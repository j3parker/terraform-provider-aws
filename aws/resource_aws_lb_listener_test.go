package aws

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAWSLBListener_basic(t *testing.T) {
	var conf elbv2.Listener
	resourceName := "aws_lb_listener.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAWSLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSLBListenerConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSLBListenerExists(resourceName, &conf),
					resource.TestCheckResourceAttrPair(resourceName, "load_balancer_arn", "aws_lb.test", "arn"),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile("listener/.+$")),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "port", "80"),
					resource.TestCheckResourceAttr(resourceName, "default_action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.type", "forward"),
					resource.TestCheckResourceAttrPair(resourceName, "default_action.0.target_group_arn", "aws_lb_target_group.test", "arn"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.redirect.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.fixed_response.#", "0"),
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

func TestAccAWSLBListener_forwardWeighted(t *testing.T) {
	var conf elbv2.Listener
	resourceName := "aws_lb_listener.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	targetGroupName1 := acctest.RandomWithPrefix("tf-acc-test")
	targetGroupName2 := acctest.RandomWithPrefix("tf-acc-test2")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAWSLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSLBListenerConfig_forwardWeighted(rName, targetGroupName1, targetGroupName2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSLBListenerExists(resourceName, &conf),
					resource.TestCheckResourceAttrPair(resourceName, "load_balancer_arn", "aws_lb.test", "arn"),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile("listener/.+$")),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "port", "80"),
					resource.TestCheckResourceAttr(resourceName, "default_action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.type", "forward"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.forward.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.forward.0.target_group.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.forward.0.stickiness.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.forward.0.stickiness.0.duration", "0"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.redirect.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.fixed_response.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSLBListenerConfig_changeForwardWeightedStickiness(rName, targetGroupName1, targetGroupName2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSLBListenerExists(resourceName, &conf),
					resource.TestCheckResourceAttrPair(resourceName, "load_balancer_arn", "aws_lb.test", "arn"),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile("listener/.+$")),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "port", "80"),
					resource.TestCheckResourceAttr(resourceName, "default_action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.type", "forward"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.forward.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.forward.0.target_group.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.forward.0.stickiness.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.forward.0.stickiness.0.duration", "3600"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.redirect.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.fixed_response.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSLBListenerConfig_changeForwardWeightedToBasic(rName, targetGroupName1, targetGroupName2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSLBListenerExists(resourceName, &conf),
					resource.TestCheckResourceAttrPair(resourceName, "load_balancer_arn", "aws_lb.test", "arn"),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile("listener/.+$")),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "port", "80"),
					resource.TestCheckResourceAttr(resourceName, "default_action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.type", "forward"),
					resource.TestCheckResourceAttrPair(resourceName, "default_action.0.target_group_arn", "aws_lb_target_group.test1", "arn"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.redirect.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.fixed_response.#", "0"),
				),
			},
		},
	})
}

