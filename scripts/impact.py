#!/usr/bin/env python
import argparse
import urllib2
import json

class GitHub(object):
    BASE = "https://api.github.com"
    def __init__(self):
        pass
    def _req(self, path, headers={}):
        url = self.BASE+path
        print "url = "+str(url)
        req = urllib2.Request(url, headers={})
        response = urllib2.urlopen(req)
        return json.loads(response.read())
    def getRepos(self, user):
        repos = self._req("/users/"+user+"/repos")
        return repos

github = GitHub()

def refresh(args):
    # Get all repositories associated with the modelica user
    print "Do refresh"
    repos = github.getRepos("modelica")
    print "repos = "+str(repos)

def install(args):
    print "Install"

def search(args):
    print "Search"

parser = argparse.ArgumentParser(prog='impact')
subparsers = parser.add_subparsers(help='command help')

parser_refresh = subparsers.add_parser('refresh',
                                       help='Refresh package cache')
parser_refresh.set_defaults(func=refresh)

parser_search = subparsers.add_parser('search',
                                      help="Search for term in package")
parser_search.set_defaults(func=search)

parser_install = subparsers.add_parser('install',
                                       help="Install a named package")
parser_install.set_defaults(func=install)

args = parser.parse_args()
args.func(args)

