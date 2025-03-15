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
 * Filename: lists.go
 * Last Modified: 11/4/22, 11:34 AM
 * Modified By: newellj
 *
 *
 */

package golang_utils

import (
	"strings"
)

type Comparable interface {
	Compare(value any) bool
}

func ContainsElement[T Comparable](list []T, value T) bool {
	for _, itemValue := range list {
		if itemValue.Compare(value) {
			return true
		}
	}
	return false
}

func Contains[T comparable](list []T, value T) bool {
	for _, itemValue := range list {
		if itemValue == value {
			return true
		}
	}
	return false
}

func ContainsItemsWithPrefix(list []string, value string) bool {
	for _, itemValue := range list {
		if strings.HasPrefix(value, itemValue) {
			return true
		}
	}
	return false
}

func SliceIsEmpty[T comparable](list []T) bool {
	return len(list) == 0
}

func Reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
