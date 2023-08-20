package config

import (
	"github.com/fireinrain/cf-speedtester/cdn"
	"github.com/fireinrain/cf-speedtester/entity"
	"github.com/fireinrain/cf-speedtester/utils"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	// DefaultRoutines 默认线程数
	DefaultRoutines int = 200
	// DefaultPingTimes 默认ping次数
	DefaultPingTimes int = 4
	// DefaultTestCount 下载测速节点数量
	DefaultTestCount int = 10
	// DefaultDownloadTime 下载测速时间
	DefaultDownloadTime time.Duration = 10 * time.Second
	// DefaultTimeout 下载超时时间
	DefaultTimeout = DefaultDownloadTime
	// DefaultTCPPort 测速端口
	DefaultTCPPort int = 443
	// DefaultDownloadURL 测速地址
	DefaultDownloadURL string = "https://cloudflarest.gssmc.cf/100mb.zip"
	// DefaultHttpingMode 切换测速模式
	DefaultHttpingMode bool = false
	// DefaultHttpingStatusCode 有效状态代码
	DefaultHttpingStatusCode int = 0
	// DefaultHttpingCFColo 匹配指定区域
	DefaultHttpingCFColo string = ""
	// DefaultMaxDelay 平均延迟上限
	DefaultMaxDelay time.Duration = 9999 * time.Millisecond
	// DefaultMinDelay 平均延迟下限
	DefaultMinDelay time.Duration = 0 * time.Millisecond
	// DefaultMaxLossRate 丢包几率上限
	DefaultMaxLossRate float32 = 1
	// DefaultMinSpeed 下载速度下限
	DefaultMinSpeed float64 = 0
	// DefaultIPsArev6 测试的ip是v6
	DefaultIPsArev6 bool = false
	// DefaultDisableDownload 禁用下载测速
	DefaultDisableDownload bool = false
	// DefaultTestAllIP 测试所有ip
	DefaultTestAllIP bool = false
	// DefaultEnableIPBanCheck  默认不开启ip ban检测
	DefaultEnableIPBanCheck bool = false
)

var DefaultIpBanChecker = func(some any) any {
	return some
}

// DefaultIPListForTest 测试IP列表
//var DefaultIPListForTest []*net.IPAddr

type TestOptionFunc func(*entity.TestOptions)

func NewTestOptions(opt ...TestOptionFunc) entity.TestOptions {
	opts := entity.TestOptions{}

	for _, o := range opt {
		o(&opts)
	}
	defaultValues := entity.SpeedTestDefaultValues{
		DefaultRoutines:          DefaultRoutines,
		DefaultPingTimes:         DefaultPingTimes,
		DefaultTestCount:         DefaultTestCount,
		DefaultDownloadTime:      DefaultDownloadTime,
		DefaultTimeout:           DefaultTimeout,
		DefaultTCPPort:           DefaultTCPPort,
		DefaultDownloadURL:       DefaultDownloadURL,
		DefaultHttpingMode:       DefaultHttpingMode,
		DefaultHttpingStatusCode: DefaultHttpingStatusCode,
		DefaultHttpingCFColo:     DefaultHttpingCFColo,
		DefaultMaxDelay:          DefaultMinDelay,
		DefaultMinDelay:          DefaultMaxDelay,
		DefaultMaxLossRate:       DefaultMaxLossRate,
		DefaultMinSpeed:          DefaultMinSpeed,
		DefaultIPsArev6:          DefaultIPsArev6,
		DefaultDisableDownload:   DefaultDisableDownload,
		DefaultTestAllIP:         DefaultTestAllIP,
		DefaultEnableIPBanCheck:  DefaultEnableIPBanCheck,
		DefaultIPBanChecker:      DefaultIpBanChecker,
	}
	opts.DefaultValues = defaultValues

	//默认线程
	if opts.Routines == 0 {
		opts.Routines = opts.DefaultValues.DefaultRoutines
	}
	//默认ping 次数
	if opts.PingTimes == 0 {
		opts.PingTimes = opts.DefaultValues.DefaultPingTimes
	}
	if opts.TestCount == 0 {
		opts.TestCount = opts.DefaultValues.DefaultTestCount
	}
	if opts.DownloadTime == time.Duration(0) {
		opts.DownloadTime = opts.DefaultValues.DefaultDownloadTime
	}
	if opts.Timeout == time.Duration(0) {
		opts.Timeout = opts.DefaultValues.DefaultTimeout
	}
	if opts.TCPPort == 0 {
		opts.TCPPort = opts.DefaultValues.DefaultTCPPort
	}
	if opts.DownloadUrl == "" {
		opts.DownloadUrl = opts.DefaultValues.DefaultDownloadURL
	}
	if opts.HttpingMode == false {
		opts.HttpingMode = opts.DefaultValues.DefaultHttpingMode
	}
	if opts.HttpingStatusCode == 0 {
		opts.HttpingStatusCode = opts.DefaultValues.DefaultHttpingStatusCode
	}
	if opts.HttpingCFColo == "" {
		opts.HttpingCFColo = opts.DefaultValues.DefaultHttpingCFColo
	} else {
		tempStr := opts.HttpingCFColo
		opts.HttpingCFColoMap = MapColoMap(tempStr)
	}
	if opts.MaxDelay == time.Duration(0) {
		opts.MaxDelay = opts.DefaultValues.DefaultMaxDelay
	}
	if opts.MinDelay == time.Duration(0) {
		opts.MinDelay = opts.DefaultValues.DefaultMinDelay
	}
	if opts.MaxLossRate == 0 {
		opts.MaxLossRate = opts.DefaultValues.DefaultMaxLossRate
	}
	if opts.MinSpeed == 0 {
		opts.MinSpeed = opts.DefaultValues.DefaultMinSpeed
	}
	if opts.DisableDownload == false {
		opts.DisableDownload = opts.DefaultValues.DefaultDisableDownload
	}
	if opts.TestAllIP == false {
		opts.TestAllIP = opts.DefaultValues.DefaultTestAllIP
	}
	if opts.IPsArev6 == false {
		opts.IPsArev6 = opts.DefaultValues.DefaultIPsArev6
	}
	if len(opts.IPListForTest) <= 0 {
		if !opts.IPsArev6 {
			//初始化v4 ips
			globalCFIPs := cdn.GlobalCFIPs
			ranges := utils.LoadIPRanges(globalCFIPs.Ipv4Range, opts.TestAllIP)
			opts.IPListForTest = ranges
		} else {
			//初始化v6 ips
			globalCFIPs := cdn.GlobalCFIPs
			ranges := utils.LoadIPRanges(globalCFIPs.Ipv6Range, opts.TestAllIP)
			opts.IPListForTest = ranges
		}
	}
	if opts.EnableIPBanCheck == false {
		opts.EnableIPBanCheck = DefaultTestAllIP
	}
	if opts.IPBanChecker == nil {
		opts.IPBanChecker = DefaultIpBanChecker
	}

	return opts
}

