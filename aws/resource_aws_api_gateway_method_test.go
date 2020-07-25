package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAWSAPIGatewayMethod_basic(t *testing.T) {
	var conf apigateway.Method
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_api_gateway_method.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAPIGatewayMethodDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAPIGatewayMethodConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAPIGatewayMethodExists(resourceName, &conf),
					testAccCheckAWSAPIGatewayMethodAttributes(&conf),
					resource.TestCheckResourceAttr(resourceName, "http_method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "authorization", "NONE"),
					resource.TestCheckResourceAttr(resourceName, "request_models.application/json", "Error"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSAPIGatewayMethodImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},

			{
				Config: testAccAWSAPIGatewayMethodConfigUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAPIGatewayMethodExists(resourceName, &conf),
					testAccCheckAWSAPIGatewayMethodAttributesUpdate(&conf),
				),
			},
		},
	})
}

func TestAccAWSAPIGatewayMethod_customauthorizer(t *testing.T) {
	var conf apigateway.Method
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_api_gateway_method.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAPIGatewayMethodDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAPIGatewayMethodConfigWithCustomAuthorizer(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAPIGatewayMethodExists(resourceName, &conf),
					testAccCheckAWSAPIGatewayMethodAttributes(&conf),
					resource.TestCheckResourceAttr(resourceName, "http_method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "authorization", "CUSTOM"),
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_id", "aws_api_gateway_authorizer.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "request_models.application/json", "Error"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSAPIGatewayMethodImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},

			{
				Config: testAccAWSAPIGatewayMethodConfigUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAPIGatewayMethodExists(resourceName, &conf),
					testAccCheckAWSAPIGatewayMethodAttributesUpdate(&conf),
					resource.TestCheckResourceAttr(resourceName, "authorization", "NONE"),
					resource.TestCheckResourceAttr(resourceName, "authorizer_id", ""),
				),
			},
		},
	})
}