func TestAccAWSLBListener_basicUdp(t *testing.T) {
	var conf elbv2.Listener
	resourceName := "aws_lb_listener.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAWSLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSLBListenerConfig_basicUdp(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSLBListenerExists(resourceName, &conf),
					resource.TestCheckResourceAttrPair(resourceName, "load_balancer_arn", "aws_lb.test", "arn"),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile("listener/.+$")),
					resource.TestCheckResourceAttr(resourceName, "protocol", "UDP"),
					resource.TestCheckResourceAttr(resourceName, "port", "514"),
					resource.TestCheckResourceAttr(resourceName, "default_action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.type", "forward"),
					resource.TestCheckResourceAttrPair(resourceName, "default_action.0.target_group_arn", "aws_lb_target_group.test", "arn"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.redirect.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.fixed_response.#", "0"),
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

func TestAccAWSLBListener_BackwardsCompatibility(t *testing.T) {
	var conf elbv2.Listener
	resourceName := "aws_alb_listener.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAWSLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSLBListenerConfigBackwardsCompatibility(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSLBListenerExists(resourceName, &conf),
					resource.TestCheckResourceAttrPair(resourceName, "load_balancer_arn", "aws_alb.test", "arn"),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile("listener/.+$")),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "port", "80"),
					resource.TestCheckResourceAttr(resourceName, "default_action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.type", "forward"),
					resource.TestCheckResourceAttrPair(resourceName, "default_action.0.target_group_arn", "aws_alb_target_group.test", "arn"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.redirect.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.fixed_response.#", "0"),
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

func TestAccAWSLBListener_https(t *testing.T) {
	var conf elbv2.Listener
	key := tlsRsaPrivateKeyPem(2048)
	resourceName := "aws_lb_listener.test"
	certificate := tlsRsaX509SelfSignedCertificatePem(key, "example.com")
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAWSLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSLBListenerConfig_https(rName, key, certificate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSLBListenerExists(resourceName, &conf),
					resource.TestCheckResourceAttrPair(resourceName, "load_balancer_arn", "aws_lb.test", "arn"),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile("listener/.+$")),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "port", "443"),
					resource.TestCheckResourceAttr(resourceName, "default_action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.type", "forward"),
					resource.TestCheckResourceAttrPair(resourceName, "default_action.0.target_group_arn", "aws_lb_target_group.test", "arn"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.redirect.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.fixed_response.#", "0"),
					resource.TestCheckResourceAttrPair(resourceName, "certificate_arn", "aws_iam_server_certificate.test", "arn"),
					resource.TestCheckResourceAttr(resourceName, "ssl_policy", "ELBSecurityPolicy-2016-08"),
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

func TestAccAWSLBListener_Protocol_Tls(t *testing.T) {
	var listener1 elbv2.Listener
	key := tlsRsaPrivateKeyPem(2048)
	certificate := tlsRsaX509SelfSignedCertificatePem(key, "example.com")
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_lb_listener.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSLBListenerConfig_Protocol_Tls(rName, key, certificate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSLBListenerExists(resourceName, &listener1),
					resource.TestCheckResourceAttr(resourceName, "protocol", "TLS"),
					resource.TestCheckResourceAttrPair(resourceName, "certificate_arn", "aws_acm_certificate.test", "arn"),
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

func TestAccAWSLBListener_redirect(t *testing.T) {
	var conf elbv2.Listener
	resourceName := "aws_lb_listener.test"
	lbName := fmt.Sprintf("testlistener-redirect-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAWSLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSLBListenerConfig_redirect(lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSLBListenerExists(resourceName, &conf),
					resource.TestCheckResourceAttrPair(resourceName, "load_balancer_arn", "aws_lb.test", "arn"),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile("listener/.+$")),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "port", "80"),
					resource.TestCheckResourceAttr(resourceName, "default_action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.type", "redirect"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.target_group_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.redirect.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.redirect.0.host", "#{host}"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.redirect.0.path", "/#{path}"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.redirect.0.port", "443"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.redirect.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.redirect.0.query", "#{query}"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.redirect.0.status_code", "HTTP_301"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.fixed_response.#", "0"),
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

func TestAccAWSLBListener_fixedResponse(t *testing.T) {
	var conf elbv2.Listener
	resourceName := "aws_lb_listener.test"
	lbName := fmt.Sprintf("testlistener-fixedresponse-%s", acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAWSLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSLBListenerConfig_fixedResponse(lbName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSLBListenerExists(resourceName, &conf),
					resource.TestCheckResourceAttrPair(resourceName, "load_balancer_arn", "aws_lb.test", "arn"),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile("listener/.+$")),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "port", "80"),
					resource.TestCheckResourceAttr(resourceName, "default_action.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.type", "fixed-response"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.target_group_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.redirect.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.fixed_response.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.fixed_response.0.content_type", "text/plain"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.fixed_response.0.message_body", "Fixed response content"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.fixed_response.0.status_code", "200"),
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

func TestAccAWSLBListener_cognito(t *testing.T) {
	var conf elbv2.Listener
	key := tlsRsaPrivateKeyPem(2048)
	resourceName := "aws_lb_listener.test"
	certificate := tlsRsaX509SelfSignedCertificatePem(key, "example.com")
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAWSLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSLBListenerConfig_cognito(rName, key, certificate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSLBListenerExists(resourceName, &conf),
					resource.TestCheckResourceAttrPair(resourceName, "load_balancer_arn", "aws_lb.test", "arn"),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile("listener/.+$")),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "port", "443"),
					resource.TestCheckResourceAttr(resourceName, "default_action.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.type", "authenticate-cognito"),
					resource.TestCheckResourceAttrSet(resourceName, "default_action.0.authenticate_cognito.0.user_pool_arn"),
					resource.TestCheckResourceAttrSet(resourceName, "default_action.0.authenticate_cognito.0.user_pool_client_id"),
					resource.TestCheckResourceAttrSet(resourceName, "default_action.0.authenticate_cognito.0.user_pool_domain"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.authenticate_cognito.0.authentication_request_extra_params.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.authenticate_cognito.0.authentication_request_extra_params.param", "test"),
					resource.TestCheckResourceAttr(resourceName, "default_action.1.type", "forward"),
					resource.TestCheckResourceAttr(resourceName, "default_action.1.order", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "default_action.1.target_group_arn", "aws_lb_target_group.test", "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "certificate_arn", "aws_iam_server_certificate.test", "arn"),
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

func TestAccAWSLBListener_oidc(t *testing.T) {
	var conf elbv2.Listener
	key := tlsRsaPrivateKeyPem(2048)
	resourceName := "aws_lb_listener.test"
	certificate := tlsRsaX509SelfSignedCertificatePem(key, "example.com")
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAWSLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSLBListenerConfig_oidc(rName, key, certificate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSLBListenerExists(resourceName, &conf),
					resource.TestCheckResourceAttrPair(resourceName, "load_balancer_arn", "aws_lb.test", "arn"),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "elasticloadbalancing", regexp.MustCompile("listener/.+$")),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "port", "443"),
					resource.TestCheckResourceAttr(resourceName, "default_action.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.type", "authenticate-oidc"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.authenticate_oidc.0.authorization_endpoint", "https://example.com/authorization_endpoint"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.authenticate_oidc.0.client_id", "s6BhdRkqt3"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.authenticate_oidc.0.client_secret", "7Fjfp0ZBr1KtDRbnfVdmIw"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.authenticate_oidc.0.issuer", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.authenticate_oidc.0.token_endpoint", "https://example.com/token_endpoint"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.authenticate_oidc.0.user_info_endpoint", "https://example.com/user_info_endpoint"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.authenticate_oidc.0.authentication_request_extra_params.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.authenticate_oidc.0.authentication_request_extra_params.param", "test"),
					resource.TestCheckResourceAttr(resourceName, "default_action.1.order", "2"),
					resource.TestCheckResourceAttr(resourceName, "default_action.1.type", "forward"),
					resource.TestCheckResourceAttrPair(resourceName, "default_action.1.target_group_arn", "aws_lb_target_group.test", "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "certificate_arn", "aws_iam_server_certificate.test", "arn"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"default_action.0.authenticate_oidc.0.client_secret"},
			},
		},
	})
}

func TestAccAWSLBListener_DefaultAction_Order(t *testing.T) {
	var listener elbv2.Listener
	key := tlsRsaPrivateKeyPem(2048)
	certificate := tlsRsaX509SelfSignedCertificatePem(key, "example.com")
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_lb_listener.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSLBListenerConfig_DefaultAction_Order(rName, key, certificate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSLBListenerExists(resourceName, &listener),
					resource.TestCheckResourceAttr(resourceName, "default_action.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.1.order", "2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"default_action.0.authenticate_oidc.0.client_secret"},
			},
		},
	})
}

// Reference: https://github.com/terraform-providers/terraform-provider-aws/issues/6171
func TestAccAWSLBListener_DefaultAction_Order_Recreates(t *testing.T) {
	var listener elbv2.Listener
	key := tlsRsaPrivateKeyPem(2048)
	certificate := tlsRsaX509SelfSignedCertificatePem(key, "example.com")
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_lb_listener.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSLBListenerConfig_DefaultAction_Order(rName, key, certificate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSLBListenerExists(resourceName, &listener),
					resource.TestCheckResourceAttr(resourceName, "default_action.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "default_action.0.order", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_action.1.order", "2"),
					testAccCheckAWSLBListenerDefaultActionOrderDisappears(&listener, 1),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckAWSLBListenerDefaultActionOrderDisappears(listener *elbv2.Listener, actionOrderToDelete int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var newDefaultActions []*elbv2.Action

		for i, action := range listener.DefaultActions {
			if int(aws.Int64Value(action.Order)) == actionOrderToDelete {
				newDefaultActions = append(listener.DefaultActions[:i], listener.DefaultActions[i+1:]...)
				break
			}
		}

		if len(newDefaultActions) == 0 {
			return fmt.Errorf("Unable to find default action order %d from default actions: %#v", actionOrderToDelete, listener.DefaultActions)
		}

		conn := testAccProvider.Meta().(*AWSClient).elbv2conn

		input := &elbv2.ModifyListenerInput{
			DefaultActions: newDefaultActions,
			ListenerArn:    listener.ListenerArn,
		}

		_, err := conn.ModifyListener(input)

		return err
	}
}

func testAccCheckAWSLBListenerExists(n string, res *elbv2.Listener) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Listener ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).elbv2conn

		describe, err := conn.DescribeListeners(&elbv2.DescribeListenersInput{
			ListenerArns: []*string{aws.String(rs.Primary.ID)},
		})

		if err != nil {
			return err
		}

		if len(describe.Listeners) != 1 ||
			aws.StringValue(describe.Listeners[0].ListenerArn) != rs.Primary.ID {
			return errors.New("Listener not found")
		}

		*res = *describe.Listeners[0]
		return nil
	}
}

func testAccCheckAWSLBListenerDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).elbv2conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_lb_listener" && rs.Type != "aws_alb_listener" {
			continue
		}

		describe, err := conn.DescribeListeners(&elbv2.DescribeListenersInput{
			ListenerArns: []*string{aws.String(rs.Primary.ID)},
		})

		if err == nil {
			if len(describe.Listeners) != 0 &&
				aws.StringValue(describe.Listeners[0].ListenerArn) == rs.Primary.ID {
				return fmt.Errorf("Listener %q still exists", rs.Primary.ID)
			}
		}

		// Verify the error
		if isAWSErr(err, elbv2.ErrCodeListenerNotFoundException, "") {
			return nil
		} else {
			return fmt.Errorf("Unexpected error checking LB Listener destroyed: %s", err)
		}
	}

	return nil
}

