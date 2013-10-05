#!/usr/bin/env python
import argparse
from fnmatch import fnmatch
import zipfile
import StringIO
import urllib2
import base64
import json
import sys
import re
import os

try:
    import colorama
    from colorama import Fore, Back, Style
    colorama.init()
    use_color = True
except:
    use_color = False

class GitHub(object):
    BASE = "https://api.github.com"
    def __init__(self, username=None, password=None, token=None):
        self.username = username
        self.password = password
        self.token = token
    def _req(self, path, headers={}, raw=False, isurl=False):
        # Construct base URL
        if isurl:
            url = path
        else:
            url = self.BASE+path

        # If we have an OAuth token, add it to the URL
        if self.token!=None:
            url = url+"?access_token="+str(self.token)

        # If we have a username and password, create the appropriate
        # Basic Authorization header
        if self.username!=None and self.password!=None:
            base64string = base64.encodestring("%s:%s" % (self.username,
                                                          self.password))
            base64string.replace("\n", "")
            headers["Authorization"] = "Basic %s" % (base64string,)

        # Formulate request
        req = urllib2.Request(url, headers=headers)

        # Get response
        response = urllib2.urlopen(req)

        if raw:
            # If the request is for the raw response, return the
            # (file-like) response object
            return response
        else:
            # Convert reponse (which should be JSON) into a python dictionary
            # and return it
            return json.loads(response.read())

    def getRepos(self, user):
        try:
            repos = self._req("/users/"+user+"/repos")
            return repos
        except Exception as e:
            print "Error fetching repositories: "+str(e)
            sys.exit(1)
    def getTags(self, user, repo):
        try:
            tags = self._req("/repos/"+user+"/"+repo+"/tags")
            return tags
        except Exception as e:
            print "Error accessing repository tags: "+str(e)
            sys.exit(1)
    def getFile(self, url):
        try:
            return self._req(url, isurl=True, raw=True)
        except Exception as e:
            print "Error downloading file: "+str(e)
            sys.exit(1)

def cache_file_name():
    return os.path.expanduser("~/.impact_cache")

def load_cache_file():
    with open(cache_file_name(), "r") as fp:
        return json.load(fp)

def parse_semver(tag):
    pat = """v?([0-9]+)\.([0-9]+)\.([0-9]+)"""
    c = re.compile(pat)
    m = c.match(tag)
    if m==None:
        return None
    major = m.group(1)
    minor = m.group(2)
    patch = m.group(3)
    return {"version": "%s.%s.%s" % (major, minor, patch),
            "major": major, "minor": minor, "patch": patch}

def process_user(repo_data, user, github):
    # Get a list of repositories
    repos = github.getRepos(user)

    # Iterate over each repository
    for repo in repos:
        # Extract the repository name
        name = repo["name"]
        if args.verbose:
            print "Repository: "+name

        # Initialize data for current repository
        data = {}

        # Pull out various pieces of information about the repository
        # and store it.
        data["description"] = repo["description"]

        # Prepare to extract all versions of this library
        data["versions"] = {}

        # Get the list of tags from GitHub
        tags = github.getTags(user, name)

        # Iterate over each tag
        for tag in tags:
            # Get the name of the tag
            tagname = tag["name"]
            if args.verbose:
                print "  Tag: "+tagname
            
            # Parse the tag to see if it is a semantic version number
            ver = parse_semver(tagname)
            if ver==None:
                print "    '"+tagname+"' is not a semantic version number"
                continue
            else:
                print "    Semantic version info: "+str(ver)

            # TODO: extract dependency information

            # Create a data structure for information related to this version
            tagdata = ver
            tagdata["zipball_url"] = tag["zipball_url"]
            tagdata["tarball_url"] = tag["tarball_url"]
            data["versions"][ver["version"]] = tagdata

        # Add data for this repository to master data structure
        repo_data[name] = data

def refresh(args):
    # Setup connection to github
    github = GitHub(username=args.username, password=args.password,
                    token=args.token)

    # Initialize respository data.  This is what we are refreshing
    # and what we will store eventually.
    repo_data = {}

    # Process all 3rd party libraries
    process_user(repo_data, "modelica-3rdparty", github)

    # This gives the "modelica" user priority over "modelica-3rdparty"
    # in case of naming conflict
    process_user(repo_data, "modelica", github)
    
    # Write out repository data collected
    cache_file = cache_file_name()
    if args.verbose:
        print "Cache file: "+cache_file
    with open(cache_file, "w") as fp:
        json.dump(repo_data, fp, indent=4)
    if args.verbose:
        print "Refresh completed"

