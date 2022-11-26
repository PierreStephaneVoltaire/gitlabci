package gitlabci

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"strconv"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &FileDataSource{}
)

type fileDataSourceModel struct {
	FileLocation types.String `tfsdk:"file_location"`
	Id           types.String `tfsdk:"id"`
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
			"id": {
				MarkdownDescription: "identifier",
				Type:                types.StringType,
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
	//read the file
	maxSz, _ := strconv.Atoi("10")
	file, err := os.Open(data.FileLocation.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not open the file",
			err.Error(),
		)
		return
	}
	defer file.Close()
	b := make([]byte, maxSz)
	fileContent := ""
	for {
		// read content to buffer
		readTotal, err := file.Read(b)
		if err != nil {
			if err != io.EOF {
				resp.Diagnostics.AddError(
					"Could not read the file",
					err.Error(),
				)
				return
			}
			break
		}

		fileContent = string(b[:readTotal])
	}
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
	diags := resp.State.Set(ctx, &y)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
