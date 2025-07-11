package provider

import (
	"context"
	"terraform-provider-google-tag-manager/internal/api"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/api/tagmanager/v2"
)

var _ resource.ResourceWithConfigure = (*triggerResource)(nil)

type triggerResource struct {
	client *api.ClientInWorkspace
}

func NewTriggerResource() resource.Resource {
	return &triggerResource{}
}

// Configure adds the provider configured client to the resource.
func (r *triggerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*api.ClientInWorkspace)
}

// Metadata returns the resource type name.
func (r *triggerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trigger"
}

var triggerResourceSchemaAttributes = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Description: "The name of the trigger.",
		Required:    true,
	},
	"type": schema.StringAttribute{
		Description: "The type of the trigger.",
		Required:    true,
	},
	"id": schema.StringAttribute{
		Description: "The ID of the trigger.",
		Computed:    true,
	},
	"notes": schema.StringAttribute{
		Description: "The notes of the trigger.",
		Optional:    true,
	},
	"custom_event_filter": conditionSchema,
}

// Schema defines the schema for the resource.
func (r *triggerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: triggerResourceSchemaAttributes}
}

type resourceTriggerModel struct {
	Name              types.String             `tfsdk:"name"`
	Type              types.String             `tfsdk:"type"`
	Id                types.String             `tfsdk:"id"`
	Notes             types.String             `tfsdk:"notes"`
	CustomEventFilter []ResourceConditionModel `tfsdk:"custom_event_filter"`
}

// Create creates the resource and sets the initial Terraform state.
func (r *triggerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceTriggerModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	trigger, err := r.client.CreateTrigger(toApiTrigger(plan))
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Trigger", err.Error())
		return
	}

	plan.Id = types.StringValue(trigger.TriggerId)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *triggerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceTriggerModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	trigger, err := r.client.Trigger(state.Id.ValueString())
	if err == api.ErrNotExist {
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error Reading Trigger", err.Error())
		return
	}

	var resource = toResourceTrigger(trigger)

	diags = resp.State.Set(ctx, &resource)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *triggerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state resourceTriggerModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	trigger, err := r.client.UpdateTrigger(state.Id.ValueString(), toApiTrigger(plan))
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Trigger", err.Error())
		return
	}

	plan.Id = types.StringValue(trigger.TriggerId)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *triggerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceTriggerModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTrigger(state.Id.ValueString())
	if err == api.ErrNotExist {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error Deleting Trigger", err.Error())
		return
	}
}

func (r *triggerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	triggerId := req.ID

	if triggerId == "" {
		resp.Diagnostics.AddError("Error Importing Trigger", "Trigger ID cannot be empty")
		return
	}

	trigger, err := r.client.Trigger(triggerId)
	if err != nil {
		resp.Diagnostics.AddError("Error Importing Trigger", err.Error())
		return
	}

	resource := toResourceTrigger(trigger)

	diags := resp.State.Set(ctx, &resource)
	resp.Diagnostics.Append(diags...)
}

// Equal compares the trigger resource model with the given resource model
func (m resourceTriggerModel) Equal(o resourceTriggerModel) bool {
	if !m.Name.Equal(o.Name) ||
		!m.Type.Equal(o.Type) ||
		(!m.Id.IsUnknown() && !m.Id.Equal(o.Id)) ||
		!m.Notes.Equal(o.Notes) {
		return false
	}

	if len(m.CustomEventFilter) != len(o.CustomEventFilter) {
		return false
	}

	for i := range m.CustomEventFilter {
		if !m.CustomEventFilter[i].Equal(o.CustomEventFilter[i]) {
			return false
		}
	}

	return true
}

func toResourceTrigger(trigger *tagmanager.Trigger) resourceTriggerModel {
	return resourceTriggerModel{
		Name:              types.StringValue(trigger.Name),
		Type:              types.StringValue(trigger.Type),
		Id:                types.StringValue(trigger.TriggerId),
		Notes:             nullableStringValue(trigger.Notes),
		CustomEventFilter: toResourceCondition(trigger.CustomEventFilter),
	}
}

func toApiTrigger(resource resourceTriggerModel) *tagmanager.Trigger {
	return &tagmanager.Trigger{
		Name:              resource.Name.ValueString(),
		Type:              resource.Type.ValueString(),
		TriggerId:         resource.Id.ValueString(),
		Notes:             resource.Notes.ValueString(),
		CustomEventFilter: toApiCondition(resource.CustomEventFilter),
	}
}
