##
# JWT authentication method using Gitlab CI token and environment variable lookups
auth "gitlab" {
  method = "jwt"
  path = "gitlab"
  priority = 100

  jwt {
    role = "terraform-ci-role"
    token = env("GITLAB_CI_TOKEN")
  }

  when {
    env_present = "GITLAB_CI_TOKEN"
  }
}

## 
# Active Directory

ad "vsphere" {
   role = "vsphere-svc"
   target_provider = "vsphere"
   extra_env_vars = {
       "VSPHERE_SERVER" = "vcenter.example.com"
   }
}

ad "generic" {
   role = "tf-svc"
   target_provider = "generic"
   username_env_var = "TF_VAR_AD_USERNAME"
   password_env_var = "TF_VAR_AD_PASSWORD"
}

##
# Kv2

# kv2 secrets engine, map default username/password attributes to infoblox
# provider environment variables

kv_secret "infoblox" {
   path = "dns/infoblox/infoblox-svc"
   target_provider = "infoblox"
   extra_env_vars = {
       "INFOBLOX_SERVER" = "infoblox.example.com"
   }
}

# reads from vault kv2 secret at db/pgsql and outputs environment variables
# mapped to
#  PGSQL=<value of psql_user attribute>
#  PGPASS=<value of psql_pass attribute>

kv_secret "generic" {
   path = "db/pgsql"
   target_provider = "generic"
   attribute_map = {
       "PGUSER" = "psql_user"
       "PGPASSWORD" = "psql_pass"
   }
   extra_env_vars = {
       "PGHOST" = "foo.bar.com"
       "PGPORT" = "12345"
   } 
}

##
# AWS STS

aws "sts" {
   method = "assumed_role"
   role = "terraform"
   role_arn = "arn:aws:iam::0000000000:role/Terraform"
   ttl = 900
}

##
# Microsoft Azure

azure "packer" {
   role = "packer-builder-ci"
   extra_env_vars = {
       "ARM_TENANT_ID" = "a9075e1c-ab41-4329-a5af-1b83cefa6e5e"
       "ARM_SUBSCRIPTION_ID" = "f442bd50-200b-4d57-ab86-5ac778f8d101"
   }
}
