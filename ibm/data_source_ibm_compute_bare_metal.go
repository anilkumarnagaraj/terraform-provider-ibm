package ibm

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/sl"
)

func dataSourceIBMComputeBareMetal() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMComputeBareMetalRead,

		Schema: map[string]*schema.Schema{

			"hostname": &schema.Schema{
				Description: "The hostname of the bare metal server",
				Type:        schema.TypeString,
				Required:    true,
			},

			"domain": &schema.Schema{
				Description: "The domain of the bare metal server",
				Type:        schema.TypeString,
				Required:    true,
			},

			"datacenter": &schema.Schema{
				Description: "Datacenter in which the bare metal is deployed",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"network_speed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The connection speed, expressed in Mbps,  for the server network components.",
			},

			"public_bandwidth": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The amount of public network traffic, allowed per month.",
			},

			"public_ipv4_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The public IPv4 address of the bare metal server.",
			},

			"private_ipv4_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The private IPv4 address of the bare metal server.",
			},

			"public_vlan_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The public VLAN used for the public network interface of the server.",
			},

			"public_subnet": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The public subnet used for the public network interface of the server.",
			},

			"private_vlan_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The private VLAN used for the private network interface of the server.",
			},

			"private_subnet": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The private subnet used for the private network interface of the server.",
			},

			"hourly_billing": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The billing type of the server.",
			},

			"private_network_only": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Specifies whether the server only has access to the private network.",
			},

			"user_metadata": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Arbitrary data available to the computing server.",
			},

			"notes": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Notes associated with the server.",
			},

			"memory": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The amount of memory in gigabytes, for the server.",
			},

			"redundant_power_supply": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "When the value is `true`, it indicates additional power supply is provided.",
			},

			"redundant_network": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "When the value is `true`, two physical network interfaces are provided with a bonding configuration.",
			},

			"unbonded_network": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "When the value is `true`, two physical network interfaces are provided without a bonding configuration.",
			},

			"os_reference_code": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Tags associated with this bare metal server.",
			},

			"block_storage_ids": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Block storage to which this computing server have access.",
			},

			"file_storage_ids": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "File storage to which this computing server have access.",
			},

			"ipv6_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the public IPv6 address enabled or not",
			},

			"ipv6_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The public IPv6 address of the bare metal server ",
			},

			"secondary_ip_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of secondary IPv4 addresses of the bare metal server.",
			},

			"secondary_ip_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: " The public secondary IPv4 addresses of the bare metal server.",
			},

			"most_recent": &schema.Schema{
				Description: "If true and multiple entries are found, the most recently created bare metal is used. " +
					"If false, an error is returned",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceIBMComputeBareMetalRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetAccountService(sess)

	hostname := d.Get("hostname").(string)
	domain := d.Get("domain").(string)
	mostRecent := d.Get("most_recent").(bool)

	bms, err := service.
		Filter(filter.Build(filter.Path("hardware.hostname").Eq(hostname),
			filter.Path("hardware.domain").Eq(domain))).Mask(
		"hostname,domain,bandwidthAllocation,provisionDate,id," +
			"primaryIpAddress,primaryBackendIpAddress,privateNetworkOnlyFlag," +
			"notes,userData[value],tagReferences[id,tag[name]]," +
			"allowedNetworkStorage[id,nasType]," +
			"hourlyBillingFlag," +
			"datacenter[id,name,longName]," +
			"primaryNetworkComponent[primarySubnet[networkVlan[id,primaryRouter,vlanNumber],id],maxSpeed," +
			"primaryVersion6IpAddressRecord[subnet,guestNetworkComponentBinding[ipAddressId]]]," +
			"primaryBackendNetworkComponent[primarySubnet[networkVlan[id,primaryRouter,vlanNumber],id],maxSpeed,redundancyEnabledFlag]," +
			"memoryCapacity,powerSupplyCount," +
			"operatingSystem[softwareLicense[softwareDescription[referenceCode]]]",
	).GetHardware()

	if err != nil {
		return fmt.Errorf("Error retrieving virtual guest details for host %s: %s", hostname, err)
	}
	if len(bms) == 0 {
		return fmt.Errorf("No virtual guest with hostname %s and domain  %s", hostname, domain)
	}
	var bm datatypes.Hardware

	if len(bms) > 1 {
		if mostRecent {
			bm = mostRecentBareMetal(bms)
		} else {
			return fmt.Errorf(
				"More than one bare metals found with host matching [%s] and domain "+
					"matching [%s]. Set 'most_recent' to true in your configuration to force the most recent bare metal "+
					"to be used", hostname, domain)
		}
	} else {
		bm = bms[0]
	}

	d.SetId(fmt.Sprintf("%d", *bm.Id))
	d.Set("hostname", bm.Hostname)
	d.Set("domain", bm.Domain)

	if bm.Datacenter != nil {
		d.Set("datacenter", bm.Datacenter.Name)
	}

	d.Set("network_speed", bm.PrimaryNetworkComponent.MaxSpeed)
	d.Set("public_bandwidth", bm.BandwidthAllocation)
	if bm.PrimaryIpAddress != nil {
		d.Set("public_ipv4_address", bm.PrimaryIpAddress)
	}
	d.Set("private_ipv4_address", bm.PrimaryBackendIpAddress)

	d.Set("private_network_only", bm.PrivateNetworkOnlyFlag)
	d.Set("hourly_billing", bm.HourlyBillingFlag)

	if bm.PrimaryNetworkComponent.PrimarySubnet != nil {
		d.Set("public_vlan_id", bm.PrimaryNetworkComponent.PrimarySubnet.NetworkVlan.Id)
		d.Set("public_subnet", bm.PrimaryNetworkComponent.PrimarySubnet.Id)
	}

	if bm.PrimaryBackendNetworkComponent.PrimarySubnet != nil {
		d.Set("private_vlan_id", bm.PrimaryBackendNetworkComponent.PrimarySubnet.NetworkVlan.Id)
		d.Set("private_subnet", bm.PrimaryBackendNetworkComponent.PrimarySubnet.Id)
	}

	userData := bm.UserData
	if len(userData) > 0 && userData[0].Value != nil {
		d.Set("user_metadata", userData[0].Value)
	}

	d.Set("notes", sl.Get(bm.Notes, nil))
	d.Set("memory", bm.MemoryCapacity)

	d.Set("redundant_power_supply", false)

	if *bm.PowerSupplyCount == 2 {
		d.Set("redundant_power_supply", true)
	}

	d.Set("redundant_network", false)
	d.Set("unbonded_network", false)

	bareMetalService := services.GetHardwareService(meta.(ClientSession).SoftLayerSession())
	backendNetworkComponent, err := bareMetalService.Filter(
		filter.Build(
			filter.Path("backendNetworkComponents.status").Eq("ACTIVE"),
		),
	).Id(*bm.Id).GetBackendNetworkComponents()

	if err != nil {
		return fmt.Errorf("Error retrieving bare metal server network: %s", err)
	}

	if len(backendNetworkComponent) > 2 && bm.PrimaryBackendNetworkComponent != nil {
		if *bm.PrimaryBackendNetworkComponent.RedundancyEnabledFlag {
			d.Set("redundant_network", true)
		} else {
			d.Set("unbonded_network", true)
		}
	}

	if bm.OperatingSystem != nil &&
		bm.OperatingSystem.SoftwareLicense != nil &&
		bm.OperatingSystem.SoftwareLicense.SoftwareDescription != nil &&
		bm.OperatingSystem.SoftwareLicense.SoftwareDescription.ReferenceCode != nil {
		d.Set("os_reference_code", bm.OperatingSystem.SoftwareLicense.SoftwareDescription.ReferenceCode)
	}

	tagReferences := bm.TagReferences
	tagReferencesLen := len(tagReferences)
	if tagReferencesLen > 0 {
		tags := make([]string, 0, tagReferencesLen)
		for _, tagRef := range tagReferences {
			tags = append(tags, *tagRef.Tag.Name)
		}
		d.Set("tags", tags)
	}

	storages := bm.AllowedNetworkStorage
	if len(storages) > 0 {
		d.Set("block_storage_ids", flattenBlockStorageID(storages))
		d.Set("file_storage_ids", flattenFileStorageID(storages))
	}

	connInfo := map[string]string{"type": "ssh"}
	if !*bm.PrivateNetworkOnlyFlag && bm.PrimaryIpAddress != nil {
		connInfo["host"] = *bm.PrimaryIpAddress
	} else {
		connInfo["host"] = *bm.PrimaryBackendIpAddress
	}
	d.SetConnInfo(connInfo)

	d.Set("ipv6_enabled", false)
	if bm.PrimaryNetworkComponent.PrimaryVersion6IpAddressRecord != nil {
		d.Set("ipv6_enabled", true)
		d.Set("ipv6_address", bm.PrimaryNetworkComponent.PrimaryVersion6IpAddressRecord.IpAddress)
	}
	err = readSecondaryIPAddresses(d, meta, bm.PrimaryIpAddress)
	if err != nil {
		return err
	}

	return nil
}

type bareMetal []datatypes.Hardware

func (k bareMetal) Len() int { return len(k) }

func (k bareMetal) Swap(i, j int) { k[i], k[j] = k[j], k[i] }

func (k bareMetal) Less(i, j int) bool {
	return k[i].ProvisionDate.Before(k[j].ProvisionDate.Time)
}

func mostRecentBareMetal(keys bareMetal) datatypes.Hardware {
	sortedKeys := keys
	sort.Sort(sortedKeys)
	return sortedKeys[len(sortedKeys)-1]
}