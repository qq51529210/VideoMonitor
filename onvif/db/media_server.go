package db

type MediaServer struct {
	// 流媒体服务的 id
	ID string `json:"-" gorm:"type:varchar(32);primaryKey"`
	// 流媒体服务访问密钥
	Secret *string `json:"-" gorm:"type:varchar(64)"`
	// 名称，方便记忆
	Name *string `json:"name" gorm:"type:varchar(32);not null"`
	// API 地址 (http|https)://ip:port
	APIBaseURL *string `json:"apiBaseURL" gorm:"type:varchar(128);not null"`
	// 外网访问的 ip ，生成播放地址时使用
	PublicIP *string `json:"publicIP" gorm:"type:varchar(40);not null"`
	// 内网访问的 ip ，生成播放地址时使用
	PrivateIP *string `json:"privateIP" gorm:"type:varchar(40)"`
	// 请求服务接口超时时间，单位，毫秒，默认 5000
	APICallTimeout *uint32 `json:"apiCallTimeout" gorm:"not null;default:5000"`
	// 心跳时间戳
	KeepaliveTime int64 `json:"keepaliveTime" gorm:"-"`
	// 是否在线，0/1
	Online *int8 `json:"online" gorm:"-"`
}
