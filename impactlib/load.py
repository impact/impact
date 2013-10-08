import urllib2
import json

from impactlib import config

def load_repo_data():
    urls = config.get_repos()
    ret = {}
    # Ideally, we could cache this data if we cannot immediately
    # fetch it.  But it isn't too much good since it would only
    # allow you to search.  Install will require network connectivity
    # anyway.
    for url in urls:
        print "Loading "+str(url)
        req = urllib2.Request(url)
        response = urllib2.urlopen(req)
        data = json.loads(response.read())
        ret.update(data)

    return ret
