package provider_ghost

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	test "github.com/Lenstra/terraform-plugin-test"
	"github.com/Lenstra/terraform-provider-ghost/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	helper "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestDocumentation(t *testing.T) {
	ctx := context.Background()
	provider := New("test")()

	examples := []string{
		"provider/provider.tf",
	}
	for _, f := range provider.Resources(ctx) {
		r := f()
		resp := &resource.MetadataResponse{}
		r.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "ghost"}, resp)
		examples = append(examples, fmt.Sprintf("resources/%s/resource.tf", resp.TypeName))
	}

	for _, f := range provider.DataSources(ctx) {
		d := f()
		resp := &datasource.MetadataResponse{}
		d.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "ghost"}, resp)
		examples = append(examples, fmt.Sprintf("data-sources/%s/data-source.tf", resp.TypeName))
	}

	for _, path := range examples {
		t.Run(path, func(t *testing.T) {
			require.FileExists(t, "../../../examples/"+path)
		})
	}
}

func TestAccSiteDataSource(t *testing.T) {
	client.TestServer(t)

	path := "../../../tests"
	entries, err := os.ReadDir(path)
	require.NoError(t, err)

	cwd, err := os.Getwd()
	require.NoError(t, err)

	var testCases []string
	for _, e := range entries {
		if e.IsDir() {
			testCases = append(testCases, filepath.Join(path, e.Name()))
		}
	}

	client.TestServer(t)

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			test.Test(t, tc, func(t *testing.T, dir string, tc *helper.TestCase) {
				for i, step := range tc.Steps {
					step.Config = strings.ReplaceAll(step.Config, "cwd", cwd)
					tc.Steps[i] = step
				}
				tc.ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
					"ghost": providerserver.NewProtocol6WithError(New("test")()),
				}
			}, nil)
		})
	}
}
