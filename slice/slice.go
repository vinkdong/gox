package slice

import (
	"fmt"
	"reflect"
	"sort"
)

type Runtime struct {
}

func Difference(sliceA []string, sliceB []string) []string {
	diff := make([]string, 0)
	diffMap := make(map[string]int)

	for _, v := range sliceA {
		diffMap[v] = 1
	}
	for _, v := range sliceB {
		diffMap[v] = diffMap[v] - 1
	}

	for k, v := range diffMap {
		if v > 0 {
			diff = append(diff, k)
		}
	}
	return diff
}

func UnionString(slices ...[]string) []string {
	m := make(map[string]bool)
	for _, s := range slices {
		for _, val := range s {
			m[val] = true
		}
	}
	s := make([]string, len(m))
	i := 0
	for k, _ := range m {
		s[i] = k
		i++
	}
	return s
}

func IsContain(sliceA []string, element string) bool {
	for _, v := range sliceA {
		if v == element {
			return true
		}
	}
	return false
}

type KeyValue struct {
	Key   string
	Value string
}

func StructToFlatSlice(prefix string, obj interface{}, result *[]KeyValue) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		typ := val.Type()
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			tag := field.Tag.Get("json")
			if tag == "" || tag == "-" {
				continue
			}
			fieldValue := val.Field(i)

			newPrefix := tag
			if prefix != "" {
				newPrefix = prefix + "." + tag
			}

			StructToFlatSlice(newPrefix, fieldValue.Interface(), result)
		}
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface()
			newPrefix := fmt.Sprintf("%s[%d]", prefix, i)
			StructToFlatSlice(newPrefix, item, result)
		}
	default:
		*result = append(*result, KeyValue{Key: prefix, Value: fmt.Sprintf("%v", val.Interface())})
	}
}

func StructToSortedFlatSlice(prefix string, obj interface{}, result *[]KeyValue) {
	StructToFlatSlice(prefix, obj, result)
	sort.Slice(*result, func(i, j int) bool {
		return (*result)[i].Key < (*result)[j].Key
	})
}
