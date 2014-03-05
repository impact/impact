from setuptools import setup
from os.path import join, dirname

setup(name="impact",
      version="0.5.2",
      description="Modelica package manager",
      long_description=open(join(dirname(__file__), 'README.md')).read(),
      author="Michael Tiller",
      author_email="michael.tiller@gmail.com",
      license="MIT",
      url="https://github.com/xogeny/impact",
      scripts=['scripts/impact'],
      packages=['impactlib'],
      include_package_data=True,
      zip_safe=False)
