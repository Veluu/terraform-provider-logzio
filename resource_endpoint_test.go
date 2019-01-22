package main

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jonboydell/logzio_client/endpoints"
	"os"
	"regexp"
	"strconv"
	"testing"
)

func TestAccLogzioEndpoint_Slack_HappyPath(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccLogzioEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLogzioEndpointConfig("slackHappyPath"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogzioEndpointExists("logzio_endpoint.slack"),
					resource.TestCheckResourceAttr(
						"logzio_endpoint.slack", "title", "my_slack_title"),
					testAccCheckOutputExists("logzio_endpoint.slack"),
				),
			},
		},
	})
}

func TestAccLogzioEndpoint_Slack_BadUrl(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccLogzioEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckLogzioEndpointConfig("slackBadUrl"),
				ExpectError: regexp.MustCompile("Bad URL provided"),
			},
		},
	})
}

func TestAccLogzioEndpoint_Slack_UpdateHappyPath(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccLogzioEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLogzioEndpointConfig("slackHappyPath"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogzioEndpointExists("logzio_endpoint.slack"),
					resource.TestCheckResourceAttr(
						"logzio_endpoint.slack", "title", "my_slack_title"),
				),
			},
			{
				Config: testAccCheckLogzioEndpointConfig("slackUpdateHappyPath"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogzioEndpointExists("logzio_endpoint.slack"),
					resource.TestCheckResourceAttr(
						"logzio_endpoint.slack", "title", "my_updated_slack_title"),
				),
			},
		},
	})
}

func TestAccLogzioEndpoint_Custom_HappyPath(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccLogzioEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLogzioEndpointConfig("customHappyPath"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogzioEndpointExists("logzio_endpoint.custom"),
					resource.TestCheckResourceAttr(
						"logzio_endpoint.name", "title", "my_custom_title"),
				),
			},
		},
	})
}

func testAccCheckOutputExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		id := rs.Primary.ID
		os, ok := s.RootModule().Outputs["test"]

		if rs.Primary.ID == "" {
			return errors.New("no endpoint ID is set")
		}

		if os.Value != id {
			return fmt.Errorf("can't find resource that matches output ID")
		}

		return nil
	}
}

func testAccCheckLogzioEndpointExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no endpoint ID is set")
		}

		id, err := strconv.ParseInt(rs.Primary.ID, BASE_10, BITSIZE_64)

		var client *endpoints.Endpoints
		client, _ = endpoints.New(os.Getenv(envLogzioApiToken))

		_, err = client.GetEndpoint(int64(id))

		if err != nil {
			return fmt.Errorf("endpoint doesn't exist")
		}

		return nil
	}
}

func testAccLogzioEndpointDestroy(s *terraform.State) error {
	for _, r := range s.RootModule().Resources {
		id, err := strconv.ParseInt(r.Primary.ID, BASE_10, BITSIZE_64)
		if err != nil {
			return err
		}

		var client *endpoints.Endpoints
		client, _ = endpoints.New(os.Getenv(envLogzioApiToken))

		_, err = client.GetEndpoint(int64(id))
		if err == nil {
			return fmt.Errorf("endpoint still exists")
		}
	}
	return nil
}

func testAccCheckLogzioEndpointConfig(key string) string {
	templates := map[string]string{
		"slackHappyPath": `
resource "logzio_endpoint" "slack" {
  title = "my_slack_title"
  endpoint_type = "slack"
  description = "this_is_my_description"
  slack {
	url = "https://www.test.com"
  }
}

output "test" {
	value = "${logzio_endpoint.slack.endpoint_id}"
}
`,
		"slackBadUrl": `
resource "logzio_endpoint" "slack" {
  title = "my_slack_title"
  endpoint_type = "slack"
  description = "this_is_my_description"
  slack {
	url = "https://not_a_url"
  }
}
`,
		"slackUpdateHappyPath": `
resource "logzio_endpoint" "slack" {
  title = "my_updated_slack_title"
  endpoint_type = "slack"
  description = "this_is_my_description"
  slack {
	url = "https://www.test.com"
  }
}
`,
		"customHappyPath": `
resource "logzio_endpoint" "custom" {
  title = "my_custom_title"
  endpoint_type = "custom"
  description = "this_is_my_description"
  custom {
	url = "https://www.test.com"
	method = "POST"
	headers = {
		"this" = "is"
		"a" = "header"
	}
	body_template = "this_is_my_template"
  }
}
`,
	}
	return templates[key]
}
