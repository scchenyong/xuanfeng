package main

import (
	"../../etl"
)

type ZIPfileConfig struct {
	etl.AppConfig
	Datatime etl.DataTime
	Datazip  etl.DataZip
}

var ETLApp etl.AppInfo
var ETLConfig ZIPfileConfig

func init() {
	ETLApp = etl.AppInfo{}
	ETLApp.Home = etl.GetParentDir(etl.PWD())
	ETLApp.Name = "file2zip"
	ETLApp.Version = "1.0"
	ETLApp.Init()
	ETLApp.Load(&ETLConfig)
	ETLApp.PrintInfo()
	ETLApp.PrintArgs()
}

func main() {
	ETLApp.Start()
	// 处理过程位置

	ETLApp.End()
	ETLApp.Close(true)
}
