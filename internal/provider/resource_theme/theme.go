package resource_theme

import (
	"context"
	"os"

	"github.com/Lenstra/terraform-provider-ghost/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// Create implements resource.Resource.
func (s *themeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var theme ThemeModel
	diags := req.Plan.Get(ctx, &theme)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(update(ctx, s.client, theme))
	resp.State = tfsdk.State(req.Config)
}

// Update implements resource.Resource.
func (s *themeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var theme ThemeModel
	diags := req.Plan.Get(ctx, &theme)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(update(ctx, s.client, theme))
	resp.State = tfsdk.State(req.Config)
}

// Delete implements resource.Resource.
func (s *themeResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {}

// Read implements resource.Resource.
func (s *themeResource) Read(context.Context, resource.ReadRequest, *resource.ReadResponse) {}

func update(ctx context.Context, client *client.Client, theme ThemeModel) diag.Diagnostic {
	name := theme.Name.ValueString()
	source := theme.Source.ValueString()
	f, err := os.Open(source)
	if err != nil {
		return diag.NewErrorDiagnostic("failed to read source", err.Error())
	}
	defer f.Close()
	if _, err := client.Themes().Upload(ctx, name, f); err != nil {
		return diag.NewErrorDiagnostic("failed to upload theme", err.Error())
	}
	if _, err := client.Themes().Activate(ctx, name); err != nil {
		return diag.NewErrorDiagnostic("failed to activate theme", err.Error())
	}
	return nil
}
