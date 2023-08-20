package cf_speedtester

import (
	"errors"
	"github.com/fireinrain/cf-speedtester/config"
	"github.com/fireinrain/cf-speedtester/entity"
	"github.com/fireinrain/cf-speedtester/task"
	"github.com/fireinrain/cf-speedtester/utils"
	"sync"
	"time"
)

type CFSpeedTester struct {
	Mux          sync.RWMutex
	TestOpts     entity.TestOptions
	SpeedResults []entity.SpeedResult
}

func NewCFSpeedTestClient(testOpts ...config.TestOptionFunc) *CFSpeedTester {
	options := config.NewTestOptions(testOpts...)
	speedTester := CFSpeedTester{
		Mux:          sync.RWMutex{},
		TestOpts:     entity.TestOptions{},
		SpeedResults: nil,
	}
	speedTester.Mux.RLock()
	speedTester.TestOpts = options
	speedTester.Mux.RUnlock()
	return &speedTester
}

// DoSpeedTest
//
//	@Description: 执行cloudflare ip 速度测试
//	@receiver s
func (s *CFSpeedTester) DoSpeedTest() {
	// 开始延迟测速 + 过滤延迟/丢包
	pingData := task.NewPing(&s.TestOpts).Run().FilterDelay(&s.TestOpts).
		FilterLossRate(&s.TestOpts).
		FilterIPBan(&s.TestOpts)
	// 开始下载测速
	speedData := task.TestDownloadSpeed(pingData, &s.TestOpts)
	//格式化输出结果
	speedData.PrettyPrint()
	var speedResults []entity.SpeedResult
	for _, data := range speedData {
		result := entity.SpeedResult{
			IPAddress:      data.IP.String(),
			Sent:           data.Sended,
			Received:       data.Received,
			PacketLossRate: data.GetLossRate(),
			AvgLatency:     data.Delay.Seconds() * 1000,
			DownloadSpeed:  data.DownloadSpeed,
		}
		speedResults = append(speedResults, result)
	}
	s.SpeedResults = speedResults
}

// DoSpeedTestForResult
//
//	@Description: 执行ip优选测速
//	@receiver s
//	@return []SpeedResult
func (s *CFSpeedTester) DoSpeedTestForResult() []entity.SpeedResult {
	s.DoSpeedTest()
	return s.SpeedResults
}

// ExportToCSV
//
//	@Description: 导出为csv文件
//	@receiver s
//	@param filePath
func (s *CFSpeedTester) ExportToCSV(filePath string) error {
	if s.SpeedResults == nil || len(s.SpeedResults) <= 0 {
		return errors.New("当前未进行ip测速，暂无结果导出")
	}
	var cloudflareIPDatas []utils.CloudflareIPData
	for _, result := range s.SpeedResults {
		addr := utils.IPStrToIPAddr(result.IPAddress)
		pingData := &utils.PingData{
			IP:       addr,
			Sended:   result.Sent,
			Received: result.Received,
			Delay:    time.Duration(result.AvgLatency / 1000),
		}
		data := utils.CloudflareIPData{
			PingData:      pingData,
			DownloadSpeed: result.DownloadSpeed,
		}
		cloudflareIPDatas = append(cloudflareIPDatas, data)

	}
	utils.ExportCSV(cloudflareIPDatas, filePath)
	return nil
}
