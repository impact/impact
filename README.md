Background
==========

This is an attempt to create a static executable that is compatible
with the `impact` scheme for Modelica package management.

This implementation is an experiment to understand whether it is
possible to create an implementation of `impact` without the need for
any additional runtime support (*e.g.,* Node.js, Python).

Status
======

So far, this implementation can read the `impact` library data and
implements the `search` subcommand.

Building
========

To build this, you need create the proper build environment.  This means you need to
create the following directory structure somewhere on your computer:

SomeDir/
  bin/
  pkg/
  src/
    xogeny/

Inside the `xogeny` directory, you need to do:

`$ git clone https://github.com/xogeny/gimpact`

Finally, it is essential to set the GOPATH environment variable to the
full name of `SomeDir`.

Once this is setup, you can go to the `gimpact` directory and do:

`$ go get`

...to install all the dependencies.  To run the client, you can do

`$ go run client.go [options]`

To build an static executable, just run:

`$ go install`

This will create a static executable of called `gimpact` in `SomeDir/bin`.
