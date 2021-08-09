package main

import "regexp"

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
	return false
}

func checkTableItemKey(key string,tType tableType)  bool{
	return true
}