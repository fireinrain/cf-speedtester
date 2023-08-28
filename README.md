# cf-speedtester
A golang library to use cloudflare speed test function. use it in your own golang project.

# why you write this golang library
I want to use cloudflare cdn ip speed test in my own project.
there is a command-line tool called [CloudflareSpeedTest](https://github.com/XIU2/CloudflareSpeedTest), but if 
you want to use in your own project,you need wrap it as a bash call.
so, here is the project.


# how to use
## add this lib to you golang project
```bash
# in terminal
go get github.com/fireinrain/cf-speedtester@v1.0.6
# in your project
import (
	"github.com/fireinrain/cf-speedtester"
)

```

## simple use(use cloudflare official cdn ips)
```go

client := cf_speedtester.NewCFSpeedTestClient(
		config.WithMaxDelay(300*time.Millisecond),
		config.WithMinSpeed(2),
		config.WithTestCount(5),
	)
client.DoSpeedTest()
results := client.SpeedResults
fmt.Println(results)


```


## export data to cvs
```go

client := cf_speedtester.NewCFSpeedTestClient(
		config.WithMaxDelay(300*time.Millisecond),
		config.WithMinSpeed(2),
		config.WithTestCount(5),
	)
client.DoSpeedTest()
results := client.SpeedResults
fmt.Println(results)
client.ExportToCSV("results.csv")


```

## use custom download url
```go
client := NewCFSpeedTestClient(
    config.WithMaxDelay(300*time.Millisecond),
    config.WithMinSpeed(2),
    config.WithTestCount(5),
    config.WithDownloadUrl("https://youself-download-url.com"),
)
result := client.DoSpeedTestForResult()
fmt.Println(result)


```


## use with self-find ips that proxied to cloudflare cdn
```go
//replace ips that you find proxied for cloudflare cdn
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
    addr := utils.IPStrToIPAddr(ip)
    ipList = append(ipList, addr)
}

client := cf_speedtester.NewCFSpeedTestClient(
    config.WithMaxDelay(300*time.Millisecond),
    config.WithMinSpeed(2),
    config.WithTestCount(1),
    config.WithIPListForTest(ipList),
)
result := client.DoSpeedTestForResult()
fmt.Println(result)



```
## get the ips banned in mainland china?
if you find you get the cloudflare ip is banned in china with this library, take it easy and here
is the solution.

```go
//develop yourself ip ban check function, and inject to the config

func YouselfIPBanChecker(some any) any{
	//do you check logic
	//notice: you need convert any to handler.PingDelaySet
	//and with PingDelaySet,do your check logic and return
    //checked PingDelaySet
	return some
}
//here is an example of a IPBanChecker
func DoIPBanCheck(someData any) any {
    var result []handler.CloudflareIPData
    //转型
    if pingDelaySetValue, ok := someData.(handler.PingDelaySet); ok {
        log.Println("Convert someData to PingDelaySet type success,size is :", len(pingDelaySetValue))
        //do ip banned check
		//DoIPBanCheckInPool if self write check logic, replaced with your check logic
        checkerResults := DoIPBanCheckInPool(pingDelaySetValue, 3)
        for _, checkerResult := range checkerResults {
            if checkerResult.IsBanned == false {
                result = append(result, *checkerResult.CheckIPAddr)
            }
        }
        log.Println("Do ip banned check finished, result size is :", len(result))
        return result
    } else {
        log.Println("Convert someData to PingDelaySet type failed :", someData)
    }
    return someData
}
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
    addr := utils.IPStrToIPAddr(ip)
    ipList = append(ipList, addr)
}

client := cf_speedtester.NewCFSpeedTestClient(
    config.WithMaxDelay(300*time.Millisecond),
    config.WithMinSpeed(2),
    config.WithTestCount(1),
    config.WithIPListForTest(ipList),
    config.WithEnableIPBanCheck(true),
    config.WithIPBanChecker(YouselfIPBanChecker),
)
result := client.DoSpeedTestForResult()
fmt.Println(result)

```

## if i want to get specific country ip, what should i do?
```go
// you can do like this
// the lib use geoip2 golang, use Country.mmdb.
// it may not exactly, but seems work well. 

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
	// the priority is lowered sequentially
    config.WithWantedISOIP([]string{"NL","US"}),
    config.WithEnableIPBanCheck(true),
    config.WithIPBanChecker(YouselfIPBanChecker),
)
result := client.DoSpeedTestForResult()
fmt.Println(result)


```
# issues?
- download speed always 0?
you should change the custom download url, and have a test.

# special thanks 
- CloudflareSpeedTest,Thanks very much !
- Jetbrains Goland IDE, Thanks very much !





