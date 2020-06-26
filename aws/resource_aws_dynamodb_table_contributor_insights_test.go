package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAWSDynamoDbTableContributorInsights_basic(t *testing.T) {
	var conf dynamodb.DescribeContributorInsightsOutput
	resourceName := "aws_dynamodb_table_contributor_insights.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	hashKey := "hashKey"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDynamoDbContributorInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDynamoDbTableContributorInsightsConfigBasic(rName, hashKey, "ENABLE"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDynamoDbTableContributorInsightsExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "table_name", rName),
					resource.TestCheckResourceAttr(resourceName, "contributor_insights_action", "ENABLE"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSDynamoDbTableContributorInsightsConfigBasic(rName, hashKey, "DISABLE"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDynamoDbTableContributorInsightsExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "table_name", rName),
					resource.TestCheckResourceAttr(resourceName, "contributor_insights_action", "DISABLE"),
				),
			},
			{
				Config: testAccAWSDynamoDbTableContributorInsightsConfigBasic(rName, hashKey, "ENABLE"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDynamoDbTableContributorInsightsExists(resourceName, &conf),
					resource.TestCheckResourceAttr(resourceName, "table_name", rName),
					resource.TestCheckResourceAttr(resourceName, "contributor_insights_action", "ENABLE"),
				),
			},
		},
	})
}

func TestAccAWSDynamoDbTableContributorInsights_disappears(t *testing.T) {
	var conf dynamodb.DescribeContributorInsightsOutput
	resourceName := "aws_dynamodb_table_contributor_insights.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	hashKey := "hashKey"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSDynamoDbContributorInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSDynamoDbTableContributorInsightsConfigBasic(rName, hashKey, "ENABLE"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSDynamoDbTableContributorInsightsExists(resourceName, &conf),
					testAccCheckResourceDisappears(testAccProvider, resourceAwsDynamoDbTableContributorInsights(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckAWSDynamoDbContributorInsightsDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).dynamodbconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_dynamodb_table_contributor_insights" {
			continue
		}

		_, err := conn.DescribeContributorInsights(&dynamodb.DescribeContributorInsightsInput{
			TableName: aws.String(rs.Primary.ID),
		})
		if err != nil {
			if isAWSErr(err, dynamodb.ErrCodeResourceNotFoundException, "") {
				return nil
			}
			return fmt.Errorf("Error retrieving DynamoDB table contributor insights: %s", err)
		}

		return fmt.Errorf("DynamoDB table contributor insights %s still exists.", rs.Primary.ID)
	}

	return nil
}

func testAccCheckAWSDynamoDbTableContributorInsightsExists(n string, insights *dynamodb.DescribeContributorInsightsOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DynamoDB table contributor insights specified!")
		}

		conn := testAccProvider.Meta().(*AWSClient).dynamodbconn

		result, err := conn.DescribeContributorInsights(&dynamodb.DescribeContributorInsightsInput{
			TableName: aws.String(rs.Primary.ID),
		})
		if err != nil {
			return fmt.Errorf("Problem getting table item '%s': %s", rs.Primary.ID, err)
		}

		*insights = *result

		return nil
	}
}

func testAccAWSDynamoDbTableContributorInsightsConfigBasic(rName, hashKey, action string) string {
	return fmt.Sprintf(`
resource "aws_dynamodb_table" "test" {
  name           = %[1]q
  read_capacity  = 10
  write_capacity = 10
  hash_key       = %[2]q

  attribute {
    name = %[2]q
    type = "S"
  }
}

resource "aws_dynamodb_table_contributor_insights" "test" {
  table_name                  = "${aws_dynamodb_table.test.name}"
  contributor_insights_action = %[3]q
}
`, rName, hashKey, action)
}
