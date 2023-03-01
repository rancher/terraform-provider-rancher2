package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/rancher/terraform-provider-rancher2/rancher2"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: rancher2.Provider,
		Debug:        true,
		ProviderAddr: "registry.terraform.io/rancher/rancher2",
	}

	plugin.Serve(opts)
}
