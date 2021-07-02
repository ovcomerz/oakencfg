package main

import (
	"flag"
	"fmt"
	"strings"
)
type multiStringArg []string

func newMultiStringArg(vals []string,p *[]string) *multiStringArg {
	*p = vals
	return (*multiStringArg)(p)
}

func (a *multiStringArg) Set(val string) error{
	*a = multiStringArg(strings.Split(val,","))
	return nil
}

func (a *multiStringArg)Get() interface{}  {
	return []string(*a)
}

func (a *multiStringArg)String() string  {
	return strings.Join(*a,",")
}

//导出模式
const(
	DirectoryMode string = "dir"
	SingleFileMode string = "single"
)

//目标格式
const(
	CSharp string = "cs"
	Lua string = "lua"
	Json string = "json"
)

var (
	//是否输出日志
	log bool
	//导出模式
	mode  = SingleFileMode
	//输出路径
	outputPath string
	//目标格式
	targets  []string
	//目标格式代码路径，顺序与目标语言一一对应
	targetCodePaths []string
)

func Export(){
	//解析命令行参数
	flag.BoolVar(&log,"log",false,"开启日志打印")
	flag.StringVar(&mode,"mode","dir","")
	flag.Var(newMultiStringArg([]string{CSharp},&targets),"targets","目标语言，多个目标用逗号分隔")
	flag.Parse()


	fmt.Printf("log:%t mode:%s taget:%s",log,mode,targets)
}