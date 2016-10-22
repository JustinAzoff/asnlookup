from setuptools import setup, find_packages

version = '0.0.1'
long_description = ""

setup(name='asnlookup',
      version=version,
      description="ASN Lookup",
      long_description=long_description,
      classifiers=[], # Get strings from http://pypi.python.org/pypi?%3Aaction=list_classifiers
      keywords='ASN',
      author='Justin Azoff',
      author_email='justin@bouncybouncy.net',
      url='',
      license='MIT',
      packages=find_packages(exclude=['ez_setup', 'examples', 'tests']),
      include_package_data=True,
      install_requires=[
          # -*- Extra requirements: -*-
          "pyasn",
      ],
      entry_points = {
        'console_scripts': [
            'asnlookup          = asnlookup.main:main',
            'asnlookup-server   = asnlookup.server:main',
            'asnlookup-client   = asnlookup.client:main',
        ]
      },
  )
