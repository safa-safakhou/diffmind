#!/usr/bin/env python3
"""Exercise background refresh, concurrent reads and trusted-proxy roles.
Usage: background-matrix.py BINARY SIX_REPO_SOURCE NEW_OUTPUT [SECONDS]
The proxy is a controlled identity fixture, not production SSO.
"""
import concurrent.futures,hashlib,json,os,secrets,shutil,socket,subprocess,sys,threading,time
from pathlib import Path
from http.server import BaseHTTPRequestHandler,ThreadingHTTPServer
from urllib.request import Request,urlopen
from urllib.error import HTTPError
binary,source_s,out_s=sys.argv[1:4];duration=int(sys.argv[4]) if len(sys.argv)>4 else 120
out=Path(out_s);out.mkdir(mode=0o700);rows=[]
clean={k:v for k,v in os.environ.items() if not k.startswith("DIFFMIND_") and k not in ["GITHUB_TOKEN","GH_TOKEN"]}
for backend in ["json","sqlite"]:
 root=out/backend;root.mkdir();home=root/"home";repos=root/"repositories"
 shutil.copytree(source_s,repos,ignore=shutil.ignore_patterns(".git","__pycache__"))
 for repo in repos.iterdir():
  for args in [["init","-q"],["add","."],["-c","user.name=Validation","-c","user.email=validation@example.test","commit","-qm","Frozen background source"]]:
   subprocess.run(["git","-C",str(repo),*args],check=True,stdout=subprocess.DEVNULL)
 with socket.socket() as s:s.bind(("127.0.0.1",0));port=s.getsockname()[1]
 url=f"http://127.0.0.1:{port}";admin=secrets.token_hex(24);proxy_secret=secrets.token_hex(24);identity=secrets.token_hex(24)
 env={**clean,"DIFFMIND_HOME":str(home),"DIFFMIND_BINARY":binary,"DIFFMIND_AUTH_TOKEN":admin,"DIFFMIND_TRUSTED_PROXY_SECRET":proxy_secret}
 log=(root/"private-server.log").open("w");server=None;proxy=None
 def request(route,method="GET",body=None,headers=None,base=url):
  h={"Content-Type":"application/json","Authorization":"Bearer "+admin} if headers is None else headers
  req=Request(base+route,method=method,data=None if body is None else json.dumps(body).encode(),headers=h)
  try:r=urlopen(req,timeout=30)
  except HTTPError as e:r=e
  with r:
   raw=r.read().decode()
   if raw.startswith("event:") or raw.startswith("data:"):raw=next(x[5:].strip() for x in raw.splitlines() if x.startswith("data:"))
   return r.status,json.loads(raw) if raw else None,dict(r.headers)
 def api(route,method="GET",body=None):
  code,data,_=request(route,method,body);assert code<300,(route,code,data);return data
 def start(interval):
  global server
  server=subprocess.Popen([binary,"ui","--no-spa-rebuild","--host","127.0.0.1","--port",str(port),"--project-access","scoped","--refresh-interval",interval,"--refresh-on-start=false"],env=env,stdout=log,stderr=log)
  for _ in range(100):
   try:
    with urlopen(url+"/healthz",timeout=1):return
   except OSError:time.sleep(.1)
  raise RuntimeError("server startup failed")
 def stop():
  global server
  if server:server.terminate();server.wait(timeout=20);server=None
 class Proxy(BaseHTTPRequestHandler):
  def log_message(self,*args):pass
  def do_GET(self):self.forward()
  def do_POST(self):self.forward()
  def forward(self):
   if self.headers.get("Authorization")!="Bearer "+identity:
    self.send_response(401);self.end_headers();return
   # Ignore incoming identity/role/secret headers. Only fixture authentication
   # decides the forwarded subject. Header spoofing cannot elevate the viewer.
   headers={"X-DiffMind-Proxy-Secret":proxy_secret,"X-DiffMind-User":"joined-viewer","X-DiffMind-Role":"viewer"}
   for key in ["Content-Type","Accept","MCP-Protocol-Version","Mcp-Session-Id"]:
    if self.headers.get(key):headers[key]=self.headers[key]
   body=self.rfile.read(int(self.headers.get("Content-Length",0)))
   req=Request(url+self.path,data=body or None,method=self.command,headers=headers)
   try:r=urlopen(req,timeout=30)
   except HTTPError as e:r=e
   with r:
    payload=r.read();self.send_response(r.status)
    for key in ["Content-Type","Mcp-Session-Id"]:
     if r.headers.get(key):self.send_header(key,r.headers[key])
    self.send_header("Content-Length",str(len(payload)));self.end_headers();self.wfile.write(payload)
 try:
  start("0");pid=api("/api/projects","POST",{"name":"Background "+backend})["id"];base="/api/projects/"+pid
  for repo in sorted(repos.iterdir()):api(base+"/repos","POST",{"name":repo.name,"path":str(repo),"source_type":"local","kind":"service_repo"})
  api(base+"/ingestion","POST",{"concurrency":2})
  for _ in range(300):
   state=api(base+"/ingestion")
   if state["status"]!="running":break
   time.sleep(.1)
  assert state["status"]=="completed",state
  old_run=state["graph_run_id"];old_file=home/"projects"/pid/"runs"/old_run/"graph.json";old_hash=hashlib.sha256(old_file.read_bytes()).hexdigest()
  access="/api/v1/projects/"+pid+"/access";policy=api(access)
  api(access,"PUT",{"revision":policy["revision"],"members":{"joined-viewer":"viewer"}})
  proxy=ThreadingHTTPServer(("127.0.0.1",0),Proxy);threading.Thread(target=proxy.serve_forever,daemon=True).start()
  proxy_url=f"http://127.0.0.1:{proxy.server_port}"
  viewer={"Authorization":"Bearer "+identity,"X-DiffMind-Role":"admin","Content-Type":"application/json","Accept":"application/json, text/event-stream","MCP-Protocol-Version":"2025-11-25"}
  assert request(base,headers=viewer,base=proxy_url)[0]==200
  assert request("/api/v1/projects/"+pid+"/refresh-jobs","POST",{},viewer,proxy_url)[0]==403
  assert request(base,headers={"X-DiffMind-User":"joined-viewer","X-DiffMind-Role":"admin"})[0]==401
  code,_,headers=request("/mcp","POST",{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"proxy-matrix","version":"1"}}},viewer,proxy_url)
  assert code==200
  sid=next(v for k,v in headers.items() if k.lower()=="mcp-session-id");viewer["Mcp-Session-Id"]=sid
  request("/mcp","POST",{"jsonrpc":"2.0","method":"notifications/initialized"},viewer,proxy_url)
  call={"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"get_readiness","arguments":{"project":pid}}}
  code,read,_=request("/mcp","POST",call,viewer,proxy_url);assert code==200 and "result" in read and not read["result"].get("isError"),read
  policy=api(access);api(access,"PUT",{"revision":policy["revision"],"members":{}})
  assert request(base,headers=viewer,base=proxy_url)[0]==404
  code,denied,_=request("/mcp","POST",call,viewer,proxy_url)
  assert code==200 and (denied["result"].get("isError") or denied["result"].get("structuredContent",{}).get("status")==404),denied
  stop()
  if backend=="sqlite":subprocess.run([binary,"storage","migrate","--offline","--json"],env=env,check=True,stdout=subprocess.DEVNULL)
  start("2s")
  latencies=[];rss=[];restarted=False;edited=False;peak_running=0;peak_queued=0
  start_time=time.monotonic()
  def read(_):
   t=time.monotonic();data=api("/api/v1/projects/"+pid+"/graph/summary");return round((time.monotonic()-t)*1000,3)
  with concurrent.futures.ThreadPoolExecutor(max_workers=8) as pool:
   while time.monotonic()-start_time<duration:
    latencies.extend(pool.map(read,range(16)))
    jobs=api("/api/v1/jobs?project="+pid+"&limit=100")["jobs"]
    running=sum(j["status"]=="running" for j in jobs);queued=sum(j["status"]=="queued" for j in jobs)
    peak_running=max(peak_running,running);peak_queued=max(peak_queued,queued);assert running<=1,(backend,jobs)
    status=Path("/proc")/str(server.pid)/"status"
    rss.append(int(next(l.split()[1] for l in status.read_text().splitlines() if l.startswith("VmRSS:"))))
    elapsed=time.monotonic()-start_time
    if not edited and elapsed>duration/4:
     p=repos/"catalog/app.py";p.write_text(p.read_text()+"\n# Background local change, preserve me.\n");edited=True
    if not restarted and elapsed>duration/2:stop();start("2s");restarted=True
    time.sleep(.5)
  jobs=api("/api/v1/jobs?project="+pid+"&limit=100")["jobs"]
  successes=[j for j in jobs if j["status"]=="succeeded"];assert len(successes)>=5,(backend,jobs)
  attempts=[a for j in successes for a in j["attempts"]]
  assert any(a["analyzed"]==1 and a["reused"]==5 for a in attempts),attempts
  assert any(a["analyzed"]==0 and a["reused"]==6 for a in attempts),attempts
  assert hashlib.sha256(old_file.read_bytes()).hexdigest()==old_hash
  assert (repos/"catalog/app.py").read_text().endswith("# Background local change, preserve me.\n")
  assert len(api(base+"/repos")["repos"])==6
  rows.append({"backend":backend,"seconds":round(time.monotonic()-start_time,2),"successful_jobs":len(successes),"read_requests":len(latencies),"read_ms_p95":sorted(latencies)[int(len(latencies)*.95)],"rss_kib_min":min(rss),"rss_kib_max":max(rss),"peak_running":peak_running,"peak_queued":peak_queued,"restart_recovered":restarted,"automatic_one_changed_five_reused":True,"automatic_six_reused":True,"scope_and_edit_preserved":True,"old_graph_hash":old_hash,"proxy_viewer_read":200,"proxy_spoofed_admin_write":403,"direct_untrusted_identity":401,"membership_removed_browser":404,"existing_mcp_membership_rechecked":True})
  print(json.dumps(rows[-1]),flush=True)
 finally:
  stop();log.close()
  if proxy:proxy.shutdown();proxy.server_close()
  (out/"observations.json").write_text(json.dumps({"rows":rows},indent=2))
(out/"results.json").write_text(json.dumps({"passed":True,"binary_sha256":hashlib.sha256(Path(binary).read_bytes()).hexdigest(),"rows":rows,"limits":"Bounded local endurance with eight readers, six small repositories, controlled proxy identity and graceful restart; not production SSO, capacity or long-duration soak."},indent=2))
