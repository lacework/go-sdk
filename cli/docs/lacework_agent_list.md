---
title: "lacework agent list"
slug: lacework_agent_list
hide_title: true
---

## lacework agent list

List all hosts with a running agent

### Synopsis

List all hosts that have a running agent in your environment.

You can use 'key:value' pairs to filter the list of hosts with the --filter flag.

    lacework agent list --filter 'os:Linux' --filter 'tags.VpcId:vpc-72225916'

**NOTE:** The value can be a regular expression such as 'hostname:db-server.*'

To filter hosts with a running agent version '5.8.0'.

    lacework agent list --filter 'agentVersion:5.8.0.*' --filter 'status:ACTIVE'

The available keys for this command are:
    * agentVersion
    * hostname
    * ipAddr
    * mid
    * mode
    * os
    * status
    * tags.arch
    * tags.ExternalIp
    * tags.Hostname
    * tags.InstanceId
    * tags.InternalIp
    * tags.LwTokenShort
    * tags.os
    * tags.VmInstanceType
    * tags.VmProvider
    * tags.Zone
    * tags.Account
    * tags.AmiId
    * tags.Name
    * tags.SubnetId
    * tags.VpcId
    * tags.Cluster
    * tags.cluster-location
    * tags.cluster-name
    * tags.cluster-uid
    * tags.created-by
    * tags.enable-oslogin
    * tags.Env
    * tags.GCEtags
    * tags.gci-ensure-gke-docker
    * tags.gci-update-strategy
    * tags.google-compute-enable-pcid
    * tags.InstanceName
    * tags.InstanceTemplate
    * tags.kube-labels
    * tags.lw_KubernetesCluster
    * tags.NumericProjectId
    * tags.ProjectId

```
lacework agent list [flags]
```

### Options

```
      --filter strings   filter results by key:value pairs (e.g. 'hostname:db-server.*')
  -h, --help             help for list
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

* [lacework agent](lacework_agent.md)	 - Manage Lacework agents

