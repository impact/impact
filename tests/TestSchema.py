from nose.tools import *

import json
import jsonschema

def test_sample():
    """
    Test the index schema by running sample input
    (generated from impact) through a validation.
    """
    fp = open("index.schema", "r")
    schema = json.load(fp)
    jsonschema.Draft4Validator.check_schema(schema)
    val = jsonschema.Draft4Validator(schema)
    dfp = open("tests/sample.json", "r")
    data = json.load(dfp)
    val.validate(data)
    print "Data validated!"
