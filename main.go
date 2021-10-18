package main

import (
	"bufio"
	"fmt"
	"github.com/eddieivan01/nic"
	"github.com/json-iterator/go"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Comments struct {
	Code int `json:"code"`
	Data []struct {
		Message    string      `json:"message"`
		ID         int64       `json:"id"`
		Floor      int         `json:"floor"`
		Count      int         `json:"count"`
		Root       int         `json:"root"`
		Oid        int         `json:"oid"`
		Bvid       string      `json:"bvid"`
		Ctime      string      `json:"ctime"`
		Mtime      string      `json:"mtime"`
		State      int         `json:"state"`
		Parent     int         `json:"parent"`
		Mid        int         `json:"mid"`
		Like       int         `json:"like"`
		Replier    string      `json:"replier"`
		Uface      string      `json:"uface"`
		Cover      string      `json:"cover"`
		Title      string      `json:"title"`
		Relation   int         `json:"relation"`
		IsElec     int         `json:"is_elec"`
		Type       int         `json:"type"`
		RootInfo   interface{} `json:"root_info"`
		ParentInfo interface{} `json:"parent_info"`
		Attr       int         `json:"attr"`
		Vote       interface{} `json:"vote"`
		Action     int         `json:"action"`
	} `json:"data"`
	Message string `json:"message"`
	Pager   struct {
		Current int `json:"current"`
		Size    int `json:"size"`
		Total   int `json:"total"`
	} `json:"pager"`
}

type Config struct {
	Interval int `json:"Interval"`
	Block    struct {
		Users []string `json:"Users"`
		Video []string `json:"Video"`
		Texts []string `json:"Texts"`
		Regex string   `json:"Regex"`
	} `json:"Block"`
	WhiteList []string `json:"WhiteList"`
	Request   struct {
		Order  string `json:"Order"`
		Filter string `json:"Filter"`
		Type   string `json:"Type"`
		Page   string `json:"Page"`
		Size   string `json:"Size"`
	} `json:"Request"`
}

type ApiDeleteResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
}

var config Config
var session *nic.Session
var cookieKV nic.KV

func main() {
	session = nic.NewSession()
	log.Println("Bilibili 评论清理助手 1.0 by Akkariin")
	log.Println("--------------------------------------")
	cookie := ""
	if FileExists("config.json") {
		config = Config{}
		configText := GetFileData("config.json")
		err := jsoniter.Unmarshal([]byte(configText), &config)
		if err != nil {
			log.Fatalln(fmt.Sprintf("无法读取配置文件，错误：%s", err.Error()))
		}
		jsonStu, err := jsoniter.MarshalIndent(config, "", "    ")
		if err == nil {
			SetFileData("config.json", string(jsonStu))
		} else {
			log.Fatalln(fmt.Sprintf("尝试更新配置文件内容时发生错误：%s", err.Error()))
		}
	} else {
		config = Config{}
		configText := GenerateConfig()
		SetFileData("config.json", configText)
		err := jsoniter.Unmarshal([]byte(configText), &config)
		if err != nil {
			log.Fatalln(fmt.Sprintf("无法读取配置文件，错误：%s", err.Error()))
		}
	}
	if FileExists("cookie.txt") {
		cookie = GetFileData("cookie.txt")
	} else {
		ir := bufio.NewReader(os.Stdin)
		log.Println("当前未设置 Cookie，请打开 Bilibili 创作中心，并按下 F12 打开浏览器控制台，转到 “网络” 或者 “Network”，然后刷新页面。")
		log.Println("刷新后找到第一个请求，将 “请求标头” 或 “Request Header” 中的 cookie: 后面的内容复制，然后在这里粘贴（Win10 支持右键粘贴）")
		log.Println("如果不知道如何操作的话请看 B 站上的文档介绍 https://space.bilibili.com/21165317 （@Akkariins）")
		fmt.Print("> ")
		cookie, err := ir.ReadString('\n')
		if err == nil {
			SetFileData("cookie.txt", cookie)
			log.Println("已储存 Cookie 到文件 cookie.txt，程序开始运行...")
		} else {
			log.Fatalln(err.Error())
		}
	}
	for {
		checkComments(cookie)
		log.Println(fmt.Sprintf("等待 %s 秒后再次检查...", strconv.Itoa(config.Interval)))
		time.Sleep(time.Duration(config.Interval) * time.Second)
	}
}

