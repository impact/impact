import re

class SemanticVersion(object):
    def __init__(self, text, tolerant=False):
        self.build = None
        self.prerelease = None

        pat = """^v?([0-9]+)\.([0-9]+)\.([0-9]+)([+-].*)?$"""
        c = re.compile(pat)
        m = c.match(text)
        if m==None and tolerant:
            pat = """^v?([0-9]+)\.([0-9]+)()([+-].*)?$"""
            c = re.compile(pat)
            m = c.match(text)
        if m==None:
            msg = "Version number %s isn't a semantic version" % (text,)
            raise ValueError(msg)
        self.major = int(m.group(1))
        self.minor = int(m.group(2))
        if m.group(3)=="":
            self.patch = 0 # We get here if we are tolerant
        else:
            self.patch = int(m.group(3))
        if m.group(4)!=None and m.group(4)[0]=="+":
            self.build = m.group(4)[1:]
        elif m.group(4)!=None and m.group(4)[0]=="-":
            self.prerelease = m.group(4)[1:]
    def json(self):
        ret = {"major": self.major,
                "minor": self.minor,
                "patch": self.patch}
        if self.build!=None:
            ret["build"] = self.build
        if self.prerelease!=None:
            ret["prerelease"] = self.prerelease
        ret["version"] = str(self)
        return ret
    def __repr__(self):
        return str(self)
    def __str__(self):
        if self.build!=None:
            return "%d.%d.%d+%s" % (self.major, self.minor,
                                    self.patch, self.build)
        elif self.prerelease!=None:
            return "%d.%d.%d-%s" % (self.major, self.minor,
                                    self.patch, self.prerelease)
        else:
            return "%d.%d.%d" % (self.major, self.minor, self.patch)

    def __lt__(self, other):
        if self.major==other.major:
            if self.minor==other.minor:
                if self.prerelease==None and other.prerelease==None:
                    return self.patch < other.patch
                elif self.prerelease!=None and other.prerelease==None:
                    return -1
                elif self.prerelease==None and other.prerelease!=None:
                    return 1
                else:
                    if self.build!=other.build:
                        msg = "Identical versions with different build: "
                        raise ValueError(msg+str(self)+", "+str(other))
                    return (self.prerelease.split(".") <
                               other.prerelease.split("."))
            else:
                return self.minor < other.minor
        else:
            return self.major < other.major
