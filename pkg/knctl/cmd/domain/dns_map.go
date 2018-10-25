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

package domain

import (
	"encoding/json"

	"github.com/cppforlife/go-cli-ui/ui"
	cmdcore "github.com/cppforlife/knctl/pkg/knctl/cmd/core"
	ctling "github.com/cppforlife/knctl/pkg/knctl/ingress"
	"github.com/spf13/cobra"
)

type DNSMapOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
}

func NewDNSMapOptions(ui ui.UI, depsFactory cmdcore.DepsFactory) *DNSMapOptions {
	return &DNSMapOptions{ui: ui, depsFactory: depsFactory}
}

func NewDNSMapCmd(o *DNSMapOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dns-map",
		Short: "Print domain to IP map in JSON format",
		RunE:  func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	cmd.Hidden = true // for advanced usage
	return cmd
}

func (o *DNSMapOptions) Run() error {
	coreClient, err := o.depsFactory.CoreClient()
	if err != nil {
		return err
	}

	domains, err := NewDomains(coreClient).List()
	if err != nil {
		return err
	}

	ingSvcs, err := ctling.NewIngressServices(coreClient).List()
	if err != nil {
		return err
	}

	dnsMap := map[string][]string{}

	var addrs []string

	for _, svc := range ingSvcs {
		addrs = append(addrs, svc.Addresses()...)
	}

	for _, domain := range domains {
		dnsMap[domain.Name] = addrs
	}

	outBytes, err := json.Marshal(dnsMap)
	if err != nil {
		return err
	}

	o.ui.PrintBlock(outBytes)

	return nil
}
