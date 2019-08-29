package main

import (
	"fmt"
	"reflect"
)

func i2s(data interface{}, out interface{}) error {
	val := reflect.ValueOf(data)
	res := reflect.ValueOf(out)
	if res.Kind() != reflect.Ptr {
		return fmt.Errorf("out is not a pointer")
	}
	return set(val, res.Elem())
}

func set(val, res reflect.Value) error {
	if val.Kind() == reflect.Interface {
		val = reflect.ValueOf(val.Interface())
	}
	// fmt.Println("val", val, "kind", val.Kind(), "res", res, "kind", res.Kind())
	switch res.Kind() {
	case reflect.Struct:
		if val.Kind() != reflect.Map || val.Type().Key().Kind() != reflect.String {
			return fmt.Errorf("type mismatch")
		}
		for _, k := range val.MapKeys() {
			if f := res.FieldByName(k.String()); !f.IsValid() {
				continue
			} else if err := set(val.MapIndex(k), f); err != nil {
				return err
			}
		}
	case reflect.Slice:
		if val.Kind() != reflect.Slice {
			return fmt.Errorf("type mismatch")
		}
		l := val.Len()
		res.Set(reflect.MakeSlice(res.Type(), l, l))
		for i := 0; i < l; i++ {
			if err := set(val.Index(i), res.Index(i)); err != nil {
				return err
			}
		}
	case reflect.String:
		if val.Kind() != reflect.String {
			return fmt.Errorf("type mismatch")
		}
		res.Set(val)
	case reflect.Int:
		if val.Kind() != reflect.Float64 {
			return fmt.Errorf("type mismatch")
		}
		res.SetInt(int64(val.Float()))
	case reflect.Bool:
		if val.Kind() != reflect.Bool {
			return fmt.Errorf("type mismatch")
		}
		res.Set(val)
	default:
		return fmt.Errorf("unsupported type")
	}
	return nil
}