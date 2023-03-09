---
title: "lacework team-member"
slug: lacework_team-member
hide_title: true
---

## lacework team-member

Manage team members

### Synopsis

Manage Team Members to grant or restrict access to multiple Lacework Accounts. 
			  Team members can also be granted organization-level roles.


### Options

```
  -h, --help   help for team-member
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
* [lacework team-member create](lacework_team-member_create.md)	 - Create a new team member
* [lacework team-member delete](lacework_team-member_delete.md)	 - Delete a team member
* [lacework team-member list](lacework_team-member_list.md)	 - List all team members
* [lacework team-member show](lacework_team-member_show.md)	 - Show a team member by id

