package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/openfga/go-sdk/client"
	"github.com/openfga/go-sdk/credentials"

	"github.com/mauriceackel/terraform-provider-openfga/internal/authorizationmodel"
	"github.com/mauriceackel/terraform-provider-openfga/internal/query"
	"github.com/mauriceackel/terraform-provider-openfga/internal/relationshiptuple"
	"github.com/mauriceackel/terraform-provider-openfga/internal/store"
)

// Ensure OpenFgaProvider satisfies various provider interfaces.
var _ provider.Provider = &OpenFgaProvider{}
var _ provider.ProviderWithFunctions = &OpenFgaProvider{}

type OpenFgaProvider struct {
	version string
}

// OpenFgaProviderModel describes the provider data model.
type OpenFgaProviderModel struct {
	ApiUrl types.String `tfsdk:"api_url"`

	ApiToken         types.String `tfsdk:"api_token"`
	ClientId         types.String `tfsdk:"client_id"`
	ClientSecret     types.String `tfsdk:"client_secret"`
	TokenEndpointUrl types.String `tfsdk:"token_endpoint_url"`
}

func (p *OpenFgaProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "openfga"
	resp.Version = p.version
}

func (p *OpenFgaProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_url": schema.StringAttribute{
				MarkdownDescription: "URL of the OpenFGA server",
				Optional:            true,
			},
			"api_token": schema.StringAttribute{
				MarkdownDescription: "Access token for authentication to the OpenFGA server",
				Optional:            true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "Client ID for client credentials authentication",
				Optional:            true,
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "Client secret for client credentials authentication",
				Optional:            true,
				Sensitive:           true,
			},
			"token_endpoint_url": schema.StringAttribute{
				MarkdownDescription: "The token endpoint URL for client credentials authentication",
				Optional:            true,
			},
		},
	}
}

func (p *OpenFgaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config OpenFgaProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if config.ApiUrl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_url"),
			"Unknown OpenFGA API URL",
			"The provider cannot create the OpenFGA API client as there is an unknown configuration value for the API URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the FGA_API_URL environment variable.",
		)
	}

	if config.ApiToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Unknown OpenFGA API token",
			"The provider cannot create the OpenFGA API client as there is an unknown configuration value for the API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the FGA_API_TOKEN environment variable.",
		)
	}

	if config.ClientId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Unknown OpenFGA client ID",
			"The provider cannot create the OpenFGA API client as there is an unknown configuration value for the OpenFGA client id. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the FGA_CLIENT_ID environment variable.",
		)
	}

	if config.ClientSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Unknown OpenFGA client secret",
			"The provider cannot create the OpenFGA API client as there is an unknown configuration value for the OpenFGA client secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the FGA_CLIENT_SECRET environment variable.",
		)
	}

	if config.TokenEndpointUrl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token_endpoint_url"),
			"Unknown OpenFGA token endpoint URL",
			"The provider cannot create the OpenFGA API client as there is an unknown configuration value for the OpenFGA token endpoint URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the FGA_TOKEN_ENDPOINT_URL environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	apiUrl := os.Getenv("FGA_API_URL")
	apiToken := os.Getenv("FGA_API_TOKEN")
	clientId := os.Getenv("FGA_CLIENT_ID")
	clientSecret := os.Getenv("FGA_CLIENT_SECRET")
	tokenEndpointUrl := os.Getenv("FGA_TOKEN_ENDPOINT_URL")

	if !config.ApiUrl.IsNull() {
		apiUrl = config.ApiUrl.ValueString()
	}

	if !config.ApiToken.IsNull() {
		apiToken = config.ApiToken.ValueString()
	}

	if !config.ClientId.IsNull() {
		clientId = config.ClientId.ValueString()
	}

	if !config.ClientSecret.IsNull() {
		clientSecret = config.ClientSecret.ValueString()
	}

	if !config.TokenEndpointUrl.IsNull() {
		tokenEndpointUrl = config.TokenEndpointUrl.ValueString()
	}

	if apiUrl == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_url"),
			"Missing OpenFGA API URL",
			"The provider cannot create the OpenFGA API client as there is a missing or empty value for the OpenFGA API URL. "+
				"Set the host value in the configuration or use the FGA_API_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	tokenSpecified := apiToken != ""
	clientCredentialsSpecified := clientId != "" && clientSecret != "" && tokenEndpointUrl != ""

	if tokenSpecified && clientCredentialsSpecified {
		resp.Diagnostics.AddError(
			"Invalid Credentials",
			"Exactly one of API token or client ID, client secret, and token endpoint URL must be specified.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Build OpenFGA client
	var apiCredentials credentials.Credentials
	if tokenSpecified {
		apiCredentials = credentials.Credentials{
			Method: credentials.CredentialsMethodApiToken,
			Config: &credentials.Config{
				ApiToken: apiToken,
			},
		}
	} else if clientCredentialsSpecified {
		apiCredentials = credentials.Credentials{
			Method: credentials.CredentialsMethodClientCredentials,
			Config: &credentials.Config{
				ClientCredentialsClientId:       clientId,
				ClientCredentialsClientSecret:   clientSecret,
				ClientCredentialsApiTokenIssuer: tokenEndpointUrl,
			},
		}
	} else {
		apiCredentials = credentials.Credentials{
			Method: credentials.CredentialsMethodNone,
		}
	}

	client, err := client.NewSdkClient(&client.ClientConfiguration{
		ApiUrl:      apiUrl,
		Credentials: &apiCredentials,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create OpenFGA API client",
			"An unexpected error occurred when creating the OpenFGA API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"OpenFGA Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *OpenFgaProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		store.NewStoreResource,
		authorizationmodel.NewAuthorizationModelResource,
		relationshiptuple.NewRelationshipTupleResource,
	}
}

func (p *OpenFgaProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		store.NewStoreDataSource,
		store.NewStoresDataSource,
		authorizationmodel.NewAuthorizationModelDocumentDataSource,
		authorizationmodel.NewAuthorizationModelDataSource,
		authorizationmodel.NewAuthorizationModelsDataSource,
		relationshiptuple.NewRelationshipTupleDataSource,
		relationshiptuple.NewRelationshipTuplesDataSource,
		query.NewCheckQueryDataSource,
	}
}

func (p *OpenFgaProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OpenFgaProvider{
			version: version,
		}
	}
}
