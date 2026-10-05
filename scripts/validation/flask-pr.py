#!/usr/bin/env python3
"""Controlled PR experiment using pinned public Flask source; no remote mutations."""
import json,os,shutil,socket,subprocess,sys,threading,time
from pathlib import Path
from http.server import BaseHTTPRequestHandler,ThreadingHTTPServer
from urllib.request import Request,urlopen
binary,source,out_text=sys.argv[1:4]
host_runner=sys.argv[4] if len(sys.argv)>4 else None
out=Path(out_text);out.mkdir(mode=0o700)
repo=out/'flask';shutil.copytree(source,repo)
def git(*args):return subprocess.check_output(['git','-C',str(repo),*args],text=True).strip()
git('init','-q');git('add','.')
git('-c','user.name=Validation','-c','user.email=validation@example.test','commit','-qm','Pinned public tutorial baseline')
base=git('rev-parse','HEAD')
auth=repo/'flaskr/auth.py';auth.write_text(auth.read_text().replace('"/login"','"/sign-in"'))
git('add','.');git('-c','user.name=Validation','-c','user.email=validation@example.test','commit','-qm','Controlled public route rename')
head=git('rev-parse','HEAD');patch=git('diff','--unified=3',base,head,'--','flaskr/auth.py')
client=out/'client';client.mkdir();(client/'client.py').write_text('import requests\n\ndef sign_in():\n    return requests.get("http://flask/auth/sign-in")\n')
subprocess.run(['git','init','-q',str(client)],check=True)
subprocess.run(['git','-C',str(client),'add','.'],check=True)
subprocess.run(['git','-C',str(client),'-c','user.name=Validation','-c','user.email=validation@example.test','commit','-qm','Explicit synthetic caller'],check=True)
# Oracle frozen from source/diff before provider or DiffMind analysis starts.
(out/'source-oracle.json').write_text(json.dumps({
 'source':'pallets/flask 3.1.2 tutorial','derived_base':base,'derived_head':head,
 'operation':{'method':'GET','path':'/auth/sign-in','file':'flaskr/auth.py'},
 'expected_direct_callers':['caller'],'synthetic_caller':True,
 'reason':'Literal requests.get URL matches the renamed Blueprint route.',
 'negative_controls':['missing patch','unrelated filename','deleted line','stale head','dirty caller'],
 'runtime_claim':False,'upstream_pr':False},indent=2))
pull={'number':1,'title':'Controlled route rename','html_url':'https://controlled-pr.example.test/flask/pull/1','head':{'sha':head,'ref':'validation'},'base':{'ref':'baseline'},'changed_files':1,'additions':1,'deletions':1,'commits':1}
files=[{'filename':'flaskr/auth.py','status':'modified','patch':patch,'additions':1,'deletions':1}]
class Provider(BaseHTTPRequestHandler):
 def log_message(self,*args):pass
 def do_GET(self):
  data=files if '/files' in self.path else ([pull] if '/pulls?' in self.path else pull)
  self.send_response(200);self.send_header('Content-Type','application/json');self.end_headers();self.wfile.write(json.dumps(data).encode())
provider=ThreadingHTTPServer(('127.0.0.1',0),Provider);threading.Thread(target=provider.serve_forever,daemon=True).start()
with socket.socket() as s:s.bind(('127.0.0.1',0));port=s.getsockname()[1]
url='http://127.0.0.1:'+str(port);home=out/'home'
env={k:v for k,v in os.environ.items() if not k.startswith('DIFFMIND_') and k not in ('GITHUB_TOKEN','GH_TOKEN')}
env.update(DIFFMIND_HOME=str(home),DIFFMIND_BINARY=binary)
log=(out/'private-server.log').open('w')
server=subprocess.Popen([binary,'ui','--no-spa-rebuild','--port',str(port),'--refresh-interval','0','--refresh-on-start=false'],env=env,stdout=log,stderr=log)
mcp_session=''
def api(route,body=None):
 global mcp_session
 req=Request(url+route,data=json.dumps(body).encode() if body is not None else None,headers={'Content-Type':'application/json','Accept':'application/json, text/event-stream','MCP-Protocol-Version':'2025-11-25'})
 if route=='/mcp' and mcp_session:req.add_header('Mcp-Session-Id',mcp_session)
 with urlopen(req,timeout=30) as r:
  if r.headers.get('Mcp-Session-Id'):mcp_session=r.headers['Mcp-Session-Id']
  raw=r.read().decode()
  if not raw:return None
  if r.headers.get('Content-Type','').startswith('text/event-stream'):
   return json.loads(next(line[5:].strip() for line in raw.splitlines() if line.startswith('data:')))
  return json.loads(raw)
def host_trial(name,prompt):
 if not host_runner:return
 spec={'id':name,'binary':binary,'mode':'http','url':url+'/mcp','work':str(out/'host-work'),'prompt':prompt,'timeout_seconds':300}
 file=out/(name+'-spec.json');file.write_text(json.dumps(spec,indent=2))
 subprocess.run([sys.executable,host_runner,str(file),str(out/name)],check=True,env=env)
