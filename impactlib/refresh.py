import json
import re
import os
# urlparse is split up in python 3:
try:
    from urllib.parse import urlparse
except ImportError:
    from urlparse import urlparse

from impactlib.github import GitHub
from impactlib.semver import SemanticVersion
from impactlib import config

def extract_dependencies(fp):
    deps = {}
    try:
        contents = fp.read().decode(encoding='utf8')
    except UnicodeDecodeError:
        try:
            contents = fp.read().decode(encoding='latin1')
            print("WARNING: Library uses non-standard latin1 encoding.\n"
                  "Consider converting the library to UTF8!\n")
        except UnicodeDecodeError:
            print("ERROR: Library uses illegal text encoding.\n"
                  "Skipping check for dependencies.\n"
                  "The library needs to get updated to UTF8!\n")
            return
    contents = contents.replace("\n","")
    contents = contents.replace("\r","")
    pat = """([A-Za-z]\w*)\s*\(\s*version\s*=\s*"(\d+\.\d+)(\.\d+)?"\s*\)"""
    c = re.compile(pat)
    matches = re.findall(c, contents)
    for m in matches:
        if m[0]=="uses":
            continue
        deps[m[0]] = m[1]+m[2]
    ret = []
    for dep in deps:
        ret.append({"name": dep, "version": deps[dep]})
    return ret

def strip_extra(version):
    return (version.split("-")[0]).split("+")[0]

def get_package_details(user, repo, tag, github, ver, verbose):
    """
    This function returns a tuple.

    The first element in the tuple is that path of the actual Modelica
    package.  This tells us what, after extracting the zipball, needs
    to be moved to the install path.

    The second element in the tuple is the list of dependencies.
    """
    root = github.getRawFile(user, repo, tag, "package.mo")
    if root!=None:
        deps = extract_dependencies(root)
        root.close()
        return (".", deps)
    elif verbose:
        print("Not in root directory")
    unver_dir = github.getRawFile(user, repo, tag, repo+"/package.mo")
    if unver_dir!=None:
        deps = extract_dependencies(unver_dir)
        unver_dir.close()
        return (repo, deps)
    elif verbose:
        print("Not in unversioned directory of the same name")
    version = tag
    if version[0]=="v":
        version = version[1:]
    version = strip_extra(version)
    ver_name = repo+" "+version
    path = repo+"%20"+version+"/package.mo"
    ver_dir = github.getRawFile(user, repo, tag, path)
    if ver_dir!=None:
        deps = extract_dependencies(ver_dir)
        ver_dir.close()
        return (ver_name, deps)
    elif verbose:
        print("Not in versioned directory ("+ver_name+")")
    return (None, [])

def process_github_user(repo_data, user, pat, github, verbose,
                        tolerant, ignore_empty):
    # Get a list of repositories
    repos = github.getRepos(user)

    c = re.compile(pat)

    # Iterate over each repository
    for repo in repos:
        # Extract the repository name
        name = repo["name"]
        if c.match(name)==None:
            continue
        print("Repository: "+name)

        # Initialize data for current repository
        data = {}

        # Pull out various pieces of information about the repository
        # and store it.
        data["description"] = repo["description"]

        # If homepage field exist store this otherwise use repo home
        data["homepage"] = repo["homepage"] or repo["html_url"]

        # Prepare to extract all versions of this library
        data["versions"] = {}

        # Get the list of tags from GitHub
        tags = github.getTags(user, name)

        # Iterate over each tag
        for tag in tags:
            # Get the name of the tag
            tagname = tag["name"]
            if verbose:
                print("  Tag: "+tagname)

            # Parse the tag to see if it is a semantic version number
            try:
                ver = SemanticVersion(tagname, tolerant=tolerant)
            except ValueError as e:
                if verbose:
                    print("Exception: "+str(e))
                continue

            # TODO: extract dependency information
            (path, deps) = get_package_details(user, name, tagname,
                                               github, ver, verbose)
            if path==None:
                print("Couldn't find Modelica package root")
                continue

            print("  Semantic version info: "+str(ver))
            print("    Path: "+str(path))
            print("    Dependencies: "+str(deps))

            # Create a data structure for information related to this version
            tagurlbase = ('https://github.com/%s/%s/archive/%s'
                                      % (str(user), str(name), str(tagname)))
            tagdata = ver.json()
            tagdata["zipball_url"] = tagurlbase+".zip"
            tagdata["tarball_url"] = tagurlbase+".tar.gz"
            if "commit" in tag and "sha" in tag["commit"]:
                tagdata["sha"] = tag["commit"]["sha"]
            tagdata["path"] = path
            tagdata["dependencies"] = deps

            data["versions"][str(ver)] = tagdata
            # Useful for legacy (non-semver) versions
            tver = tagname
            if tver[0]=="v":
                tver = tver[1:]
            if str(ver)!=tver:
                if verbose:
                    print("  Also storing under version: "+tver)
                data["versions"][tver] = tagdata

        if len(data["versions"])==0:
            print("  No useable version tags found")
            if ignore_empty:
                continue

        # Add data for this repository to master data structure
        repo_data[name] = data

def refresh(output, verbose, tolerant, ignore_empty, source_list=None):
    username = config.get("Impact", "username", None)
    password = config.get("Impact", "password", None)
    token = config.get("Impact", "token", None)

    if source_list==None:
        source_list = config.get("Impact", "source_list",
                        "github://modelica-3rdparty/.*,github://modelica/.*")
        source_list = source_list.split(",")

    if verbose:
        if username!=None:
            print("Using username: "+username+" to authenticate")
        if token!=None:
            print("Using API token to authenticate")

    # Setup connection to github
    github = GitHub(username=username, password=password,
                    token=token)

    # Initialize respository data.  This is what we are refreshing
    # and what we will store eventually.
    repo_data = {}

    for source in source_list:
        if verbose:
            print("Scanning "+source)
        data = urlparse(source)
        if data.scheme=="github":
            user = data.netloc
            repo_pat = data.path[1:]
            process_github_user(repo_data, user=user, github=github, pat=repo_pat,
                                verbose=verbose, tolerant=tolerant,
                                ignore_empty=ignore_empty)
        else:
            print("Unknown scheme: "+data.scheme+" in "+source+", skipping")

    # Process all 3rd party libraries
    #process_github_user(repo_data, user="modelica-3rdparty", github=github,
    #             verbose=verbose, tolerant=tolerant, ignore_empty=ignore_empty)

    # This gives the "modelica" user priority over "modelica-3rdparty"
    # in case of naming conflict
    #process_user(repo_data, user="modelica", github=github, verbose=verbose,
    #             tolerant=tolerant, ignore_empty=ignore_empty)

    # Write out repository data collected
    if output==None:
        print(json.dumps(repo_data, indent=4))
    else:
        if verbose:
            print("Output file: "+output)
        with open(output, "w") as fp:
            json.dump(repo_data, fp, indent=4)
    if verbose:
        print("Refresh completed")
