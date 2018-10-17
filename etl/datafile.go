package etl

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// 数据处理文件结构
type DataFile struct {
	// 文件路径
	Path string
	// 文件名
	File string
}

// 数据处理文件转换处理时间
func (df *DataFile) GetTimeFile(dt DataTime, t time.Time) string {
	filedate := dt.OffsetTime(t)
	filepath := filepath.Join(df.Path, df.File)
	filepath = strings.Replace(filepath, "YYYY", fmt.Sprintf("%d", filedate.Year()), -1)
	filepath = strings.Replace(filepath, "MM", fmt.Sprintf("%02d", filedate.Month()), -1)
	filepath = strings.Replace(filepath, "DD", fmt.Sprintf("%02d", filedate.Day()), -1)
	filepath = strings.Replace(filepath, "HH24", fmt.Sprintf("%02d", filedate.Hour()), -1)
	filepath = strings.Replace(filepath, "MI", fmt.Sprintf("%02d", filedate.Minute()), -1)
	filepath = strings.Replace(filepath, "SS", fmt.Sprintf("%02d", filedate.Second()), -1)
	return filepath
}
