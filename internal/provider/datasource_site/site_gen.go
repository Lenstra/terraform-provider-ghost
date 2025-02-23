
package datasource_site


import (
	"context"

	"github.com/Lenstra/terraform-provider-ghost/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var (
	_ datasource.DataSource              = &siteDatasource{}
	_ datasource.DataSourceWithConfigure = &siteDatasource{}
)

type siteDatasource struct {
	client *client.Client
}

// Configure implements datasource.DataSourceWithConfigure.
func (s *siteDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func SiteDatasource() datasource.DataSource {
	return &siteDatasource{}
}

// Metadata implements datasource.DataSource.
func (s *siteDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

// Schema implements datasource.DataSource.
func (s *siteDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = SiteDataSourceSchema(ctx)
}

