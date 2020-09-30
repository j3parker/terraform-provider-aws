package aws

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

func init() {
	resource.AddTestSweepers("aws_sagemaker_notebook_instance", &resource.Sweeper{
		Name: "aws_sagemaker_notebook_instance",
		F:    testSweepSagemakerNotebookInstances,
	})
}

func testSweepSagemakerNotebookInstances(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.(*AWSClient).sagemakerconn

	err = conn.ListNotebookInstancesPages(&sagemaker.ListNotebookInstancesInput{}, func(page *sagemaker.ListNotebookInstancesOutput, lastPage bool) bool {
		for _, instance := range page.NotebookInstances {
			name := aws.StringValue(instance.NotebookInstanceName)
			status := aws.StringValue(instance.NotebookInstanceStatus)

			input := &sagemaker.DeleteNotebookInstanceInput{
				NotebookInstanceName: instance.NotebookInstanceName,
			}

			log.Printf("[INFO] Stopping SageMaker Notebook Instance: %s", name)
			if status != sagemaker.NotebookInstanceStatusFailed && status != sagemaker.NotebookInstanceStatusStopped {
				if err := stopSagemakerNotebookInstance(conn, name); err != nil {
					log.Printf("[ERROR] Error stopping SageMaker Notebook Instance (%s): %s", name, err)
					continue
				}
			}

			log.Printf("[INFO] Deleting SageMaker Notebook Instance: %s", name)
			if _, err := conn.DeleteNotebookInstance(input); err != nil {
				log.Printf("[ERROR] Error deleting SageMaker Notebook Instance (%s): %s", name, err)
				continue
			}

			stateConf := &resource.StateChangeConf{
				Pending: []string{sagemaker.NotebookInstanceStatusDeleting},
				Target:  []string{""},
				Refresh: sagemakerNotebookInstanceStateRefreshFunc(conn, name),
				Timeout: 10 * time.Minute,
			}

			if _, err := stateConf.WaitForState(); err != nil {
				log.Printf("[ERROR] Error waiting for SageMaker Notebook Instance (%s) deletion: %s", name, err)
			}
		}

		return !lastPage
	})

	if testSweepSkipSweepError(err) {
		log.Printf("[WARN] Skipping SageMaker Notebook Instance sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error retrieving SageMaker Notebook Instances: %s", err)
	}

	return nil
}

func TestAccAWSSagemakerNotebookInstance_basic(t *testing.T) {
	var notebook sagemaker.DescribeNotebookInstanceOutput
	rName := acctest.RandomWithPrefix("tf-acc-test")
	var resourceName = "aws_sagemaker_notebook_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSagemakerNotebookInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSagemakerNotebookInstanceBasicConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerNotebookInstanceExists(resourceName, &notebook),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "instance_type", "ml.t2.medium"),
					resource.TestCheckResourceAttrPair(resourceName, "role_arn", "aws_iam_role.test", "arn"),
					resource.TestCheckResourceAttr(resourceName, "direct_internet_access", "Enabled"),
					resource.TestCheckResourceAttr(resourceName, "root_access", "Enabled"),
					resource.TestCheckResourceAttr(resourceName, "security_groups.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSSagemakerNotebookInstance_update(t *testing.T) {
	var notebook sagemaker.DescribeNotebookInstanceOutput
	rName := acctest.RandomWithPrefix("tf-acc-test")
	var resourceName = "aws_sagemaker_notebook_instance.foo"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSagemakerNotebookInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSagemakerNotebookInstanceBasicConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerNotebookInstanceExists(resourceName, &notebook),
					resource.TestCheckResourceAttr(resourceName, "instance_type", "ml.t2.medium"),
				),
			},

			{
				Config: testAccAWSSagemakerNotebookInstanceUpdateConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerNotebookInstanceExists(resourceName, &notebook),
					resource.TestCheckResourceAttr(resourceName, "instance_type", "ml.m4.xlarge"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSSagemakerNotebookInstance_LifecycleConfigName(t *testing.T) {
	var notebook sagemaker.DescribeNotebookInstanceOutput
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_sagemaker_notebook_instance.test"
	sagemakerLifecycleConfigResourceName := "aws_sagemaker_notebook_instance_lifecycle_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSagemakerNotebookInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSagemakerNotebookInstanceConfigLifecycleConfigName(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerNotebookInstanceExists(resourceName, &notebook),
					resource.TestCheckResourceAttrPair(resourceName, "lifecycle_config_name", sagemakerLifecycleConfigResourceName, "name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSSagemakerNotebookInstance_tags(t *testing.T) {
	var notebook sagemaker.DescribeNotebookInstanceOutput
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSagemakerNotebookInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSagemakerNotebookInstanceConfigTags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerNotebookInstanceExists("aws_sagemaker_notebook_instance.foo", &notebook),
					resource.TestCheckResourceAttr("aws_sagemaker_notebook_instance.foo", "tags.%", "1"),
					resource.TestCheckResourceAttr("aws_sagemaker_notebook_instance.foo", "tags.foo", "bar"),
				),
			},

			{
				Config: testAccAWSSagemakerNotebookInstanceConfigTags2(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerNotebookInstanceExists("aws_sagemaker_notebook_instance.foo", &notebook),
					testAccCheckAWSSagemakerNotebookInstanceTags(&notebook, "foo", ""),
					testAccCheckAWSSagemakerNotebookInstanceTags(&notebook, "bar", "baz"),

					resource.TestCheckResourceAttr("aws_sagemaker_notebook_instance.foo", "tags.%", "1"),
					resource.TestCheckResourceAttr("aws_sagemaker_notebook_instance.foo", "tags.bar", "baz"),
				),
			},
		},
	})
}

