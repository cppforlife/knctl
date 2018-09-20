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

package e2e

import (
	"testing"
)

func TestHelpCmdRoot(t *testing.T) {
	env := BuildEnv(t)
	knctl := Knctl{t, env.Namespace, Logger{}}

	out, _ := knctl.RunWithOpts([]string{"-h"}, RunOpts{NoNamespace: true})

	const expectedOutput = `knctl controls Knative resources.

CLI docs: https://github.com/cppforlife/knctl#docs.
Knative docs: https://github.com/knative/docs.

Usage:
  knctl [flags]
  knctl [command]

Basic Commands:
  curl              Curl service
  deploy            Deploy service
  logs              Print service logs
  pod               Pod management (list)
  revision          Revision management (annotate, delete, list, show, tag, untag)
  service           Service management (annotate, delete, list, open, show)

Build Management Commands:
  build             Build management (create, delete, list)

Secret Management Commands:
  basic-auth-secret Basic auth secret management (create)
  service-account   Service account management (create)
  ssh-auth-secret   SSH auth secret management (create)

Route Management Commands:
  domain            Domain management (create, list)
  ingress           Ingress management (list)
  route             Route management (create, delete, list)

Other Commands:
  namespace         Namespace management (create)

System Commands:
  install           Install Knative and Istio
  uninstall         Uninstall Knative and Istio
  version           Print client version

Available Commands:
  help              Help about any command

Flags:
      --column strings              Filter to show only given columns
  -h, --help                        help for knctl
      --json                        Output as JSON
      --kubeconfig string           Path to the kubeconfig file ($KNCTL_KUBECONFIG or $KUBECONFIG)
      --kubeconfig-context string   Kubeconfig context override ($KNCTL_KUBECONFIG_CONTEXT)
      --no-color                    Disable colorized output
      --non-interactive             Don't ask for user input
      --tty                         Force TTY-like output

Use "knctl [command] --help" for more information about a command.

Succeeded
`

	if out != expectedOutput {
		t.Fatalf("Expected to find exact help content")
	}
}

func TestHelpCmdWithChildren(t *testing.T) {
	env := BuildEnv(t)
	knctl := Knctl{t, env.Namespace, Logger{}}

	out, _ := knctl.RunWithOpts([]string{"service", "-h"}, RunOpts{NoNamespace: true})

	const expectedOutput = `Service management (annotate, delete, list, open, show)

Usage:
  knctl service [flags]
  knctl service [command]

Aliases:
  service, s, svc, svc, services

Available Commands:
  annotate    Annotate service
  delete      Delete service
  list        List services
  open        Open web browser pointing at a service domain
  show        Show service

Flags:
  -h, --help   help for service

Global Flags:
      --column strings              Filter to show only given columns
      --json                        Output as JSON
      --kubeconfig string           Path to the kubeconfig file ($KNCTL_KUBECONFIG or $KUBECONFIG)
      --kubeconfig-context string   Kubeconfig context override ($KNCTL_KUBECONFIG_CONTEXT)
      --no-color                    Disable colorized output
      --non-interactive             Don't ask for user input
      --tty                         Force TTY-like output

Use "knctl service [command] --help" for more information about a command.

Succeeded
`

	if out != expectedOutput {
		t.Fatalf("Expected to find exact help content")
	}
}

func TestHelpCmdLeaf(t *testing.T) {
	env := BuildEnv(t)
	knctl := Knctl{t, env.Namespace, Logger{}}

	out, _ := knctl.RunWithOpts([]string{"service", "delete", "-h"}, RunOpts{NoNamespace: true})

	const expectedOutput = `Delete service

Usage:
  knctl service delete [flags]

Aliases:
  delete, del

Examples:

  # Delete service 'svc1' in namespace 'ns1'
  knctl service delete -s svc1 -n ns1

Flags:
  -h, --help               help for delete
  -n, --namespace string   Specified namespace ($KNCTL_NAMESPACE or default from kubeconfig)
  -s, --service string     Specified service

Global Flags:
      --column strings              Filter to show only given columns
      --json                        Output as JSON
      --kubeconfig string           Path to the kubeconfig file ($KNCTL_KUBECONFIG or $KUBECONFIG)
      --kubeconfig-context string   Kubeconfig context override ($KNCTL_KUBECONFIG_CONTEXT)
      --no-color                    Disable colorized output
      --non-interactive             Don't ask for user input
      --tty                         Force TTY-like output

Succeeded
`

	if out != expectedOutput {
		t.Fatalf("Expected to find exact help content")
	}
}
