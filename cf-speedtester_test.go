package cf_speedtester

import (
	"fmt"
	"github.com/fireinrain/cf-speedtester/config"
	"github.com/fireinrain/cf-speedtester/handler"
	"net"
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
	results := client.SpeedResults
	fmt.Println(results)
	client.ExportToCSV("results.csv")

}

func TestCFSpeedTester_DoSpeedTestForResult(t *testing.T) {
	client := NewCFSpeedTestClient(
		config.WithMaxDelay(300*time.Millisecond),
		config.WithMinSpeed(2),
		config.WithTestCount(5),
	)
	result := client.DoSpeedTestForResult()
	fmt.Println(result)
}

func TestForCustomDownloadUrl(t *testing.T) {
	client := NewCFSpeedTestClient(
		config.WithMaxDelay(300*time.Millisecond),
		config.WithMinSpeed(2),
		config.WithTestCount(5),
		config.WithDownloadUrl("https://youself-download-url.com"),
	)
	result := client.DoSpeedTestForResult()
	fmt.Println(result)
}

func TestNewCFSpeedTestClient(t *testing.T) {
	var ips = []string{
		"193.122.125.193",
		"193.122.119.93",
		"193.122.119.34",
		"193.122.108.223",
		"193.122.114.201",
		"193.122.114.63",
		"193.122.121.37",
		"193.122.113.19",
		"193.122.112.125",
		"193.122.116.161",
	}
	var ipList []*net.IPAddr
	for _, ip := range ips {
		addr := handler.IPStrToIPAddr(ip)
		ipList = append(ipList, addr)
	}

	client := NewCFSpeedTestClient(
		config.WithMaxDelay(300*time.Millisecond),
		config.WithMinSpeed(2),
		config.WithTestCount(1),
		config.WithIPListForTest(ipList),
	)
	result := client.DoSpeedTestForResult()
	fmt.Println(result)

}

func TestIPBanedInChina(t *testing.T) {
	var ips = []string{
		"193.122.125.193",
		"193.122.119.93",
		"193.122.119.34",
		"193.122.108.223",
		"193.122.114.201",
		"193.122.114.63",
		"193.122.121.37",
		"193.122.113.19",
	}
	var ipList []*net.IPAddr
	for _, ip := range ips {
		addr := handler.IPStrToIPAddr(ip)
		ipList = append(ipList, addr)
	}

	client := NewCFSpeedTestClient(
		config.WithMaxDelay(300*time.Millisecond),
		config.WithMinSpeed(2),
		config.WithTestCount(1),
		config.WithIPListForTest(ipList),
		config.WithEnableIPBanCheck(true),
		config.WithIPBanChecker(YouselfIPBanChecker),
	)
	result := client.DoSpeedTestForResult()
	fmt.Println(result)
}

func YouselfIPBanChecker(some any) any {
	//do you check logic
	return some
}

func TestWithISOIP(t *testing.T) {
	var ips = []string{
		"193.122.125.193",
		"193.122.119.93",
		"193.122.119.34",
		"193.122.108.223",
		"193.122.114.201",
		"193.122.114.63",
		"193.122.121.37",
		"193.122.113.19",
		"146.70.175.116",
	}
	var ipList []*net.IPAddr
	for _, ip := range ips {
		addr := handler.IPStrToIPAddr(ip)
		ipList = append(ipList, addr)
	}

	client := NewCFSpeedTestClient(
		config.WithMaxDelay(300*time.Millisecond),
		config.WithMinSpeed(2),
		config.WithTestCount(1),
		config.WithIPListForTest(ipList),
		// i want to get the cdn ip belongs to NL(Netherlands ISO code)
		config.WithWantedISOIP([]string{"NL"}),
		config.WithEnableIPBanCheck(true),
		config.WithIPBanChecker(YouselfIPBanChecker),
	)
	result := client.DoSpeedTestForResult()
	fmt.Println(result)
}
