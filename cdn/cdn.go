package cdn

import (
	"context"
	"github.com/carlmjohnson/requests"
	"strings"
	"sync"
)

// CdnIPFetcher
//
//	CdnIPFetcher
//	@Description:
type CdnIPFetcher struct{}

const apiUrl4 = "https://www.cloudflare.com/ips-v4"
const apiUrl6 = "https://www.cloudflare.com/ips-v6"

var GlobalCdnFetcher CdnIPFetcher
var GlobalCFIPs CdnApiResponse

func init() {
	GlobalCdnFetcher = CdnIPFetcher{}
	cloudFlare := GlobalCdnFetcher.FetchCloudFlare()
	//解析
	GlobalCFIPs = cloudFlare

}

type CdnApiResponse struct {
	Ipv4Range []string `json:"ipv_4_range"`
	Ipv6Range []string `json:"ipv_6_range"`
}

// FetchCloudFlare
//
//	@Description: 获取cloudflare家的cdn ip range
//	@receiver receiver
//	@return CdnApiResponse
func (receiver CdnIPFetcher) FetchCloudFlare() CdnApiResponse {
	ctx := context.Background()
	var data4 string
	var data6 string
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		err := requests.
			URL(apiUrl4).
			ToString(&data4).
			Fetch(ctx)
		HandleError(err)
		// fmt.Printf("%v\n", data4)
	}()

	go func() {
		defer wg.Done()
		err := requests.
			URL(apiUrl6).
			ToString(&data6).
			Fetch(ctx)
		HandleError(err)
		// fmt.Printf("%v\n", data6)
	}()
	wg.Wait()
	splitIp4 := strings.Split(data4, "\n")
	splitIp6 := strings.Split(data6, "\n")
	return CdnApiResponse{
		Ipv4Range: splitIp4,
		Ipv6Range: splitIp6,
	}
}

func HandleError(err error) {
	if err != nil {
		panic("程序当前运行出错: " + err.Error())
	}
}
