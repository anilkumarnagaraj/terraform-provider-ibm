resource "ibm_satellite_location" "location" {
        name = var.location
        location = var.zone
        script_labels = ["env=prod"]
}