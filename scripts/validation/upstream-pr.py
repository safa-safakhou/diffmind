#!/usr/bin/env python3
"""Actual upstream Flask PR 4139 at its immutable head, with a synthetic caller.
Usage: upstream-pr.py BINARY CHECKED_OUT_FLASK_HEAD NEW_OUTPUT [HOST_RUNNER]
Only reads the public API. Does not mutate upstream or the supplied checkout.
"""
import hashlib,json,os,socket,subprocess,sys,time
from pathlib import Path
from urllib.request import Request,urlopen
binary,source_s,out_s=sys.argv[1:4];source=Path(source_s);out=Path(out_s);out.mkdir(mode=0o700)
head="51196575479e34c3ea43612e1eb770db3aa5d114"
def git(path,*args):return subprocess.check_output(["git","-C",str(path),*args],text=True).strip()
assert git(source,"rev-parse","HEAD")==head and not git(source,"status","--porcelain")
def public(path):
 with urlopen(Request("https://api.github.com/repos/pallets/flask/"+path,headers={"Accept":"application/vnd.github+json","User-Agent":"diffmind-validation"}),timeout=30) as r:return json.load(r)
pull=public("pulls/4139");files=public("pulls/4139/files")
assert pull["head"]["sha"]==head
auth=next(f for f in files if f["filename"]=="examples/tutorial/flaskr/auth.py")
assert "db.IntegrityError" in auth["patch"]
caller=out/"caller";caller.mkdir();(caller/"client.py").write_text('import requests\n\ndef register():\n    return requests.post("http://flask/auth/register", data={"username": "alice", "password": "test"})\n')
git(caller,"init","-q");git(caller,"add",".");git(caller,"-c","user.name=Validation","-c","user.email=validation@example.test","commit","-qm","Synthetic literal registration caller")
oracle={"upstream_pr":pull["html_url"],"title":pull["title"],"head":head,"caller_head":git(caller,"rev-parse","HEAD"),"caller_synthetic":True,"expected_exact_callers":["caller"],"operation":"POST /auth/register","reason":"Actual PR changes registration handler implementation; the literal caller uses this method and path.","changed_files":[f["filename"] for f in files],"runtime_break_claim":False}
(out/"oracle.json").write_text(json.dumps(oracle,indent=2))
with socket.socket() as s:s.bind(("127.0.0.1",0));port=s.getsockname()[1]
url=f"http://127.0.0.1:{port}";home=out/"home"
env={k:v for k,v in os.environ.items() if not k.startswith("DIFFMIND_") and k not in ["GITHUB_TOKEN","GH_TOKEN"]}
env.update(DIFFMIND_HOME=str(home),DIFFMIND_BINARY=binary)
log=(out/"private-server.log").open("w")
server=subprocess.Popen([binary,"ui","--no-spa-rebuild","--host","127.0.0.1","--port",str(port),"--refresh-interval","0","--refresh-on-start=false"],env=env,stdout=log,stderr=log)
def api(route,body=None):
 with urlopen(Request(url+route,data=None if body is None else json.dumps(body).encode(),headers={"Content-Type":"application/json"}),timeout=90) as r:return json.load(r)
try:
 for _ in range(100):
  try:
   with urlopen(url+"/healthz",timeout=1):break
  except OSError:time.sleep(.1)
 pid=api("/api/projects",{"name":"Upstream Flask PR"})["id"];base="/api/projects/"+pid
 rid=api(base+"/repos",{"name":"flask","path":str(source),"git_url":"https://github.com/pallets/flask.git","source_type":"local","kind":"service_repo"})["id"]
 api(base+"/repos",{"name":"caller","path":str(caller),"source_type":"local","kind":"service_repo"})
 api(base+"/ingestion",{"concurrency":2})
 for _ in range(1200):
  state=api(base+"/ingestion")
  if state["status"]!="running":break
  time.sleep(.1)
 assert state["status"]=="completed",state
 impact=api(base+"/pull-requests/"+rid+"/4139/impact");(out/"impact.json").write_text(json.dumps(impact,indent=2))
 company=impact["company"]
 assert company["freshness"]=="fresh" and company["score_eligible"] and company["direct_services"]==1,company
 assert not git(source,"status","--porcelain") and not git(caller,"status","--porcelain")
 result={"passed":True,"binary_sha256":hashlib.sha256(Path(binary).read_bytes()).hexdigest(),"upstream_pr":pull["html_url"],"head":head,"freshness":company["freshness"],"exact_callers":company["direct_services"],"changed_entrypoints":company["changed_entrypoints"],"caller_synthetic":True,"source_preserved":True,"delivery":impact["delivery"],"oracle_sha256":hashlib.sha256((out/"oracle.json").read_bytes()).hexdigest()}
 (out/"results.json").write_text(json.dumps(result,indent=2));print(json.dumps(result),flush=True)
 if len(sys.argv)>4:
  spec={"id":"U03-upstream-exact","binary":binary,"mode":"http","url":url+"/mcp","work":str(out/"model-work"),"timeout_seconds":600,"prompt":"Review pallets/flask PR 4139 using Upstream Flask PR company context. Explain which other services may be affected, what the source evidence supports, its limitations, and whether your review is posted on the PR. Inspect only and preserve local source."}
  (out/"spec.json").write_text(json.dumps(spec,indent=2))
  subprocess.run([sys.executable,sys.argv[4],str(out/"spec.json"),str(out/"model")],env=env,check=True)
finally:
 server.terminate();server.wait(timeout=20);log.close()
