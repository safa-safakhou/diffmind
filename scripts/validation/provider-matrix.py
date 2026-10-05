#!/usr/bin/env python3
"""Real-binary provider failure and HTTP/MCP parity matrix; loopback fixtures only.
Usage: provider-matrix.py BINARY NEW_PRIVATE_OUTPUT
No real credentials, repositories, PR comments or remote mutations are used.
"""
import hashlib,json,os,secrets,socket,subprocess,sys,threading,time
from pathlib import Path
from http.server import BaseHTTPRequestHandler,ThreadingHTTPServer
from urllib.request import Request,urlopen
from urllib.error import HTTPError
from urllib.parse import urlsplit,parse_qs
binary,out_s=sys.argv[1:];out=Path(out_s);out.mkdir(mode=0o700)
credential=secrets.token_hex(24);admin=secrets.token_hex(24)
mode={"case":"empty","requests":0,"authenticated":0,"leaks":0}
pull={"number":1,"title":"Fixture","head":{"sha":"a"*40,"ref":"head"},"base":{"ref":"main"},"changed_files":1}
class Other(BaseHTTPRequestHandler):
 def log_message(self,*args):pass
 def do_GET(self):
  mode["leaks"]+=1;self.send_response(200);self.end_headers()
other=ThreadingHTTPServer(("127.0.0.1",0),Other)
class Provider(BaseHTTPRequestHandler):
 def log_message(self,*args):pass
 def do_GET(self):
  mode["requests"]+=1
  assert self.headers.get("Authorization")=="Bearer "+credential
  mode["authenticated"]+=1;case=mode["case"]
  if case in ["401","403","404","429","500"]:
   self.send_response(int(case));self.end_headers();self.wfile.write(json.dumps({"message":credential}).encode());return
  if case=="cross-origin":
   self.send_response(302);self.send_header("Location",f"http://127.0.0.1:{other.server_port}/capture");self.end_headers();return
  if case=="same-origin" and "/redirected" not in self.path:
   self.send_response(302);self.send_header("Location","/redirected");self.end_headers();return
  self.send_response(200);self.send_header("Content-Type","application/json");self.end_headers()
  if case=="malformed":self.wfile.write(b"{broken");return
  if "/files" in self.path:data=[{"filename":"app.py","status":"modified","patch":""}]
  elif "/pulls/1" in self.path:data=pull
  elif case=="paged":data=[{**pull,"number":n} for n in (range(1,101) if parse_qs(urlsplit(self.path).query).get("page")==["1"] else [101])]
  elif case=="empty" or case=="same-origin":data=[]
  else:data=[pull]
  self.wfile.write(json.dumps(data).encode())
provider=ThreadingHTTPServer(("127.0.0.1",0),Provider)
for service in [provider,other]:threading.Thread(target=service.serve_forever,daemon=True).start()
with socket.socket() as s:s.bind(("127.0.0.1",0));port=s.getsockname()[1]
url=f"http://127.0.0.1:{port}"
env={k:v for k,v in os.environ.items() if not k.startswith("DIFFMIND_") and k not in ["GITHUB_TOKEN","GH_TOKEN"]}
env.update(DIFFMIND_HOME=str(out/"home"),DIFFMIND_BINARY=binary,DIFFMIND_AUTH_TOKEN=admin,GITHUB_TOKEN=credential)
log=(out/"private-server.log").open("w")
server=subprocess.Popen([binary,"ui","--no-spa-rebuild","--host","127.0.0.1","--port",str(port),"--refresh-interval","0","--refresh-on-start=false"],env=env,stdout=log,stderr=log)
session="";rows=[]
def request(route,body=None):
 global session
 headers={"Authorization":"Bearer "+admin,"Content-Type":"application/json","Accept":"application/json, text/event-stream","MCP-Protocol-Version":"2025-11-25"}
 if route=="/mcp" and session:headers["Mcp-Session-Id"]=session
 try:r=urlopen(Request(url+route,data=None if body is None else json.dumps(body).encode(),headers=headers),timeout=40)
 except HTTPError as e:r=e
 with r:
  if r.headers.get("Mcp-Session-Id"):session=r.headers["Mcp-Session-Id"]
  raw=r.read().decode()
  assert credential not in raw and admin not in raw,"credential reflected"
  if raw.startswith("event:") or raw.startswith("data:"):raw=next(x[5:].strip() for x in raw.splitlines() if x.startswith("data:"))
  return r.status,json.loads(raw) if raw else None
def api(route,body=None):
 status,data=request(route,body);assert status<300,(route,status,data);return data
def inspect(operation,pid,rid=None):
 selectors={"pid":pid}
 if rid:selectors.update(repo_id=rid,number="1")
 return api("/mcp",{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"inspect_workspace","arguments":{"operation":operation,"selectors":selectors}}})
try:
 for _ in range(100):
  try:
   with urlopen(url+"/healthz",timeout=1):break
  except OSError:time.sleep(.1)
 pid=api("/api/projects",{"name":"Provider matrix"})["id"];base="/api/projects/"+pid
 rid=api(base+"/repos",{"name":"api","git_url":"https://github.fixture.test/acme/api.git","git_provider":"github","git_api_base":f"http://127.0.0.1:{provider.server_port}/api/v3","source_type":"git","kind":"service_repo"})["id"]
 api("/mcp",{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"provider-validation","version":"1"}}})
 api("/mcp",{"jsonrpc":"2.0","method":"notifications/initialized"})
 for case in ["empty","paged","same-origin","401","403","404","429","500","malformed","cross-origin"]:
  mode["case"]=case
  listing=api(base+"/pull-requests");repo=listing["repositories"][0]
  expected="ok" if case in ["empty","paged","same-origin"] else "error"
  assert repo["status"]==expected,(case,repo)
  if expected=="ok":assert repo["open_count"]==(101 if case=="paged" else 0),(case,repo)
  mcp=inspect("list_pull_requests",pid)["result"]
  assert mcp["structuredContent"]["data"]["repositories"]==listing["repositories"],(case,mcp)
  row={"case":case,"list_status":repo["status"],"count":repo["open_count"],"http_mcp_parity":True}
  if expected=="error":
   status,impact=request(base+"/pull-requests/"+rid+"/1/impact")
   assert status==502,(case,status)
   tool=inspect("pull_request_impact",pid,rid)["result"]
   assert impact["error"] in json.dumps(tool),(case,impact,tool)
   row.update(impact_http_status=status,feedback=repo.get("error"))
  rows.append(row)
 assert mode["leaks"]==0
 (out/"results.json").write_text(json.dumps({"passed":True,"binary_sha256":hashlib.sha256(Path(binary).read_bytes()).hexdigest(),"rows":rows,"credential_requests":mode["authenticated"],"cross_origin_requests":mode["leaks"],"limits":"Controlled GitHub-compatible HTTP endpoint and fake credentials; not a real Enterprise or SSO deployment."},indent=2))
 print(json.dumps(rows),flush=True)
finally:
 (out/"observations.json").write_text(json.dumps({"rows":rows},indent=2))
 server.terminate();server.wait(timeout=20);log.close()
 for service in [provider,other]:service.shutdown();service.server_close()
