#!/usr/bin/env python
import argparse

from impactlib.search import search
from impactlib.install import install
from impactlib import config

DEFAULT_REPOS = config.get_repos()
print "DEFAULT_REPOS = "+str(DEFAULT_REPOS)

parser = argparse.ArgumentParser(prog='impact')
subparsers = parser.add_subparsers(help='command help')

parser_search = subparsers.add_parser('search',
                                      help="Search for term in package")
parser_search.add_argument("term", nargs=1)
parser_search.add_argument("-v", "--verbose", action="store_true",
                           help="Verbose mode", required=False)
parser_search.add_argument("-d", "--description", action="store_true",
                           help="Include description", required=False)
parser_search.set_defaults(func=search)

parser_install = subparsers.add_parser('install',
                                       help="Install a named package")
parser_install.add_argument("pkgname", nargs=1)
parser_install.add_argument("version", nargs="?")
parser_install.add_argument("-v", "--verbose", action="store_true",
                            help="Verbose mode", required=False)
parser_install.add_argument("-d", "--dry_run", action="store_true",
                            help="Suppress installation", required=False)
parser_install.add_argument("-u", "--username", default=None,
                            help="GitHub username", required=False)
parser_install.add_argument("-p", "--password", action=None,
                            help="GitHub password", required=False)
parser_install.add_argument("-t", "--token", default=None,
                            help="GitHub OAuth token", required=False)
parser_install.set_defaults(func=install)

args = parser.parse_args()

args.func(args)
