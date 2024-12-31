/*
 * The MIT License (MIT)
 *
 * Copyright © 2024 Jonathan Newell <jonnewell@mac.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
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
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 *
 * Filename: errors.go
 * Last Modified: 2/1/24, 3:42 PM
 * Modified By: newellj
 *
 */

package golang_utils

import (
	"fmt"
	"strings"

	"github.com/apex/log"
)

func LogError(err error, msg string, args ...any) {
	if err != nil {
		log.Error(formatError(err, msg, args...))
	}
}

func LogErrorTrace(err error, msg string, args ...any) {
	if err != nil {
		log.Trace(formatError(err, msg, args...))
	}
}

func CheckError(err error, msg string, args ...any) {
	if err != nil {
		log.Error(formatError(err, msg, args...))
		SendErrorEvent("util.CheckError", err, msg, args...)
	}
}

func ThrowError(msg string, args ...any) {
	err := fmt.Errorf(msg, args...)
	log.Error(err.Error())
	SendErrorEvent("util.ThrowError", err, "Programmatically Thrown Error")
	panic(err)
}

func formatError(err error, msg string, args ...any) (errorString string) {
	if err == nil {
		errorString = fmt.Sprintf(msg, args...)
		return
	}
	if strings.Contains(msg, "%") && len(args) > 0 {
		errorString = fmt.Sprintf("%s. Details: %s", fmt.Sprintf(msg, args...), err)
	} else {
		errorString = fmt.Sprintf("%s. Details: %s", msg, err)
	}
	return
}
