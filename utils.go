package main

import (
	"log"
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



func traceInfo(info string,v ...interface{}){
	if len(v) == 0{
		log.Print(info)
	}else {
		log.Printf(info,v)
	}

}

func logInfo(info string,v ...interface{}){
	if logEnable {
		if len(v) == 0 {
			log.Print(info)
		}else {
			log.Printf(info,v)
		}
	}
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}