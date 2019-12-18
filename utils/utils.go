package utils

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

func GetApplicationDir() (ApplicationDir string) {
	ApplicationDir, _ = os.Getwd()
	ApplicationDir = strings.Replace(ApplicationDir, "\\", "/", -1)
	ApplicationDir += "/"
	return ApplicationDir
}

//MD5加密 大写
func MD5EncToUpper(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

//MD5加密 小写
func MD5EncToLower(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return strings.ToLower(hex.EncodeToString(h.Sum(nil)))
}

//排序工具
func quickSort(arr []interface{}, start, end int, compete func(a interface{}, b interface{}) int8) {
	if start < end {
		i, j := start, end
		key := arr[(start+end)/2]
		for i <= j {
			for compete(arr[i], key) == -1 {
				i++
			}
			for compete(arr[j], key) == 1 {
				j--
			}
			if i <= j {
				arr[i], arr[j] = arr[j], arr[i]
				i++
				j--
			}
		}
		if start < j {
			quickSort(arr, start, j, compete)
		}
		if end > i {
			quickSort(arr, i, end, compete)
		}
	}
}
func Sort(a []interface{}, compete func(a interface{}, b interface{}) int8) {
	if len(a) < 2 {
		return
	}
	quickSort(a, 0, len(a)-1, compete)
}

func Copy(src interface{}) interface{} {
	if src == nil {
		return nil
	}
	original := reflect.ValueOf(src)
	cpy := reflect.New(original.Type()).Elem()
	copyRecursive(original, cpy)

	return cpy.Interface()
}

func copyRecursive(src, dst reflect.Value) {
	switch src.Kind() {
	case reflect.Ptr:
		originalValue := src.Elem()

		if !originalValue.IsValid() {
			return
		}
		dst.Set(reflect.New(originalValue.Type()))
		copyRecursive(originalValue, dst.Elem())

	case reflect.Interface:
		if src.IsNil() {
			return
		}
		originalValue := src.Elem()
		copyValue := reflect.New(originalValue.Type()).Elem()
		copyRecursive(originalValue, copyValue)
		dst.Set(copyValue)

	case reflect.Struct:
		t, ok := src.Interface().(time.Time)
		if ok {
			dst.Set(reflect.ValueOf(t))
			return
		}
		for i := 0; i < src.NumField(); i++ {
			if src.Type().Field(i).PkgPath != "" {
				continue
			}
			copyRecursive(src.Field(i), dst.Field(i))
		}

	case reflect.Slice:
		if src.IsNil() {
			return
		}
		dst.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Cap()))
		for i := 0; i < src.Len(); i++ {
			copyRecursive(src.Index(i), dst.Index(i))
		}

	case reflect.Map:
		if src.IsNil() {
			return
		}
		dst.Set(reflect.MakeMap(src.Type()))
		for _, key := range src.MapKeys() {
			originalValue := src.MapIndex(key)
			copyValue := reflect.New(originalValue.Type()).Elem()
			copyRecursive(originalValue, copyValue)
			copyKey := Copy(key.Interface())
			dst.SetMapIndex(reflect.ValueOf(copyKey), copyValue)
		}

	default:
		dst.Set(src)
	}
}

func Str2Bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
