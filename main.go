package main

import (
	"flag"
	"io/ioutil"
)

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
	logEnable bool
	multiSheetMode bool
	//导出模式
	mode  = SingleFileMode
	//输入路径
	inputPath []string
	//输出路径
	outputDir string
	//目标格式
	targets  []string
	//目标格式代码路径，顺序与目标语言一一对应
	targetCodeDir []string
)

func main() {
	parseExportArgs()
	export()
}

func parseExportArgs() {
	//解析命令行参数
	flag.BoolVar(&logEnable, "log", false, "开启日志打印")
	flag.BoolVar(&multiSheetMode, "multiSheet", false, "多sheet模式，开启后每个sheet对应导出一张配置表")
	flag.StringVar(&mode, "mode", "dir", "")
	flag.Var(newMultiStringArg([]string{}, &inputPath), "inputPath", "需要导出的文件/目录路径，多个用逗号隔开")
	flag.StringVar(&outputDir, "outputDirectory", "export", "导出目录")
	flag.Var(newMultiStringArg([]string{CSharp}, &targets), "targets", "目标语言，多个目标用逗号分隔")
	flag.Var(newMultiStringArg([]string{}, &targetCodeDir), "targetCodeDir", "目标语言代码目录，与目标语言一一对应")
	flag.Parse()
	//fmt.Printf("trace:%t mode:%s target:%s",logEnable,mode,targets)
	logInfo("logEnable:%t",logEnable)
}

func export(){
	traceInfo("##开始导出...")
	if mode == SingleFileMode {
		for _,filePath := range inputPath{
			exportSingleFile(filePath)
		}
	}else{
		for _,dirPath := range inputPath{
			exportByDir(dirPath)
		}
	}
}

func exportByDir(dir string){
	rd,err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for _,fi := range rd{
		if fi.IsDir(){
			exportByDir(dir + fi.Name())
		}else {
			exportSingleFile(dir + fi.Name())
		}
	}
}

func exportSingleFile(filePath string){
	tables := readExcel(filePath)
	for _,t := range tables{
		writeToJson(t)
	}

	print(len(tables))
}
