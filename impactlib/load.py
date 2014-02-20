import urllib2
import json

from impactlib import config

cached_data = None

### I stole this from http://bugs.python.org/issue11220 ###
import httplib, ssl, socket

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
        except ssl.SSLError, e:
            print("Trying SSLv23.")
            self.sock = ssl.wrap_socket(sock, self.key_file, self.cert_file,
                                        ssl_version=ssl.PROTOCOL_SSLv23)
            
class HTTPSHandlerV3(urllib2.HTTPSHandler):
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
            urllib2.install_opener(urllib2.build_opener(HTTPSHandlerV3()))
            req = urllib2.Request(url)
            response = urllib2.urlopen(req)
            data = json.loads(response.read())
            ret.update(data)
        except Exception as e:
            print "Unable to load repo data from: "+str(url)+", skipping"
            print "Error: "+str(e)

    cached_data = ret
    return ret
