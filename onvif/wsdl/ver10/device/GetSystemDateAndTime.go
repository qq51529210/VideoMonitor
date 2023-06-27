package device

import (
	"context"
	"wsdl/soap"
	"wsdl/ver10/schema"
)

// reqGetSystemDateAndTime 用于 GetSystemDateAndTime 的请求消息
type reqGetSystemDateAndTime struct {
	envelope
	// Header soap.Header[any]
	Body soap.Body[struct {
		XMLName string `xml:"tt:GetSystemDateAndTime"`
	}]
}

// resGetSystemDateAndTime 用于 GetSystemDateAndTime 的响应消息
type resGetSystemDateAndTime struct {
	soap.Envelope
	Header soap.Header[any]
	Body   soap.Body[GetSystemDateAndTimeResponse]
}

// GetSystemDateAndTimeResponse 是 GetSystemDateAndTime 的结果
type GetSystemDateAndTimeResponse struct {
	// Contains information whether system date and time
	// are set manually or by NTP, daylight savings is
	// on or off, time zone in POSIX 1003.1 format and
	// system date and time in UTC and also local system date and time.
	schema.SystemDateTime
}

// This operation gets the device system date and time.
// The device shall support the return of the daylight
// saving setting and of the manual system date and time
// (if applicable) or indication of NTP time (if applicable)
// through the GetSystemDateAndTime command.
// A device shall provide the UTCDateTime information.
func GetSystemDateAndTime(ctx context.Context, host string) (*GetSystemDateAndTimeResponse, error) {
	// 消息结构
	var _req reqGetSystemDateAndTime
	_req.Init()
	var _res resGetSystemDateAndTime
	// 请求
	err := soap.Do(ctx, ActionURL(host, "GetSystemDateAndTime"), &_req, &_res)
	if err != nil {
		return nil, err
	}
	// 错误
	if _res.Body.Fault != nil {
		return nil, _res.Body.Fault
	}
	// 成功
	return &_res.Body.Data, nil
}
