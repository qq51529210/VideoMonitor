package discovery

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/qq51529210/log"
	"github.com/qq51529210/uuid"
)

const (
	// 发送的消息，就差 uuid 了
	msgFmt = `<?xml version="1.0" encoding="UTF-8"?>
<Envelope xmlns="http://www.w3.org/2003/05/soap-envelope" xmlns:tds="http://www.onvif.org/ver10/device/wsdl">
  <Header xmlns:a="http://schemas.xmlsoap.org/ws/2004/08/addressing">
    <a:MessageID>uuid:%s</a:MessageID>
    <a:To>urn:schemas-xmlsoap-org:ws:2005:04:discovery</a:To>
    <a:Action>http://schemas.xmlsoap.org/ws/2005/04/discovery/Probe</a:Action>
  </Header>
  <Body>
    <Probe xmlns="http://schemas.xmlsoap.org/ws/2005/04/discovery" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
		<Types>tds:Device</Types>
    </Probe>
  </Body>
</Envelope>`
	// 用于读取消息的大小，10k 足够了
	readBufLen = 1024 * 10
)

var (
	_discover Discover
)

// Discovery 使用默认的 Discover 进行探测，直接使用这个就行了
func Discovery(multicastAddr string, interval time.Duration, cb func(addr string)) error {
	return _discover.Discovery(multicastAddr, interval, cb)
}

// Discover 用于探测
type Discover struct {
	// 用于等待读写协程退出
	wg sync.WaitGroup
	// 是否正常
	ok int32
	// 退出探测协程信号
	quit chan struct{}
}

// Discovery 启动后台协程发送和接受多播消息，探测设备
// multicastAddr 是多播地址
// interval 是发送探测消息的间隔
// cb 是探测到地址后的回调
func (d *Discover) Discovery(multicastAddr string, interval time.Duration, cb func(addr string)) error {
	// 多播地址
	mAddr, err := net.ResolveUDPAddr("udp", multicastAddr)
	if err != nil {
		return err
	}
	// 监听地址
	lAddr := &net.UDPAddr{
		IP:   mAddr.IP,
		Port: 0,
	}
	// 所有网络接口
	ifis, err := net.Interfaces()
	if err != nil {
		return err
	}
	for i := 0; i < len(ifis); i++ {
		ifi := &ifis[i]
		// 过滤接口
		if ifi.Flags&net.FlagLoopback != 0 ||
			ifi.Flags&net.FlagUp != 0 ||
			ifi.Flags&net.FlagMulticast != 0 {
			continue
		}
		// 底层连接
		conn, err := net.ListenMulticastUDP("udp", ifi, lAddr)
		if err != nil {
			return err
		}
		// 启动读写
		d.wg.Add(2)
		go d.readRoutine(conn, cb)
		go d.writeRoutine(conn, mAddr, interval)
	}
	//
	d.quit = make(chan struct{})
	//
	return nil
}

// writeRoutine 协程中发送消息
func (d *Discover) writeRoutine(conn *net.UDPConn, addr *net.UDPAddr, dur time.Duration) {
	// 计时器
	timer := time.NewTimer(0)
	defer func() {
		// 计时器
		timer.Stop()
		// 协程结束
		d.wg.Done()
	}()
	buf := bytes.NewBuffer(nil)
	for {
		select {
		case <-d.quit:
			return
		case <-timer.C:
			buf.Reset()
			fmt.Fprintf(buf, msgFmt, uuid.LowerV1())
			_, err := conn.WriteTo(buf.Bytes(), addr)
			if err != nil {
				log.Error(err)
			}
		}
		timer.Reset(dur)
	}
}

// envelope 用于解析地址
type envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		ProbeMatchs struct {
			ProbeMatch struct {
				XAddr string `xml:"XAddrs"`
			} `xml:"ProbeMatch"`
		} `xml:"ProbeMatchs"`
	}
}

// readRoutine 协程中读取消息
func (d *Discover) readRoutine(conn *net.UDPConn, fn func(addr string)) {
	defer func() {
		// 协程结束
		d.wg.Done()
	}()
	buf := make([]byte, readBufLen)
	for atomic.LoadInt32(&d.ok) == 1 {
		// 读取
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Error(err)
			continue
		}
		// 解析
		var res envelope
		err = xml.Unmarshal(buf[:n], &res)
		if err != nil {
			log.Error(err)
			continue
		}
		// 地址，应该是一个 http://host/xxx 格式
		uri, err := url.Parse(res.Body.ProbeMatchs.ProbeMatch.XAddr)
		if err != nil {
			log.Error(err)
			continue
		}
		// 回调通知
		fn(uri.Host)
	}
}
