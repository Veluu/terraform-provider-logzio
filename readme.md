# Logz.io Terraform provider

The Terraform Logz.io provider offers a great way to build integrations using Logz.io APIs.

Terraform is an infrastructure orchestrator written in Hashicorp Language (HCL). It is a popular Infrastructure-as-Code (IaC) tool that does away with manual configuration processes. You can take advantage of the Terraform Logz.io provider to really streamline the process of integrating observability into your dev workflows.

This guide assumes working knowledge of HashiCorp Terraform. If you're new to Terraform, we've got a great [introduction](https://logz.io/blog/terraform-vs-ansible-vs-puppet/) if you're in for one. We also recommend the official [Terraform guides and tutorials](https://www.terraform.io/guides/index.html).

### Capabilities

You can use the Terraform Logz.io Provider to manage users and log accounts in Logz.io, create and update log-based alerts and notification channels, and more.

The following Logz.io API endpoints are supported by this provider:

* [User management](https://docs.logz.io/api/#tag/Manage-users)
* [Notification channels](https://docs.logz.io/api/#tag/Manage-notification-endpoints)
* [Log-based alerts](https://github.com/logzio/public-api/tree/master/alerts)
* [Sub accounts](https://docs.logz.io/api/#tag/Manage-sub-accounts)

#### Working with Terraform

<div class="tasklist">

**Before you begin, you'll need**:

* [Terraform CLI](https://learn.hashicorp.com/tutorials/terraform/install-cli)
* [Logz.io API token](/)

#### Get the Terraform Logz.io Provider

The easiest way to get the provider and the JetBrains IDE HCL meta-data is to run the script provided in the [Logz.io GitHub repo](https://github.com/logzio/logzio_terraform_provider/blob/master/scripts/update_plugin.sh).

The script is found under `./scripts/update_plugin.sh`. (If you ever encounter the need to update the version, you can edit the variable: `PROVIDER_VERSION`. But this shouldn't be necessary.)

Run it:

```bash
./scripts/build.sh
```


##### Configuring the provider

The provider accepts the following arguments:

* **api_token** - (Required) The API token is used for authentication. [Learn more](/user-guide/tokens/api-tokens.html).

* **region** - (Defaults to null) The 2-letter region code identifies where your Logz.io account is hosted.
Defaults to null for accounts hosted in the US East - Northern Virginia region. [Learn more](https://docs.logz.io/user-guide/accounts/account-region.html)

###### Example

You can pass the variables in a bash command for the arguments:

```bash
provider "logzio" {
  api_token = var.api_token
  region= var.your_api_region
}
```
</div>


### Example - Create a new alert and a new Slack notification endpoint

Here's a great example demonstrating how easy it is to get up and running quickly with the Terraform Logz.io Provider.

This example adds a new Slack notification channel and creates a new alert in Kibana that will send notifications to the newly-created Slack channel.

The alert in this example will trigger whenever Logz.io records 10 loglevel:ERROR messages in 10 minutes.

```
provider "logzio" {
  api_token = "8387abb8-4831-53af-91de-5cd3784d9774"
  region= "au"
}

resource "logzio_endpoint" "my_endpoint" {
  title = "my_endpoint"
  description = "my slack endpoint"
  endpoint_type = "slack"
  slack {
    url = "${var.slack_url}"
  }
}

resource "logzio_alert" "my_alert" {
  title = "my_other_title"
  query_string = "loglevel:ERROR"
  operation = "GREATER_THAN"
  notification_emails = []
  search_timeframe_minutes = 10
  value_aggregation_type = "NONE"
  alert_notification_endpoints = ["${logzio_endpoint.my_endpoint.id}"]
  suppress_notifications_minutes = 30
  severity_threshold_tiers {
      severity = "HIGH",
      threshold = 10
  }
}
```

### Example - Create user

This example will create a user in your Logz.io account.

```
variable "api_token" {
  type = "string"
  description = "Your logzio API token"
}

variable "account_id" {
  description = "The account ID where the new user will be created"
}

provider "logzio" {
  api_token = var.api_token
  region = var.region
}

resource "logzio_user" "my_user" {
  username = "user_name@fun.io"
  fullname = "John Doe"
  roles = [ 2 ]
  account_id = var.account_id
}
```

Run the above plan using the following bash script:

```
export TF_LOG=DEBUG
terraform init
TF_VAR_api_token=${LOGZIO_API_TOKEN} TF_VAR_region=${LOGZIO_REGION} terraform plan -out terraform.plan
terraform apply terraform.plan
```

Before you run the script, update the arguments to match your details.
See our [examples](https://github.com/logzio/logzio_terraform_provider/tree/master/examples) for some complete working examples. 

### Contribute
Found a bug or want to suggest a feature? [Open an issue](https://github.com/logzio/logzio_terraform_provider/issues/new) about it.
Want to do it yourself? We are more than happy to accept external contributions from the community.
Simply fork the repo, add your changes and [open a PR](https://github.com/logzio/logzio_terraform_provider/pulls).

### Changelog 
- v1.1.8
    - Update client version 
    - Fix custom endpoint headers bug
- v1.1.7
    - Published to Terraform registry    
- v1.1.5
    - Fix boolean parameters not parsed bug
    - Support import command to state
- v1.1.4
    - Support Sub Accounts resource
    - few bug fixes
    - removed circleCI  
- v1.1.3 
    - examples now use TF12
    - will now generate the meta data needed for the IntelliJ type IDE HCL plugin
    - no more travis - just circle CI
    - version bump to use the latest TF library (0.12.6), now compatible with TF12
- 1.1.2 
    - Moved some of the source code around to comply with TF provider layout convention
    - Moved the examples into an examples directory