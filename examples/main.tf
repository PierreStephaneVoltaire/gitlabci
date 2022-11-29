terraform {
  required_providers {
    gitlabci = {
      source = "hashicorp.com/barelycompetent/gitlabCi"
    }
  }
}

provider "gitlabci" {}

data "gitlabci_file" "example" {
  file_location = "gitlab-ci.yaml"
}