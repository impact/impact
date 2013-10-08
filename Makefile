run_tests:
	python setup.py install
	(cd tests; python ../test/bin/nosetests)
