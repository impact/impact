# Impact

[![Build Status](https://drone.io/github.com/xogeny/impact/status.png)](https://drone.io/github.com/xogeny/impact/latest)

*Impact* is a Modelica package manager.

![ImpactLogo](https://rawgithub.com/impact/impact/master/resources/images/logo_glossy.svg)

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

To build this, you need to have [Go](http://golang.org/) installed.
Go will create the proper build environment for you.
All you need to specify is:

`$ export GOPATH=/some/path`

and then run:

`$ go get github.com/impact/impact/impact`

which will automatically clone a copy of the git repository and all its dependencies,
compile impact and put them in a structure like

```
$GOPATH/
  bin/
  pkg/
  src/
    github.com/impact/impact
```

To build as static executable again, just run:

`$ go install`

from inside

`$GOPATH/src/github.com/impact/impact/impact`

This will create a static executable of called `impact` in `$GOPATH/bin`.

## Cross Compiling

The `impact/Makefile` includes targets to build cross-compiled executables.

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
The development takes place on https://github.com/impact/impact

 * Authors: [@mtiller](https://github.com/mtiller), [@dietmarw](https://github.com/dietmarw)
 * Contributors: See [graphs/contributors](https://github.com/impact/impact/graphs/contributors)

You may report any issues with using the [Issues](https://github.com/impact/impact/issues) button.

Contributions in shape of [Pull Requests](https://github.com/impact/impact/pulls) are always welcome.
