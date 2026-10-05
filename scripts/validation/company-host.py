#!/usr/bin/env python3
"""Run isolated installed-host company role and revocation trials."""
import hashlib,json,os,secrets,socket,subprocess,sys,time
from pathlib import Path
from urllib.request import Request,urlopen
from urllib.error import HTTPError
binary,root_s,runner=sys.argv[1:]
root=Path(root_s);out=root/("company-"+Path(binary).parent.name);out.mkdir(mode=0o700)
with socket.socket() as s:s.bind(("127.0.0.1",0));port=s.getsockname()[1]
url=f"http://127.0.0.1:{port}";admin=secrets.token_hex(32)
env={k:v for k,v in os.environ.items() if not k.startswith("DIFFMIND_") and k not in ("GITHUB_TOKEN","GH_TOKEN")}
env.update(DIFFMIND_HOME=str(out/"home"),DIFFMIND_BINARY=binary,DIFFMIND_AUTH_TOKEN=admin)
log=(out/"private-server.log").open("w")
server=subprocess.Popen([binary,"ui","--no-spa-rebuild","--host","127.0.0.1","--port",str(port),"--project-access","scoped","--refresh-interval","0","--refresh-on-start=false"],env=env,stdout=log,stderr=log)
def request(route,method="GET",body=None,token=None,extra=None):
 headers={"Authorization":"Bearer "+(token or admin),"Content-Type":"application/json","Accept":"application/json, text/event-stream","MCP-Protocol-Version":"2025-11-25"}
 headers.update(extra or {})
 try:
  with urlopen(Request(url+route,method=method,data=None if body is None else json.dumps(body).encode(),headers=headers),timeout=30) as r:
   raw=r.read().decode()
   if raw.startswith("event:") or raw.startswith("data:"):raw=next(x[5:].strip() for x in raw.splitlines() if x.startswith("data:"))
   return r.status,json.loads(raw) if raw else None,dict(r.headers)
 except HTTPError as e:return e.code,{},{}
def api(route,method="GET",body=None):
 status,data,_=request(route,method,body);assert status<300,(route,status,data);return data
def ingest(pid):
 api(f"/api/projects/{pid}/ingestion","POST",{"concurrency":2})
 for _ in range(1000):
  state=api(f"/api/projects/{pid}/ingestion")
  if state["status"]!="running":assert state["status"]=="completed",state;return state
  time.sleep(.15)
 raise RuntimeError("ingestion timeout")
def trial(name,prompt,token,approved=False):
 spec={"id":name,"binary":binary,"mode":"http","url":url+"/mcp","token_env":"DIFFMIND_TRIAL_TOKEN","work":str(out/"work"),"prompt":prompt,"timeout_seconds":300}
 if approved:spec["approved_tools"]=["manage_workspace"]
 file=out/(name+"-spec.json");file.write_text(json.dumps(spec,indent=2))
 result=subprocess.run([sys.executable,runner,str(file),str(out/name)],env={**env,"DIFFMIND_TRIAL_TOKEN":token})
 assert result.returncode==0,(name,result.returncode)
 return json.loads((out/name/"result.json").read_text())
rows=[]
try:
 for _ in range(100):
  try:
   with urlopen(url+"/healthz",timeout=1):break
  except OSError:time.sleep(.1)
 project=api("/api/projects","POST",{"name":"Demo Commerce"});pid=project["id"]
 hidden=api("/api/projects","POST",{"name":"Hidden Payroll Sentinel 7f29"})
 for repo in sorted((root/"demo/repositories").iterdir()):
  api(f"/api/projects/{pid}/repos","POST",{"name":repo.name,"path":str(repo),"source_type":"local","kind":"service_repo"})
 baseline=ingest(pid)
 viewer=api(f"/api/v1/projects/{pid}/tokens","POST",{"name":"Automated viewer","role":"viewer","expires_in_seconds":3600})
 editor=api(f"/api/v1/projects/{pid}/tokens","POST",{"name":"Automated editor","role":"editor","expires_in_seconds":3600})
 assert request("/api/projects/"+hidden["id"],token=viewer["secret"])[0]==404
 trial("U02-viewer","I joined Demo Commerce. Using the saved company context, explain who calls catalog and checkout and the orders.created queue relationship. Can this connection update the context or change repository scope? Include evidence limitations. Do not make changes.",viewer["secret"])
 rows.append({"case":"viewer","hidden_project_http_status":404,"project":pid,"baseline_run":baseline["graph_run_id"]})
 # Pre-existing persistent MCP connection must recheck token state.
 init=request("/mcp","POST",{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"revocation-control","version":"1"}}},viewer["secret"])
 assert init[0]==200,init[0]
 session=next((v for k,v in init[2].items() if k.lower()=="mcp-session-id"),None)
 assert session,"stateful MCP session absent"
 source=root/"demo/repositories/catalog/app.py"
 before=source.read_bytes();source.write_bytes(before+b"\n# Automated return-session refresh observation.\n")
 changed=hashlib.sha256(source.read_bytes()).hexdigest()
 trial("U02-editor","The approved Demo Commerce catalog repository changed since the saved context was built. Bring the existing context up to date, preserve local edits and repository scope, and report the completed result and remaining limitations. Do not change company configuration.",editor["secret"],True)
 assert hashlib.sha256(source.read_bytes()).hexdigest()==changed
 after=api(f"/api/projects/{pid}/ingestion")
 assert after["status"]=="completed" and after["graph_run_id"]!=baseline["graph_run_id"],after
 assert len(api(f"/api/projects/{pid}/repos")["repos"])==6
 assert request(f"/api/projects/{pid}","PATCH",{"instruction":"unauthorized"},editor["secret"])[0]==403
 rows.append({"case":"editor","completed_run":after["graph_run_id"],"analyzed":after.get("analyzed"),"reused":after.get("reused"),"local_edit_preserved":True,"registrations":6,"configuration_denied":403})
 policy=api(f"/api/v1/projects/{pid}/access")
 api(f"/api/v1/projects/{pid}/access","PUT",{"revision":policy["revision"],"members":{"test-joiner":"viewer"}})
 policy=api(f"/api/v1/projects/{pid}/access")
 api(f"/api/v1/projects/{pid}/access","PUT",{"revision":policy["revision"],"members":{}})
 assert request(f"/api/projects/{pid}",token=viewer["secret"])[0]==200
 api(f"/api/v1/projects/{pid}/tokens/{viewer['token']['id']}/revoke","POST",{})
 denied=request("/mcp","POST",{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"list_projects","arguments":{}}},viewer["secret"],{"Mcp-Session-Id":session})[0]
 assert denied==401,denied
 trial("U02-revoked","Use my company connection to show the saved Demo Commerce architecture. If the connection cannot retrieve it, explain what you actually verified and do not invent the architecture.",viewer["secret"])
 assert request(f"/api/projects/{pid}",token=editor["secret"])[0]==200
 rows.append({"case":"revocation","membership_removal_does_not_revoke_service_token":True,"existing_session_status":denied,"editor_independent":True})
 (out/"controls.json").write_text(json.dumps({"passed":True,"rows":rows},indent=2))
 print(json.dumps({"company_controls":"passed","rows":rows}),flush=True)
finally:
 server.terminate();server.wait(timeout=20);log.close()
