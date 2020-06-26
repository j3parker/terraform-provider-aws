package waiter

import (
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	// Maximum amount of time to wait for a ContributorInsightsUpdate to finish
	ContributorInsightsUpdateTimeout = 5 * time.Minute
)

// ContributorInsightsEnabled waits for a ContributorInsights to be Enabled
func ContributorInsightsEnabled(conn *dynamodb.DynamoDB, tableName string) (*dynamodb.DescribeContributorInsightsOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{dynamodb.ContributorInsightsStatusEnabling, dynamodb.ContributorInsightsStatusDisabled,
			dynamodb.ContributorInsightsStatusDisabling},
		Target:  []string{dynamodb.ContributorInsightsStatusEnabled},
		Refresh: ContributorInsightsStatus(conn, tableName),
		Timeout: ContributorInsightsUpdateTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if v, ok := outputRaw.(*dynamodb.DescribeContributorInsightsOutput); ok {
		return v, err
	}

	return nil, err
}

// ContributorInsightsEnabled waits for a ContributorInsights to be Enabled
func ContributorInsightsDisabled(conn *dynamodb.DynamoDB, tableName string) (*dynamodb.DescribeContributorInsightsOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{dynamodb.ContributorInsightsStatusEnabling, dynamodb.ContributorInsightsStatusDisabling,
			dynamodb.ContributorInsightsStatusEnabled},
		Target:  []string{dynamodb.ContributorInsightsStatusDisabled},
		Refresh: ContributorInsightsStatus(conn, tableName),
		Timeout: ContributorInsightsUpdateTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if v, ok := outputRaw.(*dynamodb.DescribeContributorInsightsOutput); ok {
		return v, err
	}

	return nil, err
}
