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
 * Filename: files.go
 * Last Modified: 11/8/22, 9:45 AM
 * Modified By: newellj
 *
 *
 */

package io

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/apex/log"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

const PathSeparator = string(os.PathSeparator)

func GetWorkingDir() *FileInfo {
	dir, err := os.Getwd()
	CheckError(err, "Could not determine current working directory")
	return NewFileInfoFromPath(dir)
}

func GetFilesAtPath(path string, excludeDirs bool) (results []*FileInfo) {
	var absPath = GetAbsPath(path)
	log.Debugf("Loading Files found @ [%s]", path)

	contents, err := os.ReadDir(absPath)
	CheckError(err, "Error loading files @ [%s]", absPath)
	for _, entry := range contents {
		//Ignore directories if requested
		if excludeDirs && entry.IsDir() {
			continue
		}

		var info = NewFileInfo(entry.Name(), absPath)
		info.IsDir = entry.IsDir()
		log.Debugf("Found --> %s @ path [%s]", info.Name, info.BaseAbsPath)
		results = append(results, info)
	}
	return results
}

func GetDirsAtPath(path string, createIfNotExit bool) (results []*FileInfo) {
	var absPath = GetAbsPath(path)
	log.Debugf("Loading Dirs found @ [%s]", path)

	_, err := os.Stat(absPath)
	if os.IsNotExist(err) && createIfNotExit {
		os.Mkdir(absPath, 0755)
	}

	contents, err := os.ReadDir(absPath)
	CheckError(err, "Error loading files @ [%s]", absPath)
	for _, entry := range contents {
		if entry.IsDir() {
			var info = NewFileInfo(entry.Name(), absPath)
			log.Debugf("Found Dir --> %s @ path [%s]", info.Name, info.BaseAbsPath)
			results = append(results, info)
		}
		continue
	}
	return results
}

func RemoveDir(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		log.Errorf("Error deleting directory [%s] Details: %v", path, err)
	}
}

func CreateDirIfNotExist(directoryPath string) (existed bool, err error) {
	existed = true
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		if !strings.Contains(directoryPath, PathSeparator) {
			directoryPath = path.Join(".", directoryPath)
		}
		err = os.MkdirAll(directoryPath, 0755)
		existed = false
	}
	return existed, err
}

func CreateDir(path string) *FileInfo {
	_, err := CreateDirIfNotExist(path)
	CheckError(err, "Could not create directory [%s]", path)
	return NewFileInfoFromPath(path)
}

func GetAbsPath(path string) string {
	absPath, err := filepath.Abs(path)
	CheckError(err, "Unable to determine absolute path for [%s]", path)
	return absPath
}

const JSON string = "json"
const YAML string = "yaml"
const YML string = "yml"

type FileEncoder interface {
	Encode(v any) (err error)
}

type FileDecoder interface {
	Decode(v any) error
}

func FileIsJson(filename string) bool {
	if strings.HasSuffix(strings.ToLower(filename), JSON) {
		return true
	}
	return false
}

func FileIsYaml(filename string) bool {
	if strings.HasSuffix(strings.ToLower(filename), YAML) ||
		strings.HasSuffix(strings.ToLower(filename), YML) {
		return true
	}
	return false
}

func CheckAndCreateDir(directoryPath string) {
	var err error
	if _, err = os.Stat(directoryPath); os.IsNotExist(err) {
		if !strings.Contains(directoryPath, string(os.PathSeparator)) {
			directoryPath = path.Join(".", directoryPath)
		}
		CheckError(
			os.MkdirAll(directoryPath, os.ModePerm),
			"Failed creating directory directoryPath [%s]",
			directoryPath,
		)
		return
	}
	CheckError(err, "Failed creating directory at directoryPath [%s]", directoryPath)
}

