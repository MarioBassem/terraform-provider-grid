package converter

import (
	"reflect"
)

func Encode(obj interface{}, d interface{}) {
	v := reflect.ValueOf(obj)
	dv := reflect.ValueOf(d)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		element := reflect.New(getType(field))
		snakeCase := ToSnakeCase(v.Type().Field(i).Name)
		if dv.MethodByName("Get").Call([]reflect.Value{reflect.ValueOf(snakeCase)})[0].IsNil() {
			continue
		}
		convert(element.Interface(), field.Interface())
		dv.MethodByName("Set").Call([]reflect.Value{reflect.ValueOf(snakeCase), reflect.ValueOf(reflect.Indirect(element).Interface())})

	}
}

func convert(ie interface{}, dataInterface interface{}) {
	v := reflect.Indirect(reflect.ValueOf(ie))
	dataValue := reflect.ValueOf(dataInterface)
	if dataValue.Kind() == reflect.Slice {
		for i := 0; i < dataValue.Len(); i++ {
			element := reflect.New(getType(dataValue.Index(i)))
			convert(element.Interface(), dataValue.Index(i).Interface())
			v.Set(reflect.Append(v, reflect.Indirect(element)))
		}
	} else if dataValue.Kind() == reflect.Struct {
		mp := reflect.MakeMap(getType(dataValue))
		for i := 0; i < dataValue.NumField(); i++ {
			snakeCase := ToSnakeCase(dataValue.Type().Field(i).Name)
			element := reflect.New(getType(dataValue.Field(i)))
			convert(element.Interface(), dataValue.Field(i).Interface())
			mp.SetMapIndex(reflect.ValueOf(snakeCase), reflect.Indirect(element))
		}
		v.Set(mp)
	} else if dataValue.Kind() == reflect.Map {
		mp := reflect.MakeMap(getType(dataValue))
		for _, k := range dataValue.MapKeys() {
			element := reflect.New(getType(dataValue.MapIndex(k)))
			convert(element.Interface(), dataValue.MapIndex(k).Interface())
			mp.SetMapIndex(k, reflect.Indirect(element))
		}
		v.Set(mp)
	} else {
		v.Set(dataValue)
	}

}

func getType(v reflect.Value) reflect.Type {
	if v.Kind() == reflect.Struct || v.Kind() == reflect.Map {
		return reflect.TypeOf(map[string]interface{}{})
	} else if v.Kind() == reflect.Slice {
		return reflect.TypeOf([]interface{}{})
	} else {
		return v.Type()
	}
}
