import re

class SemanticVersion(object):
    def __init__(self, text):
        self.build = None
        self.prerelease = None

        pat = """^v?([0-9]+)\.([0-9]+)\.([0-9]+)([+-].*)?$"""
        c = re.compile(pat)
        m = c.match(text)
        if m==None:
            msg = "Version number %s isn't a semantic version" % (text,)
            raise ValueError(msg)
        self.major = int(m.group(1))
        self.minor = int(m.group(2))
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
    def __str__(self):
        if self.build!=None:
            return "%d.%d.%d+%s" % (self.major, self.minor,
                                    self.patch, self.build)
        elif self.prerelease!=None:
            return "%d.%d.%d-%s" % (self.major, self.minor,
                                    self.patch, self.prerelease)
        else:
            return "%d.%d.%d" % (self.major, self.minor, self.patch)

    def __cmp__(self, other):
        if self.major==other.major:
            if self.minor==other.minor:
                if self.prerelease==None and other.prerelease==None:
                    return self.patch-other.patch
                elif self.prerelease!=None and other.prerelease==None:
                    return -1
                elif self.prerelease==None and other.prerelease!=None:
                    return 1
                else:
                    if self.build!=other.build:
                        msg = "Identical versions with different build: "
                        raise ValueError(msg+str(self)+", "+str(other))
                    return cmp(self.prerelease.split("."),
                               other.prerelease.split("."))
            else:
                return self.minor-other.minor
        else:
            return self.major-other.major
        
def parse_semver(tag):
    pat = """v?([0-9]+)\.([0-9]+)\.([0-9]+)"""
    c = re.compile(pat)
    m = c.match(tag)
    if m==None:
        return None
    major = m.group(1)
    minor = m.group(2)
    patch = m.group(3)
    return {"version": "%s.%s.%s" % (major, minor, patch),
            "major": major, "minor": minor, "patch": patch}

def semver_cmp(v1, v2, versions):
    maj1 = int(versions[v1]["major"])
    maj2 = int(versions[v2]["major"])
    if maj1==maj2:
        min1 = int(versions[v1]["minor"])
        min2 = int(versions[v2]["minor"])
        if min1==min2:
            pat1 = int(versions[v1]["patch"])
            pat2 = int(versions[v2]["patch"])
            return pat1>pat2
        else:
            return min1>min2
    else:
        return maj1>maj2