func TestAccAWSAPIGatewayMethod_cognitoauthorizer(t *testing.T) {
	var conf apigateway.Method
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_api_gateway_method.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAPIGatewayMethodDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAPIGatewayMethodConfigWithCognitoAuthorizer(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAPIGatewayMethodExists(resourceName, &conf),
					testAccCheckAWSAPIGatewayMethodAttributes(&conf),
					resource.TestCheckResourceAttr(resourceName, "http_method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "authorization", "COGNITO_USER_POOLS"),
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_id", "aws_api_gateway_authorizer.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "request_models.application/json", "Error"),
					resource.TestCheckResourceAttr(resourceName, "authorization_scopes.#", "2"),
				),
			},

			{
				Config: testAccAWSAPIGatewayMethodConfigWithCognitoAuthorizerUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAPIGatewayMethodExists(resourceName, &conf),
					testAccCheckAWSAPIGatewayMethodAttributesUpdate(&conf),
					resource.TestCheckResourceAttr(resourceName, "authorization", "COGNITO_USER_POOLS"),
					resource.TestCheckResourceAttrPair(resourceName, "authorizer_id", "aws_api_gateway_authorizer.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "request_models.application/json", "Error"),
					resource.TestCheckResourceAttr(resourceName, "authorization_scopes.#", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSAPIGatewayMethodImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSAPIGatewayMethod_customrequestvalidator(t *testing.T) {
	var conf apigateway.Method
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_api_gateway_method.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAPIGatewayMethodDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAPIGatewayMethodConfigWithCustomRequestValidator(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAPIGatewayMethodExists(resourceName, &conf),
					testAccCheckAWSAPIGatewayMethodAttributes(&conf),
					resource.TestCheckResourceAttr(resourceName, "http_method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "authorization", "NONE"),
					resource.TestCheckResourceAttr(resourceName, "request_models.application/json", "Error"),
					resource.TestCheckResourceAttrPair(resourceName, "request_validator_id", "aws_api_gateway_request_validator.test", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSAPIGatewayMethodImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},

			{
				Config: testAccAWSAPIGatewayMethodConfigWithCustomRequestValidatorUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAPIGatewayMethodExists(resourceName, &conf),
					testAccCheckAWSAPIGatewayMethodAttributesUpdate(&conf),
					resource.TestCheckResourceAttr(resourceName, "request_validator_id", ""),
				),
			},
		},
	})
}

func TestAccAWSAPIGatewayMethod_disappears(t *testing.T) {
	var conf apigateway.Method
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_api_gateway_method.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAPIGatewayMethodDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAPIGatewayMethodConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAPIGatewayMethodExists(resourceName, &conf),
					testAccCheckResourceDisappears(testAccProvider, resourceAwsApiGatewayMethod(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckAWSAPIGatewayMethodAttributes(conf *apigateway.Method) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if aws.StringValue(conf.HttpMethod) != "GET" {
			return fmt.Errorf("Wrong HttpMethod: %q", *conf.HttpMethod)
		}
		if aws.StringValue(conf.AuthorizationType) != "NONE" &&
			aws.StringValue(conf.AuthorizationType) != "CUSTOM" &&
			aws.StringValue(conf.AuthorizationType) != "COGNITO_USER_POOLS" {
			return fmt.Errorf("Wrong Authorization: %q", *conf.AuthorizationType)
		}

		if val, ok := conf.RequestParameters["method.request.header.Content-Type"]; !ok {
			return fmt.Errorf("missing Content-Type RequestParameters")
		} else {
			if aws.BoolValue(val) {
				return fmt.Errorf("wrong Content-Type RequestParameters value")
			}
		}
		if val, ok := conf.RequestParameters["method.request.querystring.page"]; !ok {
			return fmt.Errorf("missing page RequestParameters")
		} else {
			if !aws.BoolValue(val) {
				return fmt.Errorf("wrong query string RequestParameters value")
			}
		}

		return nil
	}
}

func testAccCheckAWSAPIGatewayMethodAttributesUpdate(conf *apigateway.Method) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if aws.StringValue(conf.HttpMethod) != "GET" {
			return fmt.Errorf("Wrong HttpMethod: %q", *conf.HttpMethod)
		}
		if conf.RequestParameters["method.request.header.Content-Type"] != nil {
			return fmt.Errorf("Content-Type RequestParameters shouldn't exist")
		}
		if val, ok := conf.RequestParameters["method.request.querystring.page"]; !ok {
			return fmt.Errorf("missing updated page RequestParameters")
		} else {
			if aws.BoolValue(val) {
				return fmt.Errorf("wrong query string RequestParameters updated value")
			}
		}

		return nil
	}
}

func testAccCheckAWSAPIGatewayMethodExists(n string, res *apigateway.Method) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No API Gateway Method ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).apigatewayconn

		req := &apigateway.GetMethodInput{
			HttpMethod: aws.String("GET"),
			ResourceId: aws.String(s.RootModule().Resources["aws_api_gateway_resource.test"].Primary.ID),
			RestApiId:  aws.String(s.RootModule().Resources["aws_api_gateway_rest_api.test"].Primary.ID),
		}
		describe, err := conn.GetMethod(req)
		if err != nil {
			return err
		}

		*res = *describe

		return nil
	}
}

func testAccCheckAWSAPIGatewayMethodDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).apigatewayconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_api_gateway_method" {
			continue
		}

		req := &apigateway.GetMethodInput{
			HttpMethod: aws.String("GET"),
			ResourceId: aws.String(s.RootModule().Resources["aws_api_gateway_resource.test"].Primary.ID),
			RestApiId:  aws.String(s.RootModule().Resources["aws_api_gateway_rest_api.test"].Primary.ID),
		}
		_, err := conn.GetMethod(req)

		if err == nil {
			return fmt.Errorf("API Gateway Method still exists")
		}

		if !isAWSErr(err, apigateway.ErrCodeNotFoundException, "") {
			return err
		}

		return nil
	}

	return nil
}

func testAccAWSAPIGatewayMethodImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["rest_api_id"], rs.Primary.Attributes["resource_id"], rs.Primary.Attributes["http_method"]), nil
	}
}

func testAccAWSAPIGatewayMethodConfigWithCustomAuthorizer(rName string) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = %[1]q
}

resource "aws_iam_role" "test" {
  name = %[1]q
  path = "/"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "apigateway.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "test" {
  name = %[1]q
  role = "${aws_iam_role.test.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "lambda:InvokeFunction",
      "Effect": "Allow",
      "Resource": "${aws_lambda_function.test.arn}"
    }
  ]
}
EOF
}

