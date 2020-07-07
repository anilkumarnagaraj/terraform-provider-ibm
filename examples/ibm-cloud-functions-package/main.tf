resource "ibm_function_package" "package" {
  name      = var.packageName
  namespace = var.namespace
  publish = true
}
