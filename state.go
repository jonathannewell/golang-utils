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
 * Filename: state.go
 * Last Modified: 11/14/23, 8:27 AM
 * Modified By: newellj
 *
 */

package golang_utils

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/spf13/cobra"
)

// ApplicationState - States of an application.
type ApplicationState int

const (
	// Starting - The application is in the process of bootstrapping and is not fully operational
	Starting = iota
	// Running - The application is fully loaded/bootstrapped and is executing normally
	Running
	// Stopping - The application is stopping or shutting down
	Stopping
	// Stopped - The application is fully stopped or shutdown. Are you ever going to actually be able to see this state?
	Stopped
	// Paused - The application is in a state of suspended animation. It is up and can be restarted instantly but is not consuming CPU
	Paused
	// Errored - The application has suffered an error or exception. It is not running normally.
	Errored
)

var (
	appState *State
	lock     sync.Mutex
)

type State struct {
	sync.Mutex         //Used for locking state to ensure safe concurrent access
	state              ApplicationState
	properties         Properties
	config             *Configuration
	cmd                *cobra.Command
	startTime          time.Time
	stopTime           time.Time
	verbose            bool
	useHomeDir         bool
	user               string
	homeDir            string
	workDir            string
	tempDir            string
	dataDir            string
	configFileName     string
	appName            string
	version            string
	buildDate          string
	commitSha          string
	errored            bool
	PersistenceContext *PersistenceContext
}

type Event struct {
	Name     string
	Type     string
	Category string
	Data     Properties
}

func CurrentState() *State {
	if appState == nil {
		lock.Lock()
		defer lock.Unlock()
		if appState == nil {
			appState = &State{
				state:     Starting,
				startTime: time.Now(),
			}
			appState.init()
		}
	}
	return appState
}

func (s *State) init() {
	log.SetHandler(cli.New(os.Stdout))
	log.SetLevel(log.InfoLevel)

	u, err := user.Current()
	CheckError(err, "Unable to read/access users info")
	s.user = u.Username
	s.homeDir = GetAbsPath(u.HomeDir)
	s.workDir = GetAbsPath(GetWorkingDir().AbsFilePath())
}

func (s *State) SetState(newState ApplicationState) *State {

	//Capture new state
	s.state = newState

	switch newState {
	case Errored:
		s.errored = true
	case Stopped:
		s.stopTime = time.Now()
	}
	return s
}

func (s *State) EnablePersistence(config *PersistenceConfig) *State {
	s.PersistenceContext = NewPersistenceContext(config)
	s.PersistenceContext.OpenDB()
	return s
}

func (s *State) EnableTempDir() *State {
	s.tempDir = CreateUniqueTempDir(s.appName).AbsFilePath()
	return s
}

func (s *State) EnableBaseDataDir(path string) *State {
	s.dataDir = CreateDir(path).AbsFilePath()
	return s
}

func (s *State) SetTempDir(path string) *State {
	s.tempDir = path
	return s
}

func (s *State) TempDir() string {
	return s.tempDir
}

func (s *State) StartTime() time.Time {
	return s.startTime
}

func (s *State) StopTime() time.Time {
	return s.stopTime
}

func (s *State) User() string {
	return s.user
}

func (s *State) HomeDir() string {
	return s.homeDir
}

func (s *State) WorkDir() string {
	return s.workDir
}

func (s *State) Errored() bool {
	return s.errored
}

func (s *State) Verbose() bool { return s.verbose }

func (s *State) Duration() time.Duration {
	if s.state == Stopped {
		return s.stopTime.Sub(s.startTime)
	}
	return time.Now().Sub(s.startTime)
}

func (s *State) CheckError(terminate bool, err error, msg string, args ...any) {
	if err != nil {
		s.state = Errored
		if terminate {
			CheckError(err, msg, args)
		} else {
			LogError(err, msg, args)
		}
	}
}

func (s *State) AppendDataDir(appendPath string) *State {
	s.dataDir = path.Join(s.dataDir, appendPath)
	_, err := CreateDirIfNotExist(s.dataDir)
	CheckError(err, "Failure creating data dir path [%s] for app [%s]", s.dataDir, s.appName)
	return s
}

func (s *State) DataDir() string {
	return s.dataDir
}

func (s *State) SetLogging(verbose bool) *State {
	if verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	return s
}

func (s *State) UseHomeDir(uhd bool) *State {
	s.useHomeDir = uhd
	return s
}

func (s *State) SetAppName(name string) *State {
	s.appName = name
	return s
}

func (s *State) SetVersion(version string) *State {
	s.version = version
	return s
}

func (s *State) SetCommitSha(sha string) *State {
	s.commitSha = sha
	return s
}

func (s *State) SetBuildDate(date string) *State {
	s.buildDate = date
	return s
}

func (s *State) SetConfigFileName(name string) *State {
	s.configFileName = name
	return s
}

func (s *State) SetCmd(cmd *cobra.Command) *State {
	s.cmd = cmd
	return s
}

func (s *State) TrackCmd(cmd *cobra.Command, _ []string) {
	s.SetCmd(cmd)
}

func (s *State) FullCommand() string {
	if s.cmd != nil {
		return GetFullCmdName(s.cmd)
	}

	return "???"
}

func (s *State) UpdateState(name string, value any) *State {

	//Mutate State
	s.Lock()
	_, has := s.properties[name]
	s.properties[name] = value
	s.Unlock()

	//Send State Change Events for interested parties
	if has {
		//State Update
	} else {
		//State Created
	}

	return s
}

func (s *State) RemoveFromState(name string) *State {
	s.Lock()
	if s.properties.Remove(name) {
		//Send Deleted Event
	}
	s.Unlock()
	return s
}

func (s *State) UpdateConfigProperty(property, value string) *State {
	s.config.Update(property, value)
	return s
}

func (s *State) InitConfig(defaults Properties) *State {
	s.config = NewConfiguration(s.configFile(), s.useHomeDir, defaults)
	s.config.Load("")
	return s
}

func (s *State) DefaultConfig() *State {
	s.config.Default()
	return s
}

func (s *State) PrintConfig() *State {
	s.config.Print()
	return s
}

func (s *State) Config() *Configuration {
	if s.config == nil {
		s.Lock()
		if s.config == nil {
			s.InitConfig(Properties{})
		}
		s.Unlock()
	}
	return s.config
}

func (s *State) PrintVersion() {
	fmt.Printf("Name: %s\n", s.appName)
	fmt.Printf("Version: %s\n", s.version)
	fmt.Printf("Commit SHA: %s\n", s.commitSha)
	fmt.Printf("Build Date: %s\n\n", s.buildDate)
}

func (s *State) configFile() string {
	if IsEmpty(s.configFileName) {
		//Don't double lock!
		appName := s.applicationName()
		s.Lock()
		if IsEmpty(s.configFileName) {
			s.configFileName = fmt.Sprintf(".%s", appName)
		}
		s.Unlock()
	}
	return s.configFileName
}

func (s *State) applicationName() string {
	if IsEmpty(s.appName) {
		s.Lock()
		if IsEmpty(s.appName) {
			s.appName = "unknown(???)"
		}
		s.Unlock()
	}
	return s.appName
}
