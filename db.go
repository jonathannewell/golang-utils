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

package golang_utils

import (
	"context"
	"fmt"
	"github.com/apex/log"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// LogLevel log level
type LogLevel int
type dBJournalMode struct {
	OFF      int
	WAL      int
	DELETE   int
	TRUNCATE int
	PERSIST  int
	MEMORY   int
}

var JournalMode = dBJournalMode{
	OFF:      jrnl_off,
	WAL:      jrnl_wal,
	DELETE:   jrnl_delete,
	TRUNCATE: jrnl_truncate,
	PERSIST:  jrnl_persist,
	MEMORY:   jrnl_memory,
}

const (
	jrnl_off int = iota
	jrnl_wal
	jrnl_delete
	jrnl_truncate
	jrnl_persist
	jrnl_memory
)

// Interface logger interface
type Interface interface {
	LogMode(LogLevel) Interface
	Info(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}

type PersistenceContext struct {
	DB     *gorm.DB //Do I really need this?
	DBFile *FileInfo
	config *PersistenceConfig
}

type PersistenceConfig struct {
	Name     string
	Path     string
	Entities []any
	// GORM perform single create, update, delete operations in transactions by default to ensure database data integrity
	// You can disable it by setting `SkipDefaultTransaction` to true
	SkipDefaultTransaction bool
	// FullSaveAssociations full save associations
	FullSaveAssociations bool
	// Logger
	Logger logger.Interface
	// NowFunc the function to be used when creating a new timestamp
	NowFunc func() time.Time
	// DryRun generate sql without execute
	DryRun bool
	// PrepareStmt executes the given query in cached statement
	PrepareStmt bool
	// DisableAutomaticPing
	DisableAutomaticPing bool
	// DisableForeignKeyConstraintWhenMigrating
	DisableForeignKeyConstraintWhenMigrating bool
	// IgnoreRelationshipsWhenMigrating
	IgnoreRelationshipsWhenMigrating bool
	// DisableNestedTransaction disable nested transaction
	DisableNestedTransaction bool
	// AllowGlobalUpdate allow global update
	AllowGlobalUpdate bool
	// QueryFields executes the SQL query with all fields of the table
	QueryFields bool
	// CreateBatchSize default create batch size
	CreateBatchSize int
	// TranslateError enabling error translation
	TranslateError bool
	JournalMode    int
}

func NewPersistenceContext(config *PersistenceConfig) *PersistenceContext {
	dbFileInfo := NewFileInfo(config.Name, config.Path)
	return &PersistenceContext{
		DBFile: dbFileInfo,
		config: config,
	}
}

func NewPersistenceConfig(dbName, path string, entities []any) *PersistenceConfig {
	return &PersistenceConfig{
		Name:                   dbName,
		Path:                   path,
		Entities:               entities,
		Logger:                 logger.Default.LogMode(logger.Silent),
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		DryRun:                                   false,
		IgnoreRelationshipsWhenMigrating:         false,
		DisableNestedTransaction:                 true,
		AllowGlobalUpdate:                        true,
		DisableForeignKeyConstraintWhenMigrating: true,
		DisableAutomaticPing:                     true,
		QueryFields:                              true,
		CreateBatchSize:                          100,
		TranslateError:                           false,
		FullSaveAssociations:                     true,
		JournalMode:                              jrnl_delete,
	}
}

func (c *PersistenceContext) OpenDB() {

	log.Debugf("Connecting to DB [%s] @ %s", c.config.Name, c.DBFile.BaseAbsPath)
	var err error

	c.DB, err = gorm.Open(
		sqlite.Open(c.DBFile.AbsFilePath()),
		c.config.gormConfig(),
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

func (config *PersistenceConfig) gormConfig() *gorm.Config {
	return &gorm.Config{
		SkipDefaultTransaction:                   config.SkipDefaultTransaction,
		FullSaveAssociations:                     config.FullSaveAssociations,
		Logger:                                   config.Logger,
		NowFunc:                                  config.NowFunc,
		DryRun:                                   config.DryRun,
		PrepareStmt:                              config.PrepareStmt,
		DisableAutomaticPing:                     config.DisableAutomaticPing,
		DisableForeignKeyConstraintWhenMigrating: config.DisableForeignKeyConstraintWhenMigrating,
		IgnoreRelationshipsWhenMigrating:         config.IgnoreRelationshipsWhenMigrating,
		DisableNestedTransaction:                 config.DisableNestedTransaction,
		AllowGlobalUpdate:                        config.AllowGlobalUpdate,
		QueryFields:                              config.QueryFields,
		CreateBatchSize:                          config.CreateBatchSize,
		TranslateError:                           config.TranslateError,
	}

}

func (c *PersistenceContext) setJournalMode() {
	var mode string
	switch c.config.JournalMode {
	case jrnl_off:
		mode = "OFF"
	case jrnl_wal:
		mode = "WAL"
	case jrnl_truncate:
		mode = "TRUNCATE"
	case jrnl_persist:
		mode = "PERSIST"
	case jrnl_memory:
		mode = "MEMORY"
	default:
		mode = "DELETE"
	}

	c.DB.Raw(fmt.Sprintf("PRAGMA journal_mode=%s;", mode))
}
