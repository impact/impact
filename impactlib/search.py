from impactlib.load import load_repo_data
from fnmatch import fnmatch

try:
    import colorama
    from colorama import Fore, Back, Style
    colorama.init()
    use_color = True
except:
    use_color = False

def search(term, verbose):
    repo_data = load_repo_data()

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
        for m in sorted(matches):
            if verbose:
                if len(m[2].keys())==0:
                    versions = "None"
                else:
                    versions = ", ".join(m[2].keys())
                msg = "\n  Available versions: "+versions
                if use_color:
                    print (Fore.RED+m[0]+Fore.RESET+" - "+Fore.GREEN+m[1]
                           + Fore.CYAN + msg)
                else:
                    print m[0]+" - "+m[1] + msg
            else:
                if use_color:
                    print Fore.RED + m[0]
                else:
                    print m[0]
