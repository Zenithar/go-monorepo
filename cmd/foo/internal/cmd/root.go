// Licensed to go-monorepo under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. go-monorepo licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	iconfig "github.com/Zenithar/go-monorepo/cmd/foo/internal/config"

	"github.com/Zenithar/go-monorepo/build/version"
	"github.com/Zenithar/go-monorepo/pkg/config"
	configcmd "github.com/Zenithar/go-monorepo/pkg/config/cmd"
	"github.com/Zenithar/go-monorepo/pkg/log"
)

// -----------------------------------------------------------------------------

const (
	cmdPrefix = "FOO"
)

var (
	cfgFile string
	conf    = &iconfig.Configuration{}
)

// -----------------------------------------------------------------------------

// RootCmd describes root command of the tool
var mainCmd = &cobra.Command{
	Use:   "foo",
	Short: "Foo microservice",
}

func init() {
	mainCmd.Flags().StringVar(&cfgFile, "config", "", "config file")

	mainCmd.AddCommand(version.Command())
	mainCmd.AddCommand(configcmd.NewConfigCommand(conf, cmdPrefix))
}

// -----------------------------------------------------------------------------

// Execute main command
func Execute() error {
	initConfig()
	return mainCmd.Execute()
}

// -----------------------------------------------------------------------------

func initConfig() {
	if err := config.Load(conf, cmdPrefix, cfgFile); err != nil {
		log.Bg().Fatal("Unable to load settings", zap.Error(err))
	}
}
