package gitlabci

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gopkg.in/yaml.v3"
	"os"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &FileDataSource{}
)

type fileDataSourceModel struct {
	FileLocation types.String   `tfsdk:"file_location"`
	Stages       []types.String `tfsdk:"stages"`
}
type Yaml struct {
	Stages []string `yaml:"stages,flow"`
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
		},
	}, nil
}

// Read refreshes the Terraform state with the latest data.
func (d *FileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// get value from schema
	var data fileDataSourceModel
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
	y := Yaml{}
	err = yaml.Unmarshal([]byte(fileContent), &y)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not parse the file",
			err.Error(),
		)
		return
	}
	// save state
	state := fileDataSourceModel{FileLocation: data.FileLocation}
	for _, stage := range y.Stages {
		state.Stages = append(state.Stages, types.StringValue(stage))
	}
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
