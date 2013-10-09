import StringIO
from nose.tools import *

from impactlib.refresh import extract_dependencies

def test_complex1():
    fp = StringIO.StringIO("""
package RealTimeCoordinationLibrary
annotation (uses(Modelica(version="3.2"), RealTimeCoordinationLibrary(version=
            "1.0.2"),
      Modelica_StateGraph2(version="2.0.1")),
    preferredView="info",
    version="1.0.2",
    versionBuild=1,
    versionDate="2013-04-04",
    dateModified = "2012-04-04",
    revisionId="$Id:: package.mo 1 2013-04-04 16:18:47Z #$",
    Documentation(info="<html>
<p><b>RealTimeCoordinationLibrary</b> is a <b>free</b> Modelica package providing components to model <b>real-time</b>, <b>reactive</b>, <b>hybrid</b> and, <b>asynchronous communicating</b> systems in a convenient way with <b>statecharts</b>.</p>
<p>For an introduction, have especially a look at: </p>
<p><ul>
<li><a href=\"modelica://RealTimeCoordinationLibrary.UsersGuide.Elements\">Elements</a> provide an overview of the library inside the User&apos;s Guide.</li>
<li><a href=\"modelica://RealTimeCoordinationLibrary.Examples\">Examples</a> provide simple introductory examples as well as involved application examples. </li>
</ul></p>
<p>For an application example have a look at: <a href=\"modelica://RealTimeCoordinationLibrary.Examples.Application.BeBotSystem\">BeBotSystem</a> </p>
<p><br/><b>Licensed under the Modelica License 2</b></p>
<p><i>This Modelica package is <u>free</u> software and the use is completely at <u>your own risk</u>; it can be redistributed and/or modified under the terms of the Modelica license 2, see the license conditions (including the disclaimer of warranty) <a href=\"modelica://RealTimeCoordinationLibrary.UsersGuide.ModelicaLicense2\">here</a> or at <a href=\"http://www.Modelica.org/licenses/ModelicaLicense2\">http://www.Modelica.org/licenses/ModelicaLicense2</a>.</i> </p>
</html>", revisions="<html>
<p>Name: RealTimeCoordinationLibrary</p>
<p>Path: RealTimeCoordinationLibrary</p>
<p>Version: 1.0.2, 2013-04-04, build 1 (2013-04-04)</p>
<p>Uses:Modelica (version=&QUOT;3.2&QUOT;), RealTimeCoordinationLibrary (version=&QUOT;1.0.2&QUOT;), Modelica_StateGraph2 (version=&QUOT;2.0.1&QUOT;)</p>
</html>"));
end RealTimeCoordinationLibrary;
""")
    deps = extract_dependencies(fp)
    print "deps = "+str(deps)
    assert_equal([("Modelica", "3.2"),
                  ("RealTimeCoordinationLibrary", "1.0.2"),
                  ("Modelica_StateGraph2", "2.0.1")], deps)
                   
