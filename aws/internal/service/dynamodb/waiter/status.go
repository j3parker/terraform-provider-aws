package waiter

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// ContributorInsightsStatus fetches the ContributorInsights and its Status
func ContributorInsightsStatus(conn *dynamodb.DynamoDB, tableName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		input := &dynamodb.DescribeContributorInsightsInput{
			TableName: aws.String(tableName),
		}

		output, err := conn.DescribeContributorInsights(input)

		if err != nil {
			return nil, dynamodb.ContributorInsightsStatusFailed, err
		}

		// Error messages can also be contained in the response with FAILED status

		if aws.StringValue(output.ContributorInsightsStatus) == dynamodb.ContributorInsightsStatusFailed {
			return output, dynamodb.ContributorInsightsStatusFailed, fmt.Errorf("%s", aws.StringValue(output.ContributorInsightsStatus))
		}

		return output, aws.StringValue(output.ContributorInsightsStatus), nil
	}
}
