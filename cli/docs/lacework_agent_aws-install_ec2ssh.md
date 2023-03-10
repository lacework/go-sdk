---
title: "lacework agent aws-install ec2ssh"
slug: lacework_agent_aws-install_ec2ssh
hide_title: true
---

## lacework agent aws-install ec2ssh

Use SSH to securely connect to EC2 instances

### Synopsis

This command installs the agent on all EC2 instances in an AWS account
using SSH.

To filter by one or more regions:

    lacework agent aws-install ec2ssh --include_regions us-west-2,us-east-2

To filter by instance tag:

    lacework agent aws-install ec2ssh --tag TagName,TagValue

To filter by instance tag key:

    lacework agent aws-install ec2ssh --tag_key TagName

To provide an existing access token, use the '--token' flag. This flag is required
when running non-interactively ('--noninteractive' flag). The interactive command
'lacework agent token list' can be used to query existing tokens.

    lacework agent aws-install ec2ic --token <token>

You will need to provide an SSH authentication method. This authentication method
should work for all instances that your tag or region filters select. Instances must
be routable from your local host.

To authenticate using username and password:

    lacework agent aws-install ec2ssh --ssh_username <your-user> --ssh_password <secret>

To authenticate using an identity file:

    lacework agent aws-install ec2ssh -i /path/to/your/key

The environment should contain AWS credentials in the following variables:
- AWS_ACCESS_KEY_ID
- AWS_SECRET_ACCESS_KEY
- AWS_SESSION_TOKEN (optional),
- AWS_REGION (optional)

This command will automatically add hosts with successful connections to
'~/.ssh/known_hosts' unless specified with '--trust_host_key=false'.

```
lacework agent aws-install ec2ssh [flags]
```

### Options

```
  -h, --help                      help for ec2ssh
  -i, --identity_file string      identity (private key) for public key authentication (default "~/.ssh/id_rsa")
  -r, --include_regions strings   list of regions to filter on
  -n, --max_parallelism int       maximum number of workers executing AWS API calls, set if rate limits are lower or higher than normal (default 50)
      --ssh_password string       password for authentication
      --ssh_port int              port to connect to on the remote host (default 22)
      --ssh_username string       username to login with
      --tag strings               only select instances with this tag
      --tag_key string            only install agents on infra with this tag key
      --token string              agent access token
      --trust_host_key            automatically add host keys to the ~/.ssh/known_hosts file (default true)
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

* [lacework agent aws-install](lacework_agent_aws-install.md)	 - Install the datacollector agent on all remote AWS hosts

