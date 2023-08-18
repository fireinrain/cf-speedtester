package cf_speedtester

import "sync"

type CFSpeedTester struct {
	Mux          sync.RWMutex
	TestOpts     TestOptions
	SpeedResults []SpeedResult
}

// DoSpeedTest
//
//	@Description: 执行cloudflare ip 速度测试
//	@receiver s
func (s *CFSpeedTester) DoSpeedTest() {

}

// DoSpeedTestForResult
//
//	@Description: 执行ip优选测速
//	@receiver s
//	@return []SpeedResult
func (s *CFSpeedTester) DoSpeedTestForResult() []SpeedResult {
	s.DoSpeedTest()
	return s.SpeedResults
}
