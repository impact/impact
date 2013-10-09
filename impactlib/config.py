from ConfigParser import SafeConfigParser, NoSectionError, NoOptionError
import os

def get(section, option, default=None):
    config = SafeConfigParser({})
    filename = os.path.expanduser("~/.impactrc")
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

def get_repos():
    repo_list = get("Impact", "repos",
                    "https://xogeny-public.s3.amazonaws.com/impact_data.json")
    return repo_list.split(",")