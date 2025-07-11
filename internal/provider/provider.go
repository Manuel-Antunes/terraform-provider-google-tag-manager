package provider

import (
	"context"
	"fmt"
	"os"
	"terraform-provider-google-tag-manager/internal/api"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &gtmProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &gtmProvider{}
}

// gtmProvider is the provider implementation.
type gtmProvider struct{}

// Metadata returns the provider type name.
func (p *gtmProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "gtm"
}

// Schema defines the provider-level schema for configuration data.
func (p *gtmProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"credential_file": schema.StringAttribute{
				Description: "Path to the credential file. Can also use GTM_CREDENTIAL_FILE environment variable.",
				Optional:    true},
			"account_id": schema.StringAttribute{
				Description: "GTM Account ID. Can also use GTM_ACCOUNT_ID environment variable.",
				Optional:    true},
			"container_id": schema.StringAttribute{
				Description: "GTM Container ID. Can also use GTM_CONTAINER_ID environment variable.",
				Optional:    true},
			"workspace_name": schema.StringAttribute{
				Description: "Workspace name. Can also use GTM_WORKSPACE_NAME environment variable.",
				Optional:    true},
			"retry_limit": schema.Int64Attribute{
				Description: "Number of times to retry requests when rate-limited before giving up. Can also use GTM_RETRY_LIMIT environment variable.",
				Optional:    true},
		},
	}
}

type gtmProviderModel struct {
	CredentialFile types.String `tfsdk:"credential_file"`
	AccountId      types.String `tfsdk:"account_id"`
	ContainerId    types.String `tfsdk:"container_id"`
	WorkspaceName  types.String `tfsdk:"workspace_name"`
	RetryLimit     types.Int64  `tfsdk:"retry_limit"`
}

// Configure prepares an API client for data sources and resources.
func (p *gtmProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Provider Configure starts.")
	defer tflog.Info(ctx, "Provider Configure finished.")

	// Retrieve provider data from configuration
	var config gtmProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get credentials from config or environment variables
	credFile := os.Getenv(api.EnvCredentialFile)
	if !config.CredentialFile.IsNull() && !config.CredentialFile.IsUnknown() {
		credFile = config.CredentialFile.ValueString()
	}

	accountId := os.Getenv(api.EnvAccountId)
	if !config.AccountId.IsNull() && !config.AccountId.IsUnknown() {
		accountId = config.AccountId.ValueString()
	}

	containerId := os.Getenv(api.EnvContainerId)
	if !config.ContainerId.IsNull() && !config.ContainerId.IsUnknown() {
		containerId = config.ContainerId.ValueString()
	}

	workspaceName := os.Getenv(api.EnvWorkspaceName)
	if !config.WorkspaceName.IsNull() && !config.WorkspaceName.IsUnknown() {
		workspaceName = config.WorkspaceName.ValueString()
	}

	var retryLimit = 10
	if retryLimitEnv := os.Getenv(api.EnvRetryLimit); retryLimitEnv != "" {
		// Try to parse the retry limit from environment
		var parsed int
		if _, err := fmt.Sscanf(retryLimitEnv, "%d", &parsed); err == nil && parsed > 0 {
			retryLimit = parsed
		}
	}
	if !config.RetryLimit.IsNull() && !config.RetryLimit.IsUnknown() {
		retryLimit = int(config.RetryLimit.ValueInt64())
	}

	// Validation for required fields
	if credFile == "" {
		resp.Diagnostics.AddError("Missing credential_file",
			"credential_file must be set in provider config or GTM_CREDENTIAL_FILE environment variable")
		return
	}

	if accountId == "" {
		resp.Diagnostics.AddError("Missing account_id",
			"account_id must be set in provider config or GTM_ACCOUNT_ID environment variable")
		return
	}

	if containerId == "" {
		resp.Diagnostics.AddError("Missing container_id",
			"container_id must be set in provider config or GTM_CONTAINER_ID environment variable")
		return
	}

	if workspaceName == "" {
		resp.Diagnostics.AddError("Missing workspace_name",
			"workspace_name must be set in provider config or GTM_WORKSPACE_NAME environment variable")
		return
	}

	client, err := api.NewClientInWorkspace(&api.ClientInWorkspaceOptions{
		ClientOptions: &api.ClientOptions{
			CredentialFile: credFile,
			AccountId:      accountId,
			ContainerId:    containerId,
			RetryLimit:     retryLimit,
		},
		WorkspaceName: workspaceName,
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create GTM Client", err.Error())
		return
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *gtmProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *gtmProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewWorkspaceResource,
		NewTagResource,
		NewVariableResource,
		NewTriggerResource,
	}
}
