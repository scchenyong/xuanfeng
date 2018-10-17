package etl

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// 调度参数结构
type Args struct {
	// 是否启用日志详细输出，默认为false
	Output bool
	// 调度执行时间，不允许为空
	EtlTime string
	// 调度流程编号，不允许为空
	FlowId string
	// 调度流程节点编号，不允许为空
	NodeId string
	// 执行配置编号，不允许为空
	CallId string
}

// 解析
func (args *Args) Parse() bool {
	// 解析执行参数
	flag.BoolVar(&args.Output, "o", false, "启用详细输出<选项>\n\t默认：false")
	flag.StringVar(&args.EtlTime, "t", "", "调度执行时间\n\t格式：[YYYYMMDDHH24MISS]")
	flag.StringVar(&args.FlowId, "f", "", "调度流程编号")
	flag.StringVar(&args.NodeId, "n", "", "调度流程节点编号")
	flag.StringVar(&args.CallId, "c", "", "调度执行配置编号")
	flag.Usage = func() {
		msg := "用法: %s [-options] name [args...]\n\n"
		fmt.Fprintf(os.Stderr, msg, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()
	if !args.check() {
		flag.Usage()
		return false
	}
	return true
}

// 检查执行参数
func (args *Args) check() bool {
	// 无调度执行时间
	if len(args.EtlTime) == 0 || (len(args.EtlTime) != 8 && len(args.EtlTime) != 14) {
		fmt.Println("请输入正确的参数：调度执行时间（ -t ）")
		return false
	}
	// 无调度流程编号
	if len(args.FlowId) == 0 {
		fmt.Println("请输入正确的参数：调度流程编号（ -f ）")
		return false
	}
	// 无调度流程节点编号
	if len(args.NodeId) == 0 {
		fmt.Println("请输入正确的参数：调度流程节点编号（ -n ）")
		return false
	}
	// 无调度执行配置编号
	if len(args.CallId) == 0 {
		fmt.Println("请输入正确的参数：调度流程节点编号（ -c ）")
		return false
	}
	return true
}

// 解析数据处理执行时间
func (arg *Args) GetEtlTime() time.Time {
	dl := "20060102"
	tl := "20060102150405"
	var t time.Time
	if len(arg.EtlTime) == len(dl) {
		t, _ = time.Parse(dl, arg.EtlTime)
	}
	if len(arg.EtlTime) == len(tl) {
		t, _ = time.Parse(tl, arg.EtlTime)
	}
	return t
}
