package task

import (
	"context"
	"fmt"
	"github.com/VividCortex/ewma"
	"github.com/fireinrain/cf-speedtester/entity"
	"github.com/fireinrain/cf-speedtester/utils"
	"io"
	"log"
	"net"
	"net/http"
	"sort"
	"time"
)

const (
	BufferSize = 1024
)

func checkDownloadDefault(globalConfig *entity.TestOptions) {
	if globalConfig.DownloadUrl == "" {
		globalConfig.DownloadUrl = globalConfig.DefaultValues.DefaultDownloadURL
	}
	if globalConfig.Timeout <= 0 {
		globalConfig.Timeout = globalConfig.DefaultValues.DefaultTimeout
	}
	if globalConfig.TestCount <= 0 {
		globalConfig.TestCount = globalConfig.DefaultValues.DefaultTestCount
	}
	if globalConfig.MinSpeed <= 0.0 {
		globalConfig.MinSpeed = globalConfig.DefaultValues.DefaultMinSpeed
	}
}

func TestDownloadSpeed(ipSet utils.PingDelaySet, globalConfig *entity.TestOptions) (speedSet utils.DownloadSpeedSet) {
	checkDownloadDefault(globalConfig)
	if globalConfig.DisableDownload {
		return utils.DownloadSpeedSet(ipSet)
	}
	if len(ipSet) <= 0 { // IP数组长度(IP数量) 大于 0 时才会继续下载测速
		log.Println("延迟测速结果 IP 数量为 0，跳过下载测速.")
		return
	}
	testNum := globalConfig.TestCount
	if len(ipSet) < globalConfig.TestCount || globalConfig.MinSpeed > 0 { // 如果IP数组长度(IP数量) 小于下载测速数量（-dn），则次数修正为IP数
		testNum = len(ipSet)
	}
	if testNum < globalConfig.TestCount {
		globalConfig.TestCount = testNum
	}
	log.Printf("开始下载测速（下载速度下限：%.2f MB/s，下载测速数量：%d，下载测速队列：%d\n", globalConfig.MinSpeed, globalConfig.TestCount, testNum)

	for i := 0; i < testNum; i++ {
		speed := downloadHandler(ipSet[i].IP, globalConfig)
		ipSet[i].DownloadSpeed = speed
		// 在每个 IP 下载测速后，以 [下载速度下限] 条件过滤结果
		if speed >= globalConfig.MinSpeed*1024*1024 {
			speedSet = append(speedSet, ipSet[i])        // 高于下载速度下限时，添加到新数组中
			if len(speedSet) == globalConfig.TestCount { // 凑够满足条件的 IP 时（下载测速数量 -dn），就跳出循环
				break
			}
		}
	}
	if len(speedSet) == 0 { // 没有符合速度限制的数据，返回所有测试数据
		speedSet = utils.DownloadSpeedSet(ipSet)
	}
	// 按速度排序
	sort.Sort(speedSet)
	return speedSet
}

func getDialContext(ip *net.IPAddr, globalConfig *entity.TestOptions) func(ctx context.Context, network, address string) (net.Conn, error) {
	var fakeSourceAddr string
	if utils.IsIPv4(ip.String()) {
		fakeSourceAddr = fmt.Sprintf("%s:%d", ip.String(), globalConfig.TCPPort)
	} else {
		fakeSourceAddr = fmt.Sprintf("[%s]:%d", ip.String(), globalConfig.TCPPort)
	}
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, network, fakeSourceAddr)
	}
}

// return download Speed
func downloadHandler(ip *net.IPAddr, globalConfig *entity.TestOptions) float64 {
	client := &http.Client{
		Transport: &http.Transport{DialContext: getDialContext(ip, globalConfig)},
		Timeout:   globalConfig.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > 10 { // 限制最多重定向 10 次
				return http.ErrUseLastResponse
			}
			if req.Header.Get("Referer") == globalConfig.DownloadUrl { // 当使用默认下载测速地址时，重定向不携带 Referer
				req.Header.Del("Referer")
			}
			return nil
		},
	}
	req, err := http.NewRequest("GET", globalConfig.DownloadUrl, nil)
	if err != nil {
		return 0.0
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")

	response, err := client.Do(req)
	if err != nil {
		return 0.0
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return 0.0
	}
	timeStart := time.Now()                        // 开始时间（当前）
	timeEnd := timeStart.Add(globalConfig.Timeout) // 加上下载测速时间得到的结束时间

	contentLength := response.ContentLength // 文件大小
	buffer := make([]byte, BufferSize)

	var (
		contentRead     int64 = 0
		timeSlice             = globalConfig.Timeout / 100
		timeCounter           = 1
		lastContentRead int64 = 0
	)

	var nextTime = timeStart.Add(timeSlice * time.Duration(timeCounter))
	e := ewma.NewMovingAverage()

	// 循环计算，如果文件下载完了（两者相等），则退出循环（终止测速）
	for contentLength != contentRead {
		currentTime := time.Now()
		if currentTime.After(nextTime) {
			timeCounter++
			nextTime = timeStart.Add(timeSlice * time.Duration(timeCounter))
			e.Add(float64(contentRead - lastContentRead))
			lastContentRead = contentRead
		}
		// 如果超出下载测速时间，则退出循环（终止测速）
		if currentTime.After(timeEnd) {
			break
		}
		bufferRead, err := response.Body.Read(buffer)
		if err != nil {
			if err != io.EOF { // 如果文件下载过程中遇到报错（如 Timeout），且并不是因为文件下载完了，则退出循环（终止测速）
				break
			} else if contentLength == -1 { // 文件下载完成 且 文件大小未知，则退出循环（终止测速），例如：https://speed.cloudflare.com/__down?bytes=200000000 这样的，如果在 10 秒内就下载完成了，会导致测速结果明显偏低甚至显示为 0.00（下载速度太快时）
				break
			}
			// 获取上个时间片
			lastTimeSlice := timeStart.Add(timeSlice * time.Duration(timeCounter-1))
			// 下载数据量 / (用当前时间 - 上个时间片/ 时间片)
			e.Add(float64(contentRead-lastContentRead) / (float64(currentTime.Sub(lastTimeSlice)) / float64(timeSlice)))
		}
		contentRead += int64(bufferRead)
	}
	return e.Value() / (globalConfig.Timeout.Seconds() / 120)
}
