# Release Notes
Another day, another release. These are the release notes for the version `v0.1.6`.

### New `lacework event` command
This new command includes two sub-commands, `list` and `show`, check them out!
```
$ lacework event list
  EVENT ID |                TYPE                | SEVERITY |      START TIME      |       END TIME
-----------+------------------------------------+----------+----------------------+-----------------------
        10 | NewViolations                      | High     | 2020-04-20T13:00:00Z | 2020-04-20T14:00:00Z
         4 | VPCNetworkFirewallRuleChanged      | Medium   | 2020-04-16T20:00:00Z | 2020-04-16T21:00:00Z
         8 | VPCNetworkRouteChanged             | Medium   | 2020-04-19T23:00:00Z | 2020-04-20T00:00:00Z
         1 | ProjectOwnershipAssignmentsChanged | Medium   | 2020-04-16T17:00:00Z | 2020-04-16T18:00:00Z
         6 | NewViolations                      | Medium   | 2020-04-18T13:00:00Z | 2020-04-18T14:00:00Z
         3 | VPCNetworkChanged                  | Medium   | 2020-04-16T20:00:00Z | 2020-04-16T21:00:00Z
         2 | CloudStorageIAMPermissionChanged   | Medium   | 2020-04-16T18:00:00Z | 2020-04-16T19:00:00Z
         5 | CloudStorageIAMPermissionChanged   | Low      | 2020-04-17T19:00:00Z | 2020-04-17T20:00:00Z
         9 | VPCNetworkRouteChanged             | Low      | 2020-04-20T04:00:00Z | 2020-04-20T05:00:00Z
         7 | VPCNetworkFirewallRuleChanged      | Low      | 2020-04-19T23:00:00Z | 2020-04-20T00:00:00Z
```

## Features
* feat(api): add EventsService to inspect events (Salim Afiune Maya)([533a271](https://github.com/lacework/go-sdk/commit/533a2713f5c179e50c90c63318991643f005a750))
* feat(api): add Details func to EventsService (Salim Afiune Maya)([56b95ca](https://github.com/lacework/go-sdk/commit/56b95ca2c02c2f8af24dd351e7fe6247b4da7eba))
* feat(cli): new event list command (Salim Afiune Maya)([d7c9f9e](https://github.com/lacework/go-sdk/commit/d7c9f9e2c41a1bd1411b92b6b8632aa2e32845dd))
* feat(cli): new event show command (Salim Afiune Maya)([8f75c78](https://github.com/lacework/go-sdk/commit/8f75c78d222f3fa00d24cf41e8d1e712f6600122))
* feat(cli): `--noninteractive` mode flag (Salim Afiune Maya)([10536af](https://github.com/lacework/go-sdk/commit/10536afe1d6ce76ba3391c145f3533b4d6725484))
## Bug Fixes
* fix(api): omitempty integration responses fields (Salim Afiune Maya)([44e2314](https://github.com/lacework/go-sdk/commit/44e2314f4ca02f1c0e6a134bedbc82161a81473c))
## Documentation Updates
* docs(cli): document environment variables (Salim Afiune Maya)([0012ec1](https://github.com/lacework/go-sdk/commit/0012ec14f574f6e4c1dc2a5d774e17ef038f1308))
## Other Changes
* chore(cli): update usage of commands (Salim Afiune Maya)([5dd3057](https://github.com/lacework/go-sdk/commit/5dd3057371fe87434f2da54f68cdcc3dc5fd754a))
* chore(cli): style updates to release scripts (Salim Afiune Maya)([f4355bf](https://github.com/lacework/go-sdk/commit/f4355bf481a3349fe9bac3700a1eaa8e80227238))
