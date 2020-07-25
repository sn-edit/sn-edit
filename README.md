<div style="margin: 0 auto;"><img src="assets/sn-edit.png"/></div>

# sn-edit

## Tools and Docs

- [Docs](https://docs.sn-edit.com/)
- [Config Generator](https://conf.sn-edit.com/)
- [Extensions](https://docs.sn-edit.com/#/?id=extensions-support)

## What is sn-edit?

The app has been built out of a need to have a way to develop scripts locally in your favorite editor.
With sn-edit, you are able to develop your scripts, change entries on your instance, without opening it.

You can use your favorite IDE/Editor to develop new features. We provide binaries for all the major platforms, this makes sn-edit
fully compatible with MacOS, Windows and Linux too. This is achieved due to the nature of the language sn-edit was built with.

We've built sn-edit in Go which makes it easy to support all the major platforms listed above.

> Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.

We would like to make sn-edit as minimal as possible. This means that in the long-run, we would rather have
fewer features which are stable, than to support many and have a lot of issues with the scope of it.

We are pro-simplicity. If we can get the information from a Servicenow instance, we would like to make the experience
as easy as possible. For the supported commands, we always try to minimize the impact of **you** providing the minimal amount
of flags needed to execute the command. In any case, that you find that we are asking too much, that could've been queried
from the instance and easily reusable, then let us know in an issue, and we will try to evaluate the requirement.

## Installation
Please refer to the docs for an [installation guide](https://docs.sn-edit.com/#/getting-started/README). You can find the guide for every major platform there. Please follow
the steps closely.

## Features

Here is a list of features we are supporting right now. The list is something more like a highlight of it. If you find the list incomplete
at any point, please send a pull request.

* Download an entry
* Upload fields of an entry
* Scope support
* Update sets support
* Masking the credentials (rest)
* Custom tables support
* Custom fields, saved into a file based on the configured extension (script => js, name => txt)
* Execute scripts on the instance
* A local low-profile sqlite database for metadata and usage inside of sn-edit

## Extensions support

We've built sn-edit in a way, to make it very easy to integrate anywhere. Due to it's nature of being supported on any major platform
and for the fact, that the commands are basically the same everywhere, we invite you, the community to develop extensions for
any IDE or Editor of choice. We may support some of them officially, but we are not able cover every one of them.

Official Extensions:
- [VSCode](https://github.com/sn-edit/vscode)

**Personal maintainer note**:

I invite all of you to create community built extensions to the major editors out there.
I am determined to provide a stable CLI as the building stone of the various extensions.

We would like to support every major Editor or IDE in this case too, so if you have the skill and time, please develop an extension for any
one of them. The best could be moved to an official repo under this organization. 
The idea is to have wide support for all the major platforms out there.