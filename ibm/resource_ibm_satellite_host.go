package ibm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.ibm.com/ibmcloud/kubernetesservice-go-sdk/kubernetesserviceapiv1"
)

const (
	hostCluster    = "cluster"
	hostLocName    = "location_name"
	hostName       = "host_id"
	hostLabels     = "labels"
	hostZone       = "zone"
	hostWorkerPool = "workerpool"
)

func resourceIBMSatelliteHost() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMSatelliteHostCreate,
		Read:     resourceIBMSatelliteHostRead,
		Update:   resourceIBMSatelliteHostUpdate,
		Delete:   resourceIBMSatelliteHostDelete,
		Exists:   resourceIBMSatelliteHostExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			hostCluster: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A unique name for the new Satellite location",
				//ValidateFunc: InvokeValidator("ibm_satellite_location", satLocName),
			},
			hostLocName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The IBM Cloud metro from which the Satellite location is managed",
				//ValidateFunc: InvokeValidator("ibm_satellite_location", sateLocZone),
			},
			hostName: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A description of the new Satellite location",
			},
			hostLabels: {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of labels for the host",
			},
			hostZone: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The account ID for IBM Log Analysis with LogDNA log forwarding",
			},
			hostWorkerPool: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The account ID for IBM Log Analysis with LogDNA log forwarding",
			},
		},
	}
}

func resourceIBMFuncSatHostValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 satLocName,
			ValidateFunctionIdentifier: ValidateRegexp,
			Type:                       TypeString,
			Regexp:                     `^[^/*][a-zA-Z0-9/_@.-]`,
			Required:                   true})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 sateLocZone,
			ValidateFunctionIdentifier: ValidateNoZeroValues,
			Type:                       TypeString,
			Required:                   true})

	ibmFuncSatLocResourceValidator := ResourceValidator{ResourceName: "ibm_satellite_location", Schema: validateSchema}
	return &ibmFuncSatLocResourceValidator
}

func resourceIBMSatelliteHostCreate(d *schema.ResourceData, meta interface{}) error {
	satClient, err := meta.(ClientSession).SatelliteServiceV1()
	if err != nil {
		return err
	}

	labels := make(map[string]string)
	hostAssignOptions := &kubernetesserviceapiv1.CreateSatelliteAssignmentOptions{}
	hostAssignOptions.Cluster = ptrToString(d.Get(hostCluster).(string))
	hostAssignOptions.Controller = ptrToString(d.Get(hostLocName).(string))

	if _, ok := d.GetOk(hostLabels); ok {
		l := d.Get(hostLabels).(*schema.Set)
		labels = flattenHostLabels(l.List())
		hostAssignOptions.Labels = labels
	}

	if _, ok := d.GetOk(hostWorkerPool); ok {
		hostAssignOptions.Workerpool = ptrToString(d.Get(hostWorkerPool).(string))
	}

	if _, ok := d.GetOk(hostZone); ok {
		hostAssignOptions.Zone = ptrToString(d.Get(hostZone).(string))
	}

	host, response, err := satClient.CreateSatelliteAssignment(hostAssignOptions)
	if err != nil {
		return fmt.Errorf("Error Creating Satellite Location: %s\n%s", err, response)
	}

	log.Printf("[INFO] Assigned satellite Host : %s", *host.HostID)
	d.SetId(fmt.Sprintf("%s:%s", *hostAssignOptions.Controller, *host.HostID))

	return resourceIBMSatelliteHostRead(d, meta)
}

func resourceIBMSatelliteHostRead(d *schema.ResourceData, meta interface{}) error {
	parts, err := cfIdParts(d.Id())
	if err != nil {
		return err
	}
	locationName := parts[0]
	ID := parts[1]

	satClient, err := meta.(ClientSession).SatelliteServiceV1()
	if err != nil {
		return err
	}

	hostOptions := &kubernetesserviceapiv1.GetSatelliteHostsOptions{}
	hostList, _, err := satClient.GetSatelliteHosts(hostOptions)
	if err != nil {
		return err
	}

	for _, h := range hostList {
		if ID == *h.ID {
			d.SetId(fmt.Sprintf("%s:%s", locationName, ID))
			d.Set(hostCluster, *h.Assignment.ClusterName)
			d.Set(hostName, *h.Name)
			d.Set(hostLocName, locationName)

			labelRefs := h.Labels
			labelRefsLen := len(labelRefs)
			if labelRefsLen > 0 {
				tags := make([]string, labelRefsLen)
				for k, v := range labelRefs {
					tags = append(tags, fmt.Sprintf("%s=%s", k, v))
				}
				d.Set(hostLabels, tags)
			}
			d.Set(hostWorkerPool, *h.Assignment.WorkerPoolName)
			d.Set(hostZone, *h.Assignment.Zone)
		}
	}

	return nil
}

func resourceIBMSatelliteHostUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceIBMSatelliteHostDelete(d *schema.ResourceData, meta interface{}) error {
	satClient, err := meta.(ClientSession).SatelliteServiceV1()
	if err != nil {
		return err
	}

	removeSatLocOptions := &kubernetesserviceapiv1.RemoveSatelliteLocationOptions{}
	name := d.Get(satLocName).(string)
	removeSatLocOptions.Controller = &name

	response, err := satClient.RemoveSatelliteLocation(removeSatLocOptions)
	if err != nil && response.StatusCode != 404 {
		return fmt.Errorf("Error Deleting Satelliet Location: %s\n%s", err, response)
	}

	d.SetId("")
	return nil
}

func resourceIBMSatelliteHostExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	parts, err := cfIdParts(d.Id())
	if err != nil {
		return false, err
	}

	//locationName := parts[0]
	ID := parts[1]
	satClient, err := meta.(ClientSession).SatelliteServiceV1()
	if err != nil {
		return false, err
	}

	hostOptions := &kubernetesserviceapiv1.GetSatelliteHostsOptions{}
	hostOptions.Controller = &ID
	hostList, response, err := satClient.GetSatelliteHosts(hostOptions)
	if err != nil {
		d.SetId("")
		return false, fmt.Errorf("Error Getting Sateliite Location: %s\n%s", err, response)
	}

	for _, h := range hostList {
		if ID == *h.ID {
			return true, nil
		}
	}
	return false, nil
}
