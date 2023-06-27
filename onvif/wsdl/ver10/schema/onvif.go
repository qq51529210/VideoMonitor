package schema

import "time"

type Time struct {
	// Range is 0 to 23.
	Hour int
	// Range is 0 to 59.
	Minute int
	// Range is 0 to 61 (typically 59).
	Second int
}

type Date struct {
	Year int
	// Range is 1 to 12.
	Month int
	// Range is 1 to 31.
	Day int
}

type DateTime struct {
	Time Time
	Date Date
}

// ToTime return time
func (t *DateTime) ToTime(location *time.Location) time.Time {
	return time.Date(t.Date.Year,
		time.Month(t.Date.Month),
		t.Date.Day,
		t.Time.Hour,
		t.Time.Minute,
		t.Time.Second,
		0, location)
}

// TimeZone
// The TZ format is specified by POSIX, please refer to POSIX 1003.1 section 8.3
// Example: Europe, Paris TZ=CET-1CEST,M3.5.0/2,M10.5.0/3
// CET = designation for standard time when daylight saving is not in force
// -1 = offset in hours = negative so 1 hour east of Greenwich meridian
// CEST = designation when daylight saving is in force ("Central European Summer Time")
// , = no offset number between code and comma, so default to one hour ahead for daylight saving
// M3.5.0 = when daylight saving starts = the last Sunday in March (the "5th" week means the last in the month)
// /2, = the local time when the switch occurs = 2 a.m. in this case
// M10.5.0 = when daylight saving ends = the last Sunday in October.
// /3, = the local time when the switch occurs = 3 a.m. in this case
type TimeZone struct {
	// Posix timezone string.
	TZ string `xml:"tt:TZ"`
}

// SetDateTimeType
// enumeration
type SetDateTimeType string

const (
	// SetDateTimeTypeManual
	// Indicates that the date and time are set manually.
	SetDateTimeTypeManual = "Manual"
	// SetDateTimeTypeNTP
	// Indicates that the date and time are set through NTP
	SetDateTimeTypeNTP = "NTP"
)

// SystemDateTime
// General date time inforamtion returned
// by the GetSystemDateTime method.
type SystemDateTime struct {
	// Indicates if the time is set manully or through NTP.
	DateTimeType SetDateTimeType
	// Informative indicator whether daylight savings is currently on/off.
	DaylightSavings bool
	// Timezone information in Posix format.
	TimeZone TimeZone
	// Current system date and time in UTC format.
	// This field is mandatory since version 2.0.
	UTCDateTime DateTime
	// Date and time in local format.
	LocalDateTime DateTime
	//
	Extension any
	// Extension     SystemDateTimeExtension
}

// UTC return utc time
func (t *SystemDateTime) UTC() time.Time {
	return t.UTCDateTime.ToTime(time.UTC)
}

// Local return local time
func (t *SystemDateTime) Local() time.Time {
	return t.UTCDateTime.ToTime(time.Local)
}

type Capabilities struct {
	// Analytics capabilities
	Analytics *AnalyticsCapabilities
	// Device capabilities
	Device *DeviceCapabilities
	// Events capabilities
	Events *EventCapabilities
	// Imaging capabilities
	Imaging *ImagingCapabilities
	// Media capabilities
	Media *MediaCapabilities
	// PTZ capabilities
	PTZ *PTZCapabilities
	//
	Extension *CapabilitiesExtension
}

type AnalyticsCapabilities struct {
	// Analytics service URI.
	XAddr string
	// Indicates whether or not rules are supported.
	RuleSupport bool
	// Indicates whether or not modules are supported.
	AnalyticsModuleSupport bool
}

type DeviceCapabilities struct {
	// Device service URI.
	XAddr string
	// Network capabilities.
	Network *NetworkCapabilities
	// System capabilities.
	System    *SystemCapabilities
	IO        *IOCapabilities
	Security  *SecurityCapabilities
	Extension *DeviceCapabilitiesExtension
}

type NetworkCapabilities struct {
	// Indicates whether or not IP filtering is supported.
	IPFilter bool
	// Indicates whether or not zeroconf is supported.
	ZeroConfiguration bool
	// Indicates whether or not IPv6 is supported.
	IPVersion6 bool
	// Indicates whether or not DynDNS is supported.
	DynDNS bool
	//
	Extension NetworkCapabilitiesExtension
}

type NetworkCapabilitiesExtension struct {
	//
	Dot11Configuration bool
	//
	Extension NetworkCapabilitiesExtension2
}

type NetworkCapabilitiesExtension2 any

type SystemCapabilities struct {
	// Indicates whether or not WS Discovery resolve requests are supported.
	DiscoveryResolve bool
	// Indicates whether or not WS-Discovery Bye is supported.
	DiscoveryBye bool
	// Indicates whether or not remote discovery is supported.
	RemoteDiscovery bool
	// Indicates whether or not system backup is supported.
	SystemBackup bool
	// Indicates whether or not system logging is supported.
	SystemLogging bool
	// Indicates whether or not firmware upgrade is supported.
	FirmwareUpgrade bool
	// Indicates supported ONVIF version(s).
	SupportedVersions OnvifVersion
	Extension         SystemCapabilitiesExtension
}

type OnvifVersion struct {
	// Major version number.
	Major int
	// Two digit minor version number.
	// If major version number is less than "16",
	// X.0.1 maps to "01" and X.2.1 maps to "21"
	// where X stands for Major version number.
	// Otherwise, minor number is month of release,
	// such as "06" for June.
	Minor int
}

