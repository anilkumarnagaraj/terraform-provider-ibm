// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"math/rand"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccIBMKMSResource_basic(t *testing.T) {
	instanceName := fmt.Sprintf("kms_%d", acctest.RandIntRange(10, 100))
	cosInstanceName := fmt.Sprintf("cos_%d", acctest.RandIntRange(10, 100))
	bucketName := fmt.Sprintf("bucket-test77")
	keyName := fmt.Sprintf("key_%d", acctest.RandIntRange(10, 100))
	payload := "LqMWNtSi3Snr4gFNO0PsFFLFRNs57mSXCQE7O2oE+g0="

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMKmsResourceStandardConfig(instanceName, keyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_kms_key.test", "key_name", keyName),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMKmsResourceImportStandardConfig(instanceName, keyName, payload),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_kms_key.test", "key_name", keyName),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMKmsResourceRootkeyWithCOSConfig(instanceName, keyName, cosInstanceName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_kms_key.test", "key_name", keyName),
				),
			},
		},
	})
}
func TestAccIBMKMSHPCSResource_basic(t *testing.T) {
	t.Skip()
	hpcskeyName := fmt.Sprintf("hpcs_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMKmsResourceHpcsConfig(hpcsInstanceID, hpcskeyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_kms_key.hpcstest", "key_name", hpcskeyName),
				),
			},
		},
	})
}

// Test for valid expiration date for create key operation
func TestAccIBMKMSResource_ValidExpDate(t *testing.T) {

	instanceName := fmt.Sprintf("kms_%d", acctest.RandIntRange(10, 100))
	keyName := fmt.Sprintf("key_%d", acctest.RandIntRange(10, 100))

	hours := time.Duration(rand.Intn(24) + 1)
	mins := time.Duration(rand.Intn(60) + 1)
	sec := time.Duration(rand.Intn(60) + 1)
	loc, _ := time.LoadLocation("UTC")
	expirationDateValid := ((time.Now().In(loc).Add(time.Hour*hours + time.Minute*mins + time.Second*sec)).Format(time.RFC3339))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMKmsCreateStandardKeyConfig(instanceName, keyName, expirationDateValid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_kms_key.test", "key_name", keyName),
					resource.TestCheckResourceAttr("ibm_kms_key.test", "expiration_date", expirationDateValid),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMKmsCreateRootKeyConfig(instanceName, keyName, expirationDateValid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_kms_key.test", "key_name", keyName),
					resource.TestCheckResourceAttr("ibm_kms_key.test", "expiration_date", expirationDateValid),
				),
			},
		},
	})
}

// Test for invalid expiration date for create key operation
func TestAccIBMKMSResource_InvalidExpDate(t *testing.T) {
	instanceName := fmt.Sprintf("kms_%d", acctest.RandIntRange(10, 100))
	keyName := fmt.Sprintf("key_%d", acctest.RandIntRange(10, 100))

	hours := time.Duration(rand.Intn(24) + 1)
	mins := time.Duration(rand.Intn(60) + 1)
	sec := time.Duration(rand.Intn(60) + 1)
	expirationDateInvalid := (time.Now().Add(time.Hour*hours + time.Minute*mins + time.Second*sec)).String()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      testAccCheckIBMKmsCreateStandardKeyConfig(instanceName, keyName, expirationDateInvalid),
				ExpectError: regexp.MustCompile("Invalid time format"),
			},
			resource.TestStep{
				Config:      testAccCheckIBMKmsCreateRootKeyConfig(instanceName, keyName, expirationDateInvalid),
				ExpectError: regexp.MustCompile("Invalid time format"),
			},
		},
	})
}

func TestAccIBMKMSKeyPolicy_basic(t *testing.T) {
	instanceName := fmt.Sprintf("kms_%d", acctest.RandIntRange(10, 100))
	keyName := fmt.Sprintf("key_%d", acctest.RandIntRange(10, 100))
	rotation_interval := 3
	dual_auth_delete := false
	rotation_interval_new := 5
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMKmsKeyPolicyStandardConfig(instanceName, keyName, rotation_interval, dual_auth_delete),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_kms_key.test", "key_name", keyName),
					resource.TestCheckResourceAttr("ibm_kms_key.test", "policies.0.rotation.0.interval_month", "3"),
					resource.TestCheckResourceAttr("ibm_kms_key.test", "policies.0.dual_auth_delete.0.enabled", "false"),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMKmsKeyPolicyStandardConfig(instanceName, keyName, rotation_interval_new, dual_auth_delete),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_kms_key.test", "key_name", keyName),
					resource.TestCheckResourceAttr("ibm_kms_key.test", "policies.0.rotation.0.interval_month", "5"),
					resource.TestCheckResourceAttr("ibm_kms_key.test", "policies.0.dual_auth_delete.0.enabled", "false"),
				),
			},
		},
	})
}

