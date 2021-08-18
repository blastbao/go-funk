package funk

import (
	"fmt"
	"reflect"
	"strings"
)

// Filter iterates over elements of collection, returning an array of
// all elements predicate returns truthy for.
func Filter(arr interface{}, predicate interface{}) interface{} {
	// 检查 arr 是否可迭代
	if !IsIteratee(arr) {
		panic("First parameter must be an iteratee")
	}

	// 检查 predicate 是否是函数，且出/入参数目分别为 1/1
	if !IsFunction(predicate, 1, 1) {
		panic("Second argument must be function")
	}

	// 函数类型
	funcValue := reflect.ValueOf(predicate)
	funcType := funcValue.Type()

	// 检查函数出参类型是否为 bool
	if funcType.Out(0).Kind() != reflect.Bool {
		panic("Return argument should be a boolean")
	}

	arrValue := reflect.ValueOf(arr)
	arrType := arrValue.Type()

	// Get slice type corresponding to array type
	// 数组类型
	resultSliceType := reflect.SliceOf(arrType.Elem())

	// MakeSlice takes a slice kind type, and makes a slice.
	// 构造返回数组
	resultSlice := reflect.MakeSlice(resultSliceType, 0, 0)

	// 遍历 arr
	for i := 0; i < arrValue.Len(); i++ {
		// 取 arr[i]
		elem := arrValue.Index(i)
		// 执行 predicate(arr[i])
		result := funcValue.Call([]reflect.Value{elem})[0].Interface().(bool)
		// 如果返回 true 就 append 到 result 中
		if result {
			resultSlice = reflect.Append(resultSlice, elem)
		}
	}

	// 返回结果数组
	return resultSlice.Interface()
}

// Find iterates over elements of collection, returning the first
// element predicate returns truthy for.
func Find(arr interface{}, predicate interface{}) interface{} {
	_, val := FindKey(arr, predicate)
	return val
}

// FindKey iterates over elements of collection, returning the first
// element of an array and random of a map which predicate returns truthy for.
func FindKey(arr interface{}, predicate interface{}) (matchKey, matchEle interface{}) {
	// 检查 arr 是否可迭代
	if !IsIteratee(arr) {
		panic("First parameter must be an iteratee")
	}
	// 检查 predicate 是否是函数，且出/入参数目分别为 1/1
	if !IsFunction(predicate, 1, 1) {
		panic("Second argument must be function")
	}

	// 函数类型
	funcValue := reflect.ValueOf(predicate)
	funcType := funcValue.Type()

	// 检查函数出参类型是否为 bool
	if funcType.Out(0).Kind() != reflect.Bool {
		panic("Return argument should be a boolean")
	}


	arrValue := reflect.ValueOf(arr)
	var keyArrs []reflect.Value

	isMap := arrValue.Kind() == reflect.Map
	if isMap {
		keyArrs = arrValue.MapKeys()
	}

	for i := 0; i < arrValue.Len(); i++ {
		var (
			elem reflect.Value
			key  reflect.Value
		)
		if isMap {
			key = keyArrs[i]
			elem = arrValue.MapIndex(key)
		} else {
			key = reflect.ValueOf(i)
			elem = arrValue.Index(i)
		}

		// 执行 predicate(arr[i]) 返回 true/false
		result := funcValue.Call([]reflect.Value{elem})[0].Interface().(bool)
		if result {
			return key.Interface(), elem.Interface()
		}
	}

	return nil, nil
}

// IndexOf gets the index at which the first occurrence of value is found in array or return -1
// if the value cannot be found
func IndexOf(in interface{}, elem interface{}) int {
	inValue := reflect.ValueOf(in)

	elemValue := reflect.ValueOf(elem)

	inType := inValue.Type()

	if inType.Kind() == reflect.String {
		return strings.Index(inValue.String(), elemValue.String())
	}

	if inType.Kind() == reflect.Slice {
		equalTo := equal(elem)
		for i := 0; i < inValue.Len(); i++ {
			if equalTo(reflect.Value{}, inValue.Index(i)) {
				return i
			}
		}
	}

	return -1
}

// LastIndexOf gets the index at which the last occurrence of value is found in array or return -1
// if the value cannot be found
func LastIndexOf(in interface{}, elem interface{}) int {
	inValue := reflect.ValueOf(in)

	elemValue := reflect.ValueOf(elem)

	inType := inValue.Type()

	if inType.Kind() == reflect.String {
		return strings.LastIndex(inValue.String(), elemValue.String())
	}

	if inType.Kind() == reflect.Slice {
		length := inValue.Len()

		equalTo := equal(elem)
		for i := length - 1; i >= 0; i-- {
			if equalTo(reflect.Value{}, inValue.Index(i)) {
				return i
			}
		}
	}

	return -1
}

// Contains returns true if an element is present in a iteratee.
func Contains(in interface{}, elem interface{}) bool {
	inValue := reflect.ValueOf(in)
	elemValue := reflect.ValueOf(elem)
	inType := inValue.Type()

	switch inType.Kind() {
	case reflect.String:
		return strings.Contains(inValue.String(), elemValue.String())
	case reflect.Map:
		equalTo := equal(elem, true)
		for _, key := range inValue.MapKeys() {
			if equalTo(key, inValue.MapIndex(key)) {
				return true
			}
		}
	case reflect.Slice, reflect.Array:
		equalTo := equal(elem)
		for i := 0; i < inValue.Len(); i++ {
			if equalTo(reflect.Value{}, inValue.Index(i)) {
				return true
			}
		}
	default:
		panic(fmt.Sprintf("Type %s is not supported by Contains, supported types are String, Map, Slice, Array", inType.String()))
	}

	return false
}

// Every returns true if every element is present in a iteratee.
func Every(in interface{}, elements ...interface{}) bool {
	for _, elem := range elements {
		if !Contains(in, elem) {
			return false
		}
	}
	return true
}

// Some returns true if atleast one element is present in an iteratee.
func Some(in interface{}, elements ...interface{}) bool {
	for _, elem := range elements {
		if Contains(in, elem) {
			return true
		}
	}
	return false
}