func checkComments(cookie string) {
	if cookie != "" {
		cookieKV = nic.KV{}
		cookieList := strings.Split(cookie, ";")
		csrfToken := ""
		for _, element := range cookieList {
			cookiePart := strings.Split(element, "=")
			if len(cookiePart) == 2 {
				key := strings.TrimSpace(cookiePart[0])
				value := strings.TrimSpace(cookiePart[1])
				cookieKV[key] = value
				if key == "bili_jct" {
					csrfToken = value
				}
			}
		}
		resp, err := session.Get("https://member.bilibili.com/x/web/replies", nic.H{
			Cookies: cookieKV,
			Params: nic.KV{
				"order":           config.Request.Order,
				"filter":          config.Request.Filter,
				"is_hidden":       "0",
				"type":            config.Request.Type,
				"comment_display": "0",
				"pn":              config.Request.Page,
				"ps":              config.Request.Size,
			},
		})
		if err == nil {
			comments := &Comments{}
			err := resp.JSON(&comments)
			if err == nil {
				for _, comment := range comments.Data {
					// 判断是否在白名单列表
					log.Println(fmt.Sprintf("[%s] %s", comment.Replier, comment.Message))
					if !IsWhiteListUser(comment.Replier, comment.Mid) {
						// 判断此视频是否需要监控
						if IsVideoNeedBlock(comment.Bvid) {
							if IsBlockedUser(comment.Replier, comment.Mid) { // 判断此用户是否在黑名单内
								log.Println(fmt.Sprintf("发现黑名单用户 [ %s ] 已删除该用户评论！", comment.Replier))
								DeleteComment(comment.Oid, comment.Type, comment.ID, csrfToken)
							} else if IsBlockedText(comment.Message) { // 判断是否触发违禁词内容
								log.Println(fmt.Sprintf("触发黑名单词库 [ %s ] 已删除该用户评论！", comment.Replier))
								DeleteComment(comment.Oid, comment.Type, comment.ID, csrfToken)
							} else if IsBlockedRegex(comment.Message) { // 判断是否触发正则表达式判定
								log.Println(fmt.Sprintf("触发黑名单正则 [ %s ] 已删除该用户评论！", comment.Replier))
								DeleteComment(comment.Oid, comment.Type, comment.ID, csrfToken)
							}
						}
					}
				}
			} else {
				log.Fatalln("无法获取评论，错误：", err.Error())
			}
		} else {
			log.Fatalln("无法获取评论，错误：", err.Error())
		}
	}
}

func DeleteComment(oid int, ctype int, rpid int64, csrfToken string) bool {
	resp, err := session.Post("https://api.bilibili.com/x/v2/reply/del", nic.H{
		Cookies: cookieKV,
		Data: nic.KV{
			"oid":   strconv.Itoa(oid),
			"type":  strconv.Itoa(ctype),
			"rpid":  strconv.FormatInt(rpid, 10),
			"jsonp": "jsonp",
			"csrf":  csrfToken,
		},
	})
	if err == nil {
		apiJson := &ApiDeleteResponse{}
		err = jsoniter.Unmarshal([]byte(resp.Text), &apiJson)
		if err == nil {
			if apiJson.Code == 0 {
				return true
			} else {
				log.Println(fmt.Sprintf("无法删除评论：%s", apiJson.Message))
				return false
			}
		} else {
			log.Println(fmt.Sprintf("无法删除评论：%s", err.Error()))
			return false
		}
	} else {
		log.Println(fmt.Sprintf("无法删除评论：%s", err.Error()))
		return false
	}
}

func IsWhiteListUser(userName string, userId int) bool {
	for _, user := range config.WhiteList {
		if userName == user {
			return true
		} else {
			intId, err := strconv.Atoi(user)
			if err == nil {
				if userId == intId {
					return true
				}
			}
		}
	}
	return false
}

func IsBlockedUser(userName string, userId int) bool {
	for _, user := range config.Block.Users {
		if userName == user {
			return true
		} else {
			intId, err := strconv.Atoi(user)
			if err == nil {
				if userId == intId {
					return true
				}
			}
		}
	}
	return false
}

func IsBlockedText(message string) bool {
	for _, block := range config.Block.Texts {
		if strings.Contains(message, block) {
			return true
		}
	}
	return false
}

func IsBlockedRegex(message string) bool {
	if config.Block.Regex == "" {
		return false
	}
	match, err := regexp.MatchString(config.Block.Regex, message)
	if err == nil {
		return match
	}
	return false
}

func IsVideoNeedBlock(bvid string) bool {
	if len(config.Block.Video) == 0 {
		return true
	}
	for _, id := range config.Block.Video {
		if bvid == id {
			return true
		}
	}
	return false
}

func GenerateConfig() string {
	config := Config{
		Interval: 30,
	}
	config.Request.Order = "ctime"
	config.Request.Filter = "-1"
	config.Request.Type = "1"
	config.Request.Page = "1"
	config.Request.Size = "10"
	jsonStu, err := jsoniter.MarshalIndent(config, "", "    ")
	if err == nil {
		return string(jsonStu)
	} else {
		return "{}"
	}
}

// 读文件
func GetFileData(fileName string) string {
	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln("Failed to read file", err)
	}
	return string(f)
}

// 写文件
func SetFileData(fileName string, data string) bool {
	var bytes = []byte(data)
	err := ioutil.WriteFile(fileName, bytes, 0666)
	return err == nil
}

// 判断所给路径文件/文件夹是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path) // os.Stat 获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}
