import zipfile
import io
import tempfile
import shutil
import os

from impactlib.load import load_repo_data
from impactlib.refresh import strip_extra
from impactlib.github import GitHub
from impactlib.semver import SemanticVersion
from impactlib import config

try:
    import colorama
    colorama.init()
    use_color = True
except:
    use_color = False

def get_package(pkg):
    repo_data = load_repo_data()
    if not pkg in repo_data:
        msg = "No package named '"+pkg+"' found"
        if use_color:
            print((colorama.Fore.RED+msg))
        else:
            print(msg)
        return None
    return repo_data[pkg]

def latest_version(versions):
    if len(versions)==0:
        return None
    keys = list(versions.keys())
    svs = [(SemanticVersion(x, tolerant=True), x) for x in keys]
    sorted_versions = sorted(svs, reverse=True)
    return sorted_versions[0][1]

def install_version(pkg, version, github, dryrun, verbose, target):
    # repo_data = load_repo_data()

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
            print(colorama.Fore.RED+msg)
        else:
            print(msg)
        return

    zipurl = vdata["zipball_url"]
    vpath = vdata["path"]
    if verbose:
        print("  URL: "+zipurl)
    if not dryrun:
        zfp = io.StringIO(github.getDownload(zipurl).read())
        zf = zipfile.ZipFile(zfp)
        root = zf.infolist()[0].filename
        dst = os.path.join(target, str(pkg)+" "+str(strip_extra(version)))
        if os.path.exists(dst):
            print("  Directory "+dst+" already exists, skipping")
        else:
            td = tempfile.mkdtemp()
            zf.extractall(td)
            src = os.path.join(td, root, vpath)
            if verbose:
                print("  Root zip directory: "+root)
                print("  Temp directory: "+str(td))
                print("  Version path: "+str(vpath))
                print("  Source: "+str(src))
                print("  Destination: "+str(dst))
            shutil.copytree(src,dst)
            shutil.rmtree(td)

def elaborate_dependencies(pkgname, version, current):
    repo_data = load_repo_data()
    if not pkgname in repo_data:
        print("  No information for package "+pkgname+", skipping")
        return current
    if not version in repo_data[pkgname]["versions"]:
        print("  No version "+version+" of package "+pkgname+" found, skipping")
        return current
    ret = current.copy()
    ret[pkgname] = version
    vdata = repo_data[pkgname]["versions"][version]
    deps = vdata["dependencies"]
    for dep in deps:
        dname = dep["name"]
        dver = dep["version"]
        if dname in ret:
            if dver==ret[dname]:
                # This could avoid circular dependencies?
                continue
            else:
                raise NameError("Dependency on version %s and %s of %s" % \
                                    (ret[dname], dver, dname))
        subs = elaborate_dependencies(dname, dver, ret)
        for sub in subs:
            if sub in ret:
                if subs[sub]==ret[sub]:
                    continue
                else:
                    raise NameError("Dependency on version %s and %s of %s" % \
                                        (sub[sub], ret[sub], sub))
            ret[sub] = subs[sub]
    return ret

def install(pkgname, verbose, dry_run, target):
    username = config.get("Impact", "username", None)
    password = config.get("Impact", "password", None)
    token = config.get("Impact", "token", None)

    if "#" in pkgname:
        pkg_data = pkgname.split("#")
    else:
        pkg_data = pkgname.split(" ")
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
            print("  Choosing latest version: "+version)
        if version==None:
            msg = "No (semantic) versions found for package '"+pkg+"'"
            if use_color:
                print(colorama.Fore.RED+msg)
            else:
                print(msg)
            return

    msg = "Installing version '"+version+"' of package '"+pkg+"'"
    if use_color:
        print(colorama.Fore.GREEN+msg)
    else:
        print(msg)

    # Setup connection to github
    github = GitHub(username=username, password=password,
                    token=token)

    pkgversions = elaborate_dependencies(pkg, version, current={})

    if verbose:
        print("Libraries to install:")
        for pkgname in pkgversions:
            print("  "+pkgname+" version "+pkgversions[pkgname])
        print("Installation...")

    for pkgname in pkgversions:
        install_version(pkgname, pkgversions[pkgname], github,
                        dryrun=dry_run, verbose=verbose, target=target)
