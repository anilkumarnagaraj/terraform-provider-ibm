resource "ibm_function_action" "action" {
  name      = var.actionName
  namespace = var.namespace

  exec {
    kind = "nodejs:10"
    code = file("hello.js")
  }

}

