package airtable

import (
	"log"
	"reflect"
	"strings"
)

type TableNamer interface {
	TableName() string
}

func tableName(entity any) string {
	t := reflect.ValueOf(entity).Type()

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}

	if namer, ok := reflect.New(t).Interface().(TableNamer); ok {
		return namer.TableName()
	}

	return t.Name() + "s"
}

const RecordIDField = "RecordID"

func fields(entity any) map[string]any {
	v := reflect.ValueOf(entity)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		log.Fatal("entity must be struct type or ptr")
	}

	m := map[string]any{}

	for i := 0; i < v.NumField(); i++ {
		fType := v.Type().Field(i)
		fName := fType.Name
		fValue := v.Field(i)

		if fName == RecordIDField {
			continue
		}

		omitempty := strings.Contains(fType.Tag.Get("json"), "omitempty")
		if omitempty && fValue.IsZero() {
			continue
		}

		m[fName] = fValue.Interface()

		if fType.Type.Kind() == reflect.Ptr {
			fType.Type = fType.Type.Elem()
			fValue = fValue.Elem()
		}
	}

	return m
}

func nonZeroFields(entity any) map[string]any {
	f := fields(entity)

	for k, v := range f {
		if reflect.ValueOf(v).IsZero() {
			delete(f, k)
		}
	}

	return f
}

func id(entity any) string {
	return reflect.ValueOf(entity).FieldByName(RecordIDField).String()
}
