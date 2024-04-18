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
 * Filename: fileinfo.go
 * Last Modified: 10/25/23, 9:06 AM
 * Modified By: newellj
 *
 */

package golang_utils

import (
	"io"
	"os"
	"path/filepath"
)

type FileInfo struct {
	Name        string
	BaseAbsPath string
	Info        os.FileInfo
	IsDir       bool
	FileHandle  *os.File
}

func NewFileInfo(name string, absPath string) *FileInfo {
	return &FileInfo{
		Name:        name,
		BaseAbsPath: absPath,
	}
}

func NewFileInfoFromPath(path string) *FileInfo {
	absPath := GetAbsPath(path)
	return &FileInfo{
		Name:        filepath.Base(absPath),
		BaseAbsPath: filepath.Dir(absPath),
	}
}

func (fi *FileInfo) GetFileInfo() os.FileInfo {
	info, err := os.Stat(filepath.Join(fi.BaseAbsPath, fi.Name))
	CheckError(err, "Unable to read FileInfo for File [%s] @ [%s]", fi.Name, fi.BaseAbsPath)
	return info
}

func (fi *FileInfo) ReadFully() []byte {

	fi.Open()
	defer fi.Close()

	var size int
	size64 := fi.GetFileInfo().Size()
	if int64(int(size64)) == size64 {
		size = int(size64)
	}
	size++ // one byte for final read at EOF

	// If a file claims a small size, read at least 512 bytes.
	// In particular, files in Linux's /proc claim size 0 but
	// then do not work right if read in small pieces,
	// so an initial read of 1 byte would not work correctly.
	if size < 512 {
		size = 512
	}

	data := make([]byte, 0, size)
	for {
		if len(data) >= cap(data) {
			d := append(data[:cap(data)], 0)
			data = d[:len(data)]
		}
		n, err := fi.FileHandle.Read(data[len(data):cap(data)])
		data = data[:len(data)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			CheckError(err, "Unable to read file contents from [%s]", fi.AbsFilePath())
			return data
		}
	}
}

func (fi *FileInfo) Exists() bool {
	return PathExists(fi.AbsFilePath())
}

func (fi *FileInfo) Create() *FileInfo {
	fi.FileHandle = CreateFile(fi.Name, fi.BaseAbsPath)
	return fi
}

func (fi *FileInfo) Open() {
	var err error
	if fi.FileHandle == nil {
		flags := os.O_CREATE | os.O_RDWR
		fi.FileHandle, err = os.OpenFile(fi.AbsFilePath(), flags, 0755)
		CheckError(err, "Error Opening File [%s]", fi.AbsFilePath())
	}
}

func (fi *FileInfo) OpenForWriting(truncate bool) {
	var err error
	if fi.FileHandle == nil {
		flags := os.O_CREATE | os.O_RDWR
		if truncate {
			flags |= os.O_TRUNC
		}
		fi.FileHandle, err = os.OpenFile(fi.AbsFilePath(), flags, 0755)
		CheckError(err, "Error Opening File [%s] for writing", fi.AbsFilePath())
	}
}

func (fi *FileInfo) Close() {
	if fi.FileHandle != nil {
		LogError(fi.FileHandle.Close(), "Error closing file [%s]", fi.AbsFilePath())
		fi.FileHandle = nil
	}
}

func (fi *FileInfo) MoveToPath(path string) {
	targetPath := filepath.Join(fi.BaseAbsPath, path)
	CheckError(
		os.Rename(fi.AbsFilePath(), targetPath),
		"Error Moving/Renaming [%s] to [%s]",
		fi.AbsFilePath(),
		targetPath,
	)
}

func (fi *FileInfo) AbsFilePath() string {
	return filepath.Join(fi.BaseAbsPath, fi.Name)
}

func (fi *FileInfo) WriteFile(data []byte) {
	if fi.FileHandle == nil {
		ThrowError(
			"Unable to write file [%s]. File does not exists or no handle has been established",
			fi.AbsFilePath(),
		)
	}
	CheckError(os.WriteFile(fi.FileHandle.Name(), data, 0644), "Unable to write file [%s]", fi.AbsFilePath())
}
