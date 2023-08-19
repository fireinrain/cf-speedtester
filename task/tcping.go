package task

import (
	"fmt"
	"github.com/fireinrain/cf-speedtester/entity"
	"github.com/fireinrain/cf-speedtester/utils"
	"log"
	"net"
	"sort"
	"sync"
	"time"
)

const (
	TcpConnectTimeout = time.Second * 1
	// MaxRoutine 最先线程数
	MaxRoutine = 1000
)

type Ping struct {
	wg           *sync.WaitGroup
	m            *sync.Mutex
	ips          []*net.IPAddr
	pingResults  utils.PingDelaySet
	control      chan bool
	globalConfig *entity.TestOptions
}

func checkPingDefault(testConfig *entity.TestOptions) {
	if testConfig.Routines <= 0 || testConfig.Routines > MaxRoutine {
		testConfig.Routines = testConfig.DefaultValues.DefaultRoutines
	}
	if testConfig.TCPPort <= 0 || testConfig.TCPPort >= 65535 {
		testConfig.TCPPort = testConfig.DefaultValues.DefaultTCPPort
	}
	if testConfig.PingTimes <= 0 {
		testConfig.PingTimes = testConfig.DefaultValues.DefaultPingTimes
	}

}

func NewPing(testConfig *entity.TestOptions) *Ping {
	checkPingDefault(testConfig)
	ips := testConfig.IPListForTest
	return &Ping{
		wg:           &sync.WaitGroup{},
		m:            &sync.Mutex{},
		ips:          ips,
		pingResults:  make(utils.PingDelaySet, 0),
		control:      make(chan bool, testConfig.Routines),
		globalConfig: testConfig,
	}
}

func (p *Ping) Run() utils.PingDelaySet {
	if len(p.ips) == 0 {
		return p.pingResults
	}
	if p.globalConfig.HttpingMode {
		p.printInitTestInfo("HTTP")
	} else {
		p.printInitTestInfo("TCP")
	}
	for _, ip := range p.ips {
		p.wg.Add(1)
		p.control <- false
		go p.start(ip)
	}
	p.wg.Wait()
	sort.Sort(p.pingResults)
	return p.pingResults
}

func (p *Ping) start(ip *net.IPAddr) {
	defer p.wg.Done()
	p.tcpingHandler(ip)
	<-p.control
}

// bool connectionSucceed float32 time
func (p *Ping) tcping(ip *net.IPAddr) (bool, time.Duration) {
	startTime := time.Now()
	var fullAddress string
	if utils.IsIPv4(ip.String()) {
		fullAddress = fmt.Sprintf("%s:%d", ip.String(), p.globalConfig.TCPPort)
	} else {
		fullAddress = fmt.Sprintf("[%s]:%d", ip.String(), p.globalConfig.TCPPort)
	}
	conn, err := net.DialTimeout("tcp", fullAddress, TcpConnectTimeout)
	if err != nil {
		return false, 0
	}
	defer conn.Close()
	duration := time.Since(startTime)
	return true, duration
}

// pingReceived pingTotalTime
func (p *Ping) checkConnection(ip *net.IPAddr) (recv int, totalDelay time.Duration) {
	if p.globalConfig.HttpingMode {
		recv, totalDelay = p.httping(ip)
		return
	}
	for i := 0; i < p.globalConfig.PingTimes; i++ {
		if ok, delay := p.tcping(ip); ok {
			recv++
			totalDelay += delay
		}
	}
	return
}

func (p *Ping) appendIPData(data *utils.PingData) {
	p.m.Lock()
	defer p.m.Unlock()
	p.pingResults = append(p.pingResults, utils.CloudflareIPData{
		PingData: data,
	})
}

// handle tcping
func (p *Ping) tcpingHandler(ip *net.IPAddr) {
	recv, totalDlay := p.checkConnection(ip)
	nowAble := len(p.pingResults)
	if recv != 0 {
		nowAble++
	}
	if recv == 0 {
		return
	}
	data := &utils.PingData{
		IP:       ip,
		Sended:   p.globalConfig.PingTimes,
		Received: recv,
		Delay:    totalDlay / time.Duration(recv),
	}
	p.appendIPData(data)
}

func (p *Ping) printInitTestInfo(httpingModeStr string) {
	log.Printf("开始延迟测速（模式：%s，端口：%d，平均延迟上限：%v ms，平均延迟下限：%v ms，丢包几率上限：%.2f )\n", httpingModeStr, p.globalConfig.TCPPort, p.globalConfig.MaxDelay.Milliseconds(), p.globalConfig.MinDelay.Milliseconds(), p.globalConfig.MaxLossRate)
}
