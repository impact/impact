# Contributor Guide

If you want to contribute, here is how the code is organized:

	* `config` - Types and functions related to settings.
	* `crawl` - This package is responsible for crawling through
	  a given collection of repositories and collecting information
	  about the libraries it finds.
	* `dirinfo` - This package holds the data structures that represent
	  information about a library that is found while crawling.
	* `graph` - This package contains types and functions related to
	  the dependency resolution functionality.  The types allow you to
	  build the graph of libraries and their dependencies and the
	  functions operate on those graphs.
	* `impact` - This package contains all the functionality related
	  to the command line tool.  This includes all the various
	  "sub-commands" as well as the "main" routine.
	* `index` - This package contains information about the complete
	  `impact` index.  This is different from the `graph`.  The
	  `graph` is essentially a subset of the information in the
	  complete index.  The `index` contains additional information
	  like what repository the source came from, who the author is,
	  etc.
	* `install` - Functions related installing libraries and interacting
	  with already installed libraries.
	* `parsing` - Functions and types related to parsing of versions,
	  normalizing versions and extracting uses annotations from Modelica
	  code.
	* `recorder` - This package contains the interfaces associated
	  with the implementations that record information about
	  repositories.  This is in a separate package from `crawl` in
	  order to avoid circular dependencies among Go packages.

## Testing

I've switched this project to use GoConvey.  In addition, it depends
on a small package I've created called `xconvey` which includes some
very simple shorthand functions.

To run tests from the command line, you can simply do the normal:

```
$ go test ./...
```

This will recurse into all the directories and give a reasonably nice
textual summary.  But I normally run:

```
$ goconvey
```

In the root directory.  This provides a web based UI at
http://localhost:8080.  Futhermore, I recommend turning on
notifications.  This means that you'll get a simply summary
notification every time the tests are run.  This is an easy way to
keep an eye on any regressions without having to have the full web
interface open all the time.  When a regression is found, it is then a
simple matter to open the web browser and see a much more detailed
report.