func TestAccAWSSagemakerNotebookInstance_disappears(t *testing.T) {
	var notebook sagemaker.DescribeNotebookInstanceOutput
	rName := acctest.RandomWithPrefix("tf-acc-test")
	var resourceName = "aws_sagemaker_notebook_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSagemakerNotebookInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSagemakerNotebookInstanceBasicConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerNotebookInstanceExists(resourceName, &notebook),
					testAccCheckResourceDisappears(testAccProvider, resourceAwsSagemakerNotebookInstance(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckAWSSagemakerNotebookInstanceDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).sagemakerconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_sagemaker_notebook_instance" {
			continue
		}

		describeNotebookInput := &sagemaker.DescribeNotebookInstanceInput{
			NotebookInstanceName: aws.String(rs.Primary.ID),
		}
		notebookInstance, err := conn.DescribeNotebookInstance(describeNotebookInput)
		if err != nil {
			return nil
		}

		if *notebookInstance.NotebookInstanceName == rs.Primary.ID {
			return fmt.Errorf("sagemaker notebook instance %q still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckAWSSagemakerNotebookInstanceExists(n string, notebook *sagemaker.DescribeNotebookInstanceOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No sagmaker Notebook Instance ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).sagemakerconn
		opts := &sagemaker.DescribeNotebookInstanceInput{
			NotebookInstanceName: aws.String(rs.Primary.ID),
		}
		resp, err := conn.DescribeNotebookInstance(opts)
		if err != nil {
			return err
		}

		*notebook = *resp

		return nil
	}
}

func testAccCheckAWSSagemakerNotebookInstanceName(notebook *sagemaker.DescribeNotebookInstanceOutput, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		notebookName := notebook.NotebookInstanceName
		if *notebookName != expected {
			return fmt.Errorf("Bad Notebook Instance name: %s", *notebook.NotebookInstanceName)
		}

		return nil
	}
}

func TestAccAWSSagemakerNotebookInstance_root_access(t *testing.T) {
	var notebook sagemaker.DescribeNotebookInstanceOutput
	notebookName := resource.PrefixedUniqueId(sagemakerTestAccSagemakerNotebookInstanceResourceNamePrefix)
	var resourceName = "aws_sagemaker_notebook_instance.foo"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSagemakerNotebookInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSagemakerNotebookInstanceConfigRootAccess(notebookName, "Disabled"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerNotebookInstanceExists(resourceName, &notebook),
					testAccCheckAWSSagemakerNotebookRootAccess(&notebook, "Disabled"),
				),
			},
			{
				Config: testAccAWSSagemakerNotebookInstanceConfigRootAccess(notebookName, "Enabled"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerNotebookInstanceExists(resourceName, &notebook),
					testAccCheckAWSSagemakerNotebookRootAccess(&notebook, "Enabled"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSSagemakerNotebookInstance_direct_internet_access(t *testing.T) {
	var notebook sagemaker.DescribeNotebookInstanceOutput
	notebookName := resource.PrefixedUniqueId(sagemakerTestAccSagemakerNotebookInstanceResourceNamePrefix)
	var resourceName = "aws_sagemaker_notebook_instance.foo"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSSagemakerNotebookInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSagemakerNotebookInstanceConfigDirectInternetAccess(notebookName, "Disabled"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerNotebookInstanceExists(resourceName, &notebook),
					testAccCheckAWSSagemakerNotebookDirectInternetAccess(&notebook, "Disabled"),
				),
			},
			{
				Config: testAccAWSSagemakerNotebookInstanceConfigDirectInternetAccess(notebookName, "Enabled"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSagemakerNotebookInstanceExists(resourceName, &notebook),
					testAccCheckAWSSagemakerNotebookDirectInternetAccess(&notebook, "Enabled"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAWSSagemakerNotebookRootAccess(notebook *sagemaker.DescribeNotebookInstanceOutput, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rootAccess := notebook.RootAccess
		if *rootAccess != expected {
			return fmt.Errorf("root_access setting is incorrect: %s", *notebook.RootAccess)
		}

		return nil
	}
}

func testAccCheckAWSSagemakerNotebookDirectInternetAccess(notebook *sagemaker.DescribeNotebookInstanceOutput, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		directInternetAccess := notebook.DirectInternetAccess
		if *directInternetAccess != expected {
			return fmt.Errorf("direct_internet_access setting is incorrect: %s", *notebook.DirectInternetAccess)
		}

		return nil
	}
}

func testAccCheckAWSSagemakerNotebookInstanceTags(notebook *sagemaker.DescribeNotebookInstanceOutput, key string, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*AWSClient).sagemakerconn

		tags, err := keyvaluetags.SagemakerListTags(conn, aws.StringValue(notebook.NotebookInstanceArn))
		if err != nil {
			return err
		}

		m := tags.IgnoreAws().Map()
		v, ok := m[key]
		if value != "" && !ok {
			return fmt.Errorf("Missing tag: %s", key)
		} else if value == "" && ok {
			return fmt.Errorf("Extra tag: %s", key)
		}
		if value == "" {
			return nil
		}

		if v != value {
			return fmt.Errorf("%s: bad value: %s", key, v)
		}

		return nil
	}
}

func testAccAWSSagemakerNotebookInstanceBaseConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_iam_role" "test" {
  name               = %[1]q
  path               = "/"
  assume_role_policy = data.aws_iam_policy_document.test.json
}

data "aws_iam_policy_document" "test" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["sagemaker.amazonaws.com"]
    }
  }
}
`, rName)
}

func testAccAWSSagemakerNotebookInstanceBasicConfig(rName string) string {
	return testAccAWSSagemakerNotebookInstanceBaseConfig(rName) + fmt.Sprintf(`
