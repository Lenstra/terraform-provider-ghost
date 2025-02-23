package provider_ghost

import (
	"context"
	"fmt"
	"os"

	"github.com/Lenstra/terraform-provider-ghost/client"
	"github.com/Lenstra/terraform-provider-ghost/internal/provider/datasource_site"
	"github.com/Lenstra/terraform-provider-ghost/internal/provider/datasource_users"
	"github.com/Lenstra/terraform-provider-ghost/internal/provider/resource_theme"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/rs/zerolog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &ghostProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ghostProvider{
			version: version,
		}
	}
}

// ghostProvider is the provider implementation.
type ghostProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *ghostProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ghost"
}

// Schema defines the provider-level schema for configuration data.
func (p *ghostProvider) Schema(ctx context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = GhostProviderSchema(ctx)
}

func (p *ghostProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring ghost client")

	//Retrieve provider data from configuration
	var config GhostModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	reportUnknown := func(diags diag.Diagnostics, attr string, envName string) {
		diags.AddAttributeError(
			path.Root(attr),
			fmt.Sprintf("Unknown Ghost %s", attr),
			fmt.Sprintf("The provider cannot create the Ghost API client as there is an unknown configuration value for the Ghost %s. ", attr)+
				fmt.Sprintf("Either target apply the source of the value first, set the value statically in the configuration or set the %s environment variable.", envName),
		)
	}

	if config.Address.IsUnknown() {
		reportUnknown(resp.Diagnostics, "address", "GHOST_ADDRESS")
	}
	if config.AdminApiKey.IsUnknown() {
		reportUnknown(resp.Diagnostics, "admin_api_key", "GHOST_ADMIN_API_KEY")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	adminApiKey := os.Getenv("GHOST_ADMIN_API_KEY")
	address := os.Getenv("GHOST_ADDRESS")

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	reportUnset := func(diags diag.Diagnostics, attr string, envName string) {
		diags.AddAttributeError(
			path.Root(attr),
			fmt.Sprintf("Unknown Ghost %s", attr),
			fmt.Sprintf("The provider cannot create the Ghost API client as there is a missing or empty value for the Ghost %s. ", attr)+
				fmt.Sprintf("Set the %s value in the configuration or use the %s environment variable. ", attr, envName)+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if address == "" {
		reportUnset(resp.Diagnostics, "address", "GHOST_ADDRESS")
	}

	if adminApiKey == "" {
		reportUnset(resp.Diagnostics, "admin_api_key", "GHOST_ADMIN_API_KEY")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	conf := &client.Config{
		Address:     address,
		AdminAPIKey: adminApiKey,
	}

	if os.Getenv("TF_PROVIDER_GHOST_LOG") != "" {
		conf.Logger = zerolog.New(os.Stdout)
	}

	apiClient, err := client.NewClient(conf)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build Ghost client", err.Error())
		return
	}

	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient
}

// DataSources defines the data sources implemented in the provider.
func (p *ghostProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasource_site.SiteDatasource,
		datasource_users.UsersDatasource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *ghostProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resource_theme.ThemeResource,
	}
}
