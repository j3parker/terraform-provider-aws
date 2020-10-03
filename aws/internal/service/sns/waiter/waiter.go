package waiter

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	// Maximum amount of time to wait for an Topic Subscription to be deleted
	TopicSubscriptionDeleteTimeout = 5 * time.Minute
)

// SnsTopicSubscriptionDeleted waits for an Topic Subscription to be deleted
func SnsTopicSubscriptionDeleted(conn *sns.SNS, subscriptionArn string) (*sns.GetSubscriptionAttributesOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"available"},
		Target:  []string{},
		Refresh: SnsTopicSubscriptionStatus(conn, subscriptionArn),
		Timeout: TopicSubscriptionDeleteTimeout,
	}

	log.Printf("[DEBUG] Waiting for SNS topic subscription (%s) deletion", subscriptionArn)
	_, err := stateConf.WaitForState()

	return nil, err
}
