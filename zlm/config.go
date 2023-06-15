package zlm

import "net/url"

// Config 表示 zlm 的服务配置
type Config struct {
	// 是否调试http api,启用调试后，会打印每次http请求的内容和回复
	// apiDebug=1
	APIDebug string `json:"api.apiDebug"`
	// 一些比较敏感的http api在访问时需要提供secret，否则无权限调用
	// 如果是通过127.0.0.1访问,那么可以不提供secret
	// secret=035c73f7-bb6b-4889-a715-d9eb2d1925cc
	APISecret string `json:"api.secret"`
	// 截图保存路径根目录，截图通过http api(/index/api/getSnap)生成和获取
	// snapRoot=./www/snap/
	APISnapRoot string `json:"api.snapRoot"`
	// 默认截图图片，在启动FFmpeg截图后但是截图还未生成时，可以返回默认的预设图片
	// defaultSnap=./www/logo.png
	APIDefaultSnap string `json:"api.defaultSnap"`

	// [ffmpeg]

	// FFmpeg可执行程序路径,支持相对路径/绝对路径
	// bin=/usr/bin/ffmpeg
	FFMPEGBin string `json:"ffmpeg.bin"`
	// FFmpeg拉流再推流的命令模板，通过该模板可以设置再编码的一些参数
	// cmd=%s -fflags nobuffer -i %s -c:a aac -strict -2 -ar 44100 -ab 48k -c:v libx264  -f flv %s
	FFMPEGCmd string `json:"ffmpeg.cmd"`
	// FFmpeg生成截图的命令，可以通过修改该配置改变截图分辨率或质量
	// snap=%s -i %s -y -f mjpeg -t 0.001 %s
	FFMPEGSnap string `json:"ffmpeg.snap"`
	// FFmpeg日志的路径，如果置空则不生成FFmpeg日志
	// 可以为相对(相对于本可执行程序目录)或绝对路径
	// log=ffmpeg/ffmpeg.log
	FFMPEGLog string `json:"ffmpeg.log"`
	// 自动重启的时间(秒), 默认为0, 也就是不自动重启. 主要是为了避免长时间ffmpeg拉流导致的不同步现象
	// restart_sec=0
	FFMPEGRestartSec string `json:"ffmpeg.restart_sec"`
	// cmd_rtsp 是自定义的拉流命令名称
	// cmd_rtsp=%s -rtsp_transport tcp -i %s -c:v libx264 -c:a aac -strict -2 -ar 44100 -ab 48k -f flv %s
	FFMPEGCmdRTSP string `json:"ffmpeg.cmd_rtsp"`

	// [general]

	// 是否启用虚拟主机
	// enableVhost=0
	GeneralEnableVhost string `json:"general.enableVhost"`
	// 播放器或推流器在断开后会触发hook.on_flow_report事件(使用多少流量事件)，
	// flowThreshold参数控制触发hook.on_flow_report事件阈值，使用流量超过该阈值后才触发，单位KB
	// flowThreshold=1024
	GeneralFlowThreshold string `json:"general.flowThreshold"`
	// 播放最多等待时间，单位毫秒
	// 播放在播放某个流时，如果该流不存在，
	// ZLMediaKit会最多让播放器等待maxStreamWaitMS毫秒
	// 如果在这个时间内，该流注册成功，那么会立即返回播放器播放成功
	// 否则返回播放器未找到该流，该机制的目的是可以先播放再推流
	// maxStreamWaitMS=15000
	GeneralMaxStreamWaitMS string `json:"general.maxStreamWaitMS"`
	// 某个流无人观看时，触发hook.on_stream_none_reader事件的最大等待时间，单位毫秒
	// 在配合hook.on_stream_none_reader事件时，可以做到无人观看自动停止拉流或停止接收推流
	// streamNoneReaderDelayMS=20000
	GeneralStreamNoneReaderDelayMS string `json:"general.streamNoneReaderDelayMS"`
	// 是否全局添加静音aac音频，转协议时有效
	// 有些播放器在打开单视频流时不能秒开，添加静音音频可以加快秒开速度
	// addMuteAudio=1
	GeneralAddMuteAudio string `json:"general.addMuteAudio"`
	// 拉流代理时如果断流再重连成功是否删除前一次的媒体流数据，如果删除将重新开始，
	// 如果不删除将会接着上一次的数据继续写(录制hls/mp4时会继续在前一个文件后面写)
	// resetWhenRePlay=1
	GeneralResetWhenRePlay string `json:"general.resetWhenRePlay"`
	// 是否默认推流时转换成hls，hook接口(on_publish)中可以覆盖该设置
	// publishToHls=1
	GeneralPublishToHls string `json:"general.publishToHls"`
	// 是否默认推流时转换成mp4，hook接口(on_publish)中可以覆盖该设置
	// publishToMP4=0
	GeneralPublishToMP4 string `json:"general.publishToMP4"`
	// 合并写缓存大小(单位毫秒)，合并写指服务器缓存一定的数据后才会一次性写入socket，这样能提高性能，但是会提高延时
	// 开启后会同时关闭TCP_NODELAY并开启MSG_MORE
	// mergeWriteMS=0
	GeneralMergeWriteMS string `json:"general.mergeWriteMS"`
	// 全局的时间戳覆盖开关，在转协议时，对frame进行时间戳覆盖
	// 该开关对rtsp/rtmp/rtp推流、rtsp/rtmp/hls拉流代理转协议时生效
	// 会直接影响rtsp/rtmp/hls/mp4/flv等协议的时间戳
	// 同协议情况下不影响(例如rtsp/rtmp推流，那么播放rtsp/rtmp时不会影响时间戳)
	// modifyStamp=0
	GeneralModifyStamp string `json:"general.modifyStamp"`
	// 服务器唯一id，用于触发hook时区别是哪台服务器
	// mediaServerId=your_server_id
	GeneralMediaServerID string `json:"general.mediaServerId"`
	// 转协议是否全局开启或关闭音频
	// enable_audio=1
	GeneralEnableAudio string `json:"general.enable_audio"`
	// hls协议是否按需生成，如果hls.segNum配置为0(意味着hls录制)，那么hls将一直生成(不管此开关)
	// hls_demand=0
	GeneralHLSDemand string `json:"general.hls_demand"`
	// rtsp[s]协议是否按需生成
	// rtsp_demand=0
	GeneralRTSPDemand string `json:"general.rtsp_demand"`
	// rtmp[s]、http[s]-flv、ws[s]-flv协议是否按需生成
	// rtmp_demand=0
	GeneralRTMPDemand string `json:"general.rtmp_demand"`
	// http[s]-ts协议是否按需生成
	// ts_demand=0
	GeneralTSDemand string `json:"general.ts_demand"`
	// http[s]-fmp4协议是否按需生成
	// fmp4_demand=0
	GeneralFMP4Demand string `json:"general.fmp4_demand"`
	// 最多等待未初始化的Track时间，单位毫秒，超时之后会忽略未初始化的Track
	// wait_track_ready_ms=10000
	GeneralWaitTrackReadyMS string `json:"general.wait_track_ready_ms"`
	// 如果流只有单Track，最多等待若干毫秒，超时后未收到其他Track的数据，则认为是单Track
	// 如果协议元数据有声明特定track数，那么无此等待时间
	// wait_add_track_ms=3000
	GeneralWaitAddTrackMS string `json:"general.wait_add_track_ms"`
	// 如果track未就绪，我们先缓存帧数据，但是有最大个数限制，防止内存溢出
	// unready_frame_cache=100
	GeneralUnreadyFrameCache string `json:"general.unready_frame_cache"`
	// 推流断开后可以在超时时间内重新连接上继续推流，这样播放器会接着播放。
	// 置0关闭此特性(推流断开会导致立即断开播放器)
	// 此参数不应大于播放器超时时间
	// continue_push_ms=3000
	GeneralContinuePushMS string `json:"general.continue_push_ms"`

	// [hls]

	// hls写文件的buf大小，调整参数可以提高文件io性能
	// fileBufSize=65536
	HLSFileBufSize string `json:"hls.fileBufSize"`
	// hls保存文件路径
	// 可以为相对(相对于本可执行程序目录)或绝对路径
	// filePath=./www
	HLSFilePath string `json:"hls.filePath"`
	// hls最大切片时间
	// segDur=2
	HLSSegDur string `json:"hls.segDur"`
	// m3u8索引中,hls保留切片个数(实际保留切片个数大2~3个)
	// 如果设置为0，则不删除切片，而是保存为点播
	// segNum=3
	HLSSegNum string `json:"hls.segNum"`
	// HLS切片从m3u8文件中移除后，继续保留在磁盘上的个数
	// segRetain=5
	HLSSegRetain string `json:"hls.segRetain"`
	// 是否广播 ts 切片完成通知
	// broadcastRecordTs=0
	HLSBroadcastRecordTs string `json:"hls.broadcastRecordTs"`
	// 直播hls文件删除延时，单位秒，issue: 913
	// deleteDelaySec=0
	HLSDeleteDelaySec string `json:"hls.deleteDelaySec"`
	// 是否保留hls文件，此功能部分等效于segNum=0的情况
	// 不同的是这个保留不会在m3u8文件中体现
	// 0为不保留，不起作用
	// 1为保留，则不删除hls文件，如果开启此功能，注意磁盘大小，或者定期手动清理hls文件
	// segKeep=0
	HLSSegKeep string `json:"hls.segKeep"`

	// [hook]

	// 在推流时，如果url参数匹对admin_params，那么可以不经过hook鉴权直接推流成功，播放时亦然
	// 该配置项的目的是为了开发者自己调试测试，该参数暴露后会有泄露隐私的安全隐患
	// admin_params=secret=035c73f7-bb6b-4889-a715-d9eb2d1925cc
	HookAdminParams string `json:"hook.admin_params"`
	// 播放器或推流器使用流量事件，置空则关闭
	// on_flow_report=https://127.0.0.1/index/hook/on_flow_report
	HookOnFlowReport string `json:"hook.on_flow_report"`
	// 访问http文件鉴权事件，置空则关闭鉴权
	// on_http_access=https://127.0.0.1/index/hook/on_http_access
	HookOnHTTPAccess string `json:"hook.on_http_access"`
	// 播放鉴权事件，置空则关闭鉴权
	// on_play=https://127.0.0.1/index/hook/on_play
	HookOnPlay string `json:"hook.on_play"`
	// 推流鉴权事件，置空则关闭鉴权
	// on_publish=https://127.0.0.1/index/hook/on_publish
	HookOnPublish string `json:"hook.on_publish"`
	// 录制mp4切片完成事件
	// on_record_mp4=https://127.0.0.1/index/hook/on_record_mp4
	HookOnRecordMP4 string `json:"hook.on_record_mp4"`
	// 录制 hls ts 切片完成事件
	// on_record_ts=https://127.0.0.1/index/hook/on_record_ts
	HookOnRecordTS string `json:"hook.on_record_ts"`
	// rtsp播放鉴权事件，此事件中比对rtsp的用户名密码
	// on_rtsp_auth=https://127.0.0.1/index/hook/on_rtsp_auth
	HookOnRTSPAuth string `json:"hook.on_rtsp_auth"`
	// tsp播放是否开启专属鉴权事件，置空则关闭rtsp鉴权。rtsp播放鉴权还支持url方式鉴权
	// 建议开发者统一采用url参数方式鉴权，rtsp用户名密码鉴权一般在设备上用的比较多
	// 开启rtsp专属鉴权后，将不再触发on_play鉴权事件
	// on_rtsp_realm=https://127.0.0.1/index/hook/on_rtsp_realm
	HookOnRTSPRealm string `json:"hook.on_rtsp_realm"`
	// 远程telnet调试鉴权事件
	// on_shell_login=https://127.0.0.1/index/hook/on_shell_login
	HookOnShellLogin string `json:"hook.on_shell_login"`
	// 直播流注册或注销事件
	// on_stream_changed=https://127.0.0.1/index/hook/on_stream_changed
	HookOnStreamChanged string `json:"hook.on_stream_changed"`
	// 无人观看流事件，通过该事件，可以选择是否关闭无人观看的流。配合general.streamNoneReaderDelayMS选项一起使用
	// on_stream_none_reader=https://127.0.0.1/index/hook/on_stream_none_reader
	HookOnStreamNoneReader string `json:"hook.on_stream_none_reader"`
	// 播放时，未找到流事件，通过配合hook.on_stream_none_reader事件可以完成按需拉流
	// on_stream_not_found=https://127.0.0.1/index/hook/on_stream_not_found
	HookOnStreamNotFound string `json:"hook.on_stream_not_found"`
	// 服务器启动报告，可以用于服务器的崩溃重启事件监听
	// on_server_started=https://127.0.0.1/index/hook/on_server_started
	HookOnServerStarted string `json:"hook.on_server_started"`
	// server保活上报
	// on_server_keepalive=https://127.0.0.1/index/hook/on_server_keepalive
	HookOnServerKeepalive string `json:"hook.on_server_keepalive"`
	// 发送rtp(startSendRtp)被动关闭时回调
	// on_send_rtp_stopped=https://127.0.0.1/index/hook/on_send_rtp_stopped
	HookOnSendRTPStopped string `json:"hook.on_send_rtp_stopped"`
	// hook api最大等待回复时间，单位秒
	// timeoutSec=10
	HookTimeoutSec string `json:"hook.timeoutSec"`
	// keepalive hook触发间隔,单位秒，float类型
	// alive_interval=10.0
	HookAliveInterval string `json:"hook.alive_interval"`
	// hook通知失败重试次数,正整数。为0不重试，1时重试一次，以此类推
	// retry=1
	HookRetry string `json:"hook.retry"`
	// hook通知失败重试延时，单位秒，float型
	// retry_delay=3.0
	HookRetryDelay string `json:"hook.retry_delay"`

	// [cluster]

	// 设置源站拉流url模板, 格式跟printf类似，第一个%s指定app,第二个%s指定stream_id,
	// 开启集群模式后，on_stream_not_found和on_stream_none_reader hook将无效.
	// 溯源模式支持以下类型:
	// rtmp方式: rtmp://127.0.0.1:1935/%s/%s
	// rtsp方式: rtsp://127.0.0.1:554/%s/%s
	// hls方式: http://127.0.0.1:80/%s/%s/hls.m3u8
	// http-ts方式: http://127.0.0.1:80/%s/%s.live.ts
	// 支持多个源站，不同源站通过分号(;)分隔
	// origin_url=
	ClusterOriginURL string `json:"cluster.origin_url"`
	// 溯源总超时时长，单位秒，float型；假如源站有3个，那么单次溯源超时时间为timeout_sec除以3
	// 单次溯源超时时间不要超过general.maxStreamWaitMS配置
	// timeout_sec=15
	ClusterTimeoutSec string `json:"cluster.timeout_sec"`
	// 溯源失败尝试次数，-1时永久尝试
	// retry_count=3
	ClusterRetryCount string `json:"cluster.retry_count"`

	// [http]

	// http服务器字符编码，windows上默认gb2312
	// charSet=utf-8
	HTTPCharSet string `json:"http.charSet"`
	// http链接超时时间
	// keepAliveSecond=30
	HTTPKeepaliveSecond string `json:"http.keepAliveSecond"`
	// http请求体最大字节数，如果post的body太大，则不适合缓存body在内存
	// maxReqSize=40960
	HTTPMaxReqSize string `json:"http.maxReqSize"`
	// 404网页内容，用户可以自定义404网页
	// notFound=<html><head><title>404 Not Found</title></head><body bgcolor="white"><center><h1>您访问的资源不存在！</h1></center><hr><center>ZLMediaKit-4.0</center></body></html>
	HTTPNotFound string `json:"http.notFound"`
	// http服务器监听端口
	// port=80
	HTTPPort string `json:"http.port"`
	// http文件服务器根目录
	// 可以为相对(相对于本可执行程序目录)或绝对路径
	// rootPath=./www
	HTTPRootPath string `json:"http.rootPath"`
	// http文件服务器读文件缓存大小，单位BYTE，调整该参数可以优化文件io性能
	// sendBufSize=65536
	HTTPSendBufSize string `json:"http.sendBufSize"`
	// https服务器监听端口
	// sslport=443
	HTTPSSLPort string `json:"http.sslport"`
	// 是否显示文件夹菜单，开启后可以浏览文件夹
	// dirMenu=1
	HTTPDirMenu string `json:"http.dirMenu"`
	// 虚拟目录, 虚拟目录名和文件路径使用","隔开，多个配置路径间用";"隔开
	// 例如赋值为 app_a,/path/to/a;app_b,/path/to/b 那么
	// 访问 http://127.0.0.1/app_a/file_a 对应的文件路径为 /path/to/a/file_a
	// 访问 http://127.0.0.1/app_b/file_b 对应的文件路径为 /path/to/b/file_b
	// 访问其他http路径,对应的文件路径还是在rootPath内
	// virtualPath=
	HTTPVirtualPath string `json:"http.virtualPath"`
	// 禁止后缀的文件使用mmap缓存，使用“,”隔开
	// 例如赋值为 .mp4,.flv
	// 那么访问后缀为.mp4与.flv 的文件不缓存
	// forbidCacheSuffix=
	HTTPForbidCacheSuffix string `json:"http.forbidCacheSuffix"`
	// 可以把http代理前真实客户端ip放在http头中：https://github.com/ZLMediaKit/ZLMediaKit/issues/1388
	// 切勿暴露此key，否则可能导致伪造客户端ip
	// forwarded_ip_header=
	HTTPForwardedIPHeader string `json:"http.forwarded_ip_header"`

	// [multicast]

	// rtp组播截止组播ip地址
	// addrMax=239.255.255.255
	MulticastAddrMax string `json:"multicast.addrMax"`
	// rtp组播起始组播ip地址
	// addrMin=239.0.0.0
	MulticastAddrMin string `json:"multicast.addrMin"`
	// 组播udp ttl
	// udpTTL=64
	MulticastUDPTTL string `json:"multicast.udpTTL"`

	// [record]

	// mp4录制或mp4点播的应用名，通过限制应用名，可以防止随意点播
	// 点播的文件必须放置在此文件夹下
	// appName=record
	RecordAppName string `json:"record.appName"`
	// mp4录制写文件缓存，单位BYTE,调整参数可以提高文件io性能
	// fileBufSize=65536
	RecordFileBufSize string `json:"record.fileBufSize"`
	// mp4录制保存、mp4点播根路径
	// 可以为相对(相对于本可执行程序目录)或绝对路径
	// filePath=./www
	RecordFilePath string `json:"record.filePath"`
	// mp4录制切片时间，单位秒
	// fileSecond=3600
	RecordFileSecond string `json:"record.fileSecond"`
	// mp4点播每次流化数据量，单位毫秒，
	// 减少该值可以让点播数据发送量更平滑，增大该值则更节省cpu资源
	// sampleMS=500
	RecordSampleMS string `json:"record.sampleMS"`
	// mp4录制完成后是否进行二次关键帧索引写入头部
	// fastStart=0
	RecordFastStart string `json:"record.fastStart"`
	// MP4点播(rtsp/rtmp/http-flv/ws-flv)是否循环播放文件
	// fileRepeat=0
	RecordFileRepeat string `json:"record.fileRepeat"`
	// MP4录制是否当做播放器参与播放人数统计
	// mp4_as_player=0
	RecordMP4AsPlayer string `json:"record.mp4_as_player"`

	// [rtmp]

	// rtmp必须在此时间内完成握手，否则服务器会断开链接，单位秒
	// handshakeSecond=15
	RTMPHandshakeSecond string `json:"rtmp.handshakeSecond"`
	// rtmp超时时间，如果该时间内未收到客户端的数据，
	// 或者tcp发送缓存超过这个时间，则会断开连接，单位秒
	// keepAliveSecond=15
	RTMPKeepaliveSecond string `json:"rtmp.keepAliveSecond"`
	// 在接收rtmp推流时，是否重新生成时间戳(很多推流器的时间戳着实很烂)
	// modifyStamp=0
	RTMPModifyStamp string `json:"rtmp.modifyStamp"`
	// rtmp服务器监听端口
	// port=1935
	RTMPPort string `json:"rtmp.port"`
	// rtmps服务器监听地址
	// sslport=0
	RTMPSSLPort string `json:"rtmp.sslport"`

	// [rtp]

	// 音频mtu大小，该参数限制rtp最大字节数，推荐不要超过1400
	// 加大该值会明显增加直播延时
	// audioMtuSize=600
	RTPAudioMtuSize string `json:"rtp.audioMtuSize"`
	// 视频mtu大小，该参数限制rtp最大字节数，推荐不要超过1400
	// videoMtuSize=1400
	RTPVideoMtuSize string `json:"rtp.videoMtuSize"`
	// rtp包最大长度限制，单位KB,主要用于识别TCP上下文破坏时，获取到错误的rtp
	// rtpMaxSize=10
	RTPRTPMaxSize string `json:"rtp.rtpMaxSize"`
	// 导出调试数据(包括rtp/ps/h264)至该目录,置空则关闭数据导出

	// [rtp_proxy]

	// dumpDir=
	PTPProxyDumpDir string `json:"rtp_proxy.dumpDir"`
	// 导出调试数据(包括rtp/ps/h264)至该目录,置空则关闭数据导出
	// port=10000
	PTPProxyPort string `json:"rtp_proxy.port"`
	// rtp超时时间，单位秒
	// timeoutSec=15
	PTPProxyTimeoutSec string `json:"rtp_proxy.timeoutSec"`
	// 随机端口范围，最少确保36个端口
	// 该范围同时限制rtsp服务器udp端口范围
	// port_range=30000-35000
	PTPProxyPortRange string `json:"rtp_proxy.port_range"`
	// rtp h264 负载的pt
	// h264_pt=98
	PTPProxyH264PT string `json:"rtp_proxy.h264_pt"`
	// rtp h265 负载的pt
	// h265_pt=99
	PTPProxyH265PT string `json:"rtp_proxy.h265_pt"`
	// rtp ps 负载的pt
	// ps_pt=96
	PTPProxyPSPT string `json:"rtp_proxy.ps_pt"`
	// rtp ts 负载的pt
	// ts_pt=33
	PTPProxyTSPT string `json:"rtp_proxy.ts_pt"`
	// rtp opus 负载的pt
	// opus_pt=100
	PTPProxyOPUSPT string `json:"rtp_proxy.opus_pt"`
	// rtp g711u 负载的pt
	// g711u_pt=0
	PTPProxyG711UPT string `json:"rtp_proxy.g711u_pt"`
	// rtp g711a 负载的pt
	// g711a_pt=8
	PTPProxyG711APT string `json:"rtp_proxy.g711a_pt"`

	// [rtc]

	// rtc播放推流、播放超时时间
	// timeoutSec=15
	RTCTimeoutSec string `json:"rtc.timeoutSec"`
	// 本机对rtc客户端的可见ip，作为服务器时一般为公网ip，可有多个，用','分开，当置空时，会自动获取网卡ip
	// 同时支持环境变量，以$开头，如"$EXTERN_IP"; 请参考：https://github.com/ZLMediaKit/ZLMediaKit/pull/1786
	// externIP=
	RTCExternIP string `json:"rtc.externIP"`
	// rtc udp服务器监听端口号，所有rtc客户端将通过该端口传输stun/dtls/srtp/srtcp数据，
	// 该端口是多线程的，同时支持客户端网络切换导致的连接迁移
	// 需要注意的是，如果服务器在nat内，需要做端口映射时，必须确保外网映射端口跟该端口一致
	// port=8000
	RTCPort string `json:"rtc.port"`
	// 设置remb比特率，非0时关闭twcc并开启remb。该设置在rtc推流时有效，可以控制推流画质
	// 目前已经实现twcc自动调整码率，关闭remb根据真实网络状况调整码率
	// rembBitRate=0
	RTCRembBitRate string `json:"rtc.rembBitRate"`
	// rtc支持的音频codec类型,在前面的优先级更高
	// 以下范例为所有支持的音频codec
	// preferredCodecA=PCMU,PCMA,opus,mpeg4-generic
	RTCPreferredCodecA string `json:"rtc.preferredCodecA"`
	// rtc支持的视频codec类型,在前面的优先级更高
	// 以下范例为所有支持的视频codec
	// preferredCodecV=H264,H265,AV1X,VP9,VP8
	RTCPreferredCodecV string `json:"rtc.preferredCodecV"`

	// [srt]

	// srt播放推流、播放超时时间,单位秒
	// timeoutSec=5
	SRTTimeoutSec string `json:"srt.timeoutSec"`
	// srt udp服务器监听端口号，所有srt客户端将通过该端口传输srt数据，
	// 该端口是多线程的，同时支持客户端网络切换导致的连接迁移
	// port=9000
	SRTPort string `json:"srt.port"`
	// srt 协议中延迟缓存的估算参数，在握手阶段估算rtt ,然后latencyMul*rtt 为最大缓存时长，此参数越大，表示等待重传的时长就越大
	// latencyMul=4
	SRTLatencyMul string `json:"srt.latencyMul"`
	// 包缓存的大小
	// pktBufSize=8192
	SRTPktBufSize string `json:"srt.pktBufSize"`

	// [rtsp]

	// rtsp专有鉴权方式是采用base64还是md5方式
	// authBasic=0
	RTSPAuthBasic string `json:"rtsp.authBasic"`
	// rtsp拉流、推流代理是否是直接代理模式
	// 直接代理后支持任意编码格式，但是会导致GOP缓存无法定位到I帧，可能会导致开播花屏
	// 并且如果是tcp方式拉流，如果rtp大于mtu会导致无法使用udp方式代理
	// 假定您的拉流源地址不是264或265或AAC，那么你可以使用直接代理的方式来支持rtsp代理
	// 如果你是rtsp推拉流，但是webrtc播放，也建议关闭直接代理模式，
	// 因为直接代理时，rtp中可能没有sps pps,会导致webrtc无法播放; 另外webrtc也不支持Single NAL Unit Packets类型rtp
	// 默认开启rtsp直接代理，rtmp由于没有这些问题，是强制开启直接代理的
	// directProxy=1
	RTSPDirectProxy string `json:"rtsp.directProxy"`
	// rtsp必须在此时间内完成握手，否则服务器会断开链接，单位秒
	// handshakeSecond=15
	RTSPHandshakeSecond string `json:"rtsp.handshakeSecond"`
	// rtsp超时时间，如果该时间内未收到客户端的数据，
	// 或者tcp发送缓存超过这个时间，则会断开连接，单位秒
	// keepAliveSecond=15
	RTSPKeepaliveSecond string `json:"rtsp.keepAliveSecond"`
	// rtsp服务器监听地址
	// port=554
	RTSPPort string `json:"rtsp.port"`
	// rtsps服务器监听地址
	// sslport=0
	RTSPSSLPort string `json:"rtsp.sslport"`

	// [shell]

	// 调试telnet服务器接受最大bufffer大小
	// maxReqSize=1024
	ShellMaxReqSize string `json:"shell.maxReqSize"`
	// 调试telnet服务器监听端口
	// port=0
	ShellPort string `json:"shell.port"`
}

