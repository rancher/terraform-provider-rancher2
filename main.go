package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/ihor.kolomiiets/terraform-provider-rancher2/rancher2"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: rancher2.Provider})
}
