#!/usr/bin/env python3
"""Frozen source/runtime oracles for six controlled public tutorial PRs.

Arguments: BINARY PINNED_FLASK_REPOSITORY NEW_OUTPUT RUNTIME_PACKAGES
The provider, callers and PR heads are controlled fixtures, not upstream PRs.
"""
import hashlib,json,os,shutil,socket,subprocess,sys,threading,time
from pathlib import Path
from http.server import BaseHTTPRequestHandler,ThreadingHTTPServer
from urllib.request import Request,urlopen
binary,source_s,out_s,packages=sys.argv[1:]
source=Path(source_s);out=Path(out_s);out.mkdir(mode=0o700)
pin=subprocess.check_output(["git","-C",str(source),"rev-parse","HEAD"],text=True).strip()
assert pin=="2c1b30d0503cfb064f1cb252e6614a06915a362a"
sys.path.insert(0,packages)
import yaml
def git(repo,*args):return subprocess.check_output(["git","-C",str(repo),*args],text=True).strip()
def commit(repo,message):
 git(repo,"add",".");git(repo,"-c","user.name=Validation","-c","user.email=validation@example.test","commit","-qm",message)
 return git(repo,"rev-parse","HEAD")
runtime_code='''import json,sys
sys.path.insert(0,sys.argv[1])
from flaskr import create_app
app=create_app({"TESTING":True})
print(json.dumps(sorted({(m,str(r)) for r in app.url_map.iter_rules() if r.endpoint!="static" for m in r.methods if m not in {"HEAD","OPTIONS"}})))
'''
cases=[]
for name in ["direct","removal","internal-helper","configuration","contract","no-impact"]:
 directory=out/name;directory.mkdir();repo=directory/"flask"
 shutil.copytree(source/"examples/tutorial",repo,ignore=shutil.ignore_patterns("__pycache__","*.pyc"))
 git(repo,"init","-q");base=commit(repo,"Pinned tutorial")
 auth=repo/"flaskr/auth.py";text=auth.read_text();target="/auth/login";method="get";expected=0
 if name=="direct":
  auth.write_text(text.replace('"/login"','"/sign-in"'));target="/auth/sign-in";expected=1
 elif name=="removal":
  auth.write_text(text[:text.index('@bp.route("/logout")')]);target="/auth/logout"
 elif name=="internal-helper":
  auth.write_text(text.replace('user_id = session.get("user_id")','user_id = session.get("account_id")'))
 elif name=="configuration":
  init=repo/"flaskr/__init__.py";init.write_text(init.read_text().replace('SECRET_KEY="dev"','SECRET_KEY="corpus-dev"'))
 elif name=="contract":
  prefix,login=text.split("def login():",1)
  auth.write_text(prefix+"def login():"+login.replace('request.form["username"]','request.form["account_name"]'));method="post";expected=1
 else:
  readme=repo/"README.rst";readme.write_text(readme.read_text()+"\nControlled documentation-only test.\n")
 head=commit(repo,"Controlled "+name)
 runtime_env={**os.environ,"PYTHONPATH":packages,"PYTHONDONTWRITEBYTECODE":"1"}
 operations=json.loads(subprocess.check_output([sys.executable,"-c",runtime_code,str(repo)],env=runtime_env,text=True))
 client=directory/"caller";client.mkdir()
 payload=', data={"username": "alice", "password": "test"}' if name=="contract" else ""
 (client/"client.py").write_text(f'import requests\n\ndef call():\n    return requests.{method}("http://flask{target}"{payload})\n')
 git(client,"init","-q");caller_head=commit(client,"Synthetic literal caller")
 filenames=git(repo,"diff","--name-only",base,head).splitlines()
 files=[{"filename":f,"status":"modified","patch":git(repo,"diff","--unified=3",base,head,"--",f),"additions":1,"deletions":1} for f in filenames]
 oracle={"case":name,"public_pin":pin,"derived_base":base,"derived_head":head,"caller_head":caller_head,"runtime_operations":operations,"expected_exact_callers":expected,"expected_caller":"caller" if expected else None,"caller_method":method.upper(),"caller_path":target,"controlled":True,"runtime_traffic_claim":False,"source_files":filenames,"patches":files}
 if name=="contract":
  form_code='import json,sys;sys.path.insert(0,sys.argv[1]);from flaskr import create_app;from flaskr.db import init_db;app=create_app({"TESTING":True,"DATABASE":sys.argv[2]});ctx=app.app_context();ctx.push();init_db();print(app.test_client().post("/auth/login",data={"username":"alice","password":"test"}).status_code);ctx.pop()'
  statuses=[]
  for label,scope in [("baseline",source/"examples/tutorial"),("head",repo)]:
   statuses.append(int(subprocess.check_output([sys.executable,"-c",form_code,str(scope),str(directory/(label+".sqlite"))],env=runtime_env,text=True)))
  assert statuses==[200,400],statuses
  oracle["form_runtime_statuses"]={"baseline":statuses[0],"head":statuses[1],"payload":{"username":"alice","password":"test"},"oracle":"Flask test_client after init_db; runtime contract behavior, not network reachability"}
 receipt=directory/"oracle.json";receipt.write_text(json.dumps(oracle,indent=2))
 cases.append((name,repo,client,head,files,oracle,hashlib.sha256(receipt.read_bytes()).hexdigest()))
