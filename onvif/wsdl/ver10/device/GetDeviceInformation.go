package device

import (
	"context"
	"wsdl/soap"
)

// reqGetDeviceInformation 用于 GetDeviceInformation 的请求消息
type reqGetDeviceInformation struct {
	envelope
	// Header soap.Header[any]
	Body soap.Body[struct {
		XMLName string `xml:"tt:GetDeviceInformation"`
	}]
}

// resGetDeviceInformation 用于 GetDeviceInformation 的响应消息
type resGetDeviceInformation struct {
	soap.Envelope
	Header soap.Header[any]
	Body   soap.Body[GetDeviceInformationResponse]
}

// GetDeviceInformationResponse 是 GetDeviceInformation 的结果
type GetDeviceInformationResponse struct {
	// The manufactor of the device.
	Manufacturer string
	// The device model.
	Model string
	// The firmware version in the device.
	FirmwareVersion string
	// The serial number of the device.
	SerialNumber string
	// The hardware ID of the device.
	HardwareID string
}

// GetDeviceInformation
// This operation gets basic device information from the device.
func GetDeviceInformation(ctx context.Context, host string) (*GetDeviceInformationResponse, error) {
	// 消息结构
	var _req reqGetDeviceInformation
	_req.Init()
	var _res resGetDeviceInformation
	// 请求
	err := soap.Do(ctx, ActionURL(host, "GetDeviceInformation"), &_req, &_res)
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
