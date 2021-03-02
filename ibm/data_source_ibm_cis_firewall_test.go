// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccIBMCisFirewallLockdownDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMCisFirewallLockdownDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ibm_cis_firewall.lockdown", "firewall_type", "lockdowns"),
				),
			},
		},
	})
}

func TestAccIBMCisFirewallAccessRuleDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMCisFirewallAccessRuleDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ibm_cis_firewall.access_rule", "firewall_type", "access_rules"),
				),
			},
		},
	})
}

func TestAccIBMCisFirewallUARuleDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMCisFirewallUARuleDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ibm_cis_firewall.ua_rule", "firewall_type", "ua_rules"),
				),
			},
		},
	})
}

func testAccCheckIBMCisFirewallLockdownDataSourceConfigBasic() string {
	return testAccCheckIBMCisFirewallLockdownBasic() + fmt.Sprintf(`
	data "ibm_cis_firewall" "lockdown"{
		cis_id    = data.ibm_cis.cis.id
		domain_id = data.ibm_cis_domain.cis_domain.domain_id
		firewall_type = "lockdowns"
	  }
	`)
}

func testAccCheckIBMCisFirewallAccessRuleDataSourceConfigBasic() string {
	return testAccCheckIBMCisFirewallAccessRuleBasic() + fmt.Sprintf(`
	data "ibm_cis_firewall" "access_rule"{
		cis_id    = data.ibm_cis.cis.id
		domain_id = data.ibm_cis_domain.cis_domain.domain_id
		firewall_type = "access_rules"
	  }
	`)
}

func testAccCheckIBMCisFirewallUARuleDataSourceConfigBasic() string {
	return testAccCheckIBMCisFirewallUARuleBasic() + fmt.Sprintf(`
	data "ibm_cis_firewall" "ua_rule"{
		cis_id    = data.ibm_cis.cis.id
		domain_id = data.ibm_cis_domain.cis_domain.domain_id
		firewall_type = "ua_rules"
	  }
	`)
}
