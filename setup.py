from setuptools import setup
from os.path import join, dirname

setup(name="impact",
      version="0.5.7",
      description="Modelica package manager",
      long_description=open(join(dirname(__file__), 'README.md')).read(),
      author="Michael Tiller",
      author_email="michael.tiller@gmail.com",
      license="MIT",
      url="https://github.com/xogeny/impact",
      entry_points = {
          'console_scripts': ['impact = impactlib.cli:main']
      },
      packages=['impactlib'],
      extras_require = {
          'color': ['colorama']
      },
      include_package_data=True,
      zip_safe=False)
