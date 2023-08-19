package utils

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

type IPRanges struct {
	Ips     []*net.IPAddr
	Mask    string
	FirstIP net.IP
	IpNet   *net.IPNet
}

func IsIPv4(ip string) bool {
	return strings.Contains(ip, ".")
}

func randIPEndWith(num byte) byte {
	if num == 0 { // 对于 /32 这种单独的 IP
		return byte(0)
	}
	return byte(rand.Intn(int(num)))
}

func newIPRanges() *IPRanges {
	return &IPRanges{
		Ips: make([]*net.IPAddr, 0),
	}
}

// 如果是单独 IP 则加上子网掩码，反之则获取子网掩码(r.Mask)
func (r *IPRanges) fixIP(ip string) string {
	// 如果不含有 '/' 则代表不是 IP 段，而是一个单独的 IP，因此需要加上 /32 /128 子网掩码
	if i := strings.IndexByte(ip, '/'); i < 0 {
		if IsIPv4(ip) {
			r.Mask = "/32"
		} else {
			r.Mask = "/128"
		}
		ip += r.Mask
	} else {
		r.Mask = ip[i:]
	}
	return ip
}

// 解析 IP 段，获得 IP、IP 范围、子网掩码
func (r *IPRanges) parseCIDR(ip string) {
	var err error
	if r.FirstIP, r.IpNet, err = net.ParseCIDR(r.fixIP(ip)); err != nil {
		log.Fatalln("ParseCIDR err", err)
	}
}

func (r *IPRanges) appendIPv4(d byte) {
	r.appendIP(net.IPv4(r.FirstIP[12], r.FirstIP[13], r.FirstIP[14], d))
}

func (r *IPRanges) appendIP(ip net.IP) {
	r.Ips = append(r.Ips, &net.IPAddr{IP: ip})
}

// 返回第四段 ip 的最小值及可用数目
func (r *IPRanges) getIPRange() (minIP, hosts byte) {
	minIP = r.FirstIP[15] & r.IpNet.Mask[3] // IP 第四段最小值

	// 根据子网掩码获取主机数量
	m := net.IPv4Mask(255, 255, 255, 255)
	for i, v := range r.IpNet.Mask {
		m[i] ^= v
	}
	total, _ := strconv.ParseInt(m.String(), 16, 32) // 总可用 IP 数
	if total > 255 {                                 // 矫正 第四段 可用 IP 数
		hosts = 255
		return
	}
	hosts = byte(total)
	return
}

func (r *IPRanges) chooseIPv4(testAll bool) {
	if r.Mask == "/32" { // 单个 IP 则无需随机，直接加入自身即可
		r.appendIP(r.FirstIP)
	} else {
		minIP, hosts := r.getIPRange()    // 返回第四段 IP 的最小值及可用数目
		for r.IpNet.Contains(r.FirstIP) { // 只要该 IP 没有超出 IP 网段范围，就继续循环随机
			if testAll { // 如果是测速全部 IP
				for i := 0; i <= int(hosts); i++ { // 遍历 IP 最后一段最小值到最大值
					r.appendIPv4(byte(i) + minIP)
				}
			} else { // 随机 IP 的最后一段 0.0.0.X
				r.appendIPv4(minIP + randIPEndWith(hosts))
			}
			r.FirstIP[14]++ // 0.0.(X+1).X
			if r.FirstIP[14] == 0 {
				r.FirstIP[13]++ // 0.(X+1).X.X
				if r.FirstIP[13] == 0 {
					r.FirstIP[12]++ // (X+1).X.X.X
				}
			}
		}
	}
}

func (r *IPRanges) chooseIPv6() {
	if r.Mask == "/128" { // 单个 IP 则无需随机，直接加入自身即可
		r.appendIP(r.FirstIP)
	} else {
		var tempIP uint8                  // 临时变量，用于记录前一位的值
		for r.IpNet.Contains(r.FirstIP) { // 只要该 IP 没有超出 IP 网段范围，就继续循环随机
			r.FirstIP[15] = randIPEndWith(255) // 随机 IP 的最后一段
			r.FirstIP[14] = randIPEndWith(255) // 随机 IP 的最后一段

			targetIP := make([]byte, len(r.FirstIP))
			copy(targetIP, r.FirstIP)
			r.appendIP(targetIP) // 加入 IP 地址池

			for i := 13; i >= 0; i-- { // 从倒数第三位开始往前随机
				tempIP = r.FirstIP[i]              // 保存前一位的值
				r.FirstIP[i] += randIPEndWith(255) // 随机 0~255，加到当前位上
				if r.FirstIP[i] >= tempIP {        // 如果当前位的值大于等于前一位的值，说明随机成功了，可以退出该循环
					break
				}
			}
		}
	}
}

// LoadIPRanges
//
//	@Description: 解析ip数据段
//	@param ipList
//	@return []*net.IPAddr
func LoadIPRanges(ipList []string, testAll bool) []*net.IPAddr {
	ranges := newIPRanges()

	for _, IP := range ipList {
		IP = strings.TrimSpace(IP) // 去除首尾的空白字符（空格、制表符、换行符等）
		if IP == "" {              // 跳过空的（即开头、结尾或连续多个 ,, 的情况）
			continue
		}
		ranges.parseCIDR(IP) // 解析 IP 段，获得 IP、IP 范围、子网掩码
		if IsIPv4(IP) {      // 生成要测速的所有 IPv4 / IPv6 地址（单个/随机/全部）
			ranges.chooseIPv4(testAll)
		} else {
			ranges.chooseIPv6()
		}
	}

	return ranges.Ips
}

func IPStrToIPAddr(ipStr string) *net.IPAddr {
	ipAddr, err := net.ResolveIPAddr("ip", ipStr)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	return ipAddr
}
