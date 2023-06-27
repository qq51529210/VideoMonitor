package device

import (
	"encoding/xml"
	"os"
	"testing"
)

func print(v any) {
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", " ")
	enc.Encode(v)
	os.Stdout.WriteString("\n")
}

func Test_GetSystemDateAndTime(t *testing.T) {
	var _req reqGetSystemDateAndTime
	_req.Init()
	print(&_req)
}

func Test_GetDeviceInformation(t *testing.T) {
	var _req reqGetDeviceInformation
	_req.Init()
	print(&_req)
}

func Test_GetCapabilities(t *testing.T) {
	var _req reqGetCapabilities
	_req.Init()
	print(&_req)
}