func (c *Config) toQuery() url.Values {
	m := make(url.Values)
	if c.APIDebug != "" {
		m.Set("api.apiDebug", c.APIDebug)
	}
	if c.APISecret != "" {
		m.Set("api.secret", c.APISecret)
	}
	if c.APISnapRoot != "" {
		m.Set("api.snapRoot", c.APISnapRoot)
	}
	if c.APIDefaultSnap != "" {
		m.Set("api.defaultSnap", c.APIDefaultSnap)
	}

	if c.FFMPEGBin != "" {
		m.Set("ffmpeg.bin", c.FFMPEGBin)
	}
	if c.FFMPEGCmd != "" {
		m.Set("ffmpeg.cmd", c.FFMPEGCmd)
	}
	if c.FFMPEGSnap != "" {
		m.Set("ffmpeg.snap", c.FFMPEGSnap)
	}
	if c.FFMPEGLog != "" {
		m.Set("ffmpeg.log", c.FFMPEGLog)
	}
	if c.FFMPEGRestartSec != "" {
		m.Set("ffmpeg.restart_sec", c.FFMPEGRestartSec)
	}
	if c.FFMPEGCmdRTSP != "" {
		m.Set("ffmpeg.cmd_rtsp", c.FFMPEGCmdRTSP)
	}

	if c.GeneralEnableVhost != "" {
		m.Set("general.enableVhost", c.GeneralEnableVhost)
	}
	if c.GeneralFlowThreshold != "" {
		m.Set("general.flowThreshold", c.GeneralFlowThreshold)
	}
	if c.GeneralMaxStreamWaitMS != "" {
		m.Set("general.maxStreamWaitMS", c.GeneralMaxStreamWaitMS)
	}
	if c.GeneralStreamNoneReaderDelayMS != "" {
		m.Set("general.streamNoneReaderDelayMS", c.GeneralStreamNoneReaderDelayMS)
	}
	if c.GeneralAddMuteAudio != "" {
		m.Set("general.addMuteAudio", c.GeneralAddMuteAudio)
	}
	if c.GeneralResetWhenRePlay != "" {
		m.Set("general.resetWhenRePlay", c.GeneralResetWhenRePlay)
	}
	if c.GeneralPublishToHls != "" {
		m.Set("general.publishToHls", c.GeneralPublishToHls)
	}
	if c.GeneralPublishToMP4 != "" {
		m.Set("general.publishToMP4", c.GeneralPublishToMP4)
	}
	if c.GeneralMergeWriteMS != "" {
		m.Set("general.mergeWriteMS", c.GeneralMergeWriteMS)
	}
	if c.GeneralModifyStamp != "" {
		m.Set("general.modifyStamp", c.GeneralModifyStamp)
	}
	if c.GeneralMediaServerID != "" {
		m.Set("general.mediaServerId", c.GeneralMediaServerID)
	}
	if c.GeneralEnableAudio != "" {
		m.Set("general.enable_audio", c.GeneralEnableAudio)
	}
	if c.GeneralHLSDemand != "" {
		m.Set("general.hls_demand", c.GeneralHLSDemand)
	}
	if c.GeneralRTSPDemand != "" {
		m.Set("general.rtsp_demand", c.GeneralRTSPDemand)
	}
	if c.GeneralRTMPDemand != "" {
		m.Set("general.rtmp_demand", c.GeneralRTMPDemand)
	}
	if c.GeneralTSDemand != "" {
		m.Set("general.ts_demand", c.GeneralTSDemand)
	}
	if c.GeneralFMP4Demand != "" {
		m.Set("general.fmp4_demand", c.GeneralFMP4Demand)
	}
	if c.GeneralWaitTrackReadyMS != "" {
		m.Set("general.wait_track_ready_ms", c.GeneralWaitTrackReadyMS)
	}
	if c.GeneralWaitAddTrackMS != "" {
		m.Set("general.wait_add_track_ms", c.GeneralWaitAddTrackMS)
	}
	if c.GeneralUnreadyFrameCache != "" {
		m.Set("general.unready_frame_cache", c.GeneralUnreadyFrameCache)
	}
	if c.GeneralContinuePushMS != "" {
		m.Set("general.continue_push_ms", c.GeneralContinuePushMS)
	}

	if c.HLSFileBufSize != "" {
		m.Set("general.fileBufSize", c.HLSFileBufSize)
	}
	if c.HLSFilePath != "" {
		m.Set("general.filePath", c.HLSFilePath)
	}
	if c.HLSSegDur != "" {
		m.Set("general.segDur", c.HLSSegDur)
	}
	if c.HLSSegNum != "" {
		m.Set("general.segNum", c.HLSSegNum)
	}
	if c.HLSSegRetain != "" {
		m.Set("general.segRetain", c.HLSSegRetain)
	}
	if c.HLSBroadcastRecordTs != "" {
		m.Set("general.broadcastRecordTs", c.HLSBroadcastRecordTs)
	}
	if c.HLSDeleteDelaySec != "" {
		m.Set("general.deleteDelaySec", c.HLSDeleteDelaySec)
	}
	if c.HLSSegKeep != "" {
		m.Set("general.segKeep", c.HLSSegKeep)
	}

	if c.HookAdminParams != "" {
		m.Set("hook.admin_params", c.HookAdminParams)
	}
	if c.HookOnFlowReport != "" {
		m.Set("hook.on_flow_report", c.HookOnFlowReport)
	}
	if c.HookOnHTTPAccess != "" {
		m.Set("hook.on_http_access", c.HookOnHTTPAccess)
	}
	if c.HookOnPlay != "" {
		m.Set("hook.on_play", c.HookOnPlay)
	}
	if c.HookOnPublish != "" {
		m.Set("hook.on_publish", c.HookOnPublish)
	}
	if c.HookOnRecordMP4 != "" {
		m.Set("hook.on_record_mp4", c.HookOnRecordMP4)
	}
	if c.HookOnRecordTS != "" {
		m.Set("hook.on_record_ts", c.HookOnRecordTS)
	}
	if c.HookOnRTSPAuth != "" {
		m.Set("hook.on_rtsp_auth", c.HookOnRTSPAuth)
	}
	if c.HookOnRTSPRealm != "" {
		m.Set("hook.on_rtsp_realm", c.HookOnRTSPRealm)
	}
	if c.HookOnShellLogin != "" {
		m.Set("hook.on_shell_login", c.HookOnShellLogin)
	}
	if c.HookOnStreamChanged != "" {
		m.Set("hook.on_stream_changed", c.HookOnStreamChanged)
	}
	if c.HookOnStreamNoneReader != "" {
		m.Set("hook.on_stream_none_reader", c.HookOnStreamNoneReader)
	}
	if c.HookOnStreamNotFound != "" {
		m.Set("hook.on_stream_not_found", c.HookOnStreamNotFound)
	}
	if c.HookOnServerStarted != "" {
		m.Set("hook.on_server_started", c.HookOnServerStarted)
	}
	if c.HookOnServerKeepalive != "" {
		m.Set("hook.on_server_keepalive", c.HookOnServerKeepalive)
	}
	if c.HookOnSendRTPStopped != "" {
		m.Set("hook.on_send_rtp_stopped", c.HookOnSendRTPStopped)
	}
	if c.HookTimeoutSec != "" {
		m.Set("hook.timeoutSec", c.HookTimeoutSec)
	}
	if c.HookAliveInterval != "" {
		m.Set("hook.alive_interval", c.HookAliveInterval)
	}
	if c.HookRetry != "" {
		m.Set("hook.retry", c.HookRetry)
	}
	if c.HookRetryDelay != "" {
		m.Set("hook.retry_delay", c.HookRetryDelay)
	}

	if c.ClusterOriginURL != "" {
		m.Set("cluster.origin_url", c.ClusterOriginURL)
	}
	if c.ClusterTimeoutSec != "" {
		m.Set("cluster.timeout_sec", c.ClusterTimeoutSec)
	}
	if c.ClusterRetryCount != "" {
		m.Set("cluster.retry_count", c.ClusterRetryCount)
	}

	if c.HTTPCharSet != "" {
		m.Set("http.charSet", c.HTTPCharSet)
	}
	if c.HTTPKeepaliveSecond != "" {
		m.Set("http.keepAliveSecond", c.HTTPKeepaliveSecond)
	}
	if c.HTTPMaxReqSize != "" {
		m.Set("http.maxReqSize", c.HTTPMaxReqSize)
	}
	if c.HTTPNotFound != "" {
		m.Set("http.notFound", c.HTTPNotFound)
	}
	if c.HTTPPort != "" {
		m.Set("http.port", c.HTTPPort)
	}
	if c.HTTPRootPath != "" {
		m.Set("http.rootPath", c.HTTPRootPath)
	}
	if c.HTTPSendBufSize != "" {
		m.Set("http.sendBufSize", c.HTTPSendBufSize)
	}
	if c.HTTPSSLPort != "" {
		m.Set("http.sslport", c.HTTPSSLPort)
	}
	if c.HTTPDirMenu != "" {
		m.Set("http.dirMenu", c.HTTPDirMenu)
	}
	if c.HTTPVirtualPath != "" {
		m.Set("http.virtualPath", c.HTTPVirtualPath)
	}
	if c.HTTPForbidCacheSuffix != "" {
		m.Set("http.forbidCacheSuffix", c.HTTPForbidCacheSuffix)
	}
	if c.HTTPForwardedIPHeader != "" {
		m.Set("http.charSet", c.HTTPForwardedIPHeader)
	}

	if c.MulticastAddrMax != "" {
		m.Set("multicast.addrMax", c.MulticastAddrMax)
	}
	if c.MulticastAddrMin != "" {
		m.Set("multicast.addrMin", c.MulticastAddrMin)
	}
	if c.MulticastUDPTTL != "" {
		m.Set("multicast.udpTTL", c.MulticastUDPTTL)
	}

	if c.RecordAppName != "" {
		m.Set("record.appName", c.RecordAppName)
	}
	if c.RecordFileBufSize != "" {
		m.Set("record.fileBufSize", c.RecordFileBufSize)
	}
	if c.RecordFilePath != "" {
		m.Set("record.filePath", c.RecordFilePath)
	}
	if c.RecordFileSecond != "" {
		m.Set("record.fileSecond", c.RecordFileSecond)
	}
	if c.RecordSampleMS != "" {
		m.Set("record.sampleMS", c.RecordSampleMS)
	}
	if c.RecordFastStart != "" {
		m.Set("record.fastStart", c.RecordFastStart)
	}
	if c.RecordFileRepeat != "" {
		m.Set("record.fileRepeat", c.RecordFileRepeat)
	}
	if c.RecordMP4AsPlayer != "" {
		m.Set("record.mp4_as_player", c.RecordMP4AsPlayer)
	}

	if c.RTMPHandshakeSecond != "" {
		m.Set("rtmp.handshakeSecond", c.RTMPHandshakeSecond)
	}
	if c.RTMPKeepaliveSecond != "" {
		m.Set("rtmp.keepAliveSecond", c.RTMPKeepaliveSecond)
	}
	if c.RTMPModifyStamp != "" {
		m.Set("rtmp.modifyStamp", c.RTMPModifyStamp)
	}
	if c.RTMPPort != "" {
		m.Set("rtmp.port", c.RTMPPort)
	}
	if c.RTMPSSLPort != "" {
		m.Set("rtmp.sslport", c.RTMPSSLPort)
	}

	if c.RTPAudioMtuSize != "" {
		m.Set("rtp.audioMtuSize", c.RTPAudioMtuSize)
	}
	if c.RTPVideoMtuSize != "" {
		m.Set("rtp.videoMtuSize", c.RTPVideoMtuSize)
	}
	if c.RTPRTPMaxSize != "" {
		m.Set("rtp.rtpMaxSize", c.RTPRTPMaxSize)
	}

	if c.PTPProxyDumpDir != "" {
		m.Set("rtp_proxy.dumpDir", c.PTPProxyDumpDir)
	}
	if c.PTPProxyPort != "" {
		m.Set("rtp_proxy.port", c.PTPProxyPort)
	}
	if c.PTPProxyTimeoutSec != "" {
		m.Set("rtp_proxy.timeoutSec", c.PTPProxyTimeoutSec)
	}
	if c.PTPProxyPortRange != "" {
		m.Set("rtp_proxy.port_range", c.PTPProxyPortRange)
	}
	if c.PTPProxyH264PT != "" {
		m.Set("rtp_proxy.h264_pt", c.PTPProxyH264PT)
	}
	if c.PTPProxyH265PT != "" {
		m.Set("rtp_proxy.h265_pt", c.PTPProxyH265PT)
	}
	if c.PTPProxyPSPT != "" {
		m.Set("rtp_proxy.ps_pt", c.PTPProxyPSPT)
	}
	if c.PTPProxyTSPT != "" {
		m.Set("rtp_proxy.ts_pt", c.PTPProxyTSPT)
	}
	if c.PTPProxyOPUSPT != "" {
		m.Set("rtp_proxy.opus_pt", c.PTPProxyOPUSPT)
	}
	if c.PTPProxyG711UPT != "" {
		m.Set("rtp_proxy.g711u_pt", c.PTPProxyG711UPT)
	}
	if c.PTPProxyG711APT != "" {
		m.Set("rtp_proxy.g711a_pt", c.PTPProxyG711APT)
	}

	if c.RTCTimeoutSec != "" {
		m.Set("rtc.timeoutSec", c.RTCTimeoutSec)
	}
	if c.RTCExternIP != "" {
		m.Set("rtc.externIP", c.RTCExternIP)
	}
	if c.RTCPort != "" {
		m.Set("rtc.port", c.RTCPort)
	}
	if c.RTCRembBitRate != "" {
		m.Set("rtc.rembBitRate", c.RTCRembBitRate)
	}
	if c.RTCPreferredCodecA != "" {
		m.Set("rtc.preferredCodecA", c.RTCPreferredCodecA)
	}
	if c.RTCPreferredCodecV != "" {
		m.Set("rtc.preferredCodecV", c.RTCPreferredCodecV)
	}

	if c.SRTTimeoutSec != "" {
		m.Set("srt.timeoutSec", c.SRTTimeoutSec)
	}
	if c.SRTPort != "" {
		m.Set("srt.port", c.SRTPort)
	}
	if c.SRTLatencyMul != "" {
		m.Set("srt.latencyMul", c.SRTLatencyMul)
	}
	if c.SRTPktBufSize != "" {
		m.Set("srt.pktBufSize", c.SRTPktBufSize)
	}

	if c.RTSPAuthBasic != "" {
		m.Set("rtsp.authBasic", c.RTSPAuthBasic)
	}
	if c.RTSPDirectProxy != "" {
		m.Set("rtsp.directProxy", c.RTSPDirectProxy)
	}
	if c.RTSPHandshakeSecond != "" {
		m.Set("rtsp.handshakeSecond", c.RTSPHandshakeSecond)
	}
	if c.RTSPKeepaliveSecond != "" {
		m.Set("rtsp.keepAliveSecond", c.RTSPKeepaliveSecond)
	}
	if c.RTSPPort != "" {
		m.Set("rtsp.port", c.RTSPPort)
	}
	if c.RTSPSSLPort != "" {
		m.Set("rtsp.sslport", c.RTSPSSLPort)
	}

	if c.ShellMaxReqSize != "" {
		m.Set("shell.maxReqSize", c.ShellMaxReqSize)
	}
	if c.ShellPort != "" {
		m.Set("shell.port", c.ShellPort)
	}

	return m
}
