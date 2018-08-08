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
	"fmt"

	uitable "github.com/cppforlife/go-cli-ui/ui/table"
)

type ValueUnknownBool struct {
	B *bool
}

var _ uitable.Value = ValueUnknownBool{}

func NewValueUnknownBool(b *bool) ValueUnknownBool { return ValueUnknownBool{B: b} }

func (t ValueUnknownBool) String() string {
	if t.B != nil {
		return fmt.Sprintf("%t", *t.B)
	}
	return ""
}

func (t ValueUnknownBool) Value() uitable.Value            { return t }
func (t ValueUnknownBool) Compare(other uitable.Value) int { panic("Never called") }
