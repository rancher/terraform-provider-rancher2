package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/rancher/terraform-provider-cattle/cattle"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cattle.Provider})
}
