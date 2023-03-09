## Developer Guidelines

## Signed Commits
Signed commits are required for any contribution to this project. Please see Github's documentation on configuring signed commits, [tell git about your signing key](https://docs.github.com/en/github/authenticating-to-github/managing-commit-signature-verification/telling-git-about-your-signing-key) and [signing commits](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits)

## Commit message standard

The format is:
```
type(scope): subject

BODY

FOOTER
```

Example commit message:
```
feat(cli): add --api_token global flag

This new flag will replace the use of `api_key` and `api_secret` so that
users can run the Lacework CLI only with an access token and their account:

    lacework cloud-account list --api_token _secret123 -a mycompany

Closes https://github.com/lacework/go-sdk/issues/282
```

Each commit message consists of a header, body, and footer. The header with the type and subject are mandatory, the scope is optional.
When writing a commit message try and limit each line of the commit to a max of 80 characters, so it can be read easily.

### Type

Allowed `type` valued.

| Type      | Description                                                                                            |
| -----     | -----------                                                                                            |
| feat:     | A new feature you're adding                                                                            |
| fix:      | A bug fix                                                                                              |
| style:    | Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc) |
| refactor: | A code change that neither fixes a bug nor adds a feature                                              |
| test:     | Everything related to testing                                                                          |
| docs:     | Everything related to documentation                                                                    |
| chore:    | Regular code maintenance                                                                               |
| build:    | Changes that affect the build                                                                          |
| ci:       | Changes to our CI configuration files and scripts                                                      |
| perf:     | A code change that improves performance                                                                |
| metric:   | A change that provides better insights about the adoption of features and code statistics              |

### Scope
The optional scope refers to the section that this commit belongs to, for example, changing a specific component or service, a directive, pipes, etc.
Think about it as an indicator that will let the developers know at first glance what section of your code you are changing.

A few good examples are:

* feat(client):
* docs(cli):
* chore(tests):
* ci(directive):

### Subject
The subject should contain a short description of the change, and written in present-tense, for example, use "add" and not "added",  or "change" and not "changed".
I like to fill this sentence below to understand what should I put as my description of my change:

If applied, this commit will ________________________________________.

### Body
The body should contain a longer description of the change, try not to repeat the subject and keep it in the present tense as above.
Put as much context as you think it is needed, don’t be shy and explain your thought process, limitations, ideas for new features or fixes, etc.

### Footer
The footer is used to reference issues, pull requests or breaking changes, for example, "Fixes ticket #123".