func TestAccIBMKMSKeyPolicy_rotation(t *testing.T) {
	instanceName := fmt.Sprintf("kms_%d", acctest.RandIntRange(10, 100))
	keyName := fmt.Sprintf("key_%d", acctest.RandIntRange(10, 100))
	rotation_interval := 3
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMKmsKeyPolicyRotation(instanceName, keyName, rotation_interval),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_kms_key.test", "key_name", keyName),
					resource.TestCheckResourceAttr("ibm_kms_key.test", "policies.0.rotation.0.interval_month", "3"),
				),
			},
		},
	})
}

func TestAccIBMKMSKeyPolicy_dualAuth(t *testing.T) {
	instanceName := fmt.Sprintf("kms_%d", acctest.RandIntRange(10, 100))
	keyName := fmt.Sprintf("key_%d", acctest.RandIntRange(10, 100))
	dual_auth_delete := false
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMKmsKeyPolicyDualAuth(instanceName, keyName, dual_auth_delete),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_kms_key.test", "key_name", keyName),
					resource.TestCheckResourceAttr("ibm_kms_key.test", "policies.0.dual_auth_delete.0.enabled", "false"),
				),
			},
		},
	})
}

func TestAccIBMKMSKeyPolicy_invalid_interval(t *testing.T) {
	instanceName := fmt.Sprintf("kms_%d", acctest.RandIntRange(10, 100))
	keyName := fmt.Sprintf("key_%d", acctest.RandIntRange(10, 100))
	rotation_interval := 13
	dual_auth_delete := false
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      testAccCheckIBMKmsKeyPolicyStandardConfig(instanceName, keyName, rotation_interval, dual_auth_delete),
				ExpectError: regexp.MustCompile("config is invalid:"),
			},
		},
	})
}

func testAccCheckIBMKmsResourceStandardConfig(instanceName, KeyName string) string {
	return fmt.Sprintf(`
	resource "ibm_resource_instance" "kms_instance" {
		name              = "%s"
		service           = "kms"
		plan              = "tiered-pricing"
		location          = "us-south"
	  }
	  resource "ibm_kms_key" "test" {
		instance_id = "${ibm_resource_instance.kms_instance.guid}"
		key_name = "%s"
		standard_key =  true
		force_delete = true
	}
	
`, instanceName, KeyName)
}

func testAccCheckIBMKmsResourceImportStandardConfig(instanceName, KeyName, payload string) string {
	return fmt.Sprintf(`
	resource "ibm_resource_instance" "kms_instance" {
		name              = "%s"
		service           = "kms"
		plan              = "tiered-pricing"
		location          = "us-south"
	  }
	  resource "ibm_kms_key" "test" {
		instance_id = "${ibm_resource_instance.kms_instance.guid}"
		key_name = "%s"
		standard_key =  true
		payload = "%s"
		force_delete = true
	}

`, instanceName, KeyName, payload)
}

func testAccCheckIBMKmsResourceRootkeyWithCOSConfig(instanceName, KeyName, cosInstanceName, bucketName string) string {
	return fmt.Sprintf(`
	provider "ibm" {
		region = "us-south"
	}	
	resource "ibm_resource_instance" "kms_instance1" {
		name              = "%s"
		service           = "kms"
		plan              = "tiered-pricing"
		location          = "us-south"
	  }
	  resource "ibm_kms_key" "test" {
		instance_id = "${ibm_resource_instance.kms_instance1.guid}"
		key_name = "%s"
		standard_key =  false
		force_delete = true
	}

	resource "ibm_resource_instance" "cos_instance" {
		name     = "%s"
		service  = "cloud-object-storage"
		plan     = "standard"
		location = "global"
	}
	
	resource "ibm_iam_authorization_policy" "policy" {
		source_service_name = "cloud-object-storage"
		target_service_name = "kms"
		roles               = ["Reader"]
	}
	
	resource "ibm_cos_bucket" "smart-us-south" {
		depends_on           = [ibm_iam_authorization_policy.policy]
		bucket_name          = "%s"
		resource_instance_id = ibm_resource_instance.cos_instance.id
		region_location      = "us-south"
		storage_class        = "smart"
		key_protect          = ibm_kms_key.test.id
	}
	
`, instanceName, KeyName, cosInstanceName, bucketName)
}

