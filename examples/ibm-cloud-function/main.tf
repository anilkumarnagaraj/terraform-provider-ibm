
resource "ibm_function_package" "package" {
  name = var.packageName
  namespace = var.namespace

  user_defined_parameters = <<EOF
        [
    {
        "key":"name",
        "value":"terraform"
    },
    {
        "key":"place",
        "value":"India"
    }
]
EOF

}


