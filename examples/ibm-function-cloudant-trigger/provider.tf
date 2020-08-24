variable "function_namespace" {
  description = "Your Cloud Functions namespace"
}

provider "ibm" {
  function_namespace = var.function_namespace
}