func MapColoMap(coloStrs string) *sync.Map {
	if coloStrs == "" {
		return nil
	}
	// 将参数指定的地区三字码转为大写并格式化
	colos := strings.Split(strings.ToUpper(coloStrs), ",")
	colomap := &sync.Map{}
	for _, colo := range colos {
		colomap.Store(colo, colo)
	}
	return colomap
}

// WithRoutines
//
//	@Description: 设置延迟测速线程
//	@param routines
//	@return TestOptionFunc
func WithRoutines(routines int) TestOptionFunc {
	return func(o *entity.TestOptions) {
		o.Routines = routines
	}
}

func WithPingTimes(pingTimes int) TestOptionFunc {
	return func(o *entity.TestOptions) {
		o.PingTimes = pingTimes
	}
}

func WithTestCount(testCount int) TestOptionFunc {
	return func(o *entity.TestOptions) {
		o.TestCount = testCount
	}
}

func WithDownloadTime(downloadTime time.Duration) TestOptionFunc {
	return func(o *entity.TestOptions) {
		o.DownloadTime = downloadTime
	}
}

func WithTimeout(timeout time.Duration) TestOptionFunc {
	return func(o *entity.TestOptions) {
		o.Timeout = timeout
	}
}

func WithTCPPort(tcpPort int) TestOptionFunc {
	return func(o *entity.TestOptions) {
		o.TCPPort = tcpPort
	}
}

func WithDownloadUrl(downloadUrl string) TestOptionFunc {
	return func(o *entity.TestOptions) {
		o.DownloadUrl = downloadUrl
	}
}

func WithHttpingMode(httpingMode bool) TestOptionFunc {
	return func(o *entity.TestOptions) {
		o.HttpingMode = httpingMode
	}
}

func WithHttpingStatusCode(httpingStatusCode int) TestOptionFunc {
	return func(o *entity.TestOptions) {
		o.HttpingStatusCode = httpingStatusCode
	}
}

func WithHttpingCFColo(httpingCFColo string) TestOptionFunc {
	return func(o *entity.TestOptions) {
		o.HttpingCFColo = httpingCFColo
	}
}

//func WithHttpingCFColoMap(httpingCFColoMap *sync.Map) TestOptionFunc {
//	return func(o *entity.TestOptions) {
//		o.HttpingCFColoMap = httpingCFColoMap
//	}
//}

func WithMaxDelay(maxDelay time.Duration) TestOptionFunc {
	return func(o *entity.TestOptions) {
		o.MaxDelay = maxDelay
	}
}

func WithMinDelay(minDelay time.Duration) TestOptionFunc {
	return func(o *entity.TestOptions) {
		o.MinDelay = minDelay
	}
}

func WithMaxLossRate(maxLossRate float32) TestOptionFunc {
	return func(o *entity.TestOptions) {
		o.MaxLossRate = maxLossRate
	}
}

func WithMinSpeed(minSpeed float64) TestOptionFunc {
	return func(o *entity.TestOptions) {
		o.MinSpeed = minSpeed
	}
}

func WithIPListForTest(ipListForTest []*net.IPAddr) TestOptionFunc {
	return func(o *entity.TestOptions) {
		o.IPListForTest = ipListForTest
	}
}

func WithIPsArev6(ipsArev6 bool) TestOptionFunc {
	return func(o *entity.TestOptions) {
		o.IPsArev6 = ipsArev6
	}
}

func WithDisableDownload(disableDownload bool) TestOptionFunc {
	return func(o *entity.TestOptions) {
		o.DisableDownload = disableDownload
	}
}

func WithTestAllIP(testAllIP bool) TestOptionFunc {
	return func(o *entity.TestOptions) {
		o.TestAllIP = testAllIP
	}
}

func WithEnableIPBanCheck(enableIPBanCheck bool) TestOptionFunc {
	return func(o *entity.TestOptions) {
		o.EnableIPBanCheck = enableIPBanCheck
	}
}

func WithIPBanChecker(ipBanCheckerFunc func(some any) any) TestOptionFunc {
	return func(o *entity.TestOptions) {
		o.IPBanChecker = ipBanCheckerFunc
	}
}
