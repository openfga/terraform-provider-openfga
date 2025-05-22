# Terraform provider for OpenFGA

[![Go Reference](https://pkg.go.dev/badge/github.com/openfga/terraform-provider-openfga.svg)](https://pkg.go.dev/github.com/openfga/terraform-provider-openfga)
[![Release](https://img.shields.io/github/v/release/openfga/terraform-provider-openfga?sort=semver&color=green)](https://github.com/openfga/terraform-provider-openfga/releases)
[![Go Report](https://goreportcard.com/badge/github.com/openfga/openfga)](https://goreportcard.com/report/github.com/openfga/openfga)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](./LICENSE)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fopenfga%2Fterraform-provider-openfga.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fopenfga%2Fterraform-provider-openfga?ref=badge_shield)
[![Join our community](https://img.shields.io/badge/slack-cncf_%23openfga-40abb8.svg?logo=slack)](https://openfga.dev/community)
[![Twitter](https://img.shields.io/twitter/follow/openfga?color=%23179CF0&logo=twitter&style=flat-square "@openfga on Twitter")](https://twitter.com/openfga)

This is a Terraform/ OpenTofu provider for OpenFGA. It allows to manage OpenFGA resources with code. for more details, check the [provider documentation](https://registry.terraform.io/providers/openfga/openfga/latest/docs).

## Table of Contents

- [About OpenFGA](#about)
- [Resources](#resources)
- [Installation](#installation)
- [Getting Started](#getting-started)
  - [Initializing the Provider](#initializing-the-provider)
  - [Using the Provider](#using-the-provider)
    - [Stores](#stores)
      - [Create Store](#create-store)
      - [Get Store](#get-store)
      - [List Stores](#list-stores)
    - [Authorization Models](#authorization-models)
      - [Authorization Model Documents](#authorization-model-documents)
      - [Create Authorization Model](#create-authorization-model)
      - [Get Authorization Model](#get-authorization-model)
      - [Get Latest Authorization Model](#get-latest-authorization-model)
      - [List Authorization Models](#list-authorization-models)
    - [Relationship Tuples](#relationship-tuples)
      - [Create Relationship Tuple](#create-relationship-tuple)
      - [Get Relationship Tuple](#get-relationship-tuple)
      - [List Relationship Tuples](#list-relationship-tuples)
      - [Query Relationship Tuples](#query-relationship-tuples)
    - [Relationship Queries](#relationship-queries)
      - [Check](#check)
      - [List Objects](#list-objects)
      - [List Users](#list-users)
- [Contributing](#contributing)
- [Author](#author)
- [License](#license)

## About

[OpenFGA](https://openfga.dev) is an open source Fine-Grained Authorization solution inspired by [Google's Zanzibar paper](https://research.google/pubs/pub48190/). It was created by the FGA team at [Auth0](https://auth0.com) based on [Auth0 Fine-Grained Authorization (FGA)](https://fga.dev), available under [a permissive license (Apache-2)](https://github.com/openfga/rfcs/blob/main/LICENSE) and welcomes community contributions.

OpenFGA is designed to make it easy for application builders to model their permission layer, and to add and integrate fine-grained authorization into their applications. OpenFGA’s design is optimized for reliability and low latency at a high scale.

## Resources

- [OpenFGA Documentation](https://openfga.dev/docs)
- [OpenFGA API Documentation](https://openfga.dev/api/service)
- [Twitter](https://twitter.com/openfga)
- [OpenFGA Community](https://openfga.dev/community)
- [Zanzibar Academy](https://zanzibar.academy)
- [Google's Zanzibar Paper (2019)](https://research.google/pubs/pub48190/)

## Installation

To install, add the provider to your configuration:

```terraform
terraform {
  required_providers {
    openfga = {
      source  = "openfga/openfga"
      version = ">=0.4.0"
    }
  }
}
```

Then run terraform init:

```shell
terraform init
```

## Getting Started

### Initializing the Provider

After installation, configure the provider to connect to your OpenFGA server.

#### No Credentials

```terraform
provider "openfga" {
  api_url = "http://openfga:8080" # or use FGA_API_URL
}
```

#### API Token

```terraform
provider "openfga" {
  api_url   = "http://openfga:8080" # or use FGA_API_URL
  api_token = var.api_token         # or use FGA_API_TOKEN
}
```

#### OAuth2 Client Credentials

```terraform
provider "openfga" {
  api_url            = "http://openfga:8080" # or use FGA_API_URL
  client_id          = "..."                 # or use FGA_CLIENT_ID
  client_secret      = var.client_secret     # or use FGA_CLIENT_SECRET
  token_endpoint_url = "http://example.com"  # or use FGA_TOKEN_ENDPOINT_URL
  audience           = "..."                 # or use FGA_AUDIENCE
  scopes             = "..."                 # or use FGA_SCOPES
}
```

#### Environment Variables

You can also use environment variables to configure the provider. In this case, you can leave the provider block empty. If both environment variable and provider config a specified, the provider config takes precedence.

```terraform
provider "openfga" {}
```

The available environment variables are:
- `FGA_API_URL`
- `FGA_API_TOKEN`
- `FGA_CLIENT_ID`
- `FGA_CLIENT_SECRET`
- `FGA_SCOPES`
- `FGA_AUDIENCE`
- `FGA_TOKEN_ENDPOINT_URL`

### Using the Provider

#### Stores

##### Create Store

Create and initialize a store.

[Terraform Documentation](https://registry.terraform.io/providers/openfga/openfga/latest/docs/resources/store)

```terraform
resource "openfga_store" "example" {
  name = "FGA Demo"
}
```

##### Get Store

Get information about a store by ID.

[Terraform Documentation](https://registry.terraform.io/providers/openfga/openfga/latest/docs/data-sources/store)

```terraform
data "openfga_store" "example" {
  id = "01FQH7V8BEG3GPQW93KTRFR8JB"
}
```

##### List Stores

Get a list of stores.

[Terraform Documentation](https://registry.terraform.io/providers/openfga/openfga/latest/docs/data-sources/stores)

```terraform
data "openfga_stores" "example" {}
```

#### Authorization Models

##### Authorization Model Documents

Create a stable JSON representation of an authorization model.

[Terraform Documentation](https://registry.terraform.io/providers/openfga/openfga/latest/docs/data-sources/authorization_model_document)

This data source takes authorization models in different formats as an input and produces a semantiaclly equal JSON output for the use in a `openfga_authorization_model` resource. The output of this data source will only change if there are semantic changes to a model (i.e., the output won't change for formatting changes, etc.)

> Note: To learn how to build your authorization model, check the Docs at https://openfga.dev/docs.

> Learn more about [the OpenFGA configuration language](https://openfga.dev/docs/configuration-language).

```terraform
data "openfga_authorization_model_document" "dsl" {
  dsl = file("path/to/model.fga")
}

data "openfga_authorization_model_document" "json" {
  json = file("path/to/model.json")
}

data "openfga_authorization_model_document" "mod" {
  mod_file_path = "path/to/fga.mod"
}

data "openfga_authorization_model_document" "model" {
  model = {
    schema_version = "1.1"
    type_definitions = [{
      type = "user"
    }]
  }
}
```

##### Create Authorization Model

Create a new authorization model.

[Terraform Documentation](https://registry.terraform.io/providers/openfga/openfga/latest/docs/resources/authorization_model)

> Note: You should use the `openfga_authorization_model_document` data source when when creating an authoriuation model.

```terraform
resource "openfga_authorization_model" "example" {
  store_id = "01FQH7V8BEG3GPQW93KTRFR8JB"

  model_json = data.openfga_authorization_model_document.example.result
}
```

##### Get Authorization Model

Get an authorization model in a store by ID.

[Terraform Documentation](https://registry.terraform.io/providers/openfga/openfga/latest/docs/data-sources/authorization_model)

```terraform
data "openfga_authorization_model" "specific" {
  store_id = "01FQH7V8BEG3GPQW93KTRFR8JB"

  id = "01GXSA8YR785C4FYS3C0RTG7B1"
}
```

##### Get Latest Authorization Model

Get latest authorization model in a store.

[Terraform Documentation](https://registry.terraform.io/providers/openfga/openfga/latest/docs/data-sources/authorization_model)

```terraform
data "openfga_authorization_model" "example" {
  store_id = "01FQH7V8BEG3GPQW93KTRFR8JB"
}
```

##### List Authorization Models

Get a list of authorization models in a store.

[Terraform Documentation](https://registry.terraform.io/providers/openfga/openfga/latest/docs/data-sources/authorization_models)

```terraform
data "openfga_authorization_models" "example" {
  store_id = "01FQH7V8BEG3GPQW93KTRFR8JB"
}
```

#### Relationship Tuples

##### Create Relationship Tuple

Create a new relationship tuple.

[Terraform Documentation](https://registry.terraform.io/providers/openfga/openfga/latest/docs/resources/relationship_tuple)

```terraform
resource "openfga_relationship_tuple" "example" {
  store_id               = "01FQH7V8BEG3GPQW93KTRFR8JB"
  authorization_model_id = "01GXSA8YR785C4FYS3C0RTG7B1" # optional

  user     = "user:81684243-9356-4421-8fbf-a4f8d36aa31b"
  relation = "viewer"
  object   = "document:0192ab2a-d83f-756d-9397-c5ed9f3cb69a"
}
```

##### Get Relationship Tuple

Get a relationship tuple in a store by attributes.

[Terraform Documentation](https://registry.terraform.io/providers/openfga/openfga/latest/docs/data-sources/relationship_tuple)

```terraform
data "openfga_relationship_tuple" "example" {
  store_id = "01FQH7V8BEG3GPQW93KTRFR8JB"

  user     = "user:81684243-9356-4421-8fbf-a4f8d36aa31b"
  relation = "viewer"
  object   = "document:0192ab2a-d83f-756d-9397-c5ed9f3cb69a"
}
```

##### List Relationship Tuples

Get all relationship tuple in a store.

[Terraform Documentation](https://registry.terraform.io/providers/openfga/openfga/latest/docs/data-sources/relationship_tuples)

```terraform
data "openfga_relationship_tuples" "example" {
  store_id = "01FQH7V8BEG3GPQW93KTRFR8JB"
}
```

##### Query Relationship Tuples

Get a list of relationship tuple in a store based on a query.

[Terraform Documentation](https://registry.terraform.io/providers/openfga/openfga/latest/docs/data-sources/relationship_tuples)

```terraform
data "openfga_relationship_tuples" "query" {
  store_id = "01FQH7V8BEG3GPQW93KTRFR8JB"

  query = {
    user     = "user:81684243-9356-4421-8fbf-a4f8d36aa31b"
    relation = "viewer"
    object   = "document:"
  }
}
```

#### Relationship Queries

##### Check

Check if a user has a particular relation with an object.

[Terraform Documentation](https://registry.terraform.io/providers/openfga/openfga/latest/docs/data-sources/check_query)

```terraform
data "openfga_check_query" "example" {
  store_id = "01FQH7V8BEG3GPQW93KTRFR8JB"

  user     = "user:81684243-9356-4421-8fbf-a4f8d36aa31b"
  relation = "viewer"
  object   = "document:0192ab2a-d83f-756d-9397-c5ed9f3cb69a"
}
```

You can also add contextual tuples and context to the query.

```terraform
data "openfga_check_query" "example" {
  store_id = "01FQH7V8BEG3GPQW93KTRFR8JB"

  user     = "user:81684243-9356-4421-8fbf-a4f8d36aa31b"
  relation = "viewer"
  object   = "document:0192ab2a-d83f-756d-9397-c5ed9f3cb69a"

  contextual_tuples = [
    {
      user     = "user:81684243-9356-4421-8fbf-a4f8d36aa31b"
      relation = "viewer"
      object   = "document:0192ab2a-d83f-756d-9397-c5ed9f3cb69a"
    }
  ]

  context_json = jsonencode({
    time = timestamp()
  })
}
```

##### List Objects

List the objects of a particular type a user has access to.

[Terraform Documentation](https://registry.terraform.io/providers/openfga/openfga/latest/docs/data-sources/list_objects_query)

```
data "openfga_list_objects_query" "example" {
  store_id = "01FQH7V8BEG3GPQW93KTRFR8JB"

  user     = "user:81684243-9356-4421-8fbf-a4f8d36aa31b"
  relation = "viewer"
  type     = "document"
}
```

You can also add contextual tuples and context to the query.

```terraform
data "openfga_list_objects_query" "example" {
  store_id = "01FQH7V8BEG3GPQW93KTRFR8JB"

  user     = "user:81684243-9356-4421-8fbf-a4f8d36aa31b"
  relation = "viewer"
  type     = "document"

  contextual_tuples = [
    {
      user     = "user:81684243-9356-4421-8fbf-a4f8d36aa31b"
      relation = "viewer"
      object   = "document:0192ab2a-d83f-756d-9397-c5ed9f3cb69a"
    }
  ]

  context_json = jsonencode({
    time = timestamp()
  })
}
```

##### List Users

List the users who have a certain relation to a particular type.

[Terraform Documentation](https://registry.terraform.io/providers/openfga/openfga/latest/docs/data-sources/list_users_query)

```
data "openfga_list_users_query" "example" {
  store_id = "01FQH7V8BEG3GPQW93KTRFR8JB"

  type     = "user"
  relation = "viewer"
  object   = "document:0192ab2a-d83f-756d-9397-c5ed9f3cb69a"
}
```

You can also add contextual tuples and context to the query.

```terraform
data "openfga_list_users_query" "example" {
  store_id = "01FQH7V8BEG3GPQW93KTRFR8JB"

  type     = "user"
  relation = "viewer"
  object   = "document:0192ab2a-d83f-756d-9397-c5ed9f3cb69a"

  contextual_tuples = [
    {
      user     = "user:81684243-9356-4421-8fbf-a4f8d36aa31b"
      relation = "viewer"
      object   = "document:0192ab2a-d83f-756d-9397-c5ed9f3cb69a"
    }
  ]

  context_json = jsonencode({
    time = timestamp()
  })
}
```

## Contributing

See [CONTRIBUTING](https://github.com/openfga/.github/blob/main/CONTRIBUTING.md).

## Author

[OpenFGA](https://github.com/openfga)
[Maurice Ackel](https://github.com/mauriceackel)

## License

This project is licensed under the Apache-2.0 license. See the [LICENSE](https://github.com/openfga/terraform-provider-openfga/blob/main/LICENSE) file for more info.
