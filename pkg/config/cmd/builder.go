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
	"fmt"
	"sort"
	"strings"

	"github.com/Zenithar/go-monorepo/pkg/flags"
	"github.com/Zenithar/go-monorepo/pkg/log"

	defaults "github.com/mcuadros/go-defaults"
	toml "github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var configNewAsEnvFlag bool

// NewConfigCommand initialize a cobra config command tree
func NewConfigCommand(conf interface{}, envPrefix string) *cobra.Command {
	// Uppercase the prefix
	upPrefix := strings.ToUpper(envPrefix)

	// config
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage Service Configuration",
	}

	// config new
	configNewCmd := &cobra.Command{
		Use:   "new",
		Short: "Initialize a default configuration",
		Run: func(cmd *cobra.Command, args []string) {
			defaults.SetDefaults(conf)

			if !configNewAsEnvFlag {
				btes, err := toml.Marshal(conf)
				if err != nil {
					log.Bg().Fatal("Error during configuration export", zap.Error(err))
				}
				fmt.Println(string(btes))
			} else {
				m := flags.AsEnvVariables(conf, upPrefix, true)
				keys := []string{}

				for k := range m {
					keys = append(keys, k)
				}

				sort.Strings(keys)
				for _, k := range keys {
					fmt.Printf("export %s=\"%s\"\n", k, m[k])
				}
			}
		},
	}

	// flags
	configNewCmd.Flags().BoolVar(&configNewAsEnvFlag, "env", false, "Print configuration as environment variable")
	configCmd.AddCommand(configNewCmd)

	// Return base command
	return configCmd
}
