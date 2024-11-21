package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/robfig/cron/v3"
)

var lastIP string
var client = resty.New()

func init() {
	// 配置日志格式：时间 文件：行号 日志内容
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// 设置输出到标准输出
	log.SetOutput(os.Stdout)
}

func getIPbyBaidu() (string, error) {
	// 使用 Resty 发送 GET 请求
	resp, err := client.R().
		Get("https://qifu-api.baidubce.com/ip/local/geo/v1/district")
	if err != nil {
		return "", err
	}
	// 解析响应数据
	var data struct {
		Code string `json:"code"`
		Data struct {
			Continent string `json:"continent"`
			Country   string `json:"country"`
			Zipcode   string `json:"zipcode"`
			Owner     string `json:"owner"`
			ISP       string `json:"isp"`
			Adcode    string `json:"adcode"`
			Prov      string `json:"prov"`
			City      string `json:"city"`
			District  string `json:"district"`
		} `json:"data"`
		IP string `json:"ip"`
	}
	err = json.Unmarshal(resp.Body(), &data)
	if data.Code == "DailyLimited" {
		log.Printf("getIPbyBaidu: 获取 IP 失败，每日接口调用次数已达上限")
		return "", nil
	}
	if err != nil {
		return "", err
	}
	log.Printf("getIPbyBaidu: %+v", data)
	return data.IP, nil
}

func getIPbyTencent() (string, error) {
	// 使用 Resty 发送 GET 请求
	resp, err := client.R().
		Get("https://r.inews.qq.com/api/ip2city")
	if err != nil {
		return "", err
	}
	// 解析响应数据
	var data struct {
		Ret          int16  `json:"ret"`
		ErrMsg       string `json:"errMsg"`
		Country      string `json:"country"`
		ProvCode     string `json:"provcode"`
		CityCode     string `json:"citycode"`
		DistrictCode string `json:"districtCode"`
		Province     string `json:"province"`
		ISP          string `json:"isp"`
		City         string `json:"city"`
		District     string `json:"district"`
		IP           string `json:"ip"`
		Callback     string `json:"callback"`
	}
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return "", err
	}
	log.Printf("getIPbyTencent: %+v", data)
	return data.IP, nil
}

func getIPbyIpip() (string, error) {
	// 使用 Resty 发送 GET 请求
	resp, err := client.R().
		Get("https://myip.ipip.net/json")
	if err != nil {
		return "", err
	}
	// 解析响应数据
	var data struct {
		Ret  string `json:"ret"`
		Data struct {
			IP string `json:"ip"`
		} `json:"data"`
	}
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return "", err
	}
	log.Printf("getIPbyIpip: %+v", data)
	return data.Data.IP, nil
}

func getIP() (string, error) {
	// Create a slice of the getter functions with their names
	funcs := []struct {
		name string
		fn   func() (string, error)
	}{
		{"百度", getIPbyBaidu},
		{"腾讯", getIPbyTencent},
		{"IPIP", getIPbyIpip},
	}

	// Randomly shuffle the functions
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := len(funcs) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		funcs[i], funcs[j] = funcs[j], funcs[i]
	}

	// Try each function in the shuffled order
	for _, f := range funcs {
		ip, err := f.fn()
		if err != nil {
			log.Printf("通过%s获取IP失败: %v", f.name, err)
			continue
		}
		if ip != "" {
			return ip, nil
		}
	}

	return "", fmt.Errorf("无法通过任何一个方法获取 IP")
}

func speedup() (map[string]interface{}, error) {
	resp, err := client.R().
		Get("https://tisu-api.speedtest.cn/api/v2/speedup/reopen")
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return nil, err
	}
	log.Printf("speedup: %+v", data)
	return data, nil
}

func checkSpeedup() (bool, error) {
	resp, err := client.R().
		Get("https://tisu-api-v3.speedtest.cn/speedUp/query")
	if err != nil {
		return false, err
	}
	var data struct {
		Data struct {
			IndexInfo struct {
				CanSpeed int `json:"canSpeed"`
			} `json:"indexInfo"`
		} `json:"data"`
	}
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return false, err
	}
	log.Printf("checkSpeedup: %+v", data)
	return data.Data.IndexInfo.CanSpeed == 1, nil
}

func checkIPChange() {
	ip, err := getIP()
	if err != nil {
		log.Printf("获取 IP 失败：%v", err)
		return
	}
	if ip != lastIP {
		lastIP = ip
		mainFunc()
	}
}

func mainFunc() {
	ip, err := getIP()
	if err != nil {
		log.Printf("获取 IP 失败：%v", err)
		return
	}
	lastIP = ip
	speedupResult, err := speedup()
	if err != nil {
		log.Printf("提速请求异常：%v", err)
		return
	}

	if code, ok := speedupResult["code"].(float64); ok && code != 0 {
		log.Printf("提速请求异常，错误码：%v", code)
		return
	}

	isSpeedupSucceed, err := checkSpeedup()
	if err != nil {
		log.Printf("检查提速状态失败：%v", err)
		return
	}

	if !isSpeedupSucceed {
		log.Println("提速失败")
	} else {
		log.Println("提速成功")
	}
}

func main() {
	// 检查是否有 --version 参数
	for _, arg := range os.Args[1:] {
		if arg == "--version" {
			fmt.Println(GetVersionInfo())
			return
		}
	}

	c := cron.New()
	// 每 10 分钟检查 IP 变化
	if _, err := c.AddFunc("*/10 * * * *", checkIPChange); err != nil {
		log.Printf("Error adding cron function: %v", err)
	}
	// 每周一 0 点运行
	if _, err := c.AddFunc("0 0 * * 1", mainFunc); err != nil {
		log.Printf("Error adding cron function: %v", err)
	}
	c.Start()

	// 程序开始运行
	mainFunc()

	// 防止主程序退出
	select {}
}
