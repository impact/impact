from nose.tools import *

from impactlib.semver import SemanticVersion

@raises(ValueError)
def test_bogus_version1():
    SemanticVersion("1")

@raises(ValueError)
def test_bogus_version2():
    SemanticVersion("1.0")

@raises(ValueError)
def test_bogus_version3():
    SemanticVersion("v1.0")

@raises(ValueError)
def test_bogus_version4():
    SemanticVersion("v1.0+build")

@raises(ValueError)
def test_bogus_version5():
    SemanticVersion("v1.0.0build")

@raises(ValueError)
def test_bogus_version6():
    SemanticVersion("1.0+build")

@raises(ValueError)
def test_bogus_version7():
    v1 = SemanticVersion("1.0+build1")
    v2 = SemanticVersion("1.0+build2")
    assert_equal(False, v1==v2)
    
def test_parsing():
    """
    Test parsing of semantic versions
    """
    v1 = SemanticVersion("1.0.0")
    v2 = SemanticVersion("v1.0.0")
    v3 = SemanticVersion("1.0.0+build")
    v4 = SemanticVersion("v1.0.0+build")
    v5 = SemanticVersion("1.0.0-beta")
    v6 = SemanticVersion("v1.0.0-beta")
    
def test_json_rep():
    """
    Test JSON representation
    """
    v1 = SemanticVersion("1.2.3")
    assert_equal(v1.major, 1)
    assert_equal(v1.minor, 2)
    assert_equal(v1.patch, 3)
    assert_equal(v1.build, None)
    assert_equal(v1.prerelease, None)
    r1 = v1.json()
    assert_equal(r1["major"], 1)
    assert_equal(r1["minor"], 2)
    assert_equal(r1["patch"], 3)

    v2 = SemanticVersion("v1.2.3")
    assert_equal(v2.major, 1)
    assert_equal(v2.minor, 2)
    assert_equal(v2.patch, 3)
    assert_equal(v2.build, None)
    assert_equal(v2.prerelease, None)
    r2 = v2.json()
    assert_equal(r2["major"], 1)
    assert_equal(r2["minor"], 2)
    assert_equal(r2["patch"], 3)

    v3 = SemanticVersion("1.2.3+build")
    print "v3 = "+str(v3)
    assert_equal(v3.major, 1)
    assert_equal(v3.minor, 2)
    assert_equal(v3.patch, 3)
    assert_equal(v3.build, "build")
    assert_equal(v3.prerelease, None)
    r3 = v3.json()
    assert_equal(r3["major"], 1)
    assert_equal(r3["minor"], 2)
    assert_equal(r3["patch"], 3)
    assert_equal(r3["build"], "build")

    v4 = SemanticVersion("1.2.3-beta")
    assert_equal(v4.major, 1)
    assert_equal(v4.minor, 2)
    assert_equal(v4.patch, 3)
    assert_equal(v4.build, None)
    assert_equal(v4.prerelease, "beta")
    r4 = v4.json()
    assert_equal(r4["major"], 1)
    assert_equal(r4["minor"], 2)
    assert_equal(r4["patch"], 3)
    assert_equal(r4["prerelease"], "beta")

def test_comparison():
    """
    Test comparison
    """
    v1 = SemanticVersion("0.1.0")
    v2 = SemanticVersion("0.1.1")
    v3 = SemanticVersion("0.2.0")
    v4 = SemanticVersion("0.3.0-beta")
    v5 = SemanticVersion("0.3.0")
    v6 = SemanticVersion("0.3.1+build")
    v7 = SemanticVersion("0.3.2")
    v8 = SemanticVersion("1.0.0")
    v9 = SemanticVersion("1.0.1")
    v10 = SemanticVersion("1.1.0")

    assert_equal(True, v2>v1)
    assert_equal(True, v3>v2)
    assert_equal(True, v4>v3)
    assert_equal(True, v5>v4)
    assert_equal(True, v6>v5)
    assert_equal(True, v7>v6)
    assert_equal(True, v8>v7)
    assert_equal(True, v9>v8)
    assert_equal(True, v10>v9)

def test_sorting():
    v1 = SemanticVersion("1.0.0")
    v2 = SemanticVersion("0.1.0")
    v3 = SemanticVersion("0.2.0")

    x = sorted([v1, v2, v3])
    assert_equal(x, [v2, v3, v1])
