package cf_speedtester

const (
	// DefaultRoutines 默认线程数
	DefaultRoutines int = 200
	// DefaultPingTimes 默认ping次数
	DefaultPingTimes int = 4
	// DefaultTestCount 下载测速节点数量
	DefaultTestCount int = 10
	// DefaultDownloadTime 下载测速时间
	DefaultDownloadTime int = 10
	// DefaultTCPPort 测速端口
	DefaultTCPPort int = 443
	// DefaultDownloadURL 测速地址
	DefaultDownloadURL string = "https://cloudflarest.gssmc.cf/100mb.zip"
)

type TestOptionFunc func(*TestOptions)

func NewTestOptions(opt ...TestOptionFunc) TestOptions {
	opts := TestOptions{}

	for _, o := range opt {
		o(&opts)
	}

	//默认线程
	if opts.Routines == 0 {
		opts.Routines = DefaultRoutines
	}
	//默认ping 次数
	if opts.PingTimes == 0 {
		opts.PingTimes = DefaultPingTimes
	}
	if opts.TestCount == 0 {
		opts.TestCount = DefaultTestCount
	}
	if opts.DownloadTime == 0 {
		opts.DownloadTime = DefaultDownloadTime
	}
	if opts.TCPPort == 0 {
		opts.TCPPort = DefaultTCPPort
	}
	if opts.DownloadUrl == "" {
		opts.DownloadUrl = DefaultDownloadURL
	}

	return opts
}

// WithRoutines
//
//	@Description: 设置延迟测速线程
//	@param routines
//	@return TestOptionFunc
func WithRoutines(routines int) TestOptionFunc {
	return func(o *TestOptions) {
		o.Routines = routines
	}
}
