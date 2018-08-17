/*
Copyright 2018 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/cppforlife/knctl/pkg/knctl/cmd"

	// Import to initialize client auth plugins.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// TODO logs
	// TODO log flags used

	confUI := ui.NewConfUI(ui.NewNoopLogger())
	defer confUI.Flush()

	command := cmd.NewDefaultKnctlCmd(confUI)

	err := command.Execute()
	if err != nil {
		confUI.ErrorLinef("Error: %v", err)
		os.Exit(1)
	}

	confUI.PrintLinef("Succeeded")
}
