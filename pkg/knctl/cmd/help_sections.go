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

package cmd

import (
	"github.com/cppforlife/knctl/pkg/knctl/cobrautil"
)

const (
	cmdGroupKey = "knctl-group"
)

var (
	basicGroup = cobrautil.HelpSection{
		Key:   cmdGroupKey,
		Value: "basic",
		Title: "Basic Commands:",
	}
	buildMgmtGroup = cobrautil.HelpSection{
		Key:   cmdGroupKey,
		Value: "build-mgmt",
		Title: "Build Management Commands:",
	}
	secretMgmtGroup = cobrautil.HelpSection{
		Key:   cmdGroupKey,
		Value: "secret-mgmt",
		Title: "Secret Management Commands:",
	}
	routeMgmtGroup = cobrautil.HelpSection{
		Key:   cmdGroupKey,
		Value: "route-mgmt",
		Title: "Route Management Commands:",
	}
	otherGroup = cobrautil.HelpSection{
		Key:   cmdGroupKey,
		Value: "other",
		Title: "Other Commands:",
	}
	systemGroup = cobrautil.HelpSection{
		Key:   cmdGroupKey,
		Value: "system",
		Title: "System Commands:",
	}
	restOfCommandsGroup = cobrautil.HelpSection{
		Key:   cmdGroupKey,
		Value: "",
		Title: "Available Commands:",
	}
)