current={}
class Provider(BaseHTTPRequestHandler):
 def log_message(self,*args):pass
 def do_GET(self):
  data=current["files"] if "/files" in self.path else ([current["pull"]] if "/pulls?" in self.path else current["pull"])
  self.send_response(200);self.send_header("Content-Type","application/json");self.end_headers();self.wfile.write(json.dumps(data).encode())
provider=ThreadingHTTPServer(("127.0.0.1",0),Provider);threading.Thread(target=provider.serve_forever,daemon=True).start()
with socket.socket() as s:s.bind(("127.0.0.1",0));port=s.getsockname()[1]
url=f"http://127.0.0.1:{port}";home=out/"home"
env={k:v for k,v in os.environ.items() if not k.startswith("DIFFMIND_") and k not in ("GITHUB_TOKEN","GH_TOKEN")}
env.update(DIFFMIND_HOME=str(home),DIFFMIND_BINARY=binary)
log=(out/"private-server.log").open("w")
server=subprocess.Popen([binary,"ui","--no-spa-rebuild","--host","127.0.0.1","--port",str(port),"--refresh-interval","0","--refresh-on-start=false"],env=env,stdout=log,stderr=log)
def api(route,body=None):
 with urlopen(Request(url+route,data=None if body is None else json.dumps(body).encode(),headers={"Content-Type":"application/json"}),timeout=60) as r:return json.loads(r.read())
rows=[]
try:
 for _ in range(100):
  try:
   with urlopen(url+"/healthz",timeout=1):break
  except OSError:time.sleep(.1)
 for name,repo,client,head,files,oracle,digest in cases:
  current.update(files=files,pull={"number":1,"title":"Controlled "+name,"html_url":"https://controlled.example.test/flask/pull/1","head":{"sha":head,"ref":"head"},"base":{"ref":"base"},"changed_files":len(files)})
  pid=api("/api/projects",{"name":name})["id"];baseurl="/api/projects/"+pid
  r=api(baseurl+"/repos",{"name":"flask","path":str(repo),"source_type":"local","kind":"service_repo","git_url":"https://github.com/pallets/flask.git","git_api_base":f"http://127.0.0.1:{provider.server_port}"})
  api(baseurl+"/repos",{"name":"caller","path":str(client),"source_type":"local","kind":"service_repo"})
  api(baseurl+"/ingestion",{"concurrency":2})
  for _ in range(1000):
   state=api(baseurl+"/ingestion")
   if state["status"]!="running":break
   time.sleep(.1)
  assert state["status"]=="completed",state
  run=next(p["run_id"] for p in state["repo_progress"] if p["repo_id"]=="flask")
  facts=yaml.safe_load((home/"runs"/run/"diffmind.yaml").read_text())
  actual=sorted({(o["metadata"]["details"]["method"],o["metadata"]["details"]["path"]) for o in facts["objects"]["http_endpoints"]})
  assert [list(p) for p in actual]==oracle["runtime_operations"],(name,actual,oracle["runtime_operations"])
  impact=api(baseurl+"/pull-requests/"+r["id"]+"/1/impact")
  company=impact["company"]
  (out/name/"impact.json").write_text(json.dumps(impact,indent=2))
  row={"case":name,"oracle_sha256":digest,"runtime_operations":len(actual),"runtime_agreement":True,"expected_exact":oracle["expected_exact_callers"],"actual_exact":company["direct_services"],"eligible":company["score_eligible"],"freshness":company["freshness"],"changed_entrypoints":company.get("changed_entrypoints",[]),"passed":company["direct_services"]==oracle["expected_exact_callers"]}
  rows.append(row)
  assert company["direct_services"]==oracle["expected_exact_callers"],row
  assert bool(company["score_eligible"])==bool(oracle["expected_exact_callers"]),row
  if name=="contract":
   row["classification"]="handler_body_changed_surface"
   row["runtime_contract_break"]=oracle["form_runtime_statuses"]
  assert "not proof" in " ".join(company["limitations"])
  assert git(repo,"status","--porcelain")=="" and git(client,"status","--porcelain")==""
  print(json.dumps(row),flush=True)
 (out/"results.json").write_text(json.dumps({"safety_checks_passed":True,"exact_expectations_met":all(r["passed"] for r in rows),"rows":rows,"limits":"Six controlled mutations and one literal synthetic caller per case; not downstream runtime outcomes or complete graph accuracy."},indent=2))
finally:
 (out/"observations.json").write_text(json.dumps({"rows":rows},indent=2))
 server.terminate();server.wait(timeout=20);log.close();provider.shutdown();provider.server_close()
