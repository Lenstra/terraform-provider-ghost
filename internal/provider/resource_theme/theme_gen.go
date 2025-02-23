
package resource_theme


import (
	"context"

	"github.com/Lenstra/terraform-provider-ghost/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ resource.Resource              = &themeResource{}
	_ resource.ResourceWithConfigure = &themeResource{}
)

type themeResource struct {
	client *client.Client
}

// Configure implements resource.ResourceWithConfigure.
func (s *themeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func ThemeResource() resource.Resource {
	return &themeResource{}
}

// Metadata implements resource.Resource.
func (s *themeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_theme"
}

// Schema implements resource.Resource.
func (s *themeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ThemeResourceSchema(ctx)
}

