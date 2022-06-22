package converter

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
)

func Decode(i interface{}, d interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovering in Decode %+v", r)
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = fmt.Errorf("unknown panic: %+v", x)
			}
		}
	}()
	v := reflect.Indirect(reflect.ValueOf(i))
	dv := reflect.ValueOf(d)
	for i := 0; i < v.NumField(); i++ {
		fname := v.Type().Field(i).Name
		snakeCaseFname := ToSnakeCase(fname)
		fvalue := dv.MethodByName("Get").Call([]reflect.Value{reflect.ValueOf(snakeCaseFname)})[0]
		if fvalue.IsNil() {
			continue
		}
		element := reflect.New(v.Field(i).Type())
		extractElement(element.Interface(), fvalue.Interface())
		v.Field(i).Set(reflect.Indirect(element))
	}
	return err
}

func extractElement(ie interface{}, id interface{}) {
	v := reflect.Indirect(reflect.ValueOf(ie))
	dataValue := reflect.ValueOf(id)
	if v.Kind() == reflect.Slice {
		for i := 0; i < dataValue.Len(); i++ {
			element := reflect.New(reflect.TypeOf(v.Interface()).Elem())
			extractElement(element.Interface(), dataValue.Index(i).Interface())
			v.Set(reflect.Append(v, reflect.Indirect(element)))
		}
	} else if v.Kind() == reflect.Map {
		keyType := v.Type().Key()
		valType := v.Type().Elem()
		mp := reflect.MakeMap(reflect.MapOf(keyType, valType))
		mpk := dataValue.MapKeys()
		for idx := range mpk {
			element := reflect.New(valType)
			extractElement(element.Interface(), dataValue.MapIndex(mpk[idx]).Elem().Interface())
			key := reflect.New(keyType)
			getTypeFromString(key, mpk[idx].String())
			mp.SetMapIndex(reflect.Indirect(key), reflect.Indirect(element))
		}
		v.Set(mp)
	} else if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			fname := v.Type().Field(i).Name
			snakeCase := ToSnakeCase(fname)
			element := reflect.New(v.Field(i).Type())
			extractElement(element.Interface(), dataValue.MapIndex(reflect.ValueOf(snakeCase)).Interface())
			v.Field(i).Set(reflect.Indirect(element))
		}
	} else {
		switch v.Kind() {
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint8:
			v.SetInt(dataValue.Int())
		case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v.SetUint(uint64(dataValue.Int()))
		case reflect.Float32, reflect.Float64:
			v.SetFloat(dataValue.Float())
		case reflect.Bool:
			v.SetBool(dataValue.Bool())
		default:
			v.Set(dataValue)
		}
	}
}

func getTypeFromString(key reflect.Value, s string) {
	v := reflect.Indirect(key)
	switch v.Kind() {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		i, _ := strconv.ParseInt(s, 10, 64)
		v.SetInt(i)
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		i, _ := strconv.ParseUint(s, 10, 64)
		v.SetUint(i)
	case reflect.Float32, reflect.Float64:
		i, _ := strconv.ParseFloat(s, 64)
		v.SetFloat(i)
	case reflect.Bool:
		i, _ := strconv.ParseBool(s)
		v.SetBool(i)
	case reflect.String:
		v.SetString(s)
	}
}
