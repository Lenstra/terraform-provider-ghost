
package datasource_users


import (
	"context"

	"github.com/Lenstra/terraform-provider-ghost/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var (
	_ datasource.DataSource              = &usersDatasource{}
	_ datasource.DataSourceWithConfigure = &usersDatasource{}
)

type usersDatasource struct {
	client *client.Client
}

// Configure implements datasource.DataSourceWithConfigure.
func (s *usersDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func UsersDatasource() datasource.DataSource {
	return &usersDatasource{}
}

// Metadata implements datasource.DataSource.
func (s *usersDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

// Schema implements datasource.DataSource.
func (s *usersDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = UsersDataSourceSchema(ctx)
}