func CreateFile(filename string, dir string) (file *os.File) {
	var err error
	file, err = os.Create(path.Join(GetAbsPath(dir), filename))
	CheckError(err, "Failed creating file [%s]!", filename)
	return file
}

func PathExists(filePath string) (exists bool) {
	exists = true

	_, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		exists = false
	} else {
		CheckError(err, "Error checking file [%s] for existence", filePath)
	}

	return
}

func DeleteFileOrDir(path string) {
	log.Debugf("Deleting file/dir [%s]", path)
	err := os.Remove(path)
	CheckError(err, fmt.Sprintf("Error deleting file or dir @ path [%s]", path))
}

func OpenFileAtPath(path string, homepath string) (reader *os.File, err error) {

	if strings.HasPrefix(path, "~") {
		path = strings.Replace(path, "~", homepath, -1)
	}

	abspath, _ := filepath.Abs(path)

	return os.Open(abspath)
}

func WriteStructsToFile(target any, filename string, path string, createDir bool) (exportCnt int) {

	if reflect.TypeOf(target).Kind() != reflect.Slice {
		ThrowError("WriteStructsToFile() requires a target that is a slice of structs!")
	}

	slice := reflect.ValueOf(target)

	log.Debugf("Writing [%d] structs to file [%s]\n", slice.Len(), filename)

	if createDir {
		CheckAndCreateDir(path)
	}

	file := CreateFile(filename, path)
	defer file.Close()

	var encoder FileEncoder

	if filepath.Ext(filename) == JSON {
		e := json.NewEncoder(file)
		e.SetIndent("", "    ")
		encoder = FileEncoder(e)
	} else {
		encoder = yaml.NewEncoder(file)
	}

	for i := 0; i < slice.Len(); i++ {

		err := encoder.Encode(slice.Index(i).Interface())

		if err != nil {
			fmt.Printf("Error Writing To File [%s]! Details: %v", file.Name(), err)
			return
		}

		exportCnt++
	}
	return
}

func WriteStructToFile(target any, filename string, path string) {
	CheckAndCreateDir(path)
	file := CreateFile(filename, path)
	defer file.Close()

	var encoder FileEncoder

	if filepath.Ext(filename) == JSON {
		e := json.NewEncoder(file)
		e.SetIndent("", "    ")
		encoder = FileEncoder(e)
	} else {
		encoder = yaml.NewEncoder(file)
	}

	err := encoder.Encode(target)

	CheckError(err, "Failed writing struct to file @ [%s]", path)
}

func WriteStringContentsToFile(path string, contents string, force bool) error {
	if PathExists(path) && !force {
		return nil
	}
	file, err := os.Create(path)
	if err == nil {
		defer file.Close()
		_, err = file.WriteString(contents)
	}

	return err
}

func LoadYamlFileToStruct[T any](path string) (*T, error) {
	var targetStruct *T
	var err error

	fileInfo := NewFileInfoFromPath(path)
	log.Infof("Attempting to read file [%s] @ [%s]", fileInfo.Name, fileInfo.BaseAbsPath)

	fileInfo.OpenForWriting(true)
	defer fileInfo.Close()

	targetStruct = new(T)
	yamlDecoder := yaml.NewDecoder(fileInfo.FileHandle)
	err = yamlDecoder.Decode(targetStruct)

	return targetStruct, err
}

func CreateUniqueTempDir(parentDir string) *FileInfo {
	//Get Guid (unique directory for this run
	p := os.TempDir()

	if IsNotEmpty(parentDir) {
		p = path.Join(p, parentDir)
	}

	return CreateDir(path.Join(p, uuid.New().String()))
}

func CreateTempFile(dir string, filename string) *FileInfo {
	file, err := os.CreateTemp(dir, filename)
	CheckError(err, "Unable to create temporary file [%s] @ path [@s]", filename, dir)

	return &FileInfo{
		Name:        path.Base(file.Name()),
		BaseAbsPath: dir,
		FileHandle:  file,
	}
}
