package soap

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
)

type ReqHeader[Data any] struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header"`
	Data    Data
}

type ReqBody[Data any] struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	Data    Data
}

type ReqEnvelope[Header, Body any] struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Header  ReqHeader[Header]
	Body    ReqBody[Body]
}

type ResHeader[Data any] struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header"`
	Data    Data
}

type ResBody[Res any] struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	Data    Res
}

type ResEnvelope[Header, Body any] struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Header  ReqHeader[Header]
	Body    ResBody[Body]
}

// Do 发送请求, 格式化 req , 判断 status code 200 , 然后解析到 res
func Do[reqHeader, reqBody, resHeader, resBody any](ctx context.Context, url string,
	reqEnv *ReqEnvelope[reqHeader, reqBody],
	resEnv *ResEnvelope[resHeader, resBody]) error {
	// 格式化
	var data bytes.Buffer
	err := xml.NewEncoder(&data).Encode(reqEnv)
	if err != nil {
		return fmt.Errorf("encode xml data error %w", err)
	}
	// 请求
	req, err := http.NewRequest(http.MethodPost, url, &data)
	if err != nil {
		return fmt.Errorf("create request error %w", err)
	}
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req = req.WithContext(ctx)
	// 发送
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("create request error %w", err)
	}
	defer res.Body.Close()
	// 状态码
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error status code %d", res.StatusCode)
	}
	// 解析
	return xml.NewDecoder(res.Body).Decode(resEnv)
}