type SystemCapabilitiesExtension struct {
	HttpFirmwareUpgrade    bool
	HttpSystemBackup       bool
	HttpSystemLogging      bool
	HttpSupportInformation bool
	Extension              SystemCapabilitiesExtension2
}

type SystemCapabilitiesExtension2 any

type IOCapabilities struct {
	// Number of input connectors.
	InputConnectors int
	// Number of relay outputs.
	RelayOutputs int
	Extension    IOCapabilitiesExtension
}

type IOCapabilitiesExtension struct {
	Auxiliary         bool
	AuxiliaryCommands string
	Extension         IOCapabilitiesExtension2
}

type IOCapabilitiesExtension2 any

type SecurityCapabilities struct {
	// Indicates whether or not TLS 1.1 is supported.
	TLS1_1 bool
	// Indicates whether or not TLS 1.2 is supported.
	TLS1_2 bool
	// Indicates whether or not onboard key generation is supported.
	OnboardKeyGeneration bool
	// Indicates whether or not access policy configuration is supported.
	AccessPolicyConfig bool
	// Indicates whether or not WS-Security X.509 token is supported.
	X509Token bool `xml:"X.509Token"`
	// Indicates whether or not WS-Security SAML token is supported.
	SAMLToken bool
	// Indicates whether or not WS-Security Kerberos token is supported.
	KerberosToken bool
	// Indicates whether or not WS-Security REL token is supported.
	RELToken  bool
	Extension SecurityCapabilitiesExtension
}

type SecurityCapabilitiesExtension struct {
	TLS1_0    bool
	Extension SecurityCapabilitiesExtension2
}

type SecurityCapabilitiesExtension2 struct {
	// EAP Methods supported by the device.
	// The int values refer to the
	// "http://www.iana.org/assignments/eap-numbers/eap-numbers.xhtml"
	// IANA EAP Registry.
	Dot1X              bool
	SupportedEAPMethod int
	RemoteUserHandling bool
}

type DeviceCapabilitiesExtension any

type EventCapabilities struct {
	// Event service URI.
	XAddr string
	// Indicates whether or not WS Subscription policy is supported.
	WSSubscriptionPolicySupport bool
	// Indicates whether or not WS Pull Point is supported.
	WSPullPointSupport bool
	// Indicates whether or not WS Pausable Subscription Manager Interface is supported.
	WSPausableSubscriptionManagerInterfaceSupport bool
}

type ImagingCapabilities struct {
	// Imaging service URI.
	XAddr string
}

type MediaCapabilities struct {
	// Media service URI.
	XAddr string
	// Streaming capabilities.
	StreamingCapabilities RealTimeStreamingCapabilities
	//
	Extension MediaCapabilitiesExtension
}

type RealTimeStreamingCapabilities struct {
	// Indicates whether or not RTP multicast is supported.
	RTPMulticast bool
	// Indicates whether or not RTP over TCP is supported.
	RTP_TCP bool
	// Indicates whether or not RTP/RTSP/TCP is supported.
	RTP_RTSP_TCP bool
	//
	Extension RealTimeStreamingCapabilitiesExtension
}

type RealTimeStreamingCapabilitiesExtension any

type MediaCapabilitiesExtension struct {
	ProfileCapabilities ProfileCapabilities
}

type ProfileCapabilities struct {
	// Maximum number of profiles.
	MaximumNumberOfProfiles int
}

type PTZCapabilities struct {
	// PTZ service URI.
	XAddr string
}

type CapabilitiesExtension struct {
	DeviceIO        *DeviceIOCapabilities
	Display         *DisplayCapabilities
	Recording       *RecordingCapabilities
	Search          *SearchCapabilities
	Replay          *ReplayCapabilities
	Receiver        *ReceiverCapabilities
	AnalyticsDevice *AnalyticsDeviceCapabilities
	Extensions      *CapabilitiesExtension2
}

type DeviceIOCapabilities struct {
	XAddr        string
	VideoSources int
	VideoOutputs int
	AudioSources int
	AudioOutputs int
	RelayOutputs int
}

type DisplayCapabilities struct {
	XAddr       string
	FixedLayout bool
}

type RecordingCapabilities struct {
	XAddr              string
	ReceiverSource     bool
	MediaProfileSource bool
	DynamicRecordings  bool
	DynamicTracks      bool
	MaxStringLength    int
}

type SearchCapabilities struct {
	XAddr          string
	MetadataSearch bool
}

type ReplayCapabilities struct {
	// The address of the replay service.
	XAddr string
}

type ReceiverCapabilities struct {
	// The address of the receiver service.
	XAddr string
	// Indicates whether the device can receive RTP multicast streams.
	RTP_Multicast bool
	// Indicates whether the device can receive RTP/TCP streams.
	RTP_TCP bool
	// Indicates whether the device can receive RTP/RTSP/TCP streams.
	RTP_RTSP_TCP bool
	// The maximum number of receivers supported by the device.
	SupportedReceivers int
	// The maximum allowed length for RTSP URIs.
	MaximumRTSPURILength int
}

type AnalyticsDeviceCapabilities struct {
	XAddr       string
	RuleSupport bool
	Extension   *AnalyticsDeviceExtension
}

type AnalyticsDeviceExtension any

type CapabilitiesExtension2 any
