import fs from 'node:fs';
import path from 'node:path';
import net from 'node:net';
import assert from 'node:assert/strict';
import {spawn,execFileSync} from 'node:child_process';
import {randomBytes,createHash} from 'node:crypto';
import {pathToFileURL} from 'node:url';

const [binary,root,packFile]=process.argv.slice(2);
assert.ok(binary&&root&&packFile,'usage: node public-correction.mjs BINARY PRIVATE_TRIAL_ROOT PACK_JSON');
const {chromium}=await import(pathToFileURL(process.env.DIFFMIND_PLAYWRIGHT_MODULE).href);
const home=path.join(root,'live-home-'+Date.now()); fs.mkdirSync(home,{mode:0o700});
const evidence=path.join(root,'evidence'),labelsBytes=fs.readFileSync(path.join(evidence,'source-labels.json')),labels=JSON.parse(labelsBytes);
const pack=JSON.parse(fs.readFileSync(packFile));
const rows=[]; const record=(name,data={})=>{rows.push({name,...data});fs.writeFileSync(path.join(evidence,'live-results.json'),JSON.stringify({rows},null,2));console.log(name)};
const env=Object.fromEntries(Object.entries(process.env).filter(([k])=>!k.startsWith('DIFFMIND_')&&!['GITHUB_TOKEN','GH_TOKEN'].includes(k)));
const sock=net.createServer();await new Promise(r=>sock.listen(0,'127.0.0.1',r));const port=sock.address().port;await new Promise(r=>sock.close(r));
const url='http://127.0.0.1:'+port,proxy=randomBytes(32).toString('hex'),recovery=randomBytes(32).toString('hex');
const headers=(role='admin')=>({'Content-Type':'application/json','X-DiffMind-Proxy-Secret':proxy,'X-DiffMind-User':'trial-'+role,'X-DiffMind-Role':role==='ungranted'?'viewer':role});
const log=fs.createWriteStream(path.join(root,'private-server.log'));
const start=mode=>{const child=spawn(binary,['ui','--no-spa-rebuild','--host','127.0.0.1','--port',String(port),'--project-access',mode],{env:{...env,DIFFMIND_HOME:home,DIFFMIND_BINARY:binary,DIFFMIND_TRUSTED_PROXY_SECRET:proxy,DIFFMIND_AUTH_TOKEN:recovery},stdio:['ignore','pipe','pipe']});child.stdout.pipe(log,{end:false});child.stderr.pipe(log,{end:false});return child};
let server=start('legacy');
let browser;const sleep=ms=>new Promise(r=>setTimeout(r,ms));
const request=async(route,method='GET',body,role='admin')=>{const res=await fetch(url+route,{method,headers:headers(role),...(body===undefined?{}:{body:JSON.stringify(body)})});return {status:res.status,data:await res.json()}};
const api=async(...args)=>{const r=await request(...args);assert.ok(r.status<300,JSON.stringify(r));return r.data};
const tool=async(name,args,role='viewer')=>{const r=await fetch(url+'/mcp',{method:'POST',headers:{...headers(role),Accept:'application/json, text/event-stream','MCP-Protocol-Version':'2025-11-25'},body:JSON.stringify({jsonrpc:'2.0',id:1,method:'tools/call',params:{name,arguments:args}})});const d=await r.json();assert.ok(r.ok&&!d.error&&!d.result?.isError,JSON.stringify(d));return d.result.structuredContent};
const ingest=async(pid)=>{await api('/api/projects/'+pid+'/ingestion','POST',{concurrency:2});for(let i=0;i<1200;i++){const s=await api('/api/projects/'+pid+'/ingestion');if(s.status!=='running'){assert.equal(s.status,'completed',JSON.stringify(s));return s.graph_run_id}await sleep(250)}throw Error('ingestion timeout')};
const graph=(pid,run)=>api('/api/projects/'+pid+'/runs/'+run+'/archgraph?view=full');
const key=e=>[e.from,e.to,e.type].join('|');
try{
for(let i=0;i<100;i++){try{if((await fetch(url+'/healthz')).ok)break}catch{}await sleep(100)}
const p=await api('/api/projects','POST',{name:'Pinned Boutique reviewed correction'}),pid=p.id;
const policy=await api('/api/v1/projects/'+pid+'/access');
await api('/api/v1/projects/'+pid+'/access','PUT',{revision:policy.revision,members:{'trial-viewer':'viewer','trial-editor':'editor'}});
assert.equal((await request('/api/projects/'+pid,'GET',undefined,'ungranted')).status,200);
const stopping=server;const stopped=new Promise(r=>stopping.once('exit',r));stopping.kill('SIGTERM');await stopped;
server=start('scoped');for(let i=0;i<100;i++){try{if((await fetch(url+'/healthz')).ok)break}catch{}await sleep(100)}
assert.equal((await request('/api/projects/'+pid,'GET',undefined,'ungranted')).status,404);
record('reviewed_legacy_to_scoped_migration',{saved_grants:true,legacy_unchanged_before_explicit_restart:true,ungranted_denied_after_activation:true});
const base='/api/projects/'+pid;
const repos=[];for(const name of labels.mappings.map(m=>m.repository)){repos.push(await api(base+'/repos','POST',{name,path:path.join(root,'repositories',name),source_type:'local',kind:'service_repo'}))}
const baseline=await ingest(pid),before=await graph(pid,baseline),beforeBytes=fs.readFileSync(path.join(home,'projects',pid,'runs',baseline,'graph.json')),beforeHash=createHash('sha256').update(beforeBytes).digest('hex');
record('source_first_baseline',{revision:labels.revision,labels_sha256:createHash('sha256').update(labelsBytes).digest('hex'),services:before.services.length,edges:before.edges.length,expected_declarations:labels.expected_declared_edges.length,known_declared_found:labels.expected_declared_edges.filter(e=>before.edges.some(x=>key(x)===key(e))).length});
assert.equal((await request('/api/v1/projects/'+pid+'/improvement-gaps','POST',{revision:0,gap:{category:'coverage',expected:'reviewed declarations',observed:'missing'}},'viewer')).status,403);
let gap=await api('/api/v1/projects/'+pid+'/improvement-gaps','POST',{revision:0,gap:{category:'coverage',run_id:baseline,expected:'Seven frontend and other reviewed literal gRPC configuration declarations',observed:'Absent without an explicit tested pack',source_pointers:['kubernetes-manifests/frontend.yaml:79'] }},'editor');
const transition=async(status,reason)=>{gap=await api('/api/v1/projects/'+pid+'/improvement-gaps/'+gap.gap.id,'PATCH',{revision:gap.revision,status,reason,tests:['literal positive','near-match comments placeholders and unrelated variables','exact declared relationship','no phantom negative relationships']})};
await transition('reproduced','Pinned source-first labels and baseline assembled graph');
await transition('proposed','Narrow literal configuration declaration pack; no runtime inference');
await transition('tested','Candidate CLI lint and all four positive/negative assertions passed');
await transition('accepted','User authorized the correction trial; activation limited to disposable corpus');
assert.equal((await request(base+'/packs','POST',pack,'editor')).status,403);
const installed=await api(base+'/packs','POST',pack);
for(const repo of repos)await api(base+'/repos/'+repo.id,'PATCH',{pack_ids:[installed.id]});
await transition('active','Installed validated project pack and explicitly linked reviewed repositories');
const corrected=await ingest(pid),after=await graph(pid,corrected),declared=after.edges.filter(e=>e.evidence?.some(v=>v.pack_id===pack.id));
const expected=new Set(labels.expected_declared_edges.map(key)),actual=new Set(declared.map(key));
assert.deepEqual([...actual].sort(),[...expected].sort());
for(const edge of declared){const ev=edge.evidence.find(v=>v.pack_id===pack.id);assert.equal(ev.class,'pack_declared');assert.equal(ev.run_id,corrected);assert.equal(ev.coverage,'unverified');assert.ok(Array.isArray(ev.file_scope)&&ev.file_scope.length>0);const result=await tool('get_dependencies',{project:pid,run:corrected,service:edge.from,direction:'outbound'});assert.deepEqual(result.edges.find(e=>key(e)===key(edge))?.evidence,edge.evidence);}
record('corrected_declared_graph',{run:corrected,known_edges:expected.size,pack_edges:actual.size,declared_subset_precision:1,declared_subset_recall:1,http_mcp_evidence_parity:true,source_extracted_extras:'unlabeled; no complete graph precision',runtime_reachability:'unverified'});
browser=await chromium.launch({headless:true,executablePath:process.env.DIFFMIND_CHROMIUM});
const page=await browser.newPage({extraHTTPHeaders:headers(),viewport:{width:1440,height:1000}}),errors=[];page.on('pageerror',e=>errors.push(e.message));
await page.goto(url+'/#/projects/'+pid+'/access');await page.getByRole('heading',{name:'Project access',exact:true}).waitFor();
assert.match(await page.locator('main').innerText(),/Active access mode: scoped/);assert.match(await page.locator('main').innerText(),/permissions are not mirrored/);
await page.screenshot({path:path.join(evidence,'access-setup.png'),fullPage:true});
record('browser_scoped_setup',{mode:'scoped',project_visibility_boundary_visible:true});
await page.goto(url+'/#/projects/'+pid);await page.getByLabel('Workspace readiness').waitFor();
await page.locator('.compact-service-name').first().waitFor();
const sizes=await page.locator('.compact-service-name').evaluateAll(nodes=>nodes.map(n=>parseFloat(getComputedStyle(n).fontSize)*Math.abs(n.getScreenCTM()?.a||1)));assert.ok(sizes.every(n=>n>=12));
await page.screenshot({path:path.join(evidence,'corrected-public-graph.png'),fullPage:true});
record('browser_public_graph',{labels:sizes.length,minimum_label_pixels:Math.min(...sizes)});
const viewer=await browser.newPage({extraHTTPHeaders:headers('viewer')});await viewer.goto(url+'/#/projects/'+pid);await viewer.getByLabel('Workspace readiness').waitFor();assert.equal(await viewer.getByRole('button',{name:'Update context',exact:true}).count(),0);
const readiness=await tool('get_readiness',{project:pid});assert.equal(readiness.actions.refresh,false);
assert.equal((await request('/api/v1/projects/'+pid+'/readiness','GET',undefined,'ungranted')).status,404);
const recoveryResponse=await fetch(url+'/api/v1/projects/'+pid+'/readiness',{headers:{Authorization:'Bearer '+recovery}});assert.equal(recoveryResponse.status,200);
const ungranted=await browser.newPage({extraHTTPHeaders:headers('ungranted')});await ungranted.goto(url+'/#/projects/'+pid);await ungranted.getByRole('alert').filter({hasText:'unavailable'}).waitFor();
const unknownProjects=await tool('list_projects',{},'ungranted');assert.equal(unknownProjects.projects.length,0);
const editor=await tool('get_readiness',{project:pid},'editor');assert.equal(editor.actions.refresh,true);assert.equal(editor.actions.configure,false);
record('ordinary_viewer_editor_admin_and_ungranted',{viewer_readonly:true,editor_reporting_without_pack_authority:true,admin_recovery:true,ungranted_denied:true});
for(const repo of repos)await api(base+'/repos/'+repo.id,'PATCH',{pack_ids:[]});
await api(base+'/packs/'+installed.id,'DELETE');
await transition('rolled_back','Remove bindings and trial pack; empty bindings otherwise use project matching. Preserve earlier snapshots');
const rollback=await ingest(pid),rolled=await graph(pid,rollback);
assert.deepEqual(rolled.edges.map(key).sort(),before.edges.map(key).sort());
assert.equal(createHash('sha256').update(fs.readFileSync(path.join(home,'projects',pid,'runs',baseline,'graph.json'))).digest('hex'),beforeHash);
assert.deepEqual((await graph(pid,corrected)).edges,after.edges);
record('rollback_and_immutable_history',{baseline,corrected,rollback,baseline_bytes_unchanged:true,corrected_snapshot_preserved:true,baseline_topology_restored:true,gap_final_state:gap.gap.status});
assert.deepEqual(errors,[]);record('finished',{passed:true,browser_exceptions:errors.length});
}catch(e){record('failure',{error:e.stack});process.exitCode=1;}
finally{if(browser)await browser.close();server.kill('SIGTERM');await Promise.race([new Promise(r=>server.once('exit',r)),sleep(5000)]);}
