Impact
------

*Impact* is a Modelica package manager.

![ImpactLogo](https://rawgithub.com/xogeny/impact/master/images/logo_glossy.svg)

The concept was first presented in [impact - A Modelica Package Manager](docs/modelica2014/paper/impact.md)

> Michael Tiller, Dietmar Winkler (2014). impact - A Modelica Package Manager,
> Proceedings of the 10th International Modelica Conference, March 10-12, 2014,
> Lund, Sweden http://dx.doi.org/10.3384/ecp14096543

Installation and Usage
----------------------

Install by using [pip](http://www.pip-installer.org):

 * Linux/Mac:`$ pip install impact`
 * Windows:`c:\> pip.exe install impact`

and you should be able to use the command

    impact -h

to get the usage information displayed.

Conventions
-----------

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

Development
-----------

The development takes place on https://github.com/xogeny/impact
