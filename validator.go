package main

import (
	"regexp"
	"strings"
)

func checkTableName(name string) bool{
	if len(name) <= 0{
		return false
	}

	m,_ := regexp.MatchString("^[a-zA-Z]+[\\w]*$",name)
	return m
}

func checkTableNameExclude(name string ) bool {
	return false
}

func checkColExclude(comment string) bool{
	return strings.HasPrefix(comment,"#")
}

func checkTableItemKey(key string,tType tableType)  bool{
	return true
}