try:
 for _ in range(100):
  try:
   with urlopen(url+'/healthz',timeout=1):break
  except OSError:time.sleep(.1)
 p=api('/api/projects',{'name':'Public Flask controlled PR'});pid=p['id'];baseurl='/api/projects/'+pid
 registered=api(baseurl+'/repos',{'name':'flask','path':str(repo),'source_type':'local','kind':'service_repo','git_url':'https://github.com/pallets/flask.git','git_api_base':'http://127.0.0.1:'+str(provider.server_port)})
 api(baseurl+'/repos',{'name':'caller','path':str(client),'source_type':'local','kind':'service_repo'})
 api(baseurl+'/ingestion',{'concurrency':2})
 for _ in range(300):
  state=api(baseurl+'/ingestion')
  if state['status']!='running':break
  time.sleep(.1)
 assert state['status']=='completed',state
 graph=api(baseurl+'/runs/'+state['graph_run_id']+'/archgraph?view=full')
 result=api(baseurl+'/pull-requests/'+registered['id']+'/1/impact')
 (out/'graph.json').write_text(json.dumps(graph,indent=2))
 (out/'impact.json').write_text(json.dumps(result,indent=2))
 (out/'source.json').write_text(json.dumps({'controlled_pr':True,'upstream_pr':False,'public_source':'pallets/flask tutorial 3.1.2','derived_base':base,'derived_head':head,'patch':patch,'synthetic_caller':True},indent=2))
 assert result['company']['score_eligible'] and result['company']['direct_services']==1,result['company']
 api('/mcp',{'jsonrpc':'2.0','id':0,'method':'initialize','params':{'protocolVersion':'2025-11-25','capabilities':{},'clientInfo':{'name':'diffmind-validation','version':'1'}}})
 api('/mcp',{'jsonrpc':'2.0','method':'notifications/initialized'})
 mcp=api('/mcp',{'jsonrpc':'2.0','id':1,'method':'tools/call','params':{'name':'inspect_workspace','arguments':{'operation':'pull_request_impact','selectors':{'pid':pid,'repo_id':registered['id'],'number':'1'}}}})
 assert not mcp.get('error') and not mcp['result'].get('isError'),mcp
 structured=mcp['result']['structuredContent']
 assert structured['data']['company']==result['company']
 host_trial('U03-exact-pr','Review PR 1 for the configured Flask repository in Public Flask controlled PR. Explain which other services may be affected, what the source evidence supports, its limits, and whether this review is posted on the PR. Inspect only; preserve local source.')
 rows=[{'case':'matching clean route and caller','eligible':True,'exact_callers':1,'http_mcp_company_parity':True}]
 impact_path=baseurl+'/pull-requests/'+registered['id']+'/1/impact'
 for case in ['missing patch','unrelated filename','deleted line','stale head']:
  files[0]['patch']=patch;files[0]['filename']='flaskr/auth.py';pull['head']['sha']=head
  if case=='missing patch':files[0].pop('patch')
  if case=='unrelated filename':files[0]['filename']='unrelated/auth.py'
  if case=='deleted line':files[0]['patch']='@@ -84 +83,0 @@\n-removed decorator'
  if case=='stale head':pull['head']['sha']='f'*40
  trial=api(impact_path)['company']
  assert not trial['score_eligible'],(case,trial)
  rows.append({'case':case,'eligible':False,'freshness':trial['freshness'],'exact_callers':trial['direct_services']})
 files[0]['patch']=patch;files[0]['filename']='flaskr/auth.py';pull['head']['sha']=head
 (client/'client.py').write_text((client/'client.py').read_text()+'\n# uncommitted caller change\n')
 api(baseurl+'/ingestion',{'concurrency':2})
 for _ in range(300):
  state=api(baseurl+'/ingestion')
  if state['status']!='running':break
  time.sleep(.1)
 assert state['status']=='completed'
 trial=api(impact_path)['company']
 assert not trial['score_eligible'],trial
 rows.append({'case':'dirty caller analysis','eligible':False,'root_freshness':trial['freshness']})
 assert git('status','--porcelain')=='','PR read modified clean local source'
 (out/'controls.json').write_text(json.dumps({'rows':rows,'source_clean':True,'independent_human_reviewers':0,'actual_upstream_pr':False},indent=2))
 if host_runner:
  # Actual public PR is a separate stale-revision trial, never relabeled controlled data.
  patch_req=Request(url+baseurl+'/repos/'+registered['id'],method='PATCH',data=json.dumps({'git_api_base':''}).encode(),headers={'Content-Type':'application/json'})
  with urlopen(patch_req,timeout=30) as response:assert response.status==200
  public=api(baseurl+'/pull-requests/'+registered['id']+'/5918/impact')
  assert not public['company']['score_eligible'],public['company']
  (out/'public-impact.json').write_text(json.dumps(public,indent=2))
  host_trial('U03-public-stale-pr','Review pallets/flask PR 5918 using the configured company context. Explain which services may be affected, what can actually be concluded from the saved revision and available evidence, and whether your review has been delivered to the PR. Inspect only; preserve local source.')
 print(json.dumps({'services':len(graph['services']),'edges':len(graph['edges']),'rows':rows},indent=2))
finally:
 server.terminate();server.wait(timeout=20);log.close();provider.shutdown();provider.server_close()
