// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccIBMContainerVPCClusterALBDataSource_basic(t *testing.T) {
	enable := true
	name := fmt.Sprintf("tf-vpc-alb-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMContainerVPCClusterALBDataSourceConfig(enable, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_container_vpc_cluster_alb.testacc_ds_alb", "id"),
				),
			},
		},
	})
}

func testAccCheckIBMContainerVPCClusterALBDataSourceConfig(enable bool, name string) string {
	return testAccCheckIBMVpcContainerALBBasic(enable, name) + fmt.Sprintf(`
	data "ibm_container_vpc_cluster_alb" "testacc_ds_alb" {
	    alb_id = ibm_container_vpc_cluster.cluster.albs.0.id
	}
`)
}
