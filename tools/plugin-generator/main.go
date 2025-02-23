package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	sourceTemplate = `
package %[3]s_%[1]s


import (
	"context"

	"github.com/Lenstra/terraform-provider-ghost/client"
	"github.com/hashicorp/terraform-plugin-framework/%[3]s"
)

var (
	_ %[3]s.%[5]s              = &%[1]s%[4]s{}
	_ %[3]s.%[5]sWithConfigure = &%[1]s%[4]s{}
)

type %[1]s%[4]s struct {
	client *client.Client
}

// Configure implements %[3]s.%[5]sWithConfigure.
func (s *%[1]s%[4]s) Configure(ctx context.Context, req %[3]s.ConfigureRequest, resp *%[3]s.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Configuration error", "Failed to get Ghost client from provider metadata")
		return
	}
	s.client = client
}

func %[2]s%[4]s() %[3]s.%[5]s {
	return &%[1]s%[4]s{}
}

// Metadata implements %[3]s.%[5]s.
func (s *%[1]s%[4]s) Metadata(ctx context.Context, req %[3]s.MetadataRequest, resp *%[3]s.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_%[1]s"
}

// Schema implements %[3]s.%[5]s.
func (s *%[1]s%[4]s) Schema(ctx context.Context, req %[3]s.SchemaRequest, resp *%[3]s.SchemaResponse) {
	resp.Schema = %[2]s%[5]sSchema(ctx)
}

`
)

func Main() error {
	cmd := exec.Command(
		"tfplugingen-framework",
		"generate",
		"all",
		"--input=./tools/plugin-generator/provider_code_spec.json",
		"--output=internal/provider",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	if err := generateTemplates("datasource", "DataSource"); err != nil {
		return err
	}
	if err := generateTemplates("resource", "Resource"); err != nil {
		return err
	}

	return nil
}

func generateTemplates(typ string, packageType string) error {
	datasources, err := filepath.Glob(fmt.Sprintf("./internal/provider/%s_*/*_gen.go", typ))
	if err != nil {
		return err
	}
	for _, path := range datasources {
		filename := filepath.Base(path)
		name := strings.Split(filename, "_")[0]
		objectName := strings.ToUpper(string(name[0])) + name[1:]
		titledType := strings.ToUpper(string(typ[0])) + typ[1:]
		dir := filepath.Dir(path)
		content := fmt.Sprintf(sourceTemplate, name, objectName, typ, titledType, packageType)
		err := os.WriteFile(filepath.Join(dir, name+"_gen.go"), []byte(content), 0o600)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := Main(); err != nil {
		log.Fatal(err)
	}
}
