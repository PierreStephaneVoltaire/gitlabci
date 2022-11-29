package gitlabci

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gitlabci/helpers"
	"gopkg.in/yaml.v3"
	"os"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &FileDataSource{}
)

var defaultStateType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"image":         types.StringType,
		"timeout":       types.StringType,
		"tags":          types.ListType{ElemType: types.StringType},
		"retry":         types.Int64Type,
		"interruptible": types.BoolType,
		"before_script": types.ListType{ElemType: types.StringType},
		"after_script":  types.ListType{ElemType: types.StringType},
		"cache":         types.ObjectType{AttrTypes: cacheStateType},
	}}
var cacheStateType = map[string]attr.Type{
	"key":   types.StringType,
	"paths": types.ListType{ElemType: types.StringType},
}

// FileDataSource is the data source implementation.
type FileDataSource struct{}

// NewFileDataSource is a helper function to simplify the provider implementation.
func NewFileDataSource() datasource.DataSource {
	return &FileDataSource{}
}

// Metadata returns the data source type name.
func (d *FileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

// GetSchema defines the schema for the data source.
func (d *FileDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example data source",

		Attributes: map[string]tfsdk.Attribute{
			"file_location": {
				MarkdownDescription: "what file to parse",
				Optional:            true,
				Type:                types.StringType,
			},
			"stages": {
				MarkdownDescription: "The names and order of the pipeline stages",
				Type:                types.ListType{ElemType: types.StringType},
				Computed:            true,
			},
			"default": {
				Computed: true,
				Optional: true,
				Type:     &defaultStateType,
			},
		},
	}, nil
}

// Read refreshes the Terraform state with the latest data.
func (d *FileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// get value from schema
	var data schema.FileDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	file, err := os.ReadFile(data.FileLocation.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not read the file",
			err.Error(),
		)
		return
	}
	fileContent := string(file)
	// parse the file
	y := schema.Yaml{}
	err = yaml.Unmarshal([]byte(fileContent), &y)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not parse the file",
			err.Error(),
		)
		return
	}
	// save state
	state := ConvertYamlToState(ctx, data.FileLocation, y)
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func ConvertYamlToState(ctx context.Context, fileLocation types.String, y schema.Yaml) schema.FileDataSourceModel {
	tags, _ := types.ListValueFrom(ctx, types.StringType, y.Default.Tags)
	beforeScript, _ := types.ListValueFrom(ctx, types.StringType, y.Default.BeforeScript)
	afterScipt, _ := types.ListValueFrom(ctx, types.StringType, y.Default.AfterScript)
	paths, _ := types.ListValueFrom(ctx, types.StringType, y.Default.Cache.Paths)
	cache, _ := types.ObjectValue(cacheStateType, map[string]attr.Value{"key": types.StringValue(y.Default.Cache.Key), "paths": paths})
	v := map[string]attr.Value{
		"image":         types.StringValue(y.Default.Image),
		"timeout":       types.StringValue(y.Default.Timeout),
		"tags":          tags,
		"retry":         types.Int64Value(y.Default.Retry),
		"interruptible": types.BoolValue(y.Default.Interruptible),
		"before_script": beforeScript,
		"after_script":  afterScipt,
		"cache":         cache,
	}
	defaultobj, _ := types.ObjectValue(defaultStateType.AttrTypes, v)

	state := schema.FileDataSourceModel{FileLocation: fileLocation, Default: defaultobj}
	for _, stage := range y.Stages {
		state.Stages = append(state.Stages, types.StringValue(stage))
	}
	return state
}
