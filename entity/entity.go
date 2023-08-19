package entity

import (
	"net"
	"sync"
	"time"
)

// SpeedResult
//
//	SpeedResult
//	@Description: 测试结果
type SpeedResult struct {
	IPAddress      string  `json:"IPAddress"`
	Sent           int     `json:"Sent"`
	Received       int     `json:"Received"`
	PacketLossRate float32 `json:"PacketLossRate"`
	AvgLatency     float64 `json:"AvgLatency"`
	DownloadSpeed  float64 `json:"DownloadSpeed"`
}

// TestOptions
//
//	TestOptions
//	@Description: 测试配置选项
type TestOptions struct {
	//默认值
	DefaultValues SpeedTestDefaultValues `json:"defaultValues"`
	//延迟测速线程数
	Routines int `json:"routines"`
	//延迟测试次数
	PingTimes int `json:"pingTimes"`
	//下载测试数量
	TestCount int `json:"testCount"`
	//下载测速时间 单位 second
	DownloadTime time.Duration `json:"downloadTime"`
	//下载超时时间
	Timeout time.Duration `json:"timeout"`
	//指定测速端口
	TCPPort int `json:"TCPPort"`
	//指定测速地址
	DownloadUrl string `json:"downloadUrl"`
	//切换测速模式
	HttpingMode bool `json:"httpingMode"`
	//有效状态代码
	HttpingStatusCode int `json:"httpingStatusCode"`
	//匹配指定区域
	HttpingCFColo string `json:"httpingCfColo"`
	//匹配指定区域map
	HttpingCFColoMap *sync.Map `json:"httpingCfColoMap"`

	//平均延迟上限
	MaxDelay time.Duration `json:"maxDelay"`
	//平均延迟下限
	MinDelay time.Duration `json:"minDelay"`
	//丢包几率上限
	MaxLossRate float32 `json:"maxLossRate"`
	//下载速度下限
	MinSpeed float64 `json:"minSpeed"`
	//测试的IP列表
	IPListForTest []*net.IPAddr `json:"ipListForTest"`
	//测试的是IPv6
	IPsArev6 bool `json:"IPsArev6"`
	//禁用下载测速
	DisableDownload bool `json:"disableDownload"`
	//测速全部 IP
	TestAllIP bool `json:"testAllIp"`
}

type SpeedTestDefaultValues struct {
	// DefaultRoutines 默认线程数
	DefaultRoutines int
	// DefaultPingTimes 默认ping次数
	DefaultPingTimes int
	// DefaultTestCount 下载测速节点数量
	DefaultTestCount int
	// DefaultDownloadTime 下载测速时间
	DefaultDownloadTime time.Duration
	// DefaultTimeout 下载超时时间
	DefaultTimeout time.Duration
	// DefaultTCPPort 测速端口
	DefaultTCPPort int
	// DefaultDownloadURL 测速地址
	DefaultDownloadURL string
	// DefaultHttpingMode 切换测速模式
	DefaultHttpingMode bool
	// DefaultHttpingStatusCode 有效状态代码
	DefaultHttpingStatusCode int
	// DefaultHttpingCFColo 匹配指定区域
	DefaultHttpingCFColo string
	// DefaultMaxDelay 平均延迟上限
	DefaultMaxDelay time.Duration
	// DefaultMinDelay 平均延迟下限
	DefaultMinDelay time.Duration
	// DefaultMaxLossRate 丢包几率上限
	DefaultMaxLossRate float32
	// DefaultMinSpeed 下载速度下限
	DefaultMinSpeed float64
	// DefaultIPsArev6 测试的ip是v6
	DefaultIPsArev6 bool
	// DefaultDisableDownload 禁用下载测速
	DefaultDisableDownload bool
	// DefaultTestAllIP 测试所有ip
	DefaultTestAllIP bool
}
