#!/usr/bin/env python3
"""Private real-binary prior-candidate recovery drill; no ordinary ingress."""
import hashlib,json,os,secrets,socket,subprocess,sys,time
from pathlib import Path
from urllib.request import Request,urlopen
from urllib.error import HTTPError
old_binary,new_binary,root_text,source_text=sys.argv[1:]
root=Path(root_text).resolve(); root.mkdir(mode=0o700)
source=Path(source_text).resolve()
clean={k:v for k,v in os.environ.items() if not k.startswith("DIFFMIND_") and k not in ("GITHUB_TOKEN","GH_TOKEN")}
rows=[]
for backend in ("json","sqlite"):
 home=root/backend;home.mkdir(mode=0o700)
 with socket.socket() as sock:
  sock.bind(("127.0.0.1",0));port=sock.getsockname()[1]
 url="http://127.0.0.1:"+str(port); proxy=secrets.token_hex(32);recovery=secrets.token_hex(32)
 env={**clean,"DIFFMIND_HOME":str(home),"DIFFMIND_TRUSTED_PROXY_SECRET":proxy,"DIFFMIND_AUTH_TOKEN":recovery}
 process=None
 log=open(root/(backend+"-private-server.log"),"w")
 def start(binary):
  global process
  process=subprocess.Popen([binary,"ui","--no-spa-rebuild","--host","127.0.0.1","--port",str(port),"--project-access","scoped","--refresh-interval","0"],env={**env,"DIFFMIND_BINARY":binary},stdout=log,stderr=log)
  for _ in range(100):
   try:
    with urlopen(url+"/healthz",timeout=1) as response:
     if response.status==200:return
   except OSError:pass
   time.sleep(.05)
  raise RuntimeError("startup timeout")
 def stop():
  global process
  if process:
   process.terminate();process.wait(timeout=20);process=None
 def request(route,method="GET",body=None,user="admin",token=None):
  headers={"Content-Type":"application/json"}
  if token:headers["Authorization"]="Bearer "+token
  else:headers.update({"X-DiffMind-Proxy-Secret":proxy,"X-DiffMind-User":user,"X-DiffMind-Role":"admin" if user=="admin" else "viewer"})
  req=Request(url+route,data=None if body is None else json.dumps(body).encode(),headers=headers,method=method)
  try:
   with urlopen(req,timeout=30) as response:return response.status,json.load(response)
  except HTTPError as error:return error.code,json.load(error)
 def api(*args,**kwargs):
  status,data=request(*args,**kwargs);assert status<300,status;return data
 def command(binary,args,success=True):
  result=subprocess.run([binary,*args],env=env,capture_output=True,text=True,timeout=60)
  assert (result.returncode==0)==success,"unexpected CLI outcome"
  return json.loads(result.stdout) if success and "--json" in args else None
 try:
  start(old_binary)
  p=api("/api/projects","POST",{"name":"Prior candidate recovery "+backend});pid=p["id"];base="/api/projects/"+pid
  api(base+"/repos","POST",{"name":"frontend","path":str(source/"frontend"),"source_type":"local","kind":"service_repo"})
  api(base+"/ingestion","POST",{"concurrency":1})
  for _ in range(300):
   state=api(base+"/ingestion")
   if state["status"]!="running":break
   time.sleep(.1)
  assert state["status"]=="completed"
  run=state["graph_run_id"];graph_file=home/"projects"/pid/"runs"/run/"graph.json";graph_hash=hashlib.sha256(graph_file.read_bytes()).hexdigest()
  policy=api("/api/v1/projects/"+pid+"/access")
  policy=api("/api/v1/projects/"+pid+"/access","PUT",{"revision":policy["revision"],"members":{"departed":"viewer"}})
  assigned=api("/api/v1/projects/"+pid+"/tokens","POST",{"name":"assigned","role":"viewer","expires_in_seconds":3600})
  service=api("/api/v1/projects/"+pid+"/tokens","POST",{"name":"unrelated service","role":"viewer","expires_in_seconds":3600})
  stop()
  if backend=="sqlite":command(old_binary,["storage","migrate","--offline","--json"])
  archive=root/(backend+".tar.gz");t=time.monotonic()
  report=command(old_binary,["backup","create","--offline","--output",str(archive),"--json"]);create_seconds=time.monotonic()-t
  assert archive.stat().st_mode&0o777==0o600
  command(new_binary,["backup","verify","--archive",str(archive),"--sha256",report["sha256"],"--json"])
  start(old_binary)
  api("/api/v1/projects/"+pid+"/access","PUT",{"revision":policy["revision"],"members":{}})
  api("/api/v1/projects/"+pid+"/tokens/"+assigned["token"]["id"]+"/revoke","POST",{})
  assert request(base,token=assigned["secret"])[0]==401
  stop()
  rollback=home.with_name(home.name+"-rollback");home.rename(rollback)
  t=time.monotonic();command(new_binary,["backup","restore","--offline","--archive",str(archive),"--destination",str(home),"--sha256",report["sha256"],"--json"]);restore_seconds=time.monotonic()-t
  command(new_binary,["backup","restore","--offline","--archive",str(archive),"--destination",str(home)],success=False)
  if backend=="sqlite":command(new_binary,["storage","verify","--offline","--json"])
  start(new_binary)
  assert hashlib.sha256(graph_file.read_bytes()).hexdigest()==graph_hash
  assert request(base,token=assigned["secret"])[0]==200
  assert request(base,user="departed")[0]==200
  restored=api("/api/v1/projects/"+pid+"/access")
  api("/api/v1/projects/"+pid+"/access","PUT",{"revision":restored["revision"],"members":{}})
  api("/api/v1/projects/"+pid+"/tokens/"+assigned["token"]["id"]+"/revoke","POST",{})
  assert request(base,token=assigned["secret"])[0]==401
  assert request(base,user="departed")[0]==404
  assert request(base,token=service["secret"])[0]==200
  assert request(base,token=recovery)[0]==200
  command(new_binary,["backup","create","--offline","--output",str(root/(backend+"-while-live.tar.gz"))],success=False)
  rows.append({"backend":backend,"prior_graph_bytes_preserved":True,"rollback_copy_retained":rollback.exists(),"resurrection_reproduced":True,"membership_and_assigned_token_reconciled":True,"unrelated_service_and_recovery_preserved":True,"archive_private":True,"integrity_verified":True,"nonoverwrite_and_live_lock_guard":True,"observed_create_seconds":round(create_seconds,3),"observed_restore_seconds":round(restore_seconds,3),"archive_sha256":report["sha256"]})
  print("verified prior-to-new recovery:",backend,flush=True)
 finally:
  stop();log.close()
(root/"results.json").write_text(json.dumps({"passed":True,"ordinary_ingress":"never attached; isolated loopback operator probes","rows":rows,"limits":"small public fixture, prior implementation commit not a published release; no production RTO or cross-platform claim"},indent=2)+"\n")
