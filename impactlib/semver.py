import re

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
