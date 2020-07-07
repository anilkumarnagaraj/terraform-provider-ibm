variable "ibmcloud_api_key" {
  description = "Your IBM Cloud platform API key"
}

provider "ibm" {
  ibmcloud_api_key   = var.ibmcloud_api_key
}

