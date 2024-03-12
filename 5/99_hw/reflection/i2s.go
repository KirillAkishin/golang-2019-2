package main

import (
	"fmt"
	"reflect"
)

func main() {}

func i2s(data interface{}, result interface{}) error {
	outData := reflect.ValueOf(result)
	if outData.Kind() != reflect.Ptr {
		return fmt.Errorf("Error: outData is a copy, not a reference")
	}
	outData = outData.Elem()
	inData := reflect.ValueOf(data)

	inKind := inData.Kind()
	outKind := outData.Kind()
	switch outKind {
	case reflect.Slice:
		outData.Set(reflect.MakeSlice(outData.Type(), inData.Len(), inData.Len()))
		for i := 0; i < inData.Len(); i++ {
			i2s(inData.Index(i).Elem().Interface(), outData.Index(i).Addr().Interface())
		}
	case reflect.Struct:
		if inKind != reflect.Map {
			return fmt.Errorf("inData is no struct, but outData is struct")
		}
		for i := 0; i < outData.NumField(); i++ {
			outFieldName := outData.Type().Field(i).Name
			for _, inMapKey := range inData.MapKeys() {
				if inMapKey.Interface().(string) != outFieldName {
					continue
				}
				inMapValue := inData.MapIndex(inMapKey)
				outFieldValue := outData.FieldByName(outFieldName)

				inMapKind := inMapValue.Elem().Kind()
				outFieldKind := outData.FieldByName(outFieldName).Kind()
				switch inMapKind {
				case reflect.Bool:
					if outFieldKind != reflect.Bool {
						return fmt.Errorf("Incorrect struct's field. \"%v\" is not %v field",
							outFieldName, inMapKind)
					}
					value := inMapValue.Interface().(bool)
					outFieldValue.SetBool(value)
				case reflect.String:
					if outFieldKind != reflect.String {
						return fmt.Errorf("Incorrect struct's field. \"%v\" is not %v field",
							outFieldName, inMapKind)
					}
					value := inMapValue.Interface().(string)
					outFieldValue.SetString(value)
				case reflect.Float64, reflect.Float32, reflect.Int:
					if outFieldKind != reflect.Int {
						return fmt.Errorf("Incorrect struct's field. \"%v\" is not %v field",
							outFieldName, inMapKind)
					}
					value := int64(inMapValue.Interface().(float64))
					outFieldValue.SetInt(value)
				case reflect.Map:
					if outFieldKind != reflect.Struct {
						return fmt.Errorf("Incorrect struct's field: \"%v\"", outFieldName)
					}
					i2s(inMapValue.Interface(), outData.Field(i).Addr().Interface())
				case reflect.Struct:
					if outFieldKind != reflect.Struct {
						return fmt.Errorf("Incorrect struct's field. \"%v\" is not %v field",
							outFieldName, inMapKind)
					}
					i2s(inMapValue.Interface(), outData.Field(i).Addr().Interface())
				case reflect.Slice:
					if outFieldKind != reflect.Slice {
						return fmt.Errorf("Incorrect struct's field. \"%v\" is not %v field",
							outFieldName, inMapKind)
					}
					i2s(inMapValue.Interface(), outData.Field(i).Addr().Interface())
				default:
					return fmt.Errorf("Uknown kind: %v", inMapValue.Elem().Kind())
				}
			}
		}
	}
	return nil
}
