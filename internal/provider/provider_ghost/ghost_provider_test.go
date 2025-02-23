package provider_ghost

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	test "github.com/Lenstra/terraform-plugin-test"
	"github.com/Lenstra/terraform-provider-ghost/client"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

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
			test.Test(t, tc, func(t *testing.T, dir string, tc *resource.TestCase) {
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
