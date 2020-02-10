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

// +build mage

package main

import (
	"fmt"
	"runtime"

	"github.com/Zenithar/go-monorepo/build/mage/golang"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/magefile/mage/mg"
)

var Default = Build

// Build the artefact
func Build() {
	banner := figure.NewFigure("Foo", "", true)
	banner.Print()

	fmt.Println("")
	color.Red("# Build Info ---------------------------------------------------------------")
	fmt.Printf("Go version : %s\n", runtime.Version())

	color.Red("# Pipeline -----------------------------------------------------------------")
	mg.SerialDeps(golang.Vendor, golang.License, golang.Lint("../../"), Test)

	color.Red("# Artefact(s) --------------------------------------------------------------")
	mg.Deps(Compile)
}

// Test application
func Test() {
	color.Cyan("## Tests")
	mg.Deps(golang.UnitTest("./..."))
}

// Compile artefacts
func Compile() {
	mg.Deps(
		golang.Build("foo", "github.com/Zenithar/go-monorepo/cmd/foo"),
	)
}
