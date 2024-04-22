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
 * Filename: db.go
 * Last Modified: 11/14/23, 8:27 AM
 * Modified By: newellj
 *
 */

package app

import (
	"github.com/apex/log"
	"github.com/glebarez/sqlite"
	"github.com/jonathannewell/golang-utils/io"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PersistenceContext struct {
	DB     *gorm.DB //Do I really need this?
	DBFile *io.FileInfo
	config *PersistenceConfig
}

type PersistenceConfig struct {
	Name     string
	Path     string
	Entities []any
}

func NewPersistenceContext(config *PersistenceConfig) *PersistenceContext {
	dbFileInfo := io.NewFileInfo(config.Name, config.Path)
	return &PersistenceContext{
		DBFile: dbFileInfo,
		config: config,
	}
}

func NewPersistenceConfig(dbName, path string, entities []any) *PersistenceConfig {
	return &PersistenceConfig{
		Name:     dbName,
		Path:     path,
		Entities: entities,
	}
}

func (c *PersistenceContext) OpenDB() {

	log.Debugf("Connecting to DB [%s] @ %s", c.config.Name, c.DBFile.BaseAbsPath)
	var err error

	c.DB, err = gorm.Open(
		sqlite.Open(c.DBFile.AbsFilePath()),
		&gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Silent),
			PrepareStmt:            false,
			SkipDefaultTransaction: true,
		},
	)

	CheckError(err, "Error opening Database @ [%s]", c.DBFile.AbsFilePath())
	c.InitDB()
	c.PopulateReferenceData()
}

func (c *PersistenceContext) InitDB() {
	CheckError(
		c.DB.AutoMigrate(
			c.config.Entities...,
		),
		"error initializing DB schema",
	)

	//c.DB.Raw("PRAGMA journal_mode=WAL;")

}

func (c *PersistenceContext) Save(value any) {
	CheckError(c.DB.Save(value).Error, "Error storing %T information to db!", value)
}

func (c *PersistenceContext) SaveOmitting(value any, omitted ...string) {
	CheckError(
		c.DB.Session(&gorm.Session{}).Omit(omitted...).Save(value).Error,
		"Error storing %T information to db!",
		value,
	)
}

func (c *PersistenceContext) SaveFull(value any) {
	CheckError(
		c.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(value).Error,
		"Error storing %T information to db!",
		value,
	)
}

func (c *PersistenceContext) Create(value any) {
	CheckError(
		c.DB.Session(&gorm.Session{FullSaveAssociations: true}).Create(value).Error,
		"Error storing %T information to db!",
		value,
	)
}

func (c *PersistenceContext) OmitFields(omitted ...string) *gorm.DB {
	return c.DB.Omit(omitted...)
}

func (c *PersistenceContext) PopulateReferenceData() {
}
