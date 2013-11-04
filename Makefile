# This is really designed to run inside a virtual environment.
# It is strongly advised to do:
# $ virtualenv test --no-site-packages
# $ source test/bin/activate
# $ pip install nose
# $ make run_tests
install:
	python setup.py install

venv:
	virtualenv test --no-site-packages

run_tests: install
	(cd tests; nosetests)
