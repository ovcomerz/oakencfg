package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)


type tableType int
const(
	//字典类型配置表，kv结构
	intDictionaryTable tableType = iota
	stringDictionaryTable
	//静态类型配置表，直接生成对应的代码常量
	constTable
)

//字段类型
type tableFiledBaseType int
const(
	intField tableFiledBaseType = iota + 1
	byteField
	floatField
	float64Field
	int64Field
	stringField
	boolField
)

type tableFieldAdvanceType int
const (
	baseTypeField tableFieldAdvanceType = iota - 1
	listField
	twoDimensionListField
	mapField
)

const listSuffix string  = "_list"
const twoDimensionSuffix string = "_2dlist"
const mapSuffix string = "_map"

type Table struct {
	name   string
	tType  tableType
	fields []*TableField
	data   map[interface{}]*TableItem
}

type TableItem struct {
	key string
	values []string
}

type TableField struct {
	name     string
	index    int
	baseType tableFiledBaseType
	advanceType tableFieldAdvanceType
	advanceExtraType []tableFiledBaseType
	comment  string
}

func(field *TableField) parseType(typeName string) error{
	advanceType := parseFieldAdvanceType(typeName)
	var err error = nil
	field.advanceType = advanceType
	if advanceType != baseTypeField {
		baseType := getBaseTypeByAdvanceType(typeName)
		if advanceType == mapField {
			types := strings.Split(baseType,"_")
			keyType,err :=  parseFiledBaseType(types[0])
			if err != nil{
				return err
			}
			valueType,err :=  parseFiledBaseType(types[0])
			if err != nil{
				return err
			}
			field.baseType = keyType
			field.advanceExtraType = []tableFiledBaseType{valueType}
		}else {
			field.baseType,err = parseFiledBaseType(baseType)
		}
	}else {
		field.baseType,err = parseFiledBaseType(typeName)
	}
	return err
}

func (item *TableItem)valuesToMap(metadata []*TableField)  interface{} {
	values := make(map[string]interface{})
	for index,field := range metadata{
		content := ""
		if index < len(item.values) {
			content = item.values[index]
		}
		v,err := parseValue(content,field)
		if err != nil {
			traceInfo("字段:%s 解析错误",field.name,err.Error())
		}else {
			values[field.name] = v
		}
	}
	return values
}

func parseValue(raw string, field *TableField) (interface{},error) {
	if field.advanceType >= 0 {
		switch field.advanceType {
		case listField:
			return parseListValue(raw,field.baseType)
		case twoDimensionListField:
			return parseTwoDimensionListValues(raw,field.baseType)
		case mapField:
			return parseMapValues(raw,field.advanceExtraType[0],field.advanceExtraType[1])
		default:
			return nil,errors.Errorf("不支持的字段高级类型:%s",field.advanceType)
		}
	}else {
		return parseBaseTypeValue(raw,field.baseType)
	}
}

func parseBaseTypeValue(raw string,t tableFiledBaseType) (interface{},error ){
	if raw == "" {
		return getBaseTypeDefaultValue(t),nil
	}
	var value interface{}
	var err error
	switch t {
	case intField:
		value,err = strconv.Atoi(raw)
	case int64Field:
		value,err = strconv.ParseInt(raw,10,64)
	case byteField:
		var v int64
		v,err = strconv.ParseInt(raw,10,8)
		value = int8(v)
	case floatField:
		var v float64
		v,err = strconv.ParseFloat(raw,32)
		value = float32(v)
	case float64Field:
		value,err = strconv.ParseFloat(raw,64)
	case boolField:
		var v int
		v,err = strconv.Atoi(raw)
		value = v != 0
	case stringField:
		value = raw
	}
	return value,err
}

func getBaseTypeDefaultValue(t tableFiledBaseType) interface{}{
	var v interface{}
	switch t {
	case intField,int64Field,byteField:
		v = 0
	case floatField,float64Field:
		v =   0
	case boolField:
		v = false
	default:
		v = ""
	}
	return v
}

func parseListValue(raw string,baseType tableFiledBaseType) ([]interface{},error){
	if !checkListFormat(raw) {
		return nil,errors.Errorf("list类型格式错误:%s",raw)
	}

	rawLen := len(raw)
	if rawLen > 0 {
		raw = raw[1:len(raw) - 1]
		rawLen = len(raw)
	}

	if rawLen == 0 {
		return make([]interface{},0),nil
	}

	items := strings.Split(raw,",")
	var values []interface{}
	if len(items) == 0{
		return make([]interface{},0),nil
	}
	for _,item := range items{
		val,err := parseBaseTypeValue(item,baseType)
		if err != nil {
			return nil,err
		}
		values = append(values,val)
	}
	return values,nil
}

