# Impact

[![Build Status](https://drone.io/github.com/xogeny/impact/status.png)](https://drone.io/github.com/xogeny/impact/latest)

*Impact* is a Modelica package manager.

![ImpactLogo](https://rawgithub.com/xogeny/impact/master/resources/images/logo_glossy.svg)

The concept was first presented in [impact - A Modelica Package Manager](resources/docs/modelica2014/paper/impact.md)

> Michael Tiller, Dietmar Winkler (2014). impact - A Modelica Package Manager,
> Proceedings of the 10th International Modelica Conference, March 10-12, 2014,
> Lund, Sweden http://dx.doi.org/10.3384/ecp14096543


## History

*Impact* was initially development in [Python](https://www.python.org/) which is also the version
presented in the [Modelica 2014 paper](resources/docs/modelica2014/paper/impact.md).

Based on user feedback it became apparent that having to rely on a Python run-time environment
installed in order to use *impact* might become an issue.

Therefore it was decided to switch to a [Go language](https://golang.org/) implementation
in an attempt to create a static executable that is compatible with the `impact` scheme
for Modelica package management  without the need for any additional run-time support
(*e.g.,* Node.js, Python).

### Version history

Available versions can be grouped into:

 * `>= v0.6.x`: Go-lang based implementation (*active* development on the [master branch](../../tree/master))
 * `<= v0.5.x`: Python based implementation (kept on  [python-version branch](../../tree/python-version))


## Status

So far, this [Go language](https://golang.org/) implementation can read the `impact` library
data and implements the `search`, `install` and `info` sub-commands.

It currently lacks support of the `refresh` sub-command which is needed to
build the package index (you'll need to use the  [python-version](../../tree/python-version))
for that)

## Installation

Self-contained executable binaries are available under the [release section](../../releases)
for a whole range of operating systems.

Simply download the matching archive and extract its content to a place
that suits you best (preferable inside a directory which is part of your
executable system `$PATH`).

## Conventions

*Impact* follows a "convention over configuration" philosophy.  That
means that if you follow some reasonable conventions (that generally
reflect best practices), the system should work without the need for
any manual configuration.  Here are the conventions that *Impact* expects:

* The name of the repository should match (case included) the name
  of your library.

* Semantic Versioning - To identify a library release, simply
  attach a tag to the release that is a [semantic
  version](http://semver.org) (an optional "v" at the start of the
  tag name is allowed).

* Place the `package.mo` file for your library in one of the
  following locations within the repository:

  * `./package.mo` (i.e., at the root of the repository)

  * `./<LibraryName>/package.mo` (i.e., within a directory sharing
    the name of the library)

  * `./<LibraryName> <Version>/package.mo` (i.e., within a directory sharing
    the name of the library followed by a space followed by the tag name,
    without any leading `v` present)

## Building

To build this, you need create the proper build environment.  This means you need to
create the following directory structure somewhere on your computer:

```
SomeDir/
  bin/
  pkg/
  src/
    gihub.com/xogeny/
```

Inside the `xogeny` directory, you need to do:

`$ git clone https://github.com/xogeny/impact`

Finally, it is essential to set the GOPATH environment variable to the
full name of `SomeDir`.

Once this is setup, you can go to the `impact` directory and do:

`$ go get`

...to install all the dependencies.  To run the client, you can do

`$ go run client.go [options]`

To build an static executable, just run:

`$ go install`

This will create a static executable of called `impact` in `SomeDir/bin`.

## Cross Compiling

The `Makefile` includes targets to build cross-compiled executables.

In order to be able to cross-compile you need to have
built GO for all the compilation targets.

### Under Ubuntu Linux

Those are already available in the repo:

```
$ sudo apt-get install golang-$GOOS-$GOARCH
```

The available `$GOOS` and `$GOARCH` variants are documented
in the [Go-lang documentation](https://golang.org/doc/install/source#environment).

### Under OSX

The cross compiling Go compiler can be installed with the
command:

```
  $ brew install go --cross-compile-common
```

## License
See [LICENSE](LICENSE) file

## Development
The development takes place on https://github.com/xogeny/impact

 * Authors: [@mtiller](https://github.com/mtiller), [@dietmarw](https://github.com/dietmarw)
 * Contributors: See [graphs/contributors](../../graphs/contributors)

You may report any issues with using the [Issues](../../issues) button.

Contributions in shape of [Pull Requests](../../pulls) are always welcome.
