package cf_speedtester

import (
	"github.com/fireinrain/cf-speedtester/config"
	"testing"
	"time"
)

func TestCFSpeedTester_DoSpeedTest(t *testing.T) {
	client := NewCFSpeedTestClient(
		config.WithMaxDelay(300*time.Millisecond),
		config.WithMinSpeed(2),
		config.WithTestCount(5),
	)
	client.DoSpeedTest()
	client.ExportToCSV("results.csv")

}