resource "aws_iam_role" "lambda" {
  name = "%[1]s-lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_lambda_function" "test" {
  filename         = "test-fixtures/lambdatest.zip"
  source_code_hash = "${filebase64sha256("test-fixtures/lambdatest.zip")}"
  function_name    = %[1]q
  role             = "${aws_iam_role.lambda.arn}"
  handler          = "exports.example"
  runtime          = "nodejs12.x"
}

resource "aws_api_gateway_authorizer" "test" {
  name                   = %[1]q
  rest_api_id            = "${aws_api_gateway_rest_api.test.id}"
  authorizer_uri         = "${aws_lambda_function.test.invoke_arn}"
  authorizer_credentials = "${aws_iam_role.test.arn}"
}

resource "aws_api_gateway_resource" "test" {
  rest_api_id = "${aws_api_gateway_rest_api.test.id}"
  parent_id   = "${aws_api_gateway_rest_api.test.root_resource_id}"
  path_part   = "test"
}

resource "aws_api_gateway_method" "test" {
  rest_api_id   = "${aws_api_gateway_rest_api.test.id}"
  resource_id   = "${aws_api_gateway_resource.test.id}"
  http_method   = "GET"
  authorization = "CUSTOM"
  authorizer_id = "${aws_api_gateway_authorizer.test.id}"

  request_models = {
    "application/json" = "Error"
  }

  request_parameters = {
    "method.request.header.Content-Type" = false
    "method.request.querystring.page"    = true
  }
}
`, rName)
}

func testAccAWSAPIGatewayMethodConfigWithCognitoAuthorizer(rName string) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = %[1]q
}

resource "aws_iam_role" "test" {
  name = %[1]q
  path = "/"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "apigateway.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role" "lambda" {
  name = "%[1]s-lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_cognito_user_pool" "test" {
  name = %[1]q
}

resource "aws_api_gateway_authorizer" "test" {
  name            = "tf-acc-test-cognito-authorizer"
  rest_api_id     = "${aws_api_gateway_rest_api.test.id}"
  identity_source = "method.request.header.Authorization"
  provider_arns   = ["${aws_cognito_user_pool.test.arn}"]
  type            = "COGNITO_USER_POOLS"
}

resource "aws_api_gateway_resource" "test" {
  rest_api_id = "${aws_api_gateway_rest_api.test.id}"
  parent_id   = "${aws_api_gateway_rest_api.test.root_resource_id}"
  path_part   = "test"
}

resource "aws_api_gateway_method" "test" {
  rest_api_id          = "${aws_api_gateway_rest_api.test.id}"
  resource_id          = "${aws_api_gateway_resource.test.id}"
  http_method          = "GET"
  authorization        = "COGNITO_USER_POOLS"
  authorizer_id        = "${aws_api_gateway_authorizer.test.id}"
  authorization_scopes = ["test.read", "test.write"]

  request_models = {
    "application/json" = "Error"
  }

  request_parameters = {
    "method.request.header.Content-Type" = false
    "method.request.querystring.page"    = true
  }
}
`, rName)
}

func testAccAWSAPIGatewayMethodConfigWithCognitoAuthorizerUpdate(rName string) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = %[1]q
}