func parseTwoDimensionListValues(raw string,baseType tableFiledBaseType)  ([]interface{},error){
	if !checkListFormat(raw) {
		return nil,errors.Errorf("list类型格式错误:%s",raw)
	}

	rawLen := len(raw)
	if rawLen > 0 {
		raw = raw[1:len(raw) - 1]
		rawLen = len(raw)
	}

	if rawLen == 0 {
		return make([]interface{},0),nil
	}
	var values []interface{}
	var startIndex int = -1

	for index := 0;index < rawLen;index++{
		char := raw[index]
		switch char {
		case '[':
			startIndex = index
		case ']':
			sub := raw[startIndex:index + 1]
			subValue,err := parseListValue(sub,baseType)
			if err != nil {
				return nil,err
			}
			values = append(values,subValue)
			fmt.Println(sub)
			if index < rawLen - 1 {
				if raw[index + 1] != ',' || index == rawLen - 2{
					return nil,errors.Errorf("二维list解析错误!!")
				}
				index++
			}
			startIndex = -1
		default:
			if startIndex == - 1 || index == rawLen - 1{
				return nil,errors.Errorf("二维list解析错误!!")
			}
		}
	}
	return values,nil
}

func parseMapValues(raw string,itemKeyType tableFiledBaseType,itemValueType tableFiledBaseType) (map[interface{}]interface{},error){
	if !checkMapFormat(raw) {
		return nil,errors.Errorf("map类型格式错误")
	}
	raw = raw[1:len(raw) - 1]
	if len(raw)  == 0 {
		return  make(map[interface{}]interface{}),nil
	}

	var values map[interface{}]interface{}
	items := strings.Split(raw,",")
	for _,item := range items{
		kv := strings.Split(item,":")
		if len(kv) != 2 {
			return nil,errors.Errorf("map类型格式错误")
		}
		key,err := parseBaseTypeValue(kv[0],itemKeyType)
		if err != nil{
			 return nil,errors.Errorf("map类型格式错误")
		}
		value,err := parseBaseTypeValue(kv[1],itemValueType)
		if err != nil{
			return nil,errors.Errorf("map类型格式错误")
		}
		values[key] = value
	}

	return values,nil
}

func parseFiledBaseType(typeName string) (tableFiledBaseType,error) {
	var t tableFiledBaseType
	switch typeName {
	case "string":
		t = stringField
	case "int":
		t = intField
	case "int64":
		t = int64Field
	case "float":
		t = floatField
	case "double":
		t = float64Field
	case "byte":
		t = byteField
	case "bool":
		t = boolField
	default:
		return 0,errors.Errorf("不支持的基础数据类型:%s",typeName)
	}
	return t,nil
}

func parseFieldAdvanceType(typeName string)  tableFieldAdvanceType{
	if strings.Contains(typeName,listSuffix) {
		return listField
	}else if strings.Contains(typeName, twoDimensionSuffix){
		return twoDimensionListField
	}else if strings.Contains(typeName, mapSuffix) {
		return mapField
	}else{
		return baseTypeField
	}
}

func getBaseTypeByAdvanceType(typeName string) string{
	typeName = strings.ReplaceAll(typeName,listSuffix,"")
	typeName = strings.ReplaceAll(typeName,twoDimensionSuffix,"")
	typeName = strings.ReplaceAll(typeName,mapSuffix,"")
	return typeName
}

func getTableType(typeName string) tableType {
	var t tableType
	switch typeName {
		case "int":
			t = intDictionaryTable
		case "const":
			t = constTable
		case "string":
			t = stringDictionaryTable
	default:
			t  = -1
	}
	return t
}

func checkListFormat(raw string) bool {
	if  raw != "" && (!strings.HasPrefix(raw,"[") || !strings.HasSuffix(raw,"]")) {
		return false
	}
	return true
}

func checkMapFormat(raw string) bool {
	if !strings.HasPrefix(raw,"{") || !strings.HasSuffix(raw,"}") {
		return false
	}
	return true
}

func (table *Table)toJson() []byte{
	m := make(map[string]interface{})
	for key,val := range table.data{
		m[key.(string)] = val.valuesToMap(table.fields)
	}
	//b,_ := json.Marshal(m)
	//b,_ := json.MarshalIndent(m,"","    ")
	buffer := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("","    ")
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(m)
	if err != nil {
		fmt.Println("json 解析错误")
		return nil
	}
	return buffer.Bytes()

}

