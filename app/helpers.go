/*
 * The MIT License (MIT)
 *
 * Copyright © 2023 Jonathan Newell <jonnewell@mac.com>
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
 * Filename: helpers.go
 * Last Modified: 11/14/23, 9:28 AM
 * Modified By: newellj
 *
 */

package app

import (
	"math/rand"
	"time"
)

// ----------------------------------------- Helper Functions ---------------------------------------------------------

func GetPointer[V any](value V) *V {
	return &value
}

func GetConfigString(key string) string {
	return CurrentState().config.Get(key)
}

func GetConfigBool(key string) bool {
	return CurrentState().config.GetBool(key)
}

func GetConfigIntWithDefault(key string, defaultValue int) int {
	return CurrentState().config.GetIntWithDefault(key, defaultValue)
}

func GetConfigList(key string) []string {
	return CurrentState().config.GetList(key)
}

func GetConfigFloat(key string) float64 {
	return CurrentState().config.GetFloat(key)
}

func RandomWaitMillis(max int) {
	//Variable wait to start them all at slightly different times
	rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(time.Duration(rand.Intn(max)) * time.Millisecond)
}