def get_package(pkg):
    repo_data = load_cache_file()
    if not pkg in repo_data:
        msg = "No package named '"+pkg+"' found"
        if use_color:
            print Fore.RED+msg
        else:
            print msg
        return None
    return repo_data[pkg]

def semver_cmp(v1, v2, versions):
    maj1 = int(versions[v1]["major"])
    maj2 = int(versions[v2]["major"])
    if maj1==maj2:
        min1 = int(versions[v1]["minor"])
        min2 = int(versions[v2]["minor"])
        if min1==min2:
            pat1 = int(versions[v1]["patch"])
            pat2 = int(versions[v2]["patch"])
            return pat1>pat2
        else:
            return min1>min2
    else:
        return maj1>maj2

def latest_version(versions):
    if len(versions)==0:
        return None
    sorted_versions = sorted(versions,
                             cmp=lambda x, y: semver_cmp(x, y, versions))
    print "sorted_versions = "+str(sorted_versions)
    return sorted_versions[0]

def install_version(pkg, version, github):
    repo_data = load_cache_file()

    pdata = get_package(pkg)
    if pdata==None:
        return

    versions = pdata["versions"]

    vdata = None
    for ver in versions:
        if ver==version:
            vdata = versions[ver]

    if vdata==None:
        msg = "No version '"+str(version)+"' found for package '"+str(pkg)+"'"
        if use_color:
            print Fore.RED+msg
        else:
            print msg
        return

    zipurl = vdata["zipball_url"]
    if args.verbose:
        print "  URL: "+zipurl
    zfp = StringIO.StringIO(github.getFile(zipurl).read())
    zf = zipfile.ZipFile(zfp)
    zf.extractall()

def install(args):
    pkg = args.pkgname[0]
    pdata = get_package(pkg)

    if pdata==None:
        return

    version = args.version
    if version==None:
        version = latest_version(pdata["versions"])
        if args.verbose:
            print "  Choosing latest version: "+version
        if version==None:
            msg = "No (semantic) versions found for package '"+pkg+"'"
            if use_color:
                print Fore.RED+msg
            else:
                print msg
            return

    msg = "Installing version '"+version+"' of package '"+pkg+"'"
    if use_color:
        print Fore.GREEN+msg
    else:
        print msg

    # Setup connection to github
    github = GitHub(username=args.username, password=args.password,
                    token=args.token)

    install_version(pkg, version, github)

def search(args):
    repo_data = load_cache_file()

    term = args.term[0]
    matches = []
    for repo in repo_data:
        match = False
        data = repo_data[repo]
        if repo.find(term)>=0:
            match = True
        if fnmatch(repo, term):
            match = True
        if "description" in data and data["description"].find(term)>=0:
            match = True
        if match:
            matches.append((repo, data["description"], data["versions"]))
    if len(matches)==0:
        print "No matches found for search term '"+term+"'"
    else:
        for m in matches:
            if args.description:
                if use_color:
                    print Fore.RED+m[0]+Fore.RESET+" - "+Fore.GREEN+m[1]
                else:
                    print m[0]+" - "+m[1]
            else:
                if use_color:
                    print Fore.RED + m[0]
                else:
                    print m[0]
            if args.verbose:
                if len(m[2].keys())==0:
                    versions = "None"
                else:
                    versions = ", ".join(m[2].keys())
                msg = "  Available versions: "+versions
                if use_color:
                    print Fore.GREEN + msg
                else:
                    print msg

parser = argparse.ArgumentParser(prog='impact')
subparsers = parser.add_subparsers(help='command help')

parser_refresh = subparsers.add_parser('refresh',
                                       help='Refresh package cache')
parser_refresh.add_argument("-v", "--verbose", action="store_true",
                            help="Verbose mode", required=False)
parser_refresh.add_argument("-u", "--username", default=None,
                            help="GitHub username", required=False)
parser_refresh.add_argument("-p", "--password", action=None,
                            help="GitHub password", required=False)
parser_refresh.add_argument("-t", "--token", default=None,
                               help="GitHub OAuth token", required=False)
parser_refresh.set_defaults(func=refresh)

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
parser_install.add_argument("-u", "--username", default=None,
                            help="GitHub username", required=False)
parser_install.add_argument("-p", "--password", action=None,
                            help="GitHub password", required=False)
parser_install.add_argument("-t", "--token", default=None,
                            help="GitHub OAuth token", required=False)
parser_install.set_defaults(func=install)

args = parser.parse_args()

args.func(args)
