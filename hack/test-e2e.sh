#!/bin/bash

# Copyright 2018 The Knative Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e -x -u

export KNCTL_NAMESPACE=
export KNCTL_KUBECONFIG=
export KNCTL_BASIC_AUTH_SECRET_PASSWORD=
export KNCTL_SSH_AUTH_SECRET_PRIVATE_KEY=

GOCACHE=off go test ./test/e2e/ -timeout 120m -test.v $@

echo "Success"