resource "aws_iam_role" "test" {
  name = %[1]q
  path = "/"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "apigateway.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role" "lambda" {
  name = "%[1]s-lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_cognito_user_pool" "test" {
  name = %[1]q
}

resource "aws_api_gateway_authorizer" "test" {
  name            = %[1]q
  rest_api_id     = "${aws_api_gateway_rest_api.test.id}"
  identity_source = "method.request.header.Authorization"
  provider_arns   = ["${aws_cognito_user_pool.test.arn}"]
  type            = "COGNITO_USER_POOLS"
}

resource "aws_api_gateway_resource" "test" {
  rest_api_id = "${aws_api_gateway_rest_api.test.id}"
  parent_id   = "${aws_api_gateway_rest_api.test.root_resource_id}"
  path_part   = "test"
}

resource "aws_api_gateway_method" "test" {
  rest_api_id          = "${aws_api_gateway_rest_api.test.id}"
  resource_id          = "${aws_api_gateway_resource.test.id}"
  http_method          = "GET"
  authorization        = "COGNITO_USER_POOLS"
  authorizer_id        = "${aws_api_gateway_authorizer.test.id}"
  authorization_scopes = ["test.read", "test.write", "test.delete"]

  request_models = {
    "application/json" = "Error"
  }

  request_parameters = {
    "method.request.querystring.page" = false
  }
}
`, rName)
}

func testAccAWSAPIGatewayMethodConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = %[1]q
}

resource "aws_api_gateway_resource" "test" {
  rest_api_id = "${aws_api_gateway_rest_api.test.id}"
  parent_id   = "${aws_api_gateway_rest_api.test.root_resource_id}"
  path_part   = "test"
}

resource "aws_api_gateway_method" "test" {
  rest_api_id   = "${aws_api_gateway_rest_api.test.id}"
  resource_id   = "${aws_api_gateway_resource.test.id}"
  http_method   = "GET"
  authorization = "NONE"

  request_models = {
    "application/json" = "Error"
  }

  request_parameters = {
    "method.request.header.Content-Type" = false
    "method.request.querystring.page"    = true
  }
}
`, rName)
}

func testAccAWSAPIGatewayMethodConfigUpdate(rName string) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = %[1]q
}

resource "aws_api_gateway_resource" "test" {
  rest_api_id = "${aws_api_gateway_rest_api.test.id}"
  parent_id   = "${aws_api_gateway_rest_api.test.root_resource_id}"
  path_part   = "test"
}

resource "aws_api_gateway_method" "test" {
  rest_api_id   = "${aws_api_gateway_rest_api.test.id}"
  resource_id   = "${aws_api_gateway_resource.test.id}"
  http_method   = "GET"
  authorization = "NONE"

  request_models = {
    "application/json" = "Error"
  }

  request_parameters = {
    "method.request.querystring.page" = false
  }
}
`, rName)
}

func testAccAWSAPIGatewayMethodConfigWithCustomRequestValidator(rName string) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = %[1]q
}

resource "aws_api_gateway_resource" "test" {
  rest_api_id = "${aws_api_gateway_rest_api.test.id}"
  parent_id   = "${aws_api_gateway_rest_api.test.root_resource_id}"
  path_part   = "test"
}

resource "aws_api_gateway_request_validator" "test" {
  rest_api_id                 = "${aws_api_gateway_rest_api.test.id}"
  name                        = "paramsValidator"
  validate_request_parameters = true
}

resource "aws_api_gateway_method" "test" {
  rest_api_id   = "${aws_api_gateway_rest_api.test.id}"
  resource_id   = "${aws_api_gateway_resource.test.id}"
  http_method   = "GET"
  authorization = "NONE"

  request_models = {
    "application/json" = "Error"
  }

  request_parameters = {
    "method.request.header.Content-Type" = false
    "method.request.querystring.page"    = true
  }

  request_validator_id = "${aws_api_gateway_request_validator.test.id}"
}
`, rName)
}

func testAccAWSAPIGatewayMethodConfigWithCustomRequestValidatorUpdate(rName string) string {
	return fmt.Sprintf(`
resource "aws_api_gateway_rest_api" "test" {
  name = %[1]q
}

resource "aws_api_gateway_resource" "test" {
  rest_api_id = "${aws_api_gateway_rest_api.test.id}"
  parent_id   = "${aws_api_gateway_rest_api.test.root_resource_id}"
  path_part   = "test"
}

resource "aws_api_gateway_request_validator" "test" {
  rest_api_id                 = "${aws_api_gateway_rest_api.test.id}"
  name                        = "paramsValidator"
  validate_request_parameters = true
}

resource "aws_api_gateway_method" "test" {
  rest_api_id   = "${aws_api_gateway_rest_api.test.id}"
  resource_id   = "${aws_api_gateway_resource.test.id}"
  http_method   = "GET"
  authorization = "NONE"

  request_models = {
    "application/json" = "Error"
  }

  request_parameters = {
    "method.request.querystring.page" = false
  }
}
`, rName)
}
