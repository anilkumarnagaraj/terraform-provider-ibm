package ibm

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.ibm.com/ibmcloud/kubernetesservice-go-sdk/kubernetesserviceapiv1"
)

const (
	satLocName  = "name"
	sateLocZone = "location"
)

func resourceIBMSatelliteLocation() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMSatelliteLocationCreate,
		Read:     resourceIBMSatelliteLocationRead,
		Delete:   resourceIBMSatelliteLocationDelete,
		Exists:   resourceIBMSatelliteLocationExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			satLocName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A unique name for the new Satellite location",
				//ValidateFunc: InvokeValidator("ibm_satellite_location", satLocName),
			},
			sateLocZone: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The IBM Cloud metro from which the Satellite location is managed",
				//ValidateFunc: InvokeValidator("ibm_satellite_location", sateLocZone),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A description of the new Satellite location",
			},
			"logging_account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The account ID for IBM Log Analysis with LogDNA log forwarding",
			},
			"cos_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Computed:    true,
				Description: "COSBucket - IBM Cloud Object Storage bucket configuration details",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"endpoint": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"cos_credentials": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Computed:    true,
				Description: "COSAuthorization - IBM Cloud Object Storage authorization keys",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_key_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The HMAC secret access key ID",
						},
						"secret_access_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The HMAC secret access key",
						},
					},
				},
			},
			"script_labels": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of labels for the attach host",
			},
		},
	}
}

func resourceIBMFuncSatLocValidator() *ResourceValidator {

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

func resourceIBMSatelliteLocationCreate(d *schema.ResourceData, meta interface{}) error {
	satClient, err := meta.(ClientSession).SatelliteServiceV1()
	if err != nil {
		return err
	}

	labels := make(map[string]string)
	createSatLocOptions := &kubernetesserviceapiv1.CreateSatelliteLocationOptions{}
	name := d.Get(satLocName).(string)
	createSatLocOptions.Name = &name
	satLocation := d.Get(sateLocZone).(string)
	createSatLocOptions.Location = &satLocation

	if v, ok := d.GetOk("cos_config"); ok {
		createSatLocOptions.CosConfig = expandCosConfig(v.([]interface{}))
	}

	if v, ok := d.GetOk("cos_credentials"); ok {
		createSatLocOptions.CosCredentials = expandCosCredentials(v.([]interface{}))
	}

	if _, ok := d.GetOk("logging_account_id"); ok {
		logAccID := d.Get("logging_account_id").(string)
		createSatLocOptions.LoggingAccountID = &logAccID
	}

	if _, ok := d.GetOk("description"); ok {
		desc := d.Get("description").(string)
		createSatLocOptions.Description = &desc
	}

	if _, ok := d.GetOk("script_labels"); ok {
		l := d.Get("script_labels").(*schema.Set)
		labels = flattenHostLabels(l.List())
	}

	//createSatLocOptions.Headers = authMap
	location, response, err := satClient.CreateSatelliteLocation(createSatLocOptions)
	if err != nil {
		return fmt.Errorf("Error Creating Satellite Location: %s\n%s", err, response)
	}

	log.Printf("[INFO] Created satellite location : %s", *location.ID)
	time.Sleep(30 * time.Second)

	//Generate script
	createRegOptions := &kubernetesserviceapiv1.AttachSatelliteHostOptions{}
	createRegOptions.Controller = location.ID
	createRegOptions.Labels = labels

	resp, err := satClient.AttachSatelliteHost(createRegOptions)
	if err != nil {
		return fmt.Errorf("Error Generating Satellite Registration Script: %s", err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s-%s.sh", "addHost", name), resp, 0644)
	if err != nil {
		return fmt.Errorf("Error Creating Satellite Attach File: %s", err)
	}

	log.Printf("[INFO] Generated satellite location script to attach host : %s", name)

	d.SetId(*location.ID)

	return resourceIBMSatelliteLocationRead(d, meta)
}

func resourceIBMSatelliteLocationRead(d *schema.ResourceData, meta interface{}) error {
	ID := d.Id()

	satClient, err := meta.(ClientSession).SatelliteServiceV1()
	if err != nil {
		return err
	}

	bxClient, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}

	authMap := map[string]string{"Authorization": bxClient.Config.IAMAccessToken}
	getSatLocOptions := &kubernetesserviceapiv1.GetSatelliteLocationOptions{
		Headers:    authMap,
		Controller: &ID,
	}

	_, response, err := satClient.GetSatelliteLocation(getSatLocOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
	}

	return nil
}

func resourceIBMSatelliteLocationDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourceIBMSatelliteLocationExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	ID := d.Id()

	satClient, err := meta.(ClientSession).SatelliteServiceV1()
	if err != nil {
		return false, err
	}

	bxClient, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return false, err
	}

	authMap := map[string]string{"Authorization": bxClient.Config.IAMAccessToken}
	getSatLocOptions := &kubernetesserviceapiv1.GetSatelliteLocationOptions{
		Headers:    authMap,
		Controller: &ID,
	}

	_, response, err := satClient.GetSatelliteLocation(getSatLocOptions)
	if err != nil && response.StatusCode == 404 {
		d.SetId("")
		return false, fmt.Errorf("Error Getting Sateliite Location: %s\n%s", err, response)
	}

	return true, nil
}
