package main

import (
	"github.com/Luxurioust/excelize"
	"path"
	"strings"
)

//表最小列数
const tableMinColNum int = 2

//表最小行数
const tableMinRowNum int = 3

func excelToTable(filePath string) []*Table {
	logInfo("###开始读取excel文件:%s", filePath)
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		traceInfo("读取excel文件失败,路径:%s", filePath)
		logInfo("###读取结束")
		return nil
	}

	fileNameWithoutExtension := strings.TrimSuffix(filePath, path.Ext(filePath))
	var sheetCount int
	var exportSheets []string
	if !multiSheetMode {
		sheetCount = 1
		exportSheets = []string{f.GetSheetName(0)}
	} else {
		sheetCount = f.SheetCount
		exportSheets = []string{}
		allSheets := f.GetSheetList()
		for _, name := range allSheets {
			if !checkTableNameExclude(name) {
				exportSheets = append(exportSheets, name)
			} else {
				logInfo("忽略sheet:%s", name)
			}
		}
	}
	tables := make([]*Table, 0, sheetCount)
	for _, sheetName := range exportSheets {
		tableName := getTableName(fileNameWithoutExtension, sheetName)
		if checkTableNameExclude(tableName) {
			logInfo("跳过导出：%s", tableName)
			continue
		}
		if checkTableName(tableName) {
			rows, err := f.GetRows(sheetName)
			if err != nil {
				traceInfo("获取rows失败,表名:%s\n%s", tableName, err)
				continue
			}
			table := createTable(tableName,rows)
			if table != nil {
				tables = append(tables,table)
			}
		}
	}
	logInfo("###读取结束")
	return tables
}

func createTable(tableName string, rows [][]string) (table *Table) {
	if len(rows) < tableMinRowNum {
		traceInfo("%s 数据行小于%d行", tableName, tableMinColNum)
		return nil
	}

	rowIndex := 0
	commentRow := rows[rowIndex]
	rowIndex++
	nameRow := rows[rowIndex]
	if len(nameRow) < tableMinColNum {
		traceInfo("%s 数据列小于%d列", tableName, tableMinColNum)
		return nil
	}

	totalColumnNum := len(commentRow)
	rowIndex++
	typeRow := rows[rowIndex]
	rowIndex++

	col := 0
	typeName := typeRow[col] //类型行第一列类型为表的类型
	tType := getTableType(typeName)
	if tType < 0 {
		traceInfo("%s 不支持的配置表类型:%s", tableName, typeName)
		return nil
	}

	table = &Table{
		name:   tableName,
		data:   make(map[interface{}]*TableItem),
		fields: make([]*TableField, 0),
	}
	excludeColIndices := make([]int, 0, 4)

	table.tType = tType
	if tType == constTable { //常量配置表,信息直接存在data中，忽略field信息
		//valueType := parseFiledType(typeRow[col])
		table.fields = make([]*TableField, 0)
		table.data = make(map[interface{}]*TableItem)
		for ; rowIndex < len(rows); rowIndex++ {
			row := rows[rowIndex]
			field := TableField{
				name: nameRow[0],
			}
			err := field.parseType(row[col])
			if err != nil {
				traceInfo("%s字段类型解析错误:%s row:%d", tableName, err, rowIndex)
				return nil
			}
			table.fields = append(table.fields, &field)
			item := TableItem{
				values: []string{row[col+1]},
			}
			table.data[len(table.fields)+1] = &item
		}
	} else {
		col++
		fieldIndex := 0
		for ; col < totalColumnNum; col++ {
			if checkColExclude(commentRow[col]) {
				excludeColIndices = append(excludeColIndices, col)
				continue
			}
			filed := TableField{
				name:    nameRow[col],
				index:   fieldIndex,
				comment: commentRow[col],
			}

			err := filed.parseType(typeRow[col])
			if err != nil {
				traceInfo("%s字段类型解析错误:%s column:%d", tableName, err, col)
				return nil
			}
			fieldIndex++
			table.fields = append(table.fields, &filed)
		}

		for ; rowIndex < len(rows); rowIndex++ {
			col = 0
			row := rows[rowIndex]
			key := row[col]
			if !checkTableItemKey(key, tType) {
				continue
			}
			item := TableItem{
				key: key,
			}
			col++
			for ; col < len(row); col++ {
				if !contains(excludeColIndices, col) {
					item.values = append(item.values, row[col])
				}
			}
			table.data[key] = &item
		}
	}
	return table
}


func getTableName(fileName string, sheetName string) string {
	if !multiSheetMode {
		return fileName
	}
	return sheetName
}
