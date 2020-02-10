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

package golang

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/Zenithar/go-monorepo/build/mage/git"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Build the given binary usign the given package.
func Build(name, packageName string) func() error {
	return func() error {
		mg.Deps(git.CollectInfo)

		fmt.Printf(" > Building %s [%s]\n", name, packageName)

		version, err := git.TagMatch(fmt.Sprintf("cmd/%s*", name))
		if err != nil {
			return err
		}

		varsSetByLinker := map[string]string{
			"github.com/Zenithar/go-monorepo/build/version.Version":   version,
			"github.com/Zenithar/go-monorepo/build/version.Revision":  git.Revision,
			"github.com/Zenithar/go-monorepo/build/version.Branch":    git.Branch,
			"github.com/Zenithar/go-monorepo/build/version.BuildUser": os.Getenv("USER"),
			"github.com/Zenithar/go-monorepo/build/version.BuildDate": time.Now().Format(time.RFC3339),
			"github.com/Zenithar/go-monorepo/build/version.GoVersion": runtime.Version(),
		}
		var linkerArgs []string
		for name, value := range varsSetByLinker {
			linkerArgs = append(linkerArgs, "-X", fmt.Sprintf("%s=%s", name, value))
		}
		linkerArgs = append(linkerArgs, "-s", "-w")

		return sh.RunWith(map[string]string{
			"CGO_ENABLED": "0",
		}, "go", "build", "-buildmode=pie", "-ldflags", strings.Join(linkerArgs, " "), "-mod=vendor", "-o", fmt.Sprintf("../../bin/%s", name), packageName)
	}
}
