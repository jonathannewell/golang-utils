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
 * Filename: regex.go
 * Last Modified: 10/26/23, 9:53 AM
 * Modified By: newellj
 *
 */

package golang_utils

import (
	"regexp"
	"sync"

	"github.com/apex/log"
)

type Regex struct {
	Pattern string
	Regex   *regexp.Regexp
	lock    sync.Mutex
}

func NewRegex(pattern string) *Regex {
	return &Regex{
		Pattern: pattern,
	}
}

func (r *Regex) IsValid() bool {
	if r.Regex != nil {
		return true
	}
	r.lock.Lock()
	var err error
	r.Regex, err = regexp.Compile(r.Pattern)
	if err != nil {
		log.Errorf("Invalid Regex [%s]. Details: %v", r.Pattern, err)
		return false
	}
	r.lock.Unlock()
	return true
}

func (r *Regex) Init() {

	if r.Regex != nil {
		return
	}

	r.lock.Lock()

	r.Regex = regexp.MustCompile(r.Pattern) //Will throw panic!

	r.lock.Unlock()

}

func (r *Regex) Matches(target string) bool {
	r.Init()
	return r.Regex.MatchString(target)
}
