variable "namespace" {
  default = "Namespace-10"
}

variable "packageName" {
  default = "tutils1"
}

variable "actionName" {
  default = "thello1"
}

variable "boundPackageName" {
  default = "tmycloudant"
}

variable "triggerName" {
  default = "tmyCloudantTrigger"
}

variable "ruleName" {
  default = "tcloudantRule"
}

variable "dbname" {
  default = "tdatabasedemo"
}

variable "space" {
  default = "dev"
}

variable "org" {
  default = "namespacecf"
}

variable "service" {
  default = "cloudantNoSQLDB"
}

variable "plan" {
  default = "standard"
}

variable "service_instance_name" {
  default = "tmycloudantdb"
}

variable "service_key_name" {
  default = "tmycloudantdbkey"
}

variable "app_version" {
  default = "1"
}

variable "git_repo" {
  default = "https://github.com/hkantare/cf-cloudant-python.git"
}

variable "dir_to_clone" {
  default = "/tmp/my_cf_code"
}

variable "app_zip" {
  default = "/tmp/myzip.zip"
}

variable "route" {
  default = "tmy-app-cloudant"
}

variable "app_name" {
  default = "tmyapp"
}

variable "app_command" {
  default = "python app.py"
}

variable "buildpack" {
  default = ""
}

