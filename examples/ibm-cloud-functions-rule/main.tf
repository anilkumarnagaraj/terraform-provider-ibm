resource "ibm_function_action" "action" {
  name      = var.actionName
  namespace = var.namespace
  exec {
         kind = "nodejs:6"
         code = file("hellonode.js")
  }
}

resource "ibm_function_trigger" "trigger" {
  name = var.triggerName
  namespace = var.namespace
  feed {
         name       = "/whisk.system/alarms/alarm"
         parameters = <<EOF
                            [
                               {
                                  "key":"cron",
                                  "value":"0 */2 * * *"
                                }
                             ]
                  EOF

         }

         user_defined_annotations = <<EOF
                                          [
                                             {
                                                "key":"sample trigger",
                                                 "value":"Trigger for hello action"
                                              }
                                           ]

          EOF

}

resource "ibm_function_rule" "rule" {
  name         = var.ruleName
  namespace    = var.namespace
  trigger_name = ibm_function_trigger.trigger.name
  action_name  = ibm_function_action.action.name
}
