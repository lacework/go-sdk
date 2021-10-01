# Lacework Components

Proposed by: [Salim Afiune Maya](https://github.com/afiune)

As the Lacework platform grows and introduces new products and services, our security ecosystem will
have to grow at the same speed (or even faster) to adopt such new products. This design is a proposal
for adopting a new model to deliver tools and libraries to Lacework users.

## Motivation

The Lacework CLI was designed with the goal of providing fast, accurate, and actionable insights into
the Lacework platform. The first release [v0.1.0](https://github.com/lacework/go-sdk/releases/tag/v0.1.0) was on March 27, 2020 and since then, the Lacework CLI
has been broadly adopted. A year later, Lacework released a couple more binaries and today, it is clear
to us that we will adopt this model of releasing separate binaries for more products and services in the
future.

This model of modular binaries is flexible and will allow us to release features faster to our users,
but it also lacks from a cohesive ecosytem, an easy way for users to discover, install, configure and
manage all of these tools.

This design aims to solve this problem by introducing the concept of **components** into the Lacework CLI,
with these new components, the Lacework CLI will evolve into the SDK realm and therefore, we are proposing
to name it the Lacework CDK (Cloud Development Kit) since it will now provide truly a combination of
tools and libraries to build a robust and flexible ecosystem for our users.

## Scope

The main goal of this design is to introduce the concept of components to build a flexible and robust ecosystem
around the Lacework platform.

This design also aims to:

* Unify the installation and configuration of tools provided by Lacework
* Help users discover new tools from new Lacework products
* Provide the same experience for managing components
* Allows users to create new tools to enhance security workflows
* Make it easy to manage and distribute libraries and content

## What Are Components?

A component can be a command-line tool, a set of Lacework commands, or a package that contains dependencies
used by another Lacework component.

### Component Specifications

Every component should follow the following specifications:


| Name   |   Type    |   Description                                                 |
| ------ | --------- | ------------------------------------------------------------- |
| `name` | `string` | The name of the component |
| `description` | `string` | A long description of the purpose of the component |
| `type` | `enum(Type)` | The component type (read more about types here) |
| `version` | `string` | The version of the component in semantic format (`MAJOR.MINOR.PATCH`) (different from the overall components specifications) |
| `size` | `int64` | The component size in bytes |
| `checksum` | `string` | SHA256 (256-bit) checksums of the component |
| `download_url` | `string` | The URL from where to download the component |
| `dependencies` | `array(Component)` | A list of components that the component depends on |

These specifications should not be hardcoded, they should be disigned to be extensible since they will change as we
expand the usage and purpose of these components.

The specifications from all registered components will be provided by a new components service which will have
a semantic version (`MAJOR.MINOR.PATCH`) so that when the specifications change, a new version will be released
and our users will get notified.

### Components Internal Service

We should have a very lightweight service that will be the single source of truth of all available Lacework components,
this service should fulfill the following use cases:

* Provide an API to fetch the current state of all Lacework components and their specifications
* Provide an API to define (create) new components, when new components are added to this service, our users will be notified automatically
* Provide an API to deprecate a component, read more about deprecations below
* Provide an API to trigger a sync of a single component, useful for orchestrating release pipelines 
* Configure a batch process that runs every 10 minutes to verify that all components are in sync

### Component Synchronization

A component synchronization is a task that the internal components' service does to verify the latest version of one
or multiple components, when there is a new version of a component, this task updates the description, version, size,
checksum, and dependencies of the component.

Note that changing the type of the component is discouraged.

### Signature And File Checksum

As a security company, we need to ensure that any binary we install on our users' workstation is coming from us, the
installation and upgrade process will have a requirement that every component should be signed with Lacework's PGP
key, if the downloaded component doesn't match the PGP signature, we should delete the downloaded binary and
notify the user.

A second safety we should have is to check the SHA256 (256-bit) checksums of the downloaded binary or compressed file
which should match with the one provided by (APIs) the new components internal service.

### Create A New Component 

To create a new component, we need to define the following things:

* Define the component type (binary, commands, or content)
* For binary components, have cross-platform binaries (support windows, linux and osx)
* Automate the release process via CD pipelines


## Observability

End to end observability?

## Deprecations

## Attribution to third-party tools

## Open Questions:
* Will there be components that doesn't have cross-platform support? If yes, document examples and add details about how to handle these cases
* ...

