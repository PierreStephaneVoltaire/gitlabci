package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GitlabDefault struct {
	Artifacts     GitlabArtifacts `yaml:"artifacts"`
	AfterScript   []string        `yaml:"after_script,flow"`
	BeforeScript  []string        `yaml:"before_script,flow"`
	Image         string          `yaml:"image"`
	Timeout       string          `yaml:"timeout"`
	Tags          []string        `yaml:"tags"`
	Retry         int64           `yaml:"retry"`
	Interruptible bool            `yaml:"interruptible"`
	Cache         GitlabCache     `yaml:"cache"`
	Services      []GitlabService `yaml:"services"`
}
type GitlabCache struct {
	Key   string   `yaml:"key"`
	Paths []string `yaml:"paths,flow"`
}
type GitlabService struct {
	Name       string   `yaml:"name"`
	Alias      string   `yaml:"alias"`
	Entrypoint []string `yaml:"entrypoint,flow"`
	Command    []string `yaml:"command,flow"`
}
type GitlabArtifacts struct {
	Name      string   `yaml:"name"`
	When      string   `yaml:"when"`
	Public    bool     `yaml:"public"`
	Untracked bool     `yaml:"untracked"`
	ExposeAs  string   `yaml:"expose_as"`
	ExpireIn  string   `yaml:"expire_in"`
	Exclude   []string `yaml:"exclude,flow"`
	Paths     []string `yaml:"paths,flow"`
}
type GitlabDefaultDataSourceModel struct {
	//Artifacts     GitlabArtifacts `tfsdk:"artifacts"`
	AfterScript   []types.String             `tfsdk:"after_script"`
	BeforeScript  []types.String             `tfsdk:"before_script"`
	Image         types.String               `tfsdk:"image"`
	Timeout       types.String               `tfsdk:"timeout"`
	Tags          []types.String             `tfsdk:"tags"`
	Retry         types.Int64                `tfsdk:"retry"`
	Interruptible types.Bool                 `tfsdk:"interruptible"`
	Cache         GitlabCacheDataSourceModel `tfsdk:"cache"`
	//Services      []GitlabServiceDataSourceModel `tfsdk:"services"`
}
type GitlabCacheDataSourceModel struct {
	Key   types.String   `tfsdk:"key"`
	Paths []types.String `tfsdk:"paths"`
}
type GitlabServiceDataSourceModel struct {
	Name       types.String   `tfsdk:"name"`
	Alias      types.String   `tfsdk:"alias"`
	Entrypoint []types.String `tfsdk:"entrypoint"`
	Command    []types.String `tfsdk:"command"`
}
type GitlabArtifactsDataSourceModel struct {
	Name      types.String   `tfsdk:"name"`
	When      types.String   `tfsdk:"when"`
	Public    types.Bool     `tfsdk:"public"`
	Untracked types.Bool     `tfsdk:"untracked"`
	ExposeAs  types.String   `tfsdk:"expose_as"`
	ExpireIn  types.String   `tfsdk:"expire_in"`
	Exclude   []types.String `tfsdk:"exclude"`
	Paths     []types.String `tfsdk:"paths"`
}
type FileDataSourceModel struct {
	FileLocation types.String   `tfsdk:"file_location"`
	Stages       []types.String `tfsdk:"stages"`
	Default      types.Object   `tfsdk:"default"`
}
type Yaml struct {
	Stages  []string      `yaml:"stages,flow"`
	Default GitlabDefault `yaml:"default"`
}
