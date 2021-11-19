---
title: "lacework agent install"
slug: lacework_agent_install
hide_title: true
---

## lacework agent install

Install the datacollector agent on a remote host

### Synopsis

For single host installation of the Lacework agent via Secure Shell (SSH).

When this command is executed without any additional flag, an interactive prompt will be
launched to help gather the necessary authentication information to access the remote host.

To authenticate to the remote host with a username and password.

    lacework agent install <host> --ssh_username <your-user> --ssh_password <secret>

To authenticate to the remote host with an identity file instead.

    lacework agent install <user@host> -i /path/to/your/key

To provide an agent access token of your choice, use the command 'lacework agent token list',
select a token and pass it to the '--token' flag.

    lacework agent install <user@host> -i /path/to/your/key --token <token>

To authenticate to the remote host on a non-standard SSH port use the '--ssh_port' flag or
pass it directly via the argument.

    lacework agent install <user@host:port>
    

```
lacework agent install <[user@]host[:port]> [flags]
```

### Options

```
      --force                  override any pre-installed agent
  -h, --help                   help for install
  -i, --identity_file string   identity (private key) for public key authentication (default "~/.ssh/id_rsa")
      --ssh_password string    password for authentication
      --ssh_port int           port to connect to on the remote host (default 22)
      --ssh_username string    username to login with
      --token string           agent access token
      --trust_host_key         automatically add host keys to the ~/.ssh/known_hosts file
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

