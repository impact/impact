# Contributor Guide

If you want to contribute, here is how the code is organized:

    * `cmdsline` - This package contains all the functionality related
	  to the command line tool.  This includes all the various
	  "sub-commands" as well as the "main" routine.
	* `graph` - This package contains types and functions related to
	  the dependency resolution functionality.  The types allow you to
	  build the graph of libraries and their dependencies and the
	  functions operate on those graphs.
	* `index` - This pckage contains information about the complete
	  `impact` index.  This is different from the `graph`.  The
	  `graph` is essentially a subset of the information in the
	  complete index.  The `index` contains additional information
	  like what repository the source came from, who the author is,
	  etc.
