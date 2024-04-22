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
 * Filename: properties.go
 * Last Modified: 11/8/22, 5:17 PM
 * Modified By: newellj
 *
 *
 */

package collections

import (
	"fmt"
	coreio "io"
	"path/filepath"
	"strings"

	"github.com/apex/log"
	"github.com/jonathannewell/golang-utils/app"
	"github.com/jonathannewell/golang-utils/io"
	"gopkg.in/yaml.v3"
)

type Properties map[string]any

type PropertyManager struct {
	contents map[string]*OverridableProperties
}

type OverridableProperties struct {
	FileName string
	KeyPath  string
	Path     string
	Contents Properties
	Children []*OverridableProperties
	Parent   *OverridableProperties
}

func NewOverridableProperties(filename string, path string) *OverridableProperties {
	newConfig := &OverridableProperties{
		FileName: filename,
		Path:     io.GetAbsPath(path),
		Contents: make(Properties),
		Children: make([]*OverridableProperties, 0),
	}
	newConfig.loadPropertiesYamlToMap()
	return newConfig
}

func NewOverriddenProperties(properties *OverridableProperties) *OverridableProperties {

	return &OverridableProperties{
		FileName: properties.FileName,
		Path:     properties.Path,
		KeyPath:  properties.GetKeyPath(),
		Contents: properties.GetOverriddenProperties(),
	}
}

// Build Relationships
func (c *OverridableProperties) Add(child *OverridableProperties) {
	if child != nil {
		if !Contains(c.Children, child) {
			log.Debugf("Adding Child [%s] to Parent [%s]", child.GetKeyPath(), c.GetKeyPath())
			c.Children = append(c.Children, child)
			child.Parent = c
		}
	}
}

func (p Properties) Remove(key string) bool {
	return RemoveFromMap(key, p)
}

func (c *OverridableProperties) GetOverriddenProperties() Properties {
	if c.Parent != nil {
		return MapMerge(c.Parent.GetOverriddenProperties(), c.Contents)
	}
	return c.Contents
}

func (c *OverridableProperties) GetKeyPath() string {
	if c.KeyPath != "" {
		return c.KeyPath
	}
	if c.Parent == nil {
		return filepath.Base(c.Path)
	} else {
		path := c.Parent.GetKeyPath()
		return fmt.Sprintf("%s.%s", path, filepath.Base(c.Path))
	}
}

func (c *OverridableProperties) GetGrandestParent() *OverridableProperties {
	if c.Parent == nil {
		return c
	}
	return c.Parent.GetGrandestParent()
}

func (c *OverridableProperties) loadPropertiesYamlToMap() {
	LoadPropertiesMapYamlToMap(c.FileName, c.Path, &c.Contents)
}

/***********************************************************************************************************************
												PROPERTY MANAGER API
***********************************************************************************************************************/

func NewPropertyManager() *PropertyManager {
	return &PropertyManager{
		contents: make(map[string]*OverridableProperties),
	}
}

func (p *PropertyManager) Has(path string) bool {
	return MapHasKey(path, p.contents)
}

func (p *PropertyManager) Add(path string, props *OverridableProperties) {
	log.Debugf(
		"Adding props collection [core: %d total: %d path: %s",
		MapCount(props.Contents),
		MapCount(props.GetOverriddenProperties()),
		path,
	)
	p.contents[path] = props
}

func (p *PropertyManager) Get(path string) *OverridableProperties {
	if p.Has(path) {
		props, _ := p.contents[path]
		return props
	}
	return nil
}

func (p *PropertyManager) GetSetsWithSuffix(suffix string) (found []*OverridableProperties) {
	for k, v := range p.contents {
		if strings.HasSuffix(k, suffix) {
			found = append(found, v)
		}
	}
	return
}

func (p *PropertyManager) GetFirstSetWithSuffix(suffix string) *OverridableProperties {
	set := p.GetSetsWithSuffix(suffix)
	if set != nil && len(set) > 0 {
		return set[0]
	}
	return nil
}

func (p *PropertyManager) SetCount() int {
	return MapCount(p.contents)
}

func LoadPropertiesMapYamlToMap(filename, path string, properties *Properties) {
	log.Debugf("Attempting to read [%s] @ %s", filename, path)
	fileInfo := io.NewFileInfo(filename, io.GetAbsPath(path))
	if fileInfo.Exists() {
		fileInfo.Open()
		defer fileInfo.Close()
		yamlDecoder := yaml.NewDecoder(fileInfo.FileHandle)
		err := yamlDecoder.Decode(properties)
		if err != coreio.EOF { //Ignore empty files or files with nothing but comments...not real errors~!
			app.CheckError(err, "Failed Reading [%s]", fileInfo.AbsFilePath())
		}
	}
}
