package main

import (
	"context"
	"flag"
	"log"

	"github.com/Lenstra/terraform-provider-ghost/internal/provider/provider_ghost"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

//go:generate go run -- ./tools/plugin-generator

// Run "go generate" to format example terraform files and generate the docs for the registry/website
//go:generate terraform fmt -recursive .

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

var (
	version = "dev"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "github.com/Lenstra/terraform-provider-ghost",
		Debug:   debugMode,
	}

	err := providerserver.Serve(context.Background(), provider_ghost.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
