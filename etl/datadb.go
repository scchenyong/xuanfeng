package etl

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

// 数据库源配置结构（数据处理主机，数据处理文件）
type DataDB struct {
	// 数据库主机地址
	DataHost
	// 数据库名称
	DBName string
	// 执行脚本
	Script string
}

// 数据处理脚本转换处理时间
func (dd *DataDB) GetTimeScript(appinfo AppInfo, dt DataTime) (string, error) {
	datadate := dt.OffsetTime(appinfo.Args.GetEtlTime())
	YYYYMMDD := fmt.Sprintf("%d%02d%02d", datadate.Year(), datadate.Month(), datadate.Day())
	YYYYMMDDHH24MISS := fmt.Sprintf("%s%02d%02d%02d", YYYYMMDD, datadate.Hour(), datadate.Minute(), datadate.Second())
	scriptpath := filepath.Join(appinfo.Home, "scripts", dd.Script)
	sbytes, err := ioutil.ReadFile(scriptpath)
	if nil != err {
		return "", err
	}
	dbscript := string(sbytes)
	dbscript = strings.Replace(dbscript, "${YYYYMMDD}", YYYYMMDD, -1)
	dbscript = strings.Replace(dbscript, "${YYYYMMDDHH24MISS}", YYYYMMDDHH24MISS, -1)
	return dbscript, nil
}

// 数据库内容处理
type DataDBContent struct {
	// 前缀
	Prefix string
	// 后缀
	Suffix string
	// 空内容填充
	Empty string
	// 间隔
	Gap string
}

// 解析处理
func (dc *DataDBContent) Parse() {
	dc.Prefix = parseASCII(dc.Prefix)
	dc.Suffix = parseASCII(dc.Suffix)
	dc.Empty = parseASCII(dc.Empty)
	dc.Gap = parseASCII(dc.Gap)
}

// 解析ASCII字符
func parseASCII(a string) string {
	a = strings.Replace(a, "\\", "", -1)
	i, err := strconv.Atoi(a)
	if err == nil {
		a = string(i)
	}
	return a
}

// 字节内容处理
func (dc *DataDBContent) Byte2String(content []byte) string {
	if content == nil {
		content = []byte(dc.Empty)
	}
	return fmt.Sprintf("%s%s%s", dc.Prefix, string(content), dc.Suffix)
}

// 连接内容
func (dc *DataDBContent) Join(content []string) string {
	return strings.Join(content, dc.Gap)
}
