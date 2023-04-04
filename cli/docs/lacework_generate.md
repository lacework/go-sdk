---
title: "lacework generate"
slug: lacework_generate
hide_title: true
---

## lacework generate

Generate code to onboard your account

### Synopsis

Generate code to onboard your account and deploy Lacework into various cloud environments.

This command creates Infrastructure as Code (IaC) in the form of Terraform HCL, with the option of running
Terraform and deploying Lacework into AWS, Azure, or GCP.


### Options

```
  -h, --help   help for generate
```

### Options inherited from parent commands

```
  -a, --account string      account subdomain of URL (i.e. <ACCOUNT>.lacework.net)
  -k, --api_key string      access key id
  -s, --api_secret string   secret access key
      --api_token string    access token (replaces the use of api_key and api_secret)
      --debug               turn on debug logging
      --json                switch commands output from human-readable to json format
      --nocache             turn off caching
      --nocolor             turn off colors
      --noninteractive      turn off interactive mode (disable spinners, prompts, etc.)
      --organization        access organization level data sets (org admins only)
  -p, --profile string      switch between profiles configured at ~/.lacework.toml
      --subaccount string   sub-account name inside your organization (org admins only)
```

### SEE ALSO

* [lacework](lacework.md)	 - A tool to manage the Lacework cloud security platform.
* [lacework generate cloud-account](lacework_generate_cloud-account.md)	 - Generate cloud integration IaC
* [lacework generate k8s](lacework_generate_k8s.md)	 - Generate Kubernetes integration IaC

