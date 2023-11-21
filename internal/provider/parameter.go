package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/api/tagmanager/v2"
)

var parameterSchema = buildParameterSchema()

func wrapParameterSchema(list schema.ListNestedAttribute) schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"key": schema.StringAttribute{
					Description: "Parameter key.",
					Optional:    true},
				"type": schema.StringAttribute{
					Description: "Parameter type.",
					Required:    true},
				"value": schema.StringAttribute{
					Description: "Parameter value.",
					Optional:    true},
				"list": list,
				"map":  list,
			},
		},
	}
}

func buildParameterSchema() schema.ListNestedAttribute {
	var s = schema.ListNestedAttribute{
		Description:  "Parameters.",
		Optional:     true,
		NestedObject: schema.NestedAttributeObject{},
	}

	for i := 0; i < 3; i++ {
		s = wrapParameterSchema(s)
	}

	return s
}

type ResourceParameterModel struct {
	Key   types.String             `tfsdk:"key"`
	Type  types.String             `tfsdk:"type"`
	Value types.String             `tfsdk:"value"`
	List  []ResourceParameterModel `tfsdk:"list"`
	Map   []ResourceParameterModel `tfsdk:"map"`
}

func (r *ResourceParameterModel) Equal(o ResourceParameterModel) bool {
	if !r.Key.Equal(o.Key) ||
		!r.Type.Equal(o.Type) ||
		!r.Value.Equal(o.Value) ||
		len(r.List) != len(o.List) ||
		len(r.Map) != len(o.Map) {
		return false
	}

	for i := 0; i < len(r.List); i++ {
		if !r.List[i].Equal(o.List[i]) {
			return false
		}
	}

	for i := 0; i < len(r.Map); i++ {
		if !r.Map[i].Equal(o.Map[i]) {
			return false
		}
	}

	return true
}

func toApiParameter(resourceParameter []ResourceParameterModel) []*tagmanager.Parameter {
	var parameter []*tagmanager.Parameter

	for _, p := range resourceParameter {
		var list, mmap []*tagmanager.Parameter

		if p.List != nil {
			list = toApiParameter(p.List)
		}

		if p.Map != nil {
			mmap = toApiParameter(p.Map)
		}

		parameter = append(parameter, &tagmanager.Parameter{
			Key:   p.Key.ValueString(),
			Type:  p.Type.ValueString(),
			Value: p.Value.ValueString(),
			List:  list,
			Map:   mmap,
		})
	}

	return parameter
}

func toResourceParameter(parameter []*tagmanager.Parameter) []ResourceParameterModel {
	var resourceParameter []ResourceParameterModel = make([]ResourceParameterModel, len(parameter))

	for i, p := range parameter {
		var list, mmap []ResourceParameterModel

		if p.List != nil {
			list = toResourceParameter(p.List)
		}

		if p.Map != nil {
			mmap = toResourceParameter(p.Map)
		}

		resourceParameter[i] = ResourceParameterModel{
			Key:   nullableStringValue(p.Key),
			Type:  nullableStringValue(p.Type),
			Value: nullableStringValue(p.Value),
			List:  list,
			Map:   mmap,
		}
	}

	return resourceParameter
}
