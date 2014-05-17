import os
import sys
# ConfigParser has been renamed in Python 3 to configparser:
try:
    from configparser import SafeConfigParser, NoSectionError, NoOptionError
except ImportError:
    from ConfigParser import SafeConfigParser, NoSectionError, NoOptionError

ENVVAR = "IMPACT_CONFIG_FILE"

def get_config_file():
    """
    All this complexity is here because @dietmarw had to be all "you
    should really follow standards" and like "don't clutter up my home
    directory with all sorts of config files".

    I got tired of all his whining and implemented this.  If you don't
    like it, blame him.
    """
    if ENVVAR in os.environ:
        return os.environ[ENVVAR]
    if sys.platform=="win32":
        datadir = os.environ.get("APPDATA",
                                 os.path.expanduser("~/.config"))
    elif sys.platform=="darwin":
        datadir = os.path.expanduser("~/Library/Preferences")
    else:
        datadir = os.environ.get("XDG_CONFIG_HOME",
                                 os.path.expanduser("~/.config"))
    return os.path.join(datadir, "impact", "impactrc")

def get(section, option, default=None):
    config = SafeConfigParser({})
    filename = get_config_file()
    config.read([filename])
    try:
        ret = config.get(section, option)
    except NoSectionError:
        ret = None
    except NoOptionError:
        ret = None
    if ret==None:
        return default
    else:
        return ret

def get_indices():
    """
    This is a user level setting for specifying where to search
    for the *processed indices*.  When using the "refresh" command
    to **create** such an index, the `source_list` open is used.
    """
    repo_list = get("Impact", "indices",
                    "https://impact.modelica.org/impact_data.json")
    return repo_list.split(",")

##
## Configuration Options
##
## [Impact]
## indices=url1,url2
## token=API_token_from_GitHub
## username=GitHub username
## password=GitHub password
## source_list=github://user/repo_pattern,github://user/repo_pattern
##
