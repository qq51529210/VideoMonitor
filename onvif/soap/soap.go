package soap

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
)

// Do 发送请求, 格式化 req , 判断 status code 200 , 然后解析到 res
func Do[reqHeader, reqBody, resHeader, resBody any](
	ctx context.Context,
	url string,
	req *Envelope[reqHeader, reqBody],
	res *Envelope[resHeader, resBody],
) error {
	// 格式化
	var data bytes.Buffer
	err := xml.NewEncoder(&data).Encode(req)
	if err != nil {
		return fmt.Errorf("encode xml data error %w", err)
	}
	// 请求
	_req, err := http.NewRequest(http.MethodPost, url, &data)
	if err != nil {
		return fmt.Errorf("create request error %w", err)
	}
	_req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	_req = _req.WithContext(ctx)
	// 发送
	_res, err := http.DefaultClient.Do(_req)
	if err != nil {
		return fmt.Errorf("create request error %w", err)
	}
	defer _res.Body.Close()
	// 状态码
	if _res.StatusCode != http.StatusOK {
		return fmt.Errorf("error status code %d", _res.StatusCode)
	}
	// 解析
	return xml.NewDecoder(_res.Body).Decode(res)
}
