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

# download urls
here are some download speed test urls provided by [PencilNavigator](https://github.com/PencilNavigator), thanks for your efforts.
https://testfiles.blockly.cf/100mb.zip
https://testfiles.blockly.cf/200mb.zip
https://testfiles.blockly.cf/300mb.zip
https://testfiles.blockly.cf/400mb.zip
https://testfiles.blockly.cf/500mb.zip

https://testfiles.blockly.tk/100mb.zip
https://testfiles.blockly.tk/200mb.zip
https://testfiles.blockly.tk/300mb.zip
https://testfiles.blockly.tk/400mb.zip
https://testfiles.blockly.tk/500mb.zip

https://testfiles.blockly.gq/100mb.zip
https://testfiles.blockly.gq/200mb.zip
https://testfiles.blockly.gq/300mb.zip
https://testfiles.blockly.gq/400mb.zip
https://testfiles.blockly.gq/500mb.zip

https://testfiles.gssmc.cf/100mb.zip
https://testfiles.gssmc.cf/200mb.zip
https://testfiles.gssmc.cf/300mb.zip
https://testfiles.gssmc.cf/400mb.zip
https://testfiles.gssmc.cf/500mb.zip

https://testfiles.gssmc.tk/100mb.zip
https://testfiles.gssmc.tk/200mb.zip
https://testfiles.gssmc.tk/300mb.zip
https://testfiles.gssmc.tk/400mb.zip
https://testfiles.gssmc.tk/500mb.zip

https://testfiles.gssmc.gq/100mb.zip
https://testfiles.gssmc.gq/200mb.zip
https://testfiles.gssmc.gq/300mb.zip
https://testfiles.gssmc.gq/400mb.zip
https://testfiles.gssmc.gq/500mb.zip

https://testfiles.itkyou.cf/100mb.zip
https://testfiles.itkyou.cf/200mb.zip
https://testfiles.itkyou.cf/300mb.zip
https://testfiles.itkyou.cf/400mb.zip
https://testfiles.itkyou.cf/500mb.zip

https://testfiles.itkyou.tk/100mb.zip
https://testfiles.itkyou.tk/200mb.zip
https://testfiles.itkyou.tk/300mb.zip
https://testfiles.itkyou.tk/400mb.zip
https://testfiles.itkyou.tk/500mb.zip

https://testfiles.itkyou.gq/100mb.zip
https://testfiles.itkyou.gq/200mb.zip
https://testfiles.itkyou.gq/300mb.zip
https://testfiles.itkyou.gq/400mb.zip
https://testfiles.itkyou.gq/500mb.zip

https://testfiles.ityou.cf/100mb.zip
https://testfiles.ityou.cf/200mb.zip
https://testfiles.ityou.cf/300mb.zip
https://testfiles.ityou.cf/400mb.zip
https://testfiles.ityou.cf/500mb.zip

https://testfiles.ityou.tk/100mb.zip
https://testfiles.ityou.tk/200mb.zip
https://testfiles.ityou.tk/300mb.zip
https://testfiles.ityou.tk/400mb.zip
https://testfiles.ityou.tk/500mb.zip

https://testfiles.ityou.gq/100mb.zip
https://testfiles.ityou.gq/200mb.zip
https://testfiles.ityou.gq/300mb.zip
https://testfiles.ityou.gq/400mb.zip
https://testfiles.ityou.gq/500mb.zip

https://testfiles.kiring.cf/100mb.zip
https://testfiles.kiring.cf/200mb.zip
https://testfiles.kiring.cf/300mb.zip
https://testfiles.kiring.cf/400mb.zip
https://testfiles.kiring.cf/500mb.zip

https://testfiles.kiring.tk/100mb.zip
https://testfiles.kiring.tk/200mb.zip
https://testfiles.kiring.tk/300mb.zip
https://testfiles.kiring.tk/400mb.zip
https://testfiles.kiring.tk/500mb.zip

https://testfiles.kiring.gq/100mb.zip
https://testfiles.kiring.gq/200mb.zip
https://testfiles.kiring.gq/300mb.zip
https://testfiles.kiring.gq/400mb.zip
https://testfiles.kiring.gq/500mb.zip

https://testfiles.newbeer.cf/100mb.zip
https://testfiles.newbeer.cf/200mb.zip
https://testfiles.newbeer.cf/300mb.zip
https://testfiles.newbeer.cf/400mb.zip
https://testfiles.newbeer.cf/500mb.zip

https://testfiles.newbeer.gq/100mb.zip
https://testfiles.newbeer.gq/200mb.zip
https://testfiles.newbeer.gq/300mb.zip
https://testfiles.newbeer.gq/400mb.zip
https://testfiles.newbeer.gq/500mb.zip

# special thanks 
- CloudflareSpeedTest,Thanks very much !
- Jetbrains Goland IDE, Thanks very much !





