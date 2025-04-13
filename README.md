# Honest Bank Github Repository Template

* Terraform-managed template repository
* This repository is a template to be used for all Honest Bank Github repos

This template currently contains:

* [CODEOWNERS](./.github/CODEOWNERS) - currently set to `honestbank-engineers`
* [pull_request_template.md](./.github/pull_request_template.md) - as the name says
* [semantic.yaml](./.github/semantic.yaml) - settings for the [Semantic Pull Requests Github app](https://github.com/apps/semantic-pull-requests)

## Workflows

This template also contains the following Github Actions:

* [s3_upload.yaml](./.github/workflows/s3_upload.yaml) - uploads main/master to
  S3 for backup
* [trivy_scan.yaml](./.github/workflows/trivy_scan.yaml) - runs the [Trivy scanner](https://github.com/aquasecurity/trivy)
