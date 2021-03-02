// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccIBMCisDNSRecordsDataSource_basic(t *testing.T) {
	node := "data.ibm_cis_dns_records.test_dns_records"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMCisDNSRecordsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(node, "cis_dns_records.0.record_id"),
					resource.TestCheckResourceAttrSet(node, "cis_dns_records.0.name"),
					resource.TestCheckResourceAttrSet(node, "file"),
					testAccCheckIBMCisDNSRecordsExportExists("/tmp/records.txt"),
					testAccCheckIBMCisDNSRecordsExportedFileRemove("/tmp/records.txt"),
				),
			},
		},
	})
}

func testAccCheckIBMCisDNSRecordsDataSourceConfig() string {
	// status filter defaults to empty
	return testAccCheckIBMCisDNSRecordConfigCisDSBasic("test", cisDomainStatic) +
		fmt.Sprintf(`
	data "ibm_cis_dns_records" "test_dns_records" {
		cis_id    = data.ibm_cis.cis.id
		domain_id = ibm_cis_dns_record.test.domain_id
		file      = "/tmp/records.txt"
	}`)
}

func testAccCheckIBMCisDNSRecordsExportExists(file string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		testStr := fmt.Sprintf("test.%s", cisDomainStatic)

		// Splits on newlines by default.
		scanner := bufio.NewScanner(f)

		line := 1
		// https://golang.org/pkg/bufio/#Scanner.Scan
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), testStr) {
				return nil
			}
			line++
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("domain is not found in exported dns records")
		}
		return nil
	}
}

func testAccCheckIBMCisDNSRecordsExportedFileRemove(file string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		err := os.Remove(file)
		if err != nil {
			return err
		}
		return nil
	}
}
