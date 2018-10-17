package etl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// 应用程序信息
type AppInfo struct {
	Home    string
	Name    string
	Version string
	Args
	Log *logrus.Logger
}

// 程序信息
func (ai *AppInfo) PrintInfo() {
	info := "程序：%s，版本：%s"
	info = fmt.Sprintf(info, ai.Name, ai.Version)
	fmt.Println(info)
}

// 程序参数
func (ai *AppInfo) PrintArgs() {
	info := "参数：-t %s -f %s -n %s -c %s"
	info = fmt.Sprintf(info,
		ai.Args.EtlTime,
		ai.Args.FlowId,
		ai.Args.NodeId,
		ai.Args.CallId)
	fmt.Println(info)
}

// 程序初始
func (ai *AppInfo) Init() {
	ai.Args = Args{}
	if !ai.Args.Parse() {
		ai.Close(false)
	}
	ai.Log = GetLog()
	if !ai.Args.Output {
		fileLog, err := GetFileLog(*ai)
		if nil != err {
			fmt.Println(err.Error())
			ai.Close(false)
		}
		ai.Log = fileLog
	}
}

// 程序加载
func (ai *AppInfo) Load(config interface{}) {
	// 加载配置文件
	configpath := filepath.Join(ai.Home, "etc", ai.Name+".json")
	ai.Log.Debugln("正在载入配置文件", configpath)
	rf, err := ioutil.ReadFile(configpath)
	if err != nil {
		ai.Log.Errorln("载入配置文件出错", configpath)
		ai.Close(false)
	}
	var configs []interface{}
	// 解析配置文件
	err = json.Unmarshal([]byte(rf), &configs)
	if nil != err {
		ai.Log.Errorln("解析配置文件内容出错", err.Error())
		ai.Close(false)
	}

	var c map[string]interface{}
	// 匹配执行配置项
	for _, value := range configs {
		if value.(map[string]interface{})["id"] == ai.Args.CallId {
			c = value.(map[string]interface{})
			break
		}
	}
	// 未找到配置项退出
	if c == nil {
		ai.Log.Errorln("未找到配置项: ", ai.Args.CallId)
		ai.Close(false)
	}
	// 解析成结构体
	js, _ := json.Marshal(c)
	json.Unmarshal(js, &config)
}

// 程序启动
func (ai *AppInfo) Start() {

}

// 程序结束
func (ai *AppInfo) End() {

}

// 程序退出
func (ai *AppInfo) Close(success bool) {
	if success {
		os.Exit(0)
	} else {
		os.Exit(-1)
	}
}

type AppConfigMatch func(AppConfigMatch) bool

// 应用程序配置
type AppConfig struct {
	Id string
}

// 应用程序配置读取接口
type AppConfigReader interface {
	Get(id string) (AppConfig, error)
}
