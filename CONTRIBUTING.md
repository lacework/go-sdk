# Contributing to the Lacework Go-sdk

#### Table Of Contents

[Before getting started?](#before-getting-started)

[How to Contribute?](#how-can-i-contribute)
* [Reporting Bugs](#reporting-bugs)
* [Feature Requests](#feature-requests)
* [Pull Requests](#pull-requests)

[Styleguides](#styleguides)
* [Git Commit Messages](#commit-message-standard)


## Before getting started

Read the [README.md](https://github.com/lacework/go-sdk/wiki/CLI-Documentation), the [development documentation](https://github.com/lacework/go-sdk/tree/main/cli#development) 
and the Lacework CLI documentation on the [wiki](https://github.com/lacework/go-sdk/wiki/CLI-Documentation).

## Reporting Bugs

Ensure the issue you are raising has not already been created under [issues](https://github.com/lacework/go-sdk/issues).

If no current issue addresses the problem, open a new [issue](https://github.com/lacework/go-sdk/issues/new).
Include as much relevant information as possible. See the [bug template](https://github.com/lacework/go-sdk/) for help on creating a new issue.

## Feature Requests

If you wish to submit a request to add new functionality or an improvement to the go-sdk then use the the [feature request]() template to 
open a new [issue](https://github.com/lacework/go-sdk/issues/new)

## Creating a Pull Request

When submitting a pull request follow the [commit message standard](#commit-message-standard).
Reduce the likelihood of pushing breaking changes by running the go-sdk unit and integration tests, 
see [development documentation](https://github.com/lacework/go-sdk/tree/main/cli#development).

## Commit message standard

The format is:
type(scope): subject
BODY
FOOTER

Each commit message consists of a header, body, and footer. The header is mandatory, the scope is optional, the type and subject are mandatory.
When writing a commit message try and limit each line of the commit to a max of 100 hundred characters, so it can be read easily.

### Type

| Type | Description |
| ----- | ----------- |
| feat: | A new feature you're adding |
| fix: |A bug fix|
| style: | Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc) |
| refactor: | A code change that neither fixes a bug nor adds a feature |
| test: | Everything related to testing |
| docs: | Everything related to documentation |
| chore: | Regular code maintenance |
| build: | Changes that affect the build |
| ci: | Changes to our CI configuration files and scripts |
| perf: | A code change that improves performance |
| metric: | A change that provides better insights about the adoption of features and code statistics |

### Scope
The optional scope refers to the section that this commit belongs to, for example, changing a specific component or service, a directive, pipes, etc. 
Think about it as an indicator that will let the developers know at first glance what section of your code you are changing.

A few good examples are:

feat(client):
docs(cli):
chore(tests):
ci(directive):

### Subject
The subject should contain a short description of the change, and written in present-tense, for example, use “add” and not “added”,  or “change” and not “changed”. 
I like to fill this sentence below to understand what should I put as my description of my change:

If applied, this commit will ________________________________________.

### Body
The body should contain a longer description of the change, try not to repeat the subject and keep it in the present tense as above. 
Put as much context as you think it is needed, don’t be shy and explain your thought process, limitations, ideas for new features or fixes, etc.

### Footer
The footer is used to reference issues, pull requests or breaking changes, for example, “Fixes ticket #123”.

Thanks,

Project Maintainers

