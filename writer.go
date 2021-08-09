package main



func writeToJson(table *Table){
	m := make(map[interface{}]interface{})
	for key,val := range table.data{
		item := make(map[interface{}]interface{})
		for _,f := range table.fields{
			item[f.name] =
		}
		m[key] = item
	}
}