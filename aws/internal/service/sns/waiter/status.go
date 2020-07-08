package waiter

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// SNSSubscriptionStatus fetches the Operation and its Status
func SnsTopicSubscriptionStatus(conn *sns.SNS, subscriptionArn string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		output, err := conn.GetSubscriptionAttributes(&sns.GetSubscriptionAttributesInput{
			SubscriptionArn: aws.String(subscriptionArn),
		})

		if isAWSErr(err, sns.ErrCodeNotFoundException, "") {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", fmt.Errorf("error reading SNS topic subscription (%s): %s", subscriptionArn, err)
		}

		if output == nil || len(output.Attributes) == 0 {
			return nil, "", nil
		}

		return output, "available", nil
	}
}
