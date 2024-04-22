/*
 * The MIT License (MIT)
 *
 * Copyright © 2022 Jonathan Newell <jonnewell@mac.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this analyzer and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
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
 * Filename: templatefuncs.go
 * Last Modified: 11/17/22, 12:49 PM
 * Modified By: newellj
 *
 *
 */

package text

import (
	"fmt"
	"strconv"
)

const (
	QUOTE_FUNC_NAME      = "quote"
	EMPTY_FUNC_NAME      = "empty"
	NOT_EMPTY_FUNC_NAME  = "notEmpty"
	OR_DEFAULT_FUNC_NAME = "orDefault"
	OR_EMPTY_FUNC_NAME   = "orEmpty"
	TERNARY_FUNC_NAME    = "ternary"
)

func Quote(input any) string {
	switch val := input.(type) {
	case string:
		return strconv.Quote(val)
	case int:
		return strconv.Quote(strconv.Itoa(val))
	case int64:
		return strconv.Quote(strconv.FormatInt(val, 10))
	default:
		return strconv.Quote(fmt.Sprintf("%v", val))
	}
}

func Empty(input any) bool {
	switch val := input.(type) {
	case string:
		if val == "" {
			return true
		}
		return false
	default:
		if val == nil {
			return true
		}
		return false
	}
}

func NotEmpty(input any) bool {
	return !Empty(input)
}

func OrDefault(defaulVal string, input any) any {
	if Empty(input) {
		return defaulVal
	}
	return input
}

func OrEmpty(input any) any {
	return OrDefault("", input)
}

func Ternary(condition bool, trueOut any, falseOut any) any {
	if condition {
		return trueOut
	}
	return falseOut
}
