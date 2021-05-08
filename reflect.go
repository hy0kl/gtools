package gtools

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

// ChangeValueByColName 通过反射更改结构体的值
func ChangeValueByColName(aExt interface{}, colName string, dstValue interface{}) interface{} {
	v := reflect.ValueOf(aExt)
	v = v.Elem()

	col := v.FieldByName(colName)
	if !col.IsValid() {
		log.Println(fmt.Sprintf("[changeValueByColName] reflect err, aExt: %#v, colName: %s", aExt, colName))
		return aExt
	}

	col.Set(reflect.ValueOf(dstValue))
	return aExt
}

func SetByFields(aExt interface{}, colName string, dstValue interface{}) error {
	aa := reflect.ValueOf(aExt).Elem()

	field := FieldByName(aa, colName)

	if field.IsValid() {
		switch field.Type().Kind() {
		case reflect.Slice:
			//log.Printf("dstValue:%v", dstValue)
			v := reflect.Append(field, reflect.ValueOf(dstValue))
			field.Set(v)
		default:
			field.Set(reflect.ValueOf(dstValue))
		}

	} else {
		err := fmt.Errorf("[SetByFields] field is not valid, colName: %v, aExt: %v", colName, aExt)
		log.Println(err)
		return err
	}
	return nil
}

// 通过反射取得特定字段
func FieldByName(rv reflect.Value, colName string) reflect.Value {
	//log.Printf("colName:%v ", colName)
	if colName == "" {
		return reflect.Value{}
	}

	index := strings.Index(colName, ".")
	if index != -1 {
		name := colName[0:index]
		//log.Printf("[FieldByName] name:%v", name)

		field := rv.FieldByName(name)

		//log.Printf("[FieldByName] name:%v field:%#v", name, field)
		return FieldByName(field, colName[index+1:])
	}
	return rv.FieldByName(colName)
}
