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

package kube

import (
	"io"
)

type DirCp struct {
	exec Exec
}

func NewDirCp(exec Exec) DirCp {
	return DirCp{exec}
}

func (s DirCp) Execute(srcDir, dstDir string) error {
	reader, writer := io.Pipe()
	tarWritingErrCh := make(chan error)

	go func() {
		err := TarBuilder{}.Build(srcDir, "/", TarBuilderOpts{ExcludedPaths: []string{".git"}}, writer)
		writer.Close()
		tarWritingErrCh <- err
	}()

	cmd := []string{"tar", "xf", "-", "-C", dstDir}

	execErr := s.exec.Execute(cmd, reader)
	tarErr := <-tarWritingErrCh

	if tarErr != nil {
		return tarErr
	}

	return execErr
}
