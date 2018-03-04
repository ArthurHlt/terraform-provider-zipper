package main

import (
	"github.com/ArthurHlt/terraform-provider-zipper/zipper"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: zipper.Provider})
}
