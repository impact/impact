#!/usr/bin/env python

# This is a script that performs only the "refresh" functionality.  It
# is meant to be run by some "server" that will then publish the cache
# information to some public place.

import argparse
import os

from impactlib.refresh import refresh

parser = argparse.ArgumentParser(prog='impact')
parser.add_argument("-v", "--verbose", action="store_true",
                    help="Verbose mode", required=False)
parser.add_argument("-i", "--ignore", action="store_true",
                    help="Ignore packages with no versions", required=False)
parser.add_argument("-u", "--username", default=None,
                    help="GitHub username", required=False)
parser.add_argument("-p", "--password", action=None,
                    help="GitHub password", required=False)
parser.add_argument("-t", "--token", default=None,
                    help="GitHub OAuth token", required=False)
parser.add_argument("-o", "--output", default=None,
                    help="Output file", required=False)

args = parser.parse_args()

refresh(username=args.username,
        password=args.password,
        token=args.token,
        output=args.output,
        verbose=args.verbose,
        ignore_empty=args.ignore)
