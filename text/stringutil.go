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
 * Filename: stringutil.go
 * Last Modified: 11/4/22, 11:32 AM
 * Modified By: newellj
 *
 *
 */

package text

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const dashRune = rune('-')

var specials = map[string]string{
	"\"": "\\\"",
	"'":  "''",
}

func LineWrap(s string, length int) string {

	re := regexp.MustCompile("\r?\n")
	s = re.ReplaceAllString(s, "")

	var buffer bytes.Buffer
	var n_1 = length - 1
	var l_1 = len(s) - 1
	for i, rune := range s {
		buffer.WriteRune(rune)
		if i%length == n_1 && i != l_1 {
			buffer.WriteRune('\n')
		}
	}
	return buffer.String()
}

func EscapeSpecials(target string) string {

	if target != "" {

		//Replace any escapes first...
		target = strings.Replace(target, "\\", "\\\\", -1)

		for special, replacement := range specials {
			target = strings.Replace(target, special, replacement, -1)
		}

	}

	return target
}

func Escape(value string) string {
	return strings.Replace(value, "'", "''", -1)
}

func RemoveLineEndings(value string) string {
	value = strings.Replace(value, "\n", "", -1)
	value = strings.Replace(value, "\r", "", -1)
	return strings.TrimSpace(value)
}

func RemoveNonPrintingUnicode(value string) string {
	return strings.Map(
		func(r rune) rune {
			if unicode.IsPrint(r) {
				return r
			}
			return -1
		}, value,
	)
}

func IsDash(character rune) bool {
	return dashRune == character
}

func IsInt(value string) bool {
	_, err := strconv.Atoi(value)
	return err == nil
}

func ToString(value int) string {
	return strconv.Itoa(value)
}

func StripPathSeparators(value string) string {
	value = strings.Replace(value, "\\", "", -1)
	value = strings.Replace(value, "/", "", -1)
	return value
}

func IsEmpty(value string) bool {
	return strings.TrimSpace(value) == ""
}

func IsNotEmpty(value string) bool {
	return !IsEmpty(value)
}
