package geoip

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/carlmjohnson/requests"
	"github.com/fireinrain/cf-speedtester/utils"
	"log"
	"net/http"
	"os"
	"strings"
)

// geoip工具 用于过滤出想要的国家ip
// 方便使用

const ProxiedUrlPrefix = "https://ghproxy.com/"
const GeoIPGitRepoUrl = "https://github.com/Loyalsoldier/geoip"
const RepoOwner = "Loyalsoldier"
const RepoName = "geoip"
const GEOIPFileMainName = "Country.mmdb"

//"https://ghproxy.com/https://github.com/Loyalsoldier/geoip/releases/download/202308170052/Country.mmdb"

var GeoIPDBUrlTemplate = "%shttps://github.com/Loyalsoldier/geoip/releases/download/%s/Country.mmdb"

var GlobalGeoIPClient GeoIPClient

type GeoIPClient struct {
	*GeoIPGitRepo
}

type GeoIPGitRepo struct {
	GeoIPDbFileName   string `json:"geoIpDbFileName"`
	GEOIPFileMainName string `json:"geoipFileMainMame"`
	LatestTagName     string `json:"LatestTagName"`
	CurrentTagName    string `json:"currentTagName"`
}

type Release struct {
	HtmlUrl string `json:"html_url"`
	TagName string `json:"tag_name"`
}

func init() {
	repo := GeoIPGitRepo{}
	repo.GetCurrentGeoIPDbFileInfo()
	//当前没有下载geoip文件
	if repo.GeoIPDbFileName == "" {
		_, err := repo.DownloadLatestGEOIPDb()
		if err != nil {
			log.Println("Error downloading latest gepdb file: ", err.Error())
			return
		}
	}
	GlobalGeoIPClient = GeoIPClient{
		GeoIPGitRepo: &repo,
	}

	fmt.Println("GEOIP DB 初始化成功...")
}

// GetCurrentGeoIPDbFileInfo
//
//	@Description: 获取当前目录的Country.mmdb 信息
//	@receiver repo
func (repo *GeoIPGitRepo) GetCurrentGeoIPDbFileInfo() {
	//设置主要名称
	repo.GEOIPFileMainName = GEOIPFileMainName
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current directory:", err)
		return
	}

	files, err := os.ReadDir(currentDir)
	if err != nil {
		log.Println("Error reading directory:", err)
		return
	}

	prefix := GEOIPFileMainName
	var matchingFiles []string

	for _, file := range files {
		if strings.HasPrefix(file.Name(), prefix) {
			matchingFiles = append(matchingFiles, file.Name())
		}
	}

	if len(matchingFiles) > 0 {
		currentDbFile := matchingFiles[0]
		log.Println("Matching geoip file: ", currentDbFile)
		fileNameList := strings.Split(currentDbFile, "-")
		repo.CurrentTagName = strings.TrimSpace(fileNameList[1])
		repo.GeoIPDbFileName = strings.TrimSpace(currentDbFile)

	} else {
		log.Println("No matching GEOIP files found.")
	}
}

// GetRepoLatestTag
//
//	@Description: 获取最新的geoip tag
//	@receiver repo
//	@return Release
func (repo *GeoIPGitRepo) GetRepoLatestTag() Release {
	var release Release

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", RepoOwner, RepoName)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("GetRepoLatestTag Error:", err)
		return release
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		log.Println("Error decoding JSON:", err)
		return release
	}
	return release
}

func (repo *GeoIPGitRepo) DownloadLatestGEOIPDb() (dbName string, err error) {
	release := repo.GetRepoLatestTag()
	fullFileName := fmt.Sprintf("%s-%s", GEOIPFileMainName, release.TagName)
	exists := utils.FileOrDirExists(fullFileName)
	//如果不存在 就下载最新的
	if exists {
		return fullFileName, nil
	}

	proxiedDownloadUrl := fmt.Sprintf(GeoIPDBUrlTemplate, ProxiedUrlPrefix, release.TagName)
	ctx := context.Background()
	err = requests.
		URL(proxiedDownloadUrl).
		ToFile(fullFileName).
		Fetch(ctx)
	log.Printf("正在下载文件GEOIP db: %s,请稍后...\n", fullFileName)
	if err != nil {
		log.Printf("文件下载失败: %s\n", err.Error())
		return repo.GeoIPDbFileName, nil
	}
	log.Printf("GEOIP db文件下载成功: %s\n", fullFileName)
	//remove old geodb file
	err = os.Remove(repo.GeoIPDbFileName)
	if err != nil {
		log.Println("Remove old GEOIP file error: ", err)
		return
	}
	repo.LatestTagName = release.TagName
	repo.CurrentTagName = release.TagName
	repo.GeoIPDbFileName = fullFileName

	return fullFileName, nil
}
