install:
	python setup.py install

run_tests: install
	(cd tests; python ../test/bin/nosetests)
