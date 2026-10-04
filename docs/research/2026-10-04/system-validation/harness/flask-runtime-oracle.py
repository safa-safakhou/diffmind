"""Run with pinned Flask 3.1.2 and PyYAML available; never starts a web server."""
import json, subprocess, sys
from pathlib import Path
source, artifact, output = map(Path, sys.argv[1:])
revision = subprocess.check_output(["git", "-C", str(source), "rev-parse", "HEAD"], text=True).strip()
assert revision == "2c1b30d0503cfb064f1cb252e6614a06915a362a"
sys.path.insert(0, str(source))
import flask, yaml
from flaskr import create_app
app = create_app({"TESTING": True})
expected = {(method, str(rule)) for rule in app.url_map.iter_rules()
            if rule.endpoint != "static" for method in rule.methods
            if method not in {"HEAD", "OPTIONS"}}
data = yaml.safe_load(artifact.read_text())
actual = {(o["metadata"]["details"]["method"], o["metadata"]["details"]["path"])
          for o in data["objects"]["http_endpoints"]}
result = {"source_commit": revision, "oracle": "actual Flask 3.1.2 URL map; unique explicit method/path; static and implicit HEAD/OPTIONS excluded",
          "expected_unique": sorted(expected), "actual_unique": sorted(actual),
          "missing": sorted(expected-actual), "unexpected": sorted(actual-expected),
          "independent_human_reviewers": 0}
output.write_text(json.dumps(result, indent=2) + "\n")
assert expected == actual
print("12 unique explicit runtime operations match extraction")
