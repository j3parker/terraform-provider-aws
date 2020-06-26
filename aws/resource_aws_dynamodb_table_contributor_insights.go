package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/service/dynamodb/waiter"
)

func resourceAwsDynamoDbTableContributorInsights() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsDynamoDbContributorInsightsUpdate,
		Read:   resourceAwsDynamoDbContributorInsightsRead,
		Update: resourceAwsDynamoDbContributorInsightsUpdate,
		Delete: resourceAwsDynamoDbContributorInsightsDelete,

		Schema: map[string]*schema.Schema{
			"table_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"contributor_insights_action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					dynamodb.ContributorInsightsActionDisable,
					dynamodb.ContributorInsightsActionEnable,
				}, false),
			},
			"index_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"contributor_insights_rule_list": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func resourceAwsDynamoDbContributorInsightsUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).dynamodbconn

	action := d.Get("contributor_insights_action").(string)
	input := &dynamodb.UpdateContributorInsightsInput{
		TableName:                 aws.String(d.Get("table_name").(string)),
		ContributorInsightsAction: aws.String(action),
	}

	if v, ok := d.GetOk("index_name"); ok {
		input.IndexName = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Updating DynamoDB Table Contributor Insights: %#v", input)

	out, err := conn.UpdateContributorInsights(input)
	if err != nil {
		return fmt.Errorf("error updating DynamoDB Table Contributor Insights: %s", err)
	}

	d.SetId(aws.StringValue(out.TableName))

	if action == dynamodb.ContributorInsightsActionEnable {
		_, err = waiter.ContributorInsightsEnabled(conn, d.Id())
		if err != nil {
			return fmt.Errorf("error waiting for DynamoDB Table Contributor Insights (%s) to be enabled: %w", d.Id(), err)
		}
	} else {
		_, err = waiter.ContributorInsightsDisabled(conn, d.Id())
		if err != nil {
			return fmt.Errorf("error waiting for DynamoDB Table Contributor Insights (%s) to be disabled: %w", d.Id(), err)
		}
	}

	return resourceAwsDynamoDbContributorInsightsRead(d, meta)
}

func resourceAwsDynamoDbContributorInsightsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).dynamodbconn

	log.Printf("[DEBUG] Describing DynamoDB Table (%s) Contributor Insights", d.Id())

	input := &dynamodb.DescribeContributorInsightsInput{
		TableName: aws.String(d.Id()),
	}

	if v, ok := d.GetOk("index_name"); ok {
		input.IndexName = aws.String(v.(string))
	}

	result, err := conn.DescribeContributorInsights(input)
	if err != nil {
		if isAWSErr(err, dynamodb.ErrCodeResourceNotFoundException, "") {
			log.Printf("[WARN] DynamoDB Table (%s) Contributor Insights not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("table_name", result.TableName)
	d.Set("index_name", result.IndexName)
	d.Set("contributor_insights_rule_list", flattenStringSet(result.ContributorInsightsRuleList))

	if aws.StringValue(result.ContributorInsightsStatus) == dynamodb.ContributorInsightsStatusEnabled {
		d.Set("contributor_insights_action", dynamodb.ContributorInsightsActionEnable)
	} else {
		d.Set("contributor_insights_action", dynamodb.ContributorInsightsActionDisable)
	}

	return nil
}

func resourceAwsDynamoDbContributorInsightsDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).dynamodbconn

	input := &dynamodb.UpdateContributorInsightsInput{
		TableName:                 aws.String(d.Get("table_name").(string)),
		ContributorInsightsAction: aws.String(dynamodb.ContributorInsightsActionDisable),
	}

	if v, ok := d.GetOk("index_name"); ok {
		input.IndexName = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Disabling DynamoDB Table Contributor Insights: %#v", input)

	_, err := conn.UpdateContributorInsights(input)
	if err != nil {
		return fmt.Errorf("error Disabling DynamoDB Table Contributor Insights: %s", err)
	}

	_, err = waiter.ContributorInsightsDisabled(conn, d.Id())
	if err != nil {
		return fmt.Errorf("error waiting for DynamoDB Table Contributor Insights (%s) to be disabled: %w", d.Id(), err)
	}

	return nil
}