resource "aws_sagemaker_notebook_instance" "test" {
  name          = %[1]q
  role_arn      = aws_iam_role.test.arn
  instance_type = "ml.t2.medium"
}
`, rName)
}

func testAccAWSSagemakerNotebookInstanceUpdateConfig(rName string) string {
	return testAccAWSSagemakerNotebookInstanceBaseConfig(rName) + fmt.Sprintf(`
resource "aws_sagemaker_notebook_instance" "test" {
  name          = %[1]q
  role_arn      = aws_iam_role.test.arn
  instance_type = "ml.m4.xlarge"
}
`, rName)
}

func testAccAWSSagemakerNotebookInstanceConfigLifecycleConfigName(rName string) string {
	return testAccAWSSagemakerNotebookInstanceBaseConfig(rName) + fmt.Sprintf(`
resource "aws_sagemaker_notebook_instance_lifecycle_configuration" "test" {
  name = %[1]q
}

resource "aws_sagemaker_notebook_instance" "test" {
  instance_type         = "ml.t2.medium"
  lifecycle_config_name = aws_sagemaker_notebook_instance_lifecycle_configuration.test.name
  name                  = %[1]q
  role_arn              = aws_iam_role.test.arn
}
`, rName)
}

func testAccAWSSagemakerNotebookInstanceConfigTags1(rName, tagKey1, tagValue1 string) string {
	return testAccAWSSagemakerNotebookInstanceBaseConfig(rName) + fmt.Sprintf(`
resource "aws_sagemaker_notebook_instance" "test" {
  name          = %[1]q
  role_arn      = aws_iam_role.test.arn
  instance_type = "ml.t2.medium"

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccAWSSagemakerNotebookInstanceConfigTags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return testAccAWSSagemakerNotebookInstanceBaseConfig(rName) + fmt.Sprintf(`
resource "aws_sagemaker_notebook_instance" "test" {
  name          = %[1]q
  role_arn      = aws_iam_role.test.arn
  instance_type = "ml.t2.medium"

  tags = {
    %[2]q = %[3]q
    %[3]q = %[4]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}

func testAccAWSSagemakerNotebookInstanceConfigRootAccess(rName string, rootAccess string) string {
	return testAccAWSSagemakerNotebookInstanceBaseConfig(rName) + fmt.Sprintf(`
resource "aws_sagemaker_notebook_instance" "foo" {
  name          = %[1]q
  role_arn      = aws_iam_role.test.arn
  instance_type = "ml.t2.medium"
  root_access   = %[2]q
}
`, rName, rootAccess)
}

func testAccAWSSagemakerNotebookInstanceConfigDirectInternetAccess(rName string, directInternetAccess string) string {
	return testAccAWSSagemakerNotebookInstanceBaseConfig(rName) + fmt.Sprintf(`
resource "aws_sagemaker_notebook_instance" "foo" {
  name                   = %[1]q
  role_arn               = aws_iam_role.test.arn
  instance_type          = "ml.t2.medium"
  security_groups        = [aws_security_group.test.id]
  subnet_id              = aws_subnet.test.id
  direct_internet_access = %[2]q
}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_security_group" "test" {
  vpc_id = aws_vpc.test.id
}

resource "aws_subnet" "test" {
  vpc_id     = aws_vpc.test.id
  cidr_block = "10.0.0.0/24"

  tags = {
    Name = %[1]q
  }
}
`, rName, directInternetAccess)
}
