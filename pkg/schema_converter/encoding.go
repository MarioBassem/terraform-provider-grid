package converter

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
)

func Encode(obj interface{}, d interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovering in Enocde %+v", r)
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
	v := reflect.ValueOf(obj)
	dv := reflect.ValueOf(d)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		snakeCase := ToSnakeCase(v.Type().Field(i).Name)
		rd := dv.MethodByName("Get").Call([]reflect.Value{reflect.ValueOf(snakeCase)})[0]
		if rd.IsNil() {
			log.Printf("nil")
			continue
		}
		element := reflect.New(getType(field.Type()))
		convert(element.Interface(), field.Interface())
		dv.MethodByName("Set").Call([]reflect.Value{reflect.ValueOf(snakeCase), reflect.ValueOf(reflect.Indirect(element).Interface())})

	}
	return err
}

func convert(ie interface{}, dataInterface interface{}) {
	v := reflect.Indirect(reflect.ValueOf(ie))
	dataValue := reflect.ValueOf(dataInterface)
	if dataValue.Kind() == reflect.Slice {
		for i := 0; i < dataValue.Len(); i++ {
			element := reflect.New(getType(dataValue.Index(i).Type()))
			convert(element.Interface(), dataValue.Index(i).Interface())
			v.Set(reflect.Append(v, reflect.Indirect(element)))
		}
	} else if dataValue.Kind() == reflect.Struct {
		mp := reflect.MakeMap(getType(dataValue.Type()))
		for i := 0; i < dataValue.NumField(); i++ {
			snakeCase := ToSnakeCase(dataValue.Type().Field(i).Name)
			element := reflect.New(getType(dataValue.Field(i).Type()))
			convert(element.Interface(), dataValue.Field(i).Interface())
			mp.SetMapIndex(reflect.ValueOf(snakeCase), reflect.Indirect(element))
		}
		v.Set(mp)
	} else if dataValue.Kind() == reflect.Map {
		mp := reflect.MakeMap(getType(dataValue.Type()))
		for _, k := range dataValue.MapKeys() {
			element := reflect.New(getType(dataValue.MapIndex(k).Type()))
			convert(element.Interface(), dataValue.MapIndex(k).Interface())
			newK := ensureStringType(k)
			mp.SetMapIndex(newK, reflect.Indirect(element))
		}
		v.Set(mp)
	} else {

		switch dataValue.Kind() {
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			v.SetInt(dataValue.Int())
		case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v.SetInt(int64(dataValue.Uint()))
		case reflect.Float32, reflect.Float64:
			v.SetFloat(dataValue.Float())
		default:
			v.Set(dataValue)
		}
	}

}

func ensureStringType(k reflect.Value) reflect.Value {
	var newK reflect.Value
	switch k.Kind() {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		newK = reflect.ValueOf(strconv.FormatInt(k.Int(), 10))
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		newK = reflect.ValueOf(strconv.FormatUint(k.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		newK = reflect.ValueOf(strconv.FormatFloat(k.Float(), 'f', -1, 64))
	case reflect.Bool:
		newK = reflect.ValueOf(strconv.FormatBool(k.Bool()))
	case reflect.String:
		newK = k
	}
	return newK
}

func getType(v reflect.Type) reflect.Type {
	var ret reflect.Type
	switch v.Kind() {
	case reflect.Struct, reflect.Map:
		ret = reflect.TypeOf(map[string]interface{}{})
	case reflect.Slice, reflect.Array:
		ret = reflect.TypeOf([]interface{}{})
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		ret = reflect.TypeOf(int(1))
	case reflect.Float32, reflect.Float64:
		ret = reflect.TypeOf(float32(1))
	default:
		ret = v
	}
	return ret
}
