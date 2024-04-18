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
 * Filename: cmd.go
 * Last Modified: 11/14/23, 8:27 AM
 * Modified By: newellj
 *
 */

package golang_utils

import (
	"github.com/spf13/cobra"
)

type CmdFunc func(cmd *cobra.Command, args []string)

type CmdConfig struct {
	use            string
	short          string
	long           string
	aliases        []string
	pre            CmdFunc
	run            CmdFunc
	post           CmdFunc
	pPre           CmdFunc
	pPost          CmdFunc
	args           cobra.PositionalArgs
	enableTracking bool
	version        string
}

func CommandBuilder(use string) *CmdConfig {
	return &CmdConfig{
		use:            use,
		enableTracking: true,
	}
}

func (cc *CmdConfig) SetShortDescription(desc string) *CmdConfig {
	cc.short = desc
	return cc
}

func (cc *CmdConfig) SetLongDescription(desc string) *CmdConfig {
	cc.long = desc
	return cc
}

func (cc *CmdConfig) SetAliases(alias ...string) *CmdConfig {
	cc.aliases = alias
	return cc
}

func (cc *CmdConfig) SetRun(cmdFunc CmdFunc) *CmdConfig {
	cc.run = cmdFunc
	return cc
}

func (cc *CmdConfig) SetPreRun(cmdFunc CmdFunc) *CmdConfig {
	cc.pre = cmdFunc
	return cc
}

func (cc *CmdConfig) SetPersistentPreRun(cmdFunc CmdFunc) *CmdConfig {
	cc.pPre = cmdFunc
	return cc
}

func (cc *CmdConfig) SetPostRun(cmdFunc CmdFunc) *CmdConfig {
	cc.post = cmdFunc
	return cc
}

func (cc *CmdConfig) SetPersistentPostRun(cmdFunc CmdFunc) *CmdConfig {
	cc.pPost = cmdFunc
	return cc
}

func (cc *CmdConfig) SetArgValidations(argValidations cobra.PositionalArgs) *CmdConfig {
	cc.args = argValidations
	return cc
}

func (cc *CmdConfig) SetVersion(version string) *CmdConfig {
	cc.version = version
	return cc
}

func (cc *CmdConfig) DisableTracking() *CmdConfig {
	cc.enableTracking = false
	return cc
}

func (cc *CmdConfig) EnableTracking() *CmdConfig {
	cc.enableTracking = true
	return cc
}

func (cc *CmdConfig) Build() *cobra.Command {
	return newCommand(cc)
}

func newCommand(config *CmdConfig) *cobra.Command {

	newCmd := &cobra.Command{
		Use:               config.use,
		Short:             config.short,
		Aliases:           config.aliases,
		Args:              config.args,
		PreRun:            config.pre,
		PersistentPreRun:  config.pPre,
		PostRun:           config.post,
		PersistentPostRun: config.pPost,
		Run:               config.run,
		Version:           config.version,
	}

	if config.enableTracking {
		if config.pre != nil {
			newCmd.PreRun = func(cmd *cobra.Command, args []string) {
				CurrentState().TrackCmd(cmd, args)
				config.pre(cmd, args)
			}
		} else {
			newCmd.PreRun = CurrentState().TrackCmd
		}
	}

	return newCmd
}

func MakeFlagRequired(cmd *cobra.Command, flagName string) {
	CheckError(
		cmd.MarkFlagRequired(flagName),
		"Error marking flag [%s] as required for the %s cmd",
		flagName,
		GetFullCmdName(cmd),
	)
}

func GetFullCmdName(cmd *cobra.Command) string {
	if cmd == nil {
		return "????"
	}
	if cmd.Parent() == nil {
		return cmd.Name()
	}
	return GetFullCmdName(cmd.Parent()) + "." + cmd.Name()
}
