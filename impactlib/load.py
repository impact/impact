import json
# urllib2 is now split into several urllib modules in python3
try:
    from urllib.request import (urlopen, Request, HTTPSHandler,
                                install_opener, build_opener)
except ImportError:
    from urllib2 import (urlopen, Request, HTTPSHandler,
                         install_opener, build_opener)

from impactlib import config

cached_data = None

### I stole this from http://bugs.python.org/issue11220 ###
try:
    import http.client as httplib
except ImportError:
    import httplib

import ssl, socket

class HTTPSConnectionV3(httplib.HTTPSConnection):
    def __init__(self, *args, **kwargs):
        httplib.HTTPSConnection.__init__(self, *args, **kwargs)

    def connect(self):
        sock = socket.create_connection((self.host, self.port), self.timeout)
        if self._tunnel_host:
            self.sock = sock
            self._tunnel()
        try:
            self.sock = ssl.wrap_socket(sock, self.key_file, self.cert_file,
                                        ssl_version=ssl.PROTOCOL_SSLv3)
        except ssl.SSLError as e:
            print("Trying SSLv23.")
            self.sock = ssl.wrap_socket(sock, self.key_file, self.cert_file,
                                        ssl_version=ssl.PROTOCOL_SSLv23)

class HTTPSHandlerV3(HTTPSHandler):
    def https_open(self, req):
        return self.do_open(HTTPSConnectionV3, req)
######

def load_repo_data():
    global cached_data

    # If we've already called this function once in this process, use
    # the previously returned result.
    if cached_data!=None:
        return cached_data

    urls = config.get_indices()
    ret = {}
    # Ideally, we could cache this data if we cannot immediately
    # fetch it.  But it isn't too much good since it would only
    # allow you to search.  Install will require network connectivity
    # anyway.
    for url in urls:
        try:
            install_opener(build_opener(HTTPSHandlerV3()))
            req = Request(url)
            response = urlopen(req)
            data = json.loads(response.read().decode(encoding='utf8'))
            ret.update(data)
        except Exception as e:
            print("Unable to load repo data from: "+str(url)+", skipping")
            print("Error: "+str(e))

    cached_data = ret
    return ret