func testAccAWSLBListenerConfigBase(rName string) string {
	return fmt.Sprintf(`
variable "subnets" {
  default = ["10.0.1.0/24", "10.0.2.0/24"]
  type    = "list"
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  count                   = 2
  vpc_id                  = "${aws_vpc.test.id}"
  cidr_block              = "${element(var.subnets, count.index)}"
  map_public_ip_on_launch = true
  availability_zone       = "${element(data.aws_availability_zones.available.names, count.index)}"

  tags = {
    Name = "%[1]s-${count.index}"
  }
}

resource "aws_security_group" "test" {
  name        = %[1]q
  description = "Used for ALB Testing"
  vpc_id      = "${aws_vpc.test.id}"

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccAWSLBListenerConfig_basic(rName string) string {
	return testAccAWSLBListenerConfigBase(rName) + fmt.Sprintf(`
resource "aws_lb_listener" "test" {
  load_balancer_arn = "${aws_lb.test.id}"
  protocol          = "HTTP"
  port              = "80"

  default_action {
    target_group_arn = "${aws_lb_target_group.test.id}"
    type             = "forward"
  }
}

resource "aws_lb" "test" {
  name            = %[1]q
  internal        = true
  security_groups = ["${aws_security_group.test.id}"]
  subnets         = ["${aws_subnet.test.*.id[0]}", "${aws_subnet.test.*.id[1]}"]

  idle_timeout               = 30
  enable_deletion_protection = false

  tags = {
    Name = %[1]q
  }
}

resource "aws_lb_target_group" "test" {
  name     = %[1]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = "${aws_vpc.test.id}"

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccAWSLBListenerConfig_forwardWeighted(rName, targetGroupName1 string, targetGroupName2 string) string {
	return testAccAWSLBListenerConfigBase(rName) + fmt.Sprintf(`
resource "aws_lb_listener" "test" {
  load_balancer_arn = "${aws_lb.test.id}"
  protocol          = "HTTP"
  port              = "80"

  default_action {
    type = "forward"
    forward {
      target_group {
        arn    = "${aws_lb_target_group.test1.arn}"
        weight = 1
      }
      target_group {
        arn    = "${aws_lb_target_group.test2.arn}"
        weight = 1
      }
    }
  }
}

resource "aws_lb" "test" {
  name            = %[1]q
  internal        = true
  security_groups = ["${aws_security_group.test.id}"]
  subnets         = ["${aws_subnet.test.*.id[0]}", "${aws_subnet.test.*.id[1]}"]

  idle_timeout               = 30
  enable_deletion_protection = false

  tags = {
    Name = %[1]q
  }
}

resource "aws_lb_target_group" "test1" {
  name     = %[2]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = "${aws_vpc.test.id}"

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }

  tags = {
    Name = %[2]q
  }
}

resource "aws_lb_target_group" "test2" {
  name     = %[3]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = "${aws_vpc.test.id}"

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }

  tags = {
    Name = %[3]q
  }
}
`, rName, targetGroupName1, targetGroupName2)
}

func testAccAWSLBListenerConfig_changeForwardWeightedStickiness(rName, targetGroupName1 string, targetGroupName2 string) string {
	return testAccAWSLBListenerConfigBase(rName) + fmt.Sprintf(`
resource "aws_lb_listener" "test" {
  load_balancer_arn = "${aws_lb.test.id}"
  protocol          = "HTTP"
  port              = "80"

  default_action {
    type = "forward"
    forward {
      target_group {
        arn    = "${aws_lb_target_group.test1.arn}"
        weight = 1
      }
      target_group {
        arn    = "${aws_lb_target_group.test2.arn}"
        weight = 1
      }
      stickiness {
        enabled  = true
        duration = 3600
      }
    }
  }
}

resource "aws_lb" "test" {
  name            = %[1]q
  internal        = true
  security_groups = ["${aws_security_group.test.id}"]
  subnets         = ["${aws_subnet.test.*.id[0]}", "${aws_subnet.test.*.id[1]}"]

  idle_timeout               = 30
  enable_deletion_protection = false

  tags = {
    Name = %[1]q
  }
}

resource "aws_lb_target_group" "test1" {
  name     = %[2]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = "${aws_vpc.test.id}"

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }

  tags = {
    Name = %[2]q
  }
}

resource "aws_lb_target_group" "test2" {
  name     = %[3]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = "${aws_vpc.test.id}"

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }

  tags = {
    Name = %[3]q
  }
}
`, rName, targetGroupName1, targetGroupName2)
}

func testAccAWSLBListenerConfig_changeForwardWeightedToBasic(rName, targetGroupName1 string, targetGroupName2 string) string {
	return testAccAWSLBListenerConfigBase(rName) + fmt.Sprintf(`
resource "aws_lb_listener" "test" {
  load_balancer_arn = "${aws_lb.test.id}"
  protocol          = "HTTP"
  port              = "80"

  default_action {
    target_group_arn = "${aws_lb_target_group.test1.arn}"
    type             = "forward"
  }
}

resource "aws_lb" "test" {
  name            = %[1]q
  internal        = true
  security_groups = ["${aws_security_group.test.id}"]
  subnets         = ["${aws_subnet.test.*.id[0]}", "${aws_subnet.test.*.id[1]}"]

  idle_timeout               = 30
  enable_deletion_protection = false

  tags = {
    Name = %[1]q
  }
}

resource "aws_lb_target_group" "test1" {
  name     = %[2]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = "${aws_vpc.test.id}"

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }

  tags = {
    Name = %[2]q
  }
}

resource "aws_lb_target_group" "test2" {
  name     = %[3]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = "${aws_vpc.test.id}"

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }

  tags = {
    Name = %[3]q
  }
}
`, rName, targetGroupName1, targetGroupName2)
}

func testAccAWSLBListenerConfig_basicUdp(rName string) string {
	return testAccAWSLBListenerConfigBase(rName) + fmt.Sprintf(`
resource "aws_lb_listener" "test" {
  load_balancer_arn = "${aws_lb.test.id}"
  protocol          = "UDP"
  port              = "514"

  default_action {
    target_group_arn = "${aws_lb_target_group.test.id}"
    type             = "forward"
  }
}

resource "aws_lb" "test" {
  name            = %[1]q
  internal        = false
  load_balancer_type = "network"
  subnets         = ["${aws_subnet.test.*.id[0]}", "${aws_subnet.test.*.id[1]}"]

  idle_timeout               = 30
  enable_deletion_protection = false

  tags = {
    Name = %[1]q
  }
}

resource "aws_lb_target_group" "test" {
  name     = %[1]q
  port     = 514
  protocol = "UDP"
  vpc_id   = "${aws_vpc.test.id}"

  health_check {
    port                = 514
    protocol            = "TCP"
  }

  tags = {
    Name = %[1]q
  }
}

resource "aws_internet_gateway" "test" {
  vpc_id = "${aws_vpc.test.id}"

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccAWSLBListenerConfigBackwardsCompatibility(rName string) string {
	return testAccAWSLBListenerConfigBase(rName) + fmt.Sprintf(`
resource "aws_alb_listener" "test" {
  load_balancer_arn = "${aws_alb.test.id}"
  protocol          = "HTTP"
  port              = "80"

  default_action {
    target_group_arn = "${aws_alb_target_group.test.id}"
    type             = "forward"
  }
}

resource "aws_alb" "test" {
  name            = %[1]q
  internal        = true
  security_groups = ["${aws_security_group.test.id}"]
  subnets         = ["${aws_subnet.test.*.id[0]}", "${aws_subnet.test.*.id[1]}"]

  idle_timeout               = 30
  enable_deletion_protection = false

  tags = {
    Name = %[1]q
  }
}

resource "aws_alb_target_group" "test" {
  name     = %[1]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = "${aws_vpc.test.id}"

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccAWSLBListenerConfig_https(rName, key, certificate string) string {
	return testAccAWSLBListenerConfigBase(rName) + fmt.Sprintf(`
resource "aws_lb_listener" "test" {
  load_balancer_arn = "${aws_lb.test.id}"
  protocol          = "HTTPS"
  port              = "443"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = "${aws_iam_server_certificate.test.arn}"

  default_action {
    target_group_arn = "${aws_lb_target_group.test.id}"
    type             = "forward"
  }
}

resource "aws_lb" "test" {
  name            = %[1]q
  internal        = false
  security_groups = ["${aws_security_group.test.id}"]
  subnets         = ["${aws_subnet.test.*.id[0]}", "${aws_subnet.test.*.id[1]}"]

  idle_timeout               = 30
  enable_deletion_protection = false

  tags = {
    Name = %[1]q
  }
}

resource "aws_lb_target_group" "test" {
  name     = %[1]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = "${aws_vpc.test.id}"

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }

  tags = {
    Name = %[1]q
  }
}

resource "aws_iam_server_certificate" "test" {
  name             = %[1]q
  certificate_body = "%[2]s"
  private_key      = "%[3]s"
}

resource "aws_internet_gateway" "test" {
  vpc_id = "${aws_vpc.test.id}"

  tags = {
    Name = %[1]q
  }
}
`, rName, tlsPemEscapeNewlines(certificate), tlsPemEscapeNewlines(key))
}

func testAccAWSLBListenerConfig_Protocol_Tls(rName, key, certificate string) string {
	return testAccAWSLBListenerConfigBase(rName) + fmt.Sprintf(`
resource "aws_acm_certificate" "test" {
  certificate_body = "%[2]s"
  private_key      = "%[3]s"

  tags = {
    Name = %[1]q
  }
}

resource "aws_lb" "test" {
  internal           = true
  load_balancer_type = "network"
  name               = %[1]q
  subnets            = ["${aws_subnet.test.*.id[0]}", "${aws_subnet.test.*.id[1]}"]

  tags = {
    Name = %[1]q
  }
}

resource "aws_lb_target_group" "test" {
  name     = %[1]q
  port     = 443
  protocol = "TCP"
  vpc_id   = "${aws_vpc.test.id}"

  health_check {
    interval            = 10
    port                = "traffic-port"
    protocol            = "TCP"
    healthy_threshold   = 3
    unhealthy_threshold = 3
  }

  tags = {
    Name = %[1]q
  }
}

resource "aws_lb_listener" "test" {
  certificate_arn   = "${aws_acm_certificate.test.arn}"
  load_balancer_arn = "${aws_lb.test.arn}"
  port              = "443"
  protocol          = "TLS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"

  default_action {
    target_group_arn = "${aws_lb_target_group.test.arn}"
    type             = "forward"
  }
}
`, rName, tlsPemEscapeNewlines(certificate), tlsPemEscapeNewlines(key))
}

func testAccAWSLBListenerConfig_redirect(rName string) string {
	return testAccAWSLBListenerConfigBase(rName) + fmt.Sprintf(`
resource "aws_lb_listener" "test" {
  load_balancer_arn = "${aws_lb.test.id}"
  protocol          = "HTTP"
  port              = "80"

  default_action {
    type = "redirect"

    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

resource "aws_lb" "test" {
  name            = %[1]q
  internal        = true
  security_groups = ["${aws_security_group.test.id}"]
  subnets         = ["${aws_subnet.test.*.id[0]}", "${aws_subnet.test.*.id[1]}"]

  idle_timeout               = 30
  enable_deletion_protection = false

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccAWSLBListenerConfig_fixedResponse(rName string) string {
	return testAccAWSLBListenerConfigBase(rName) + fmt.Sprintf(`
resource "aws_lb_listener" "test" {
  load_balancer_arn = "${aws_lb.test.id}"
  protocol          = "HTTP"
  port              = "80"

  default_action {
    type = "fixed-response"

    fixed_response {
      content_type = "text/plain"
      message_body = "Fixed response content"
      status_code  = "200"
    }
  }
}

resource "aws_lb" "test" {
  name            = %[1]q
  internal        = true
  security_groups = ["${aws_security_group.test.id}"]
  subnets         = ["${aws_subnet.test.*.id[0]}", "${aws_subnet.test.*.id[1]}"]

  idle_timeout               = 30
  enable_deletion_protection = false

  tags = {
    Name = %[1]q
  }
}
`, rName)
}

func testAccAWSLBListenerConfig_cognito(rName, key, certificate string) string {
	return testAccAWSLBListenerConfigBase(rName) + fmt.Sprintf(`
resource "aws_lb" "test" {
  name                       = %[1]q
  internal                   = false
  security_groups            = ["${aws_security_group.test.id}"]
  subnets                    = ["${aws_subnet.test.*.id[0]}", "${aws_subnet.test.*.id[1]}"]
  enable_deletion_protection = false

  tags = {
    Name = %[1]q
  }
}

resource "aws_lb_target_group" "test" {
  name     = %[1]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = "${aws_vpc.test.id}"

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }

  tags = {
    Name = %[1]q
  }
}

resource "aws_internet_gateway" "test" {
  vpc_id = "${aws_vpc.test.id}"

  tags = {
    Name = %[1]q
  }
}

resource "aws_cognito_user_pool" "test" {
  name = %[1]q

  tags = {
    Name = %[1]q
  }
}

resource "aws_cognito_user_pool_client" "test" {
  name                                 = %[1]q
  user_pool_id                         = "${aws_cognito_user_pool.test.id}"
  generate_secret                      = true
  allowed_oauth_flows_user_pool_client = true
  allowed_oauth_flows                  = ["code", "implicit"]
  allowed_oauth_scopes                 = ["phone", "email", "openid", "profile", "aws.cognito.signin.user.admin"]
  callback_urls                        = ["https://www.example.com/callback", "https://www.example.com/redirect"]
  default_redirect_uri                 = "https://www.example.com/redirect"
  logout_urls                          = ["https://www.example.com/login"]
}

resource "aws_cognito_user_pool_domain" "test" {
  domain       = %[1]q
  user_pool_id = "${aws_cognito_user_pool.test.id}"
}

resource "aws_iam_server_certificate" "test" {
  name             = %[1]q
  certificate_body = "%[2]s"
  private_key      = "%[3]s"
}

resource "aws_lb_listener" "test" {
  load_balancer_arn = "${aws_lb.test.id}"
  protocol          = "HTTPS"
  port              = "443"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = "${aws_iam_server_certificate.test.arn}"

  default_action {
    type = "authenticate-cognito"

    authenticate_cognito {
      user_pool_arn       = "${aws_cognito_user_pool.test.arn}"
      user_pool_client_id = "${aws_cognito_user_pool_client.test.id}"
      user_pool_domain    = "${aws_cognito_user_pool_domain.test.domain}"

      authentication_request_extra_params = {
        param = "test"
      }
    }
  }

  default_action {
    target_group_arn = "${aws_lb_target_group.test.id}"
    type             = "forward"
  }
}
`, rName, tlsPemEscapeNewlines(certificate), tlsPemEscapeNewlines(key))
}

func testAccAWSLBListenerConfig_oidc(rName, key, certificate string) string {
	return testAccAWSLBListenerConfigBase(rName) + fmt.Sprintf(`
resource "aws_lb" "test" {
  name                       = %[1]q
  internal                   = false
  security_groups            = ["${aws_security_group.test.id}"]
  subnets                    = ["${aws_subnet.test.*.id[0]}", "${aws_subnet.test.*.id[1]}"]
  enable_deletion_protection = false

  tags = {
    Name = %[1]q
  }
}

resource "aws_lb_target_group" "test" {
  name     = %[1]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = "${aws_vpc.test.id}"

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }

  tags = {
    Name = %[1]q
  }
}

resource "aws_internet_gateway" "test" {
  vpc_id = "${aws_vpc.test.id}"

  tags = {
    Name = %[1]q
  }
}

resource "aws_iam_server_certificate" "test" {
  name             = %[1]q
  certificate_body = "%[2]s"
  private_key      = "%[3]s"
}

resource "aws_lb_listener" "test" {
  load_balancer_arn = "${aws_lb.test.id}"
  protocol          = "HTTPS"
  port              = "443"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = "${aws_iam_server_certificate.test.arn}"

  default_action {
    type = "authenticate-oidc"

    authenticate_oidc {
      authorization_endpoint = "https://example.com/authorization_endpoint"
      client_id              = "s6BhdRkqt3"
      client_secret          = "7Fjfp0ZBr1KtDRbnfVdmIw"
      issuer                 = "https://example.com"
      token_endpoint         = "https://example.com/token_endpoint"
      user_info_endpoint     = "https://example.com/user_info_endpoint"

      authentication_request_extra_params = {
        param = "test"
      }
    }
  }

  default_action {
    target_group_arn = "${aws_lb_target_group.test.id}"
    type             = "forward"
  }
}
`, rName, tlsPemEscapeNewlines(certificate), tlsPemEscapeNewlines(key))
}

func testAccAWSLBListenerConfig_DefaultAction_Order(rName, key, certificate string) string {
	return testAccAWSLBListenerConfigBase(rName) + fmt.Sprintf(`
resource "aws_lb_listener" "test" {
  load_balancer_arn = "${aws_lb.test.id}"
  protocol          = "HTTPS"
  port              = "443"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = "${aws_iam_server_certificate.test.arn}"

  default_action {
    order = 1
    type  = "authenticate-oidc"

    authenticate_oidc {
      authorization_endpoint = "https://example.com/authorization_endpoint"
      client_id              = "s6BhdRkqt3"
      client_secret          = "7Fjfp0ZBr1KtDRbnfVdmIw"
      issuer                 = "https://example.com"
      token_endpoint         = "https://example.com/token_endpoint"
      user_info_endpoint     = "https://example.com/user_info_endpoint"

      authentication_request_extra_params = {
        param = "test"
      }
    }
  }

  default_action {
    order            = 2
    type             = "forward"
    target_group_arn = "${aws_lb_target_group.test.arn}"
  }
}

resource "aws_iam_server_certificate" "test" {
  certificate_body = "%[2]s"
  name             = %[1]q
  private_key      = "%[3]s"
}

resource "aws_lb" "test" {
  internal        = true
  name            = %[1]q
  security_groups = ["${aws_security_group.test.id}"]
  subnets         = ["${aws_subnet.test.*.id[0]}", "${aws_subnet.test.*.id[1]}"]

  tags = {
    Name = %[1]q
  }
}

resource "aws_lb_target_group" "test" {
  name     = %[1]q
  port     = 8080
  protocol = "HTTP"
  vpc_id   = "${aws_vpc.test.id}"

  health_check {
    path                = "/health"
    interval            = 60
    port                = 8081
    protocol            = "HTTP"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-299"
  }

  tags = {
    Name = %[1]q
  }
}
`, rName, tlsPemEscapeNewlines(certificate), tlsPemEscapeNewlines(key))
}
