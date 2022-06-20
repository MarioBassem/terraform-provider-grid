package converter

import (
	"reflect"
)

func Decode(i interface{}, d interface{}) {
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
			element := reflect.New(dataValue.MapIndex(mpk[idx]).Elem().Type())
			extractElement(element.Interface(), dataValue.MapIndex(mpk[idx]).Elem().Interface())
			mp.SetMapIndex(mpk[idx], reflect.Indirect(element))
		}
		v.Set(mp)
	} else if v.Kind() == reflect.Struct {
		mpk := dataValue.MapKeys()
		for idx := range mpk {
			camelCaseFname := ToCamelCase(mpk[idx].Interface().(string))
			element := reflect.New(v.FieldByName(camelCaseFname).Type())
			extractElement(element.Interface(), dataValue.MapIndex(mpk[idx]).Interface())
			v.FieldByName(camelCaseFname).Set(reflect.Indirect(element))
		}
	} else {
		if v.Kind() == dataValue.Kind() {
			v.Set(dataValue)
			return
		}

		switch v.Kind() {
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			v.SetInt(int64(dataValue.Uint()))
		case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v.SetUint(uint64(dataValue.Int()))
		}
	}
}
