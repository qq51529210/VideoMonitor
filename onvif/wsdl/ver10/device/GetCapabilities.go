package device

import (
	"context"
	"wsdl/soap"
	"wsdl/ver10/schema"
)

// CapabilityCategory 能力
type CapabilityCategory string

// 能力列表
const (
	CapabilityCategoryAll       = "All"
	CapabilityCategoryAnalytics = "Analytics"
	CapabilityCategoryDevice    = "Device"
	CapabilityCategoryEvents    = "Events"
	CapabilityCategoryImaging   = "Imaging"
	CapabilityCategoryMedia     = "Media"
	CapabilityCategoryPTZ       = "PTZ"
)

// reqGetCapabilities 用于 GetDeviceInformation 的请求消息
type reqGetCapabilities struct {
	envelope
	Header soap.Header[soap.Security]
	Body   soap.Body[struct {
		XMLName  string               `xml:"tt:GetCapabilities"`
		Category []CapabilityCategory `xml:"tt:Category"`
	}]
}

// resGetCapabilities 用于 GetDeviceInformation 的响应消息
type resGetCapabilities struct {
	soap.Envelope
	Header soap.Header[any]
	Body   soap.Body[GetCapabilitiesResponse]
}

// GetCapabilitiesResponse 是 GetCapabilities 的结果
type GetCapabilitiesResponse struct {
	schema.Capabilities
}

// GetCapabilities
// This method has been replaced by the more generic GetServices method.
// For capabilities of individual services refer to the GetServiceCapabilities methods.
func GetCapabilities(ctx context.Context, host, username, password string, categories ...CapabilityCategory) (*GetCapabilitiesResponse, error) {
	// 消息结构
	var _req reqGetCapabilities
	_req.Init()
	_req.Header.Data.Init(username, password)
	_req.Body.Data.Category = append(_req.Body.Data.Category, categories...)
	var _res resGetCapabilities
	// 请求
	err := soap.Do(ctx, ActionURL(host, "GetCapabilities"), &_req, &_res)
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
