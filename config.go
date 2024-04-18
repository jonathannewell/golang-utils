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
 * Filename: config.go
 * Last Modified: 11/14/23, 8:27 AM
 * Modified By: newellj
 *
 */

package golang_utils

import (
	"fmt"
	"path"

	"github.com/apex/log"
	"github.com/spf13/viper"
)

type Configuration struct {
	filename    string
	LoadedFrom  string
	writeInHome bool
	RunID       uint
	defaults    Properties
}

func NewConfiguration(filename string, writeInHome bool, defaults Properties) *Configuration {
	return &Configuration{
		filename:    filename,
		writeInHome: writeInHome,
		defaults:    defaults,
	}
}

func (c *Configuration) Get(propertyName string) string {
	return viper.GetString(propertyName)
}

func (c *Configuration) GetList(propertyName string) []string {
	return viper.GetStringSlice(propertyName)
}

func (c *Configuration) GetBool(propertyName string) bool {
	return viper.GetBool(propertyName)
}

func (c *Configuration) GetFloat(propertyName string) float64 {
	return viper.GetFloat64(propertyName)
}

func (c *Configuration) GetInt(propertyName string) int {
	return viper.GetInt(propertyName)
}

func (c *Configuration) GetIntWithDefault(propertyName string, defaultValue int) int {
	val := viper.GetInt(propertyName)

	if defaultValue != 0 && val == 0 {
		return defaultValue
	}
	return val
}

func (c *Configuration) GetMap(propertyName string) map[string]string {
	return viper.GetStringMapString(propertyName)
}

func (c *Configuration) HasResource(name string, property string) bool {
	return Contains(c.GetList(property), name)
}

func (c *Configuration) Update(propertyName string, value string) {
	viper.Set(propertyName, value)
	c.Write("Error adding/updating config property [%s]", propertyName)
}

func (c *Configuration) UpdateList(propertyName string, value string) {
	current := c.GetList(propertyName)
	if current == nil {
		current = make([]string, 0)
	}
	if !Contains(current, value) {
		current = append(current, value)
	}
	viper.Set(propertyName, current)
	c.Write("Error adding/updating config property [%s]", propertyName)
}

func (c *Configuration) UpdateMap(propertyName string, key string, value string) {
	current := c.GetMap(propertyName)
	if current == nil {
		current = make(map[string]string)
	}

	current[key] = value
	viper.Set(propertyName, current)
	c.Write("Error adding/updating key [%s] in map property [%s]", key, propertyName)
}

func (c *Configuration) DeleteFromMap(propertyName string, key string) {
	current := c.GetMap(propertyName)
	if current == nil {
		return
	}
	if _, ok := current[key]; ok {
		delete(current, key)
	}
	viper.Set(propertyName, current)
	c.Write("Error deleting key [%s] map property [%s]", key, propertyName)
}

func (c *Configuration) PrintMapProperty(propertyName string) {
	current := c.GetMap(propertyName)
	if current == nil {
		fmt.Printf("no property [%s] currently defined\n", propertyName)
	} else {
		fmt.Printf("---- %s(s) ----\n", propertyName)
		for k, v := range current {
			fmt.Printf(" %s --> %s\n", k, v)
		}
	}
}

func (c *Configuration) Load(cfgFile string) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		c.setConfigPaths()
	}
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	_ = c.readConfig()
}

func (c *Configuration) Write(msg string, args ...string) {
	log.Debugf("writing configuration file")
	err := viper.WriteConfig()
	CheckError(err, msg, args)
}

func (c *Configuration) Default() {
	viper.Reset()
	c.setUpDefaults()
	c.setConfigPaths()
	c.createConfigFile()
	c.Print()
}

func (c *Configuration) Print() {
	Print(viper.AllSettings(), "------ Portfolio Viewer Config Properties [%s] ------", c.LoadedFrom)
}

func (c *Configuration) readConfig() (err error) {
	if err = viper.ReadInConfig(); err == nil {
		c.LoadedFrom = viper.ConfigFileUsed()
		log.Debugf("Loaded config from [%s]", c.LoadedFrom)
		if c.setUpDefaults() {
			c.Write("updating config with missing defaults")
		}
	} else {
		log.Debugf("No config file (%s) found in `.` or `user-home-dir`", c.filename)
		c.setUpDefaults()
		c.LoadedFrom = "defaults"
	}
	return err
}

func (c *Configuration) createConfigFile() {
	CheckError(viper.WriteConfigAs(c.configPath()), "Unable to write config")
	_ = c.readConfig()
}

func (c *Configuration) configPath() string {
	if c.writeInHome {
		return path.Join(CurrentState().homeDir, c.filename)
	} else {
		return path.Join(CurrentState().workDir, c.filename)
	}
}

func (c *Configuration) setUpDefaults() (shouldWrite bool) {
	if c.defaults != nil {
		for p, v := range c.defaults {
			defaulted := c.applyDefault(p, v)
			shouldWrite = shouldWrite || defaulted
		}
	}
	return
}

func (c *Configuration) setConfigPaths() {
	viper.AddConfigPath(".") //Look in current directory first!
	viper.AddConfigPath(CurrentState().homeDir)
	viper.SetConfigType("yaml")
	viper.SetConfigName(c.filename)
}

func (c *Configuration) applyDefault(property string, value any) bool {
	if viper.Get(property) == nil {
		viper.SetDefault(property, value)
		return true
	}
	return false
}
