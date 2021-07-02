package main

type tableType int
const(
	//字典类型配置表，kv结构
	DictionaryTable tableType = iota
	//静态类型配置表，直接生成对应的代码常量
	ConstTable
)

type tableKeyType int
const(
	intTableKey tableKeyType = iota
	stringTableKey
)

type tableFiledType int
const(
	intField tableFiledType = iota
	byteField
	floatField
	doubleField
	int64Field
	stringField
	boolField
	listField

)

type Table struct {
	name string
	tType tableType
	keyType tableKeyType
}

type TableField struct {
	name string
	fileType
}