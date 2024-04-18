/*
 * The MIT License (MIT)
 *
 * Copyright © 2022 Jonathan Newell <jonnewell@mac.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this analyzer and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, Merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NON-INFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 *
 * Filename: maps.go
 * Last Modified: 11/4/22, 9:58 AM
 * Modified By: newellj
 *
 *
 */

package golang_utils

import (
	"fmt"
	"reflect"
	"strings"
)

func HasKey[K comparable, V any](key K, target map[K]V) bool {
	if _, found := target[key]; found {
		return true
	}
	return false
}

func AddIfNotExist[K comparable, V any](key K, value V, target map[K]V) {
	if !HasKey(key, target) {
		target[key] = value
	}
}

func Get[K comparable, V any](key K, target map[K]V) V {
	if value, has := target[key]; has {
		return value
	}
	return *new(V)
}

func Remove[K comparable, V any](key K, target map[K]V) bool {
	if HasKey(key, target) {
		delete(target, key)
		return true
	}
	return false
}

func GetString(key string, target map[string]any) string {
	found := Get(key, target)
	if found != nil {
		result, ok := found.(string)
		if ok {
			return result
		} else {
			return ""
		}
	}
	return ""
}

func Merge[K comparable, V any](baseMap map[K]V, overridingMap map[K]V) map[K]V {
	if baseMap == nil {
		if overridingMap == nil {
			return make(map[K]V)
		}
		return overridingMap
	}

	if overridingMap == nil {
		return baseMap
	}

	//I have both...get to work
	for k, v := range baseMap {
		if _, found := overridingMap[k]; !found {
			overridingMap[k] = v
		}
	}
	return overridingMap
}

func AddAll[K comparable, V any](toMap map[K]V, fromMap map[K]V) {
	for k, v := range fromMap {
		toMap[k] = v
	}
}

func Count[K comparable, V any](target map[K]V) (cnt int) {
	if target == nil {
		return
	}

	for _, _ = range target {
		cnt++
	}
	return
}

func Print[K comparable, V any](target map[K]V, msg string, args ...any) {
	PrintPadded(target, 0, msg, args...)
}

func PrintPadded[K comparable, V any](target map[K]V, depth int, msg string, args ...any) {
	var builder strings.Builder
	padd := strings.Repeat(" ", depth)
	for k, v := range target {
		builder.WriteString(fmt.Sprintf("%s%v: %s", padd, k, getValue(v, depth+1)))
	}

	fmt.Printf(fmt.Sprintf("%s%s", padd, msg), args...)
	fmt.Printf("\n%s\n", builder.String())
}
func getValue(value any, depth int) string {
	actualValue := reflect.ValueOf(value)
	if actualValue.Kind() == reflect.Map {
		var builder strings.Builder
		var spacing = strings.Repeat(" ", depth)
		builder.WriteString(fmt.Sprintf("\n"))
		iter := actualValue.MapRange()
		for iter.Next() {
			builder.WriteString(
				fmt.Sprintf(
					"%s%v: %s",
					spacing,
					iter.Key().String(),
					getValue(iter.Value().Interface(), depth+1),
				),
			)
		}
		return builder.String()
	} else if actualValue.Kind() == reflect.Slice {
		var builder strings.Builder
		var spacing = strings.Repeat(" ", depth-1)
		builder.WriteString(fmt.Sprintf("\n"))
		switch actualValue.Type().Elem().Kind() {
		case reflect.String:
			for _, v := range value.([]string) {
				builder.WriteString(fmt.Sprintf("%s- %s", spacing, getValue(v, depth)))
			}
		case reflect.Interface:
			for _, v := range value.([]any) {
				builder.WriteString(fmt.Sprintf("%s- %s", spacing, getValue(v, depth)))
			}
		default:
			fmt.Printf("Attempt to iterate/print unknown slice type [%s]\n", actualValue.Elem().Kind().String())
		}

		return builder.String()
	}

	if actualValue.Kind() == reflect.String {
		return fmt.Sprintf("%s\n", value)
	} else {
		return fmt.Sprintf("%v\n", value)
	}
}
