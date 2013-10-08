import os
import json

def cache_file_name():
    return os.path.expanduser("~/.impact_cache")

def load_cache_file():
    with open(cache_file_name(), "r") as fp:
        return json.load(fp)
