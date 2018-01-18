package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/directconnect"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAwsDxGateway_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsDxGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDxGatewayConfig(acctest.RandString(5)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsDxGatewayExists("aws_dx_gateway.test"),
				),
			},
		},
	})
}

func testAccCheckAwsDxGatewayDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).dxconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_dx_gateway" {
			continue
		}

		input := &directconnect.DescribeDirectConnectGatewaysInput{
			DirectConnectGatewayId: aws.String(rs.Primary.ID),
		}

		resp, err := conn.DescribeDirectConnectGateways(input)
		if err != nil {
			return err
		}
		for _, v := range resp.DirectConnectGateways {
			if *v.DirectConnectGatewayId == rs.Primary.ID && !(*v.DirectConnectGatewayState == directconnect.GatewayStateDeleted) {
				return fmt.Errorf("[DESTROY ERROR] DX Gateway (%s) not deleted", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckAwsDxGatewayExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		return nil
	}
}

func testAccDxGatewayConfig(rName string) string {
	return fmt.Sprintf(`
    resource "aws_dx_gateway" "test" {
      name = "tf-dxg-%s"
      amazon_side_asn = "64512"
    }
    `, rName)
}
