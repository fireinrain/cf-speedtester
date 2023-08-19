# cf-speedtester
A golang library to use cloudflare speed test function. use it in your own golang project.

# why you write this golang library
I want to use cloudflare cdn ip speed test in my own project.
there are a command-line tool called [CloudflareSpeedTest](https://github.com/XIU2/CloudflareSpeedTest), but if 
you want to use in your own project,you need wrap it as a bash call.
so, here is the project.


# how to use

## simple use(use cloudflare official cdn ips)
```go

client := NewCFSpeedTestClient(
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

client := NewCFSpeedTestClient(
		config.WithMaxDelay(300*time.Millisecond),
		config.WithMinSpeed(2),
		config.WithTestCount(5),
	)
client.DoSpeedTest()
results := client.SpeedResults
fmt.Println(results)
client.ExportToCSV("results.csv")


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

client := NewCFSpeedTestClient(
    config.WithMaxDelay(300*time.Millisecond),
    config.WithMinSpeed(2),
    config.WithTestCount(1),
    config.WithIPListForTest(ipList),
)
result := client.DoSpeedTestForResult()
fmt.Println(result)



```

# Special thanks 
- CloudflareSpeedTest,Thanks very much !
- Jetbrains Goland IDE, Thanks very much !





