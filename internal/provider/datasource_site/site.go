package datasource_site

import (
	"context"

	"github.com/Lenstra/terraform-provider-ghost/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read implements datasource.DataSource.
func (s *siteDatasource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	site, err := s.client.Site().Read(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to get site information", err.Error())
		return
	}

	diags := resp.State.Set(ctx, convert(site))
	resp.Diagnostics.Append(diags...)
}

func convert(site *client.Site) *SiteModel {
	return &SiteModel{
		AccentColor:         types.StringValue(site.AccentColor),
		AllowExternalSignup: types.BoolValue(site.AllowExternalSignup),
		CoverImage:          types.StringValue(site.CoverImage),
		Description:         types.StringValue(site.Description),
		Icon:                types.StringValue(site.Icon),
		Locale:              types.StringValue(site.Locale),
		Logo:                types.StringValue(site.Logo),
		Title:               types.StringValue(site.Title),
		Url:                 types.StringValue(site.Url),
		Version:             types.StringValue(site.Version),
	}
}
