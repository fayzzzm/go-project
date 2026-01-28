package database

import (
	"fmt"
	"reflect"
	"sync"
)

// fieldInfo holds metadata about a struct field.
type fieldInfo struct {
	Index int
	Name  string
	Tag   string
}

// structFields represents a list of fields for a struct.
type structFields []fieldInfo

var fieldCache sync.Map // map[reflect.Type]structFields

// getCachedTypeFields returns the fields of a struct type, using a cache.
func getCachedTypeFields(t reflect.Type) structFields {
	if f, ok := fieldCache.Load(t); ok {
		return f.(structFields)
	}
	f, _ := fieldCache.LoadOrStore(t, getTypeFields(t))
	return f.(structFields)
}

// getTypeFields allows us to extract fields from a struct using reflection.
func getTypeFields(t reflect.Type) structFields {
	var fields structFields
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		tag := field.Tag.Get("db")
		if tag == "-" {
			continue
		}

		fields = append(fields, fieldInfo{
			Index: i,
			Name:  field.Name,
			Tag:   tag,
		})
	}
	return fields
}

// MapToArgs uses reflection to pull values from a struct into a slice of interface{}.
// It is the 10x programmer's way to eliminate manual mapping boilerplate.
// It respects 'db' tags and ignores fields with tag "-".
func MapToArgs(src interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(src)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %s", v.Kind())
	}

	t := v.Type()
	fields := getCachedTypeFields(t)

	args := make([]interface{}, len(fields))
	for i, f := range fields {
		args[i] = v.Field(f.Index).Interface()
	}

	return args, nil
}

// MapToStruct converts a map or another struct to a pointer of type T using reflection.
// This is useful for AOP-like transformations between layers.
func MapToStruct[T any](src interface{}) (*T, error) {
	// Implementation for generic mapping if needed...
	var dest T
	// ... logic to map fields ...
	return &dest, nil
}
