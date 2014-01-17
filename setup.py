from setuptools import setup
from os.path import join, dirname

setup(name="impact",
      version="0.2.2",
      description="Modelica package manager",
      long_description=open(join(dirname(__file__), 'README.md')).read(),
      author="Michael Tiller",
      author_email="michael.tiller@gmail.com",
      license="MIT",
      url="http://www.xogeny.com/",
      scripts=['scripts/impact.py'],
      packages=['impactlib'],
      include_package_data=True,
      zip_safe=False)
