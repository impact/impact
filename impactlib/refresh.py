import json

from impactlib.github import GitHub
from impactlib.semver import parse_semver
from impactlib.cache import cache_file_name

def process_user(repo_data, user, github, verbose):
    # Get a list of repositories
    repos = github.getRepos(user)

    # Iterate over each repository
    for repo in repos:
        # Extract the repository name
        name = repo["name"]
        if verbose:
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
            if verbose:
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
    process_user(repo_data, "modelica-3rdparty", github, args.verbose)

    # This gives the "modelica" user priority over "modelica-3rdparty"
    # in case of naming conflict
    process_user(repo_data, "modelica", github, args.verbose)
    
    # Write out repository data collected
    cache_file = cache_file_name()
    if args.verbose:
        print "Cache file: "+cache_file
    with open(cache_file, "w") as fp:
        json.dump(repo_data, fp, indent=4)
    if args.verbose:
        print "Refresh completed"
