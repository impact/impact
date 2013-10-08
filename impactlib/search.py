from impactlib.cache import load_cache_file
from fnmatch import fnmatch

try:
    import colorama
    from colorama import Fore, Back, Style
    colorama.init()
    use_color = True
except:
    use_color = False

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
