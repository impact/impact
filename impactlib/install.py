import zipfile
import StringIO

from impactlib.load import load_repo_data
from impactlib.github import GitHub
from impactlib.semver import SemanticVersion

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
    keys = versions.keys()
    svs = map(lambda x: (SemanticVersion(x), x), keys)
    sorted_versions = sorted(svs, cmp=lambda x, y: x[0]>y[0])
    print "sorted_versions = "+str(sorted_versions)
    return sorted_versions[0][1]

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

def install(pkgname, verbose, username, password, token, dry_run):
    pkg_data = pkgname.split("#")
    if len(pkg_data)==1:
        pkg = pkg_data[0]
        version = None
    elif len(pkg_data)==2:
        pkg = pkg_data[0]
        version = pkg_data[1]
    else:
        raise ValueError("Package name must be of the form name[#version]")
    
    pdata = get_package(pkg)

    if pdata==None:
        return

    version = version
    if version==None:
        version = latest_version(pdata["versions"])
        if verbose:
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
    github = GitHub(username=username, password=password,
                    token=token)

    install_version(pkg, version, github,
                    dryrun=dry_run, verbose=verbose)
