import zipfile
import StringIO

from impactlib.load import load_repo_data
from impactlib.github import GitHub

try:
    import colorama
    from colorama import Fore, Back, Style
    colorama.init()
    use_color = True
except:
    use_color = False

def get_package(pkg):
    repo_data = load_repo_data()
    if not pkg in repo_data:
        msg = "No package named '"+pkg+"' found"
        if use_color:
            print Fore.RED+msg
        else:
            print msg
        return None
    return repo_data[pkg]

def latest_version(versions):
    if len(versions)==0:
        return None
    sorted_versions = sorted(versions,
                             cmp=lambda x, y: semver_cmp(x, y, versions))
    print "sorted_versions = "+str(sorted_versions)
    return sorted_versions[0]

def install_version(pkg, version, github, dryrun, verbose):
    repo_data = load_repo_data()

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
    if verbose:
        print "  URL: "+zipurl
    zfp = StringIO.StringIO(github.getFile(zipurl).read())
    zf = zipfile.ZipFile(zfp)
    if not dryrun:
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

    install_version(pkg, version, github,
                    dryrun=args.dry_run, verbose=args.verbose)
