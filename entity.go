package cf_speedtester

// SpeedResult
//
//	SpeedResult
//	@Description: 测试结果
type SpeedResult struct {
	IPAddress      string  `csv:"IPAddress"`
	Sent           int     `csv:"Sent"`
	Received       int     `csv:"Received"`
	PacketLossRate float64 `csv:"PacketLossRate"`
	AvgLatency     float64 `csv:"AvgLatency"`
	DownloadSpeed  float64 `csv:"DownloadSpeed"`
}

// TestOptions
//
//	TestOptions
//	@Description: 测试配置选项
type TestOptions struct {
	//延迟测速线程数
	Routines int `json:"routines"`
	//延迟测试次数
	PingTimes int `json:"pingTimes"`
	//下载测试数量
	TestCount int `json:"testCount"`
	//下载测速时间 单位 second
	DownloadTime int `json:"downloadTime"`
	//指定测速端口
	TCPPort int `json:"TCPPort"`
	//指定测速地址
	DownloadUrl string `json:"downloadUrl"`
}