func testAccCheckIBMKmsResourceHpcsConfig(hpcsInstanceID, KeyName string) string {
	return fmt.Sprintf(`
	  resource "ibm_kms_key" "hpcstest" {
		instance_id = "%s"
		key_name = "%s"
		standard_key =  true
		force_delete = true
	}

`, hpcsInstanceID, KeyName)
}

func testAccCheckIBMKmsCreateStandardKeyConfig(instanceName, KeyName, expirationDate string) string {
	return fmt.Sprintf(`
	resource "ibm_resource_instance" "kms_instance" {
		name              = "%s"
		service           = "kms"
		plan              = "tiered-pricing"
		location          = "us-south"
	  }
	  resource "ibm_kms_key" "test" {
		instance_id = "${ibm_resource_instance.kms_instance.guid}"
		key_name = "%s"
		standard_key =  true
		force_delete = true
		expiration_date = "%s"
	}
	
`, instanceName, KeyName, expirationDate)
}

func testAccCheckIBMKmsCreateRootKeyConfig(instanceName, KeyName, expirationDate string) string {
	return fmt.Sprintf(`
	resource "ibm_resource_instance" "kms_instance" {
		name              = "%s"
		service           = "kms"
		plan              = "tiered-pricing"
		location          = "us-south"
	  }
	  resource "ibm_kms_key" "test" {
		instance_id = "${ibm_resource_instance.kms_instance.guid}"
		key_name = "%s"
		standard_key =  false
		force_delete = true
		expiration_date = "%s"
	}
	
`, instanceName, KeyName, expirationDate)
}

func testAccCheckIBMKmsKeyPolicyStandardConfig(instanceName, KeyName string, rotation_interval int, dual_auth_delete bool) string {
	return fmt.Sprintf(`
	resource "ibm_resource_instance" "kp_instance" {
		name     = "%s"
		service  = "kms"
		plan     = "tiered-pricing"
		location = "us-south"
	  }
	  
	  resource "ibm_kms_key" "test" {
		instance_id = ibm_resource_instance.kp_instance.guid
		key_name       = "%s"
		standard_key   = false
		policies {
		  rotation {
			interval_month = %d
		  }
		  dual_auth_delete {
			enabled = %t
		  }
		}
	  }
`, instanceName, KeyName, rotation_interval, dual_auth_delete)
}

func testAccCheckIBMKmsKeyPolicyRotation(instanceName, KeyName string, rotation_interval int) string {
	return fmt.Sprintf(`
	resource "ibm_resource_instance" "kp_instance" {
		name     = "%s"
		service  = "kms"
		plan     = "tiered-pricing"
		location = "us-south"
	  }
	  
	  resource "ibm_kms_key" "test" {
		instance_id = ibm_resource_instance.kp_instance.guid
		key_name       = "%s"
		standard_key   = false
		policies {
		  rotation {
			interval_month = %d
		  }
		}
	  }
`, instanceName, KeyName, rotation_interval)
}

func testAccCheckIBMKmsKeyPolicyDualAuth(instanceName, KeyName string, dual_auth_delete bool) string {
	return fmt.Sprintf(`
	resource "ibm_resource_instance" "kp_instance" {
		name     = "%s"
		service  = "kms"
		plan     = "tiered-pricing"
		location = "us-south"
	  }
	  
	  resource "ibm_kms_key" "test" {
		instance_id = ibm_resource_instance.kp_instance.guid
		key_name       = "%s"
		standard_key   = false
		policies {
		  dual_auth_delete {
			enabled = %t
		  }
		}
	  }
`, instanceName, KeyName, dual_auth_delete)
}
