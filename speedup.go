package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"

	"github.com/robfig/cron/v3"
)

var lastIP string

// 创建 Resty 客户端
var client = resty.New()

// 获取 IP 的函数
func getIP() (string, error) {
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
	if err != nil {
		return "", err
	}
	fmt.Println("getIP: ", data)
	return data.IP, nil
}

// 用来提速的函数
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
	fmt.Println("speedup: ", data)
	return data, nil
}

// 用来检查是否提速成功的函数
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
	fmt.Println("checkSpeedup: ", data)
	return data.Data.IndexInfo.CanSpeed == 1, nil
}

func checkIPChange() {
	ip, err := getIP()
	if err != nil {
		fmt.Println("获取 IP 失败：", err)
		return
	}
	if ip != lastIP {
		lastIP = ip
		mainFunc()
	}
}

// 主函数
func mainFunc() {
	ip, err := getIP()
	if err != nil {
		fmt.Println("获取 IP 失败：", err)
		return
	}
	lastIP = ip
	speedupResult, err := speedup()
	if err != nil {
		fmt.Println("提速请求异常：", err)
		return
	}

	if code, ok := speedupResult["code"].(float64); ok && code != 0 {
		fmt.Println("提速请求异常，错误码：", code)
		return
	}

	isSpeedupSucceed, err := checkSpeedup()
	if err != nil {
		fmt.Println("检查提速状态失败：", err)
		return
	}

	if !isSpeedupSucceed {
		fmt.Println("提速失败")
	} else {
		fmt.Println("提速成功")
	}
}

func main() {
	c := cron.New()
	// 每 10 分钟检查 IP 变化
	c.AddFunc("*/10 * * * *", checkIPChange)
	// 每周一 0 点运行
	c.AddFunc("0 0 * * 1", mainFunc)
	c.Start()

	// 程序开始运行
	mainFunc()

	// 防止主程序退出
	select {}
}