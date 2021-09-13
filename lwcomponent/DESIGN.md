# Lacework Components

Proposed by: [Salim Afiune Maya](https://github.com/afiune)

As the Lacework platform grows and introduces new products and services, the security ecosystem will
have to grow at the same speed (or even faster) to adopt such new products. This design is a proposal
for adopting a new model to deliver tools and libraries to Lacework users.

## Motivation

The Lacework CLI was designed with the goal of providing fast, accurate, and actionable insights into
the Lacework platform. The first release [v0.1.0](https://github.com/lacework/go-sdk/releases/tag/v0.1.0) was on March 27, 2020 and since then, the Lacework CLI
has been adopted broadly. A year later, Lacework released a couple more binaries and today, it is clear
to us that we will adopt this model for more products and services in the future.

This model of modular binaries is flexible and will allow us to release features faster to our users,
but it also lacks from a cohesive ecosytem, an easy way for users to discover, install, configure and
manage all of these tools.

This design aims to solve this problem by introducing the concept of **components** into the Lacework CLI,
with these new components, the Lacework CLI will evolve into the SDK realm and therefore, we are proposing
to name it the Lacework SDK (Software Development Kit) since it will now provide truly a combination of
tools and libraries to build a robust and flexible ecosystem for our users.

## Scope

The main goal of this design is to introduce the concept of components to build a flexible and robust ecosystem
around the Lacework platform.

This design also aims to:

* Unify the installation and configuration of tools provided by Lacework
* Help users discover new tools from new Lacework products
* Provide the same experience for upgrading components
* Allows users to create new tools to enhance Lacework workflows

## What are Components?

A component can be a command-line tool, a set of Lacework CLI commands, or a package that contains dependencies
used by a tool in the Lacework SDK.
