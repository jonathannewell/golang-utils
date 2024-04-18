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
 * Filename: maps_test.go
 * Last Modified: 11/4/22, 11:27 AM
 * Modified By: newellj
 *
 *
 */

package golang_utils

import (
	"testing"

	"github.com/tj/assert"
)

const (
	key1     = "key1"
	value1   = "value1"
	key2     = "key2"
	value2   = "value2"
	key3     = "key3"
	value3   = "value3"
	val2over = "overridden-val2"
	key4     = "key4"
	value4   = "value4"
)

func TestCountKeys(t *testing.T) {
	assert.Equal(t, 2, Count(createMap()), "Count is wrong", nil)
}

func TestMapMergeKeysInBaseAddedToResult(t *testing.T) {

	//Simple Test
	result := Merge(createMap(), nil)
	assert.Equal(t, 2, Count(result), "Merged Count is wrong", nil)
	assert.Equal(t, value1, result[key1], "Key 1 value wrong", nil)
	assert.Equal(t, value2, result[key2], "Key 2 value wrong", nil)
}

func TestMapMergeWithOverriddenValues(t *testing.T) {
	baseMap := createMap()
	baseMap[key4] = value4
	overrideMap := createMap()
	overrideMap[key3] = value3
	overrideMap[key2] = val2over

	result := Merge(baseMap, overrideMap)
	assert.Equal(t, 4, Count(result), "Merged Count is wrong", nil)
	assert.Equal(t, value1, result[key1], "Key 1 value wrong", nil)
	assert.Equal(t, val2over, result[key2], "Key 2 value wrong", nil)
	assert.Equal(t, value3, result[key3], "Key 3 value wrong", nil)
	assert.Equal(t, value4, result[key4], "Key 4 value wrong", nil)
}

func TestMapGetString(t *testing.T) {
	testMap := createMap()
	result := GetString(key1, testMap)

	assert.NotNil(t, result, "Get string should return the requested value or empty string not nil!")
	assert.Equal(t, value1, result, "Key1 should return Value1")

	result = GetString("non-existent-key", testMap)
	assert.NotNil(t, result, "Get string should return the requested value or empty string not nil!")
	assert.Equal(t, "", result, "Non existent key should return empty string")
}

func TestMapGet(t *testing.T) {
	testMap := createMap()
	result := Get(key1, testMap)

	assert.NotNil(t, result, "Get string should return the requested value or empty string not nil!")
	assert.Equal(t, value1, result, "Key1 should return Value1")

	result = GetString("non-existent-key", testMap)
	assert.NotNil(t, result, "Get string should return the requested value or empty string not nil!")
	assert.Equal(t, "", result, "Non existent key should return empty string")

}

func createMap() Properties {
	baseMap := make(Properties)
	baseMap[key1] = value1
	baseMap[key2] = value2

	return baseMap
}
