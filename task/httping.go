package task

import (
	//"crypto/tls"
	//"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"time"
)

var (
	OutRegexp = regexp.MustCompile(`[A-Z]{3}`)
)

// pingReceived pingTotalTime
func (p *Ping) httping(ip *net.IPAddr) (int, time.Duration) {
	hc := http.Client{
		Timeout: time.Second * 2,
		Transport: &http.Transport{
			DialContext: getDialContext(ip, p.globalConfig),
			//TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 跳过证书验证
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // 阻止重定向
		},
	}

	// 先访问一次获得 HTTP 状态码 及 Cloudflare Colo
	{
		requ, err := http.NewRequest(http.MethodHead, p.globalConfig.DownloadUrl, nil)
		if err != nil {
			return 0, 0
		}
		requ.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")
		resp, err := hc.Do(requ)
		if err != nil {
			return 0, 0
		}
		defer resp.Body.Close()

		//fmt.Println("IP:", ip, "StatusCode:", resp.StatusCode, resp.Request.URL)
		// 如果未指定的 HTTP 状态码，或指定的状态码不合规，则默认只认为 200、301、302 才算 HTTPing 通过
		if p.globalConfig.HttpingStatusCode == 0 || p.globalConfig.HttpingStatusCode < 100 && p.globalConfig.HttpingStatusCode > 599 {
			if resp.StatusCode != 200 && resp.StatusCode != 301 && resp.StatusCode != 302 {
				return 0, 0
			}
		} else {
			if resp.StatusCode != p.globalConfig.HttpingStatusCode {
				return 0, 0
			}
		}

		io.Copy(io.Discard, resp.Body)

		// 只有指定了地区才匹配机场三字码
		if p.globalConfig.HttpingCFColo != "" {
			// 通过头部 Server 值判断是 Cloudflare 还是 AWS CloudFront 并设置 cfRay 为各自的机场三字码完整内容
			cfRay := func() string {
				if resp.Header.Get("Server") == "cloudflare" {
					return resp.Header.Get("CF-RAY") // 示例 cf-ray: 7bd32409eda7b020-SJC
				}
				return resp.Header.Get("x-amz-cf-pop") // 示例 X-Amz-Cf-Pop: SIN52-P1
			}()
			colo := p.getColo(cfRay)
			if colo == "" { // 没有匹配到三字码或不符合指定地区则直接结束该 IP 测试
				return 0, 0
			}
		}

	}

	// 循环测速计算延迟
	success := 0
	var delay time.Duration
	for i := 0; i < p.globalConfig.PingTimes; i++ {
		requ, err := http.NewRequest(http.MethodHead, p.globalConfig.DownloadUrl, nil)
		if err != nil {
			log.Println("意外的错误：", err)
			return 0, 0
		}
		requ.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")
		if i == p.globalConfig.PingTimes-1 {
			requ.Header.Set("Connection", "close")
		}
		startTime := time.Now()
		resp, err := hc.Do(requ)
		if err != nil {
			continue
		}
		success++
		io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
		duration := time.Since(startTime)
		delay += duration

	}

	return success, delay

}

func (p *Ping) getColo(b string) string {
	if b == "" {
		return ""
	}
	// 正则匹配并返回 机场三字码
	out := OutRegexp.FindString(b)

	if p.globalConfig.HttpingCFColoMap == nil {
		return out
	}
	// 匹配 机场三字码 是否为指定的地区
	_, ok := p.globalConfig.HttpingCFColoMap.Load(out)
	if ok {
		return out
	}

	return ""
}
