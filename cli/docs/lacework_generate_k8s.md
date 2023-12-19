---
title: "lacework generate k8s"
slug: lacework_generate_k8s
hide_title: true
---

## lacework generate k8s

Generate Kubernetes integration IaC

### Synopsis

Generate IaC to deploy Lacework into a Kubernetes platform.

This command creates Infrastructure as Code (IaC) in the form of Terraform HCL, with the option of running
Terraform and deploying Lacework into GKE.


### Options

```
  -h, --help   help for k8s
```

### Options inherited from parent commands

```
  -a, --account string      account subdomain of URL (i.e. <ACCOUNT>.lacework.net)
  -k, --api_key string      access key id
  -s, --api_secret string   secret access key
      --api_token string    access token (replaces the use of api_key and api_secret)
      --apply               run terraform apply without executing plan or prompting
      --debug               turn on debug logging
      --json                switch commands output from human-readable to json format
      --nocache             turn off caching
      --nocolor             turn off colors
      --noninteractive      turn off interactive mode (disable spinners, prompts, etc.)
      --organization        access organization level data sets (org admins only)
      --output string       location to write generated content
  -p, --profile string      switch between profiles configured at ~/.lacework.toml
      --subaccount string   sub-account name inside your organization (org admins only)
```

### SEE ALSO

* [lacework generate](lacework_generate.md)	 - Generate code to onboard your account
* [lacework generate k8s eks](lacework_generate_k8s_eks.md)	 - Generate and/or execute Terraform code for EKS integration
* [lacework generate k8s gke](lacework_generate_k8s_gke.md)	 - Generate and/or execute Terraform code for GKE integration

