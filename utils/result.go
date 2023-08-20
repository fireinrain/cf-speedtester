package utils

import (
	"encoding/csv"
	"github.com/fireinrain/cf-speedtester/entity"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var (
	PrintNum = 10
)

type PingData struct {
	IP       *net.IPAddr
	Sended   int
	Received int
	Delay    time.Duration
}

type CloudflareIPData struct {
	*PingData
	lossRate      float32
	DownloadSpeed float64
}

// GetLossRate
//
//	@Description: 计算丢包率
//	@receiver cf
//	@return float32
func (cf *CloudflareIPData) GetLossRate() float32 {
	if cf.lossRate == 0 {
		pingLost := cf.Sended - cf.Received
		cf.lossRate = float32(pingLost) / float32(cf.Sended)
	}
	return cf.lossRate
}

// GetDataStrList
//
//	@Description: 获取结果字符串列表
//	@receiver cf
//	@return []string
func (cf *CloudflareIPData) GetDataStrList() []string {
	result := make([]string, 6)
	result[0] = cf.IP.String()
	result[1] = strconv.Itoa(cf.Sended)
	result[2] = strconv.Itoa(cf.Received)
	result[3] = strconv.FormatFloat(float64(cf.GetLossRate()), 'f', 2, 32)
	result[4] = strconv.FormatFloat(cf.Delay.Seconds()*1000, 'f', 2, 32)
	result[5] = strconv.FormatFloat(cf.DownloadSpeed/1024/1024, 'f', 2, 32)
	return result
}

// ConvertToString
//
//	@Description: 结构体转化为字符串
//	@param data
//	@return [][]string
func ConvertToString(data []CloudflareIPData) [][]string {
	result := make([][]string, 0)
	for _, v := range data {
		result = append(result, v.GetDataStrList())
	}
	return result
}

// PingDelaySet 延迟丢包排序
type PingDelaySet []CloudflareIPData

// FilterDelay
//
//	@Description: 延迟条件过滤
//	@receiver s
//	@param globalConfig
//	@return data
func (s PingDelaySet) FilterDelay(globalConfig *entity.TestOptions) (data PingDelaySet) {
	if globalConfig.MaxDelay > globalConfig.DefaultValues.DefaultMaxDelay || globalConfig.MinDelay < globalConfig.DefaultValues.DefaultMinDelay { // 当输入的延迟条件不在默认范围内时，不进行过滤
		return s
	}
	if globalConfig.MaxDelay == globalConfig.DefaultValues.DefaultMaxDelay && globalConfig.MinDelay < globalConfig.DefaultValues.DefaultMinDelay { // 当输入的延迟条件为默认值时，不进行过滤
		return s
	}
	for _, v := range s {
		if v.Delay > globalConfig.MaxDelay { // 平均延迟上限，延迟大于条件最大值时，后面的数据都不满足条件，直接跳出循环
			break
		}
		if v.Delay < globalConfig.MinDelay { // 平均延迟下限，延迟小于条件最小值时，不满足条件，跳过
			continue
		}
		data = append(data, v) // 延迟满足条件时，添加到新数组中
	}
	return data
}

// FilterLossRate
//
//	@Description: 丢包条件过滤
//	@receiver s
//	@param globalConfig
//	@return data
func (s PingDelaySet) FilterLossRate(globalConfig *entity.TestOptions) (data PingDelaySet) {
	if globalConfig.MaxLossRate >= globalConfig.DefaultValues.DefaultMaxLossRate { // 当输入的丢包条件为默认值时，不进行过滤
		return s
	}
	for _, v := range s {
		if v.GetLossRate() > globalConfig.MaxLossRate { // 丢包几率上限
			break
		}
		data = append(data, v) // 丢包率满足条件时，添加到新数组中
	}
	return data
}

// Len 实现sort接口
func (s PingDelaySet) Len() int {
	return len(s)
}

func (s PingDelaySet) Less(i, j int) bool {
	iRate, jRate := s[i].GetLossRate(), s[j].GetLossRate()
	if iRate != jRate {
		return iRate < jRate
	}
	return s[i].Delay < s[j].Delay
}
func (s PingDelaySet) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

//////////////////////////////// 实现sort结束

// FilterIPBan
//
//	@Description: 墙内ipban 过滤
//	@receiver s
//	@param globalConfig
//	@return data
func (s PingDelaySet) FilterIPBan(globalConfig *entity.TestOptions) (data PingDelaySet) {
	if globalConfig.EnableIPBanCheck == false {
		return s
	} else {
		if globalConfig.IPBanChecker != nil {
			checker := globalConfig.IPBanChecker
			result := checker(s)
			if pingDelaySetValue, ok := result.(PingDelaySet); ok {
				log.Println("FilterIPBan values are :", pingDelaySetValue)
				return pingDelaySetValue
			} else {
				log.Println("FilterIPBan filter failed :", s, ", discard FilterIPBan and return original values")
			}
		}
	}
	return s
}

// DownloadSpeedSet 下载速度排序
type DownloadSpeedSet []CloudflareIPData

// Len
//
//	@Description: 实现sort接口
//	@receiver s
//	@return int
func (s DownloadSpeedSet) Len() int {
	return len(s)
}
func (s DownloadSpeedSet) Less(i, j int) bool {
	return s[i].DownloadSpeed > s[j].DownloadSpeed
}
func (s DownloadSpeedSet) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

/////////////////////////实现结束

// PrettyPrint
//
//	@Description: 格式化打印输出
//	@receiver s
func (s DownloadSpeedSet) PrettyPrint() {
	if len(s) <= 0 { // IP数组长度(IP数量) 大于 0 时继续
		log.Println("完整测速结果 IP 数量为 0，跳过输出结果.")
		return
	}
	dateString := ConvertToString(s) // 转为多维数组 [][]String
	if len(dateString) < PrintNum {  // 如果IP数组长度(IP数量) 小于  打印次数，则次数改为IP数量
		PrintNum = len(dateString)
	}
	headFormat := "%-16s%-5s%-5s%-5s%-6s%-11s\n"
	dataFormat := "%-18s%-8s%-8s%-8s%-10s%-15s\n"
	for i := 0; i < PrintNum; i++ { // 如果要输出的 IP 中包含 IPv6，那么就需要调整一下间隔
		if len(dateString[i][0]) > 15 {
			headFormat = "%-40s%-5s%-5s%-5s%-6s%-11s\n"
			dataFormat = "%-42s%-8s%-8s%-8s%-10s%-15s\n"
			break
		}
	}
	log.Printf(headFormat, "IP 地址", "已发送", "已接收", "丢包率", "平均延迟", "下载速度 (MB/s)")
	for i := 0; i < PrintNum; i++ {
		log.Printf(dataFormat, dateString[i][0], dateString[i][1], dateString[i][2], dateString[i][3], dateString[i][4], dateString[i][5])
	}
}

// ExportCSV
//
//	@Description: 导出结果为csv文件
//	@param data
func ExportCSV(data []CloudflareIPData, filePath string) {
	if len(data) == 0 {
		return
	}
	fp, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("创建文件[%s]失败：%v", filePath, err)
		return
	}
	defer fp.Close()
	w := csv.NewWriter(fp) //创建一个新的写入文件流
	_ = w.Write([]string{"IP 地址", "已发送", "已接收", "丢包率", "平均延迟", "下载速度 (MB/s)"})
	_ = w.WriteAll(ConvertToString(data))
	w.Flush()
}
