package provider

import (
	"context"
	"terraform-provider-google-tag-manager/internal/api"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/api/tagmanager/v2"
)

// Interace adoption checks
var _ resource.ResourceWithConfigure = (*tagResource)(nil)

type tagResource struct {
	client *api.ClientInWorkspace
}

func NewTagResource() resource.Resource {
	return &tagResource{}
}

// Configure adds the provider configured client to the resource.
func (r *tagResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*api.ClientInWorkspace)
}

// Metadata returns the resource type name.
func (r *tagResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
}

var tagResourceSchemaAttributes = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Description: "The name of the tag.",
		Required:    true},
	"type": schema.StringAttribute{
		Description: "The type of the tag.",
		Required:    true},
	"id": schema.StringAttribute{
		Description: "The ID of the tag.",
		Computed:    true},
	"notes": schema.StringAttribute{
		Description: "The notes associated with the tag.",
		Optional:    true},
	"parameter": parameterSchema,
	"firing_trigger_id": schema.ListAttribute{
		Description: "The ID of the firing triggers associated with the tag.",
		Optional:    true,
		ElementType: types.StringType,
	},
}

// Schema defines the schema for the resource.
func (r *tagResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: tagResourceSchemaAttributes}
}

type resourceTagModel struct {
	Name            types.String             `tfsdk:"name"`
	Type            types.String             `tfsdk:"type"`
	Id              types.String             `tfsdk:"id"`
	Notes           types.String             `tfsdk:"notes"`
	Parameter       []ResourceParameterModel `tfsdk:"parameter"`
	FiringTriggerId []types.String           `tfsdk:"firing_trigger_id"`
}

// Create creates the resource and sets the initial Terraform state.
func (r *tagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceTagModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tag, err := r.client.CreateTag(toApiTag(plan, false))
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Tag", err.Error())
		return
	}

	plan.Id = types.StringValue(tag.TagId)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *tagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceTagModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tag, err := r.client.Tag(state.Id.ValueString())
	if err == api.ErrNotExist {
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error Reading Tag", err.Error())
		return
	}

	var resource = toResourceTag(tag)

	diags = resp.State.Set(ctx, &resource)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *tagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state resourceTagModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tag, err := r.client.UpdateTag(state.Id.ValueString(), toApiTag(plan, true))
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Tag", err.Error())
		return
	}

	plan.Id = types.StringValue(tag.TagId)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *tagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError(
			"Resource Import Missing ID",
			"This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				"Resource ImportState method call to ImportState with an empty ID is not allowed.",
		)
		return
	}
	tag, err := r.client.Tag(req.ID)
	if err == nil {
		resp.Diagnostics.AddError(
			"Resource Import Failed",
			"Failed to import tag with ID "+req.ID+". The tag does not exist or the ID is invalid.",
		)
		return
	}

	plan := toResourceTag(tag)

	diags := resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *tagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceTagModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if state.Id.IsNull() || state.Id.IsUnknown() {
		resp.Diagnostics.AddError("Invalid Id state", state.Id.String())
	}

	err := r.client.DeleteTag(state.Id.ValueString())
	if err == api.ErrNotExist {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error Deleting Tag", err.Error())
		return
	}
}

// Equal compares the two models and returns true if they are equal.
func (m resourceTagModel) Equal(o resourceTagModel) bool {
	if !m.Name.Equal(o.Name) ||
		!m.Type.Equal(o.Type) ||
		(!m.Id.IsUnknown() && !m.Id.Equal(o.Id)) ||
		!m.Notes.Equal(o.Notes) ||
		len(m.Parameter) != len(o.Parameter) ||
		len(m.FiringTriggerId) != len(o.FiringTriggerId) {
		return false
	}

	for i := range m.Parameter {
		if !m.Parameter[i].Equal(o.Parameter[i]) {
			return false
		}
	}

	for i := range m.FiringTriggerId {
		if !m.FiringTriggerId[i].Equal(o.FiringTriggerId[i]) {
			return false
		}
	}

	return true
}

func toResourceTag(tag *tagmanager.Tag) resourceTagModel {
	return resourceTagModel{
		Name:            types.StringValue(tag.Name),
		Type:            types.StringValue(tag.Type),
		Id:              types.StringValue(tag.TagId),
		Notes:           nullableStringValue(tag.Notes),
		Parameter:       toResourceParameter(tag.Parameter),
		FiringTriggerId: toResourceStringArray(tag.FiringTriggerId),
	}

}

func toApiTag(resource resourceTagModel, id bool) *tagmanager.Tag {
	if !id {
		return &tagmanager.Tag{
			Name:            resource.Name.ValueString(),
			Type:            resource.Type.ValueString(),
			Notes:           resource.Notes.ValueString(),
			Parameter:       toApiParameter(resource.Parameter),
			FiringTriggerId: unwrapStringArray(resource.FiringTriggerId),
		}
	}

	return &tagmanager.Tag{
		Name:            resource.Name.ValueString(),
		Type:            resource.Type.ValueString(),
		TagId:           resource.Id.String(),
		Notes:           resource.Notes.ValueString(),
		Parameter:       toApiParameter(resource.Parameter),
		FiringTriggerId: unwrapStringArray(resource.FiringTriggerId),
	}
}
