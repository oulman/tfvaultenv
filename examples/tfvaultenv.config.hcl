ad "vsphere" {
   role = "vsphere-svc"
   target_provider = "vsphere"
   extra_env_vars = {
       "VSPHERE_SERVER" = "vcenter.example.com"
   }
}

ad "infoblox" {
   role = "infoblox-svc"
   target_provider = "infoblox"
   extra_env_vars = {
       "INFOBLOX_SERVER" = "infoblox.example.com"
   }
}

aws "sts" {
   method = "assumed_role"
   role = "terraform"
   role_arn = "arn:aws:iam::0000000000:role/Terraform"
   ttl = 900
}
