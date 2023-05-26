---
title: "lacework agent aws-install ec2ic"
slug: lacework_agent_aws-install_ec2ic
hide_title: true
---

## lacework agent aws-install ec2ic

Use EC2InstanceConnect to securely connect to EC2 instances

### Synopsis

This command installs the agent on all EC2 instances in an AWS account using EC2InstanceConnect.

To filter by one or more regions:

    lacework agent aws-install ec2ic --include_regions us-west-2,us-east-2

To filter by instance tag:

    lacework agent aws-install ec2ic --tag TagName,TagValue

To filter by instance tag key:

    lacework agent aws-install ec2ic --tag_key TagName

To explicitly specify the username for all SSH logins:

    lacework agent aws-install ec2ic --ssh_username <your-user>

To provide an agent access token of your choice, use the command 'lacework agent token list',
select a token and pass it to the '--token' flag. This flag must be selected if the
'--noninteractive' flag is set.

    lacework agent aws-install ec2ic --token <token>

To explicitly specify the server URL that the agent will connect to:

    lacework agent aws-install ec2ic --server_url https://your.server.url.lacework.net

To specify an AWS credential profile other than 'default':

    lacework agent aws-install ec2ic --credential_profile aws-profile-name

AWS credentials are read from the following environment variables:
- AWS_ACCESS_KEY_ID
- AWS_SECRET_ACCESS_KEY
- AWS_SESSION_TOKEN (optional)
- AWS_REGION (optional)

This command will only install the agent on hosts that are supported by
EC2InstanceConnect. The supported AMI types are Amazon Linux 2 and Ubuntu
16.04 and later. There may also be a region restriction.

This command will automatically add hosts with successful connections to
'~/.ssh/known_hosts' unless specified with '--trust_host_key=false'.

```
lacework agent aws-install ec2ic [flags]
```

### Options

```
      --credential_profile string   AWS credential profile to use (default "default")
  -h, --help                        help for ec2ic
  -r, --include_regions strings     list of regions to filter on
  -n, --max_parallelism int         maximum number of workers executing AWS API calls, set if rate limits are lower or higher than normal (default 50)
      --server_url https://         server URL that agents will talk to, prefixed with https:// (default "https://agent.lacework.net")
      --ssh_username string         username to login with
      --tag strings                 only install agents on infra with this tag
      --tag_key string              only install agents on infra with this tag key set
      --token string                agent access token
      --trust_host_key              automatically add host keys to the ~/.ssh/known_hosts file (default true)
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

