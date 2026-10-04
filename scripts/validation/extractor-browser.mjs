import fs from 'node:fs';
import path from 'node:path';
import net from 'node:net';
import assert from 'node:assert/strict';
import {spawn} from 'node:child_process';
import {pathToFileURL} from 'node:url';

const [binary, artifacts, output]=process.argv.slice(2);
assert.ok(binary && artifacts && output,'BINARY EXISTING_ARTIFACTS NEW_EVIDENCE_DIRECTORY');
fs.mkdirSync(output,{mode:0o700});
fs.cpSync(artifacts,path.join(output,'runs'),{recursive:true});
const rows=[],errors=[];
const engines=await import(pathToFileURL(process.env.DIFFMIND_PLAYWRIGHT_MODULE).href);
const engine=process.env.DIFFMIND_BROWSER || 'chromium';
assert.ok(['chromium','firefox','webkit'].includes(engine), 'supported browser engine required');
const {default:AxeBuilder}=await import(pathToFileURL(process.env.DIFFMIND_AXE_MODULE).href);
const sock=net.createServer();
await new Promise(r=>sock.listen(0,'127.0.0.1',r));
const port=sock.address().port;
await new Promise(r=>sock.close(r));
const env=Object.fromEntries(Object.entries(process.env).filter(([k])=>!k.startsWith('DIFFMIND_')&&!['GITHUB_TOKEN','GH_TOKEN'].includes(k)));
const log=fs.openSync(path.join(output,'private-server.log'),'w',0o600);
const server=spawn(binary,['extractor-ui','--no-spa-rebuild','--host','127.0.0.1','--port',String(port),'--out',path.join(output,'runs')],{env:{...env,DIFFMIND_HOME:path.join(output,'home')},stdio:['ignore',log,log]});
let browser;
const sleep=ms=>new Promise(r=>setTimeout(r,ms)),url='http://127.0.0.1:'+port;
const save=()=>fs.writeFileSync(path.join(output,'results.json'),JSON.stringify({engine,browser_version:browser?.version(),rows,errors},null,2));
try{
 for(let i=0;i<100;i++){try{if((await fetch(url+'/api/repos')).ok)break}catch{}await sleep(100);}
 const repos=await(await fetch(url+'/api/repos')).json();
 assert.equal(repos.repos.length,1,'use one source for this observer');
 const expectedPath=repos.repos[0].path;
 browser=await engines[engine].launch({headless:true,...(engine==='chromium' && process.env.DIFFMIND_CHROMIUM ? {executablePath:process.env.DIFFMIND_CHROMIUM} : {})});
 const context=await browser.newContext();
 const page=await context.newPage();
 page.on('pageerror',e=>errors.push(e.message));
 const scan=async(stage,width)=>{
  await sleep(300);
  const a=await new AxeBuilder({page}).withTags(['wcag2a','wcag2aa','wcag21a','wcag21aa']).analyze();
  const size=await page.evaluate(()=>({width:innerWidth,document:document.documentElement.scrollWidth}));
  rows.push({stage,width,size,violations:a.violations.map(x=>({id:x.id,impact:x.impact,nodes:x.nodes.map(n=>({target:n.target,summary:n.failureSummary}))}))});
  save();
  await page.screenshot({path:path.join(output,stage+'-'+width+'.png'),fullPage:true});
 };
 for(const width of [1440,390]){
  await page.setViewportSize({width,height:1000});
  await page.goto(url);await page.getByRole('heading',{name:'Deterministic Runs',exact:true}).waitFor();
  await scan('repositories',width);
  const opener=page.getByRole('button',{name:'+ Add repository',exact:true}).first();
  await opener.focus();await page.keyboard.press('Enter');
  await page.locator('input[placeholder="/abs/path/to/repo"]').waitFor();
  await scan('add-repository',width);
  await page.getByRole('button',{name:'Cancel',exact:true}).click();
  await page.goto(url+'/#/runs/'+encodeURIComponent(repos.repos[0].last_run_id));
  await page.locator('.status-pill.completed').waitFor();
  await page.getByText('This saved run has no recorded stage history.',{exact:false}).waitFor();
  assert.match(await page.locator('.pipeline-strip').innerText(),/not recorded/);
  await scan('saved-cli-artifact',width);
  await page.goto(url);await page.getByRole('heading',{name:'Deterministic Runs',exact:true}).waitFor();
  await page.getByRole('button',{name:'Run',exact:true}).first().click();
  const initialPath=await page.locator('input[placeholder="/abs/path/to/repo"]').inputValue();
  rows.push({stage:'selected-repository',width,expected:expectedPath,actual:initialPath});save();
  assert.equal(initialPath,expectedPath,'Run must use the selected repository, including after remembered defaults');
  await scan('run-form',width);
  const dialog=page.getByRole('dialog');
  assert.equal(await dialog.count(),1);
  let trapped=true;
  for(let i=0;i<16;i++){await page.keyboard.press(i%3===0?'Shift+Tab':'Tab');trapped&&=await dialog.evaluate(d=>d.contains(document.activeElement));}
  assert.ok(trapped);await page.keyboard.press('Escape');await dialog.waitFor({state:'hidden'});
  assert.ok(await page.getByRole('button',{name:'Run',exact:true}).first().evaluate(e=>e===document.activeElement));
  rows.push({stage:'keyboard-dialog',width,trapped,restored:true});
  // Poison a remembered target; selecting this card must still choose this repository.
  await page.evaluate(()=>localStorage.setItem('diffmind:form-defaults-v7',JSON.stringify({repo_path:'/unrelated-previous-repository',runtime:{workers:2},quality:{min_confidence:0.7}})));
  await page.getByRole('button',{name:'Run',exact:true}).first().click();
  assert.equal(await page.getByLabel('Repository absolute path',{exact:true}).inputValue(),expectedPath);
  await page.getByRole('button',{name:'Run deterministic extraction',exact:true}).click();
  await page.waitForURL(/#\/runs\//);
  const runID=decodeURIComponent(page.url().split('/runs/')[1]);
  let result;
  for(let i=0;i<300;i++){
   result=await(await fetch(url+'/api/runs/'+encodeURIComponent(runID)+'/state')).json();
   if(result.state?.status==='completed'||result.state?.status==='failed')break;
   await sleep(100);
  }
  assert.equal(result.state?.status,'completed',JSON.stringify(result.state));
  assert.ok(result.events.length>0,'live run retains replayable events');
  await page.locator('.status-pill.completed').waitFor();
  await scan('completed-run',width);
  rows.push({stage:'live-extraction',width,run_id:runID,status:result.state.status,events:result.events.length});
  await page.getByRole('button',{name:'Graph',exact:true}).click();
  await sleep(400);await scan('outcome-graph',width);
  await page.keyboard.press('Escape');
 }
 assert.deepEqual(errors,[]);
 assert.equal(rows.reduce((n,r)=>n+(r.violations?.length||0),0),0,'accessibility violations');
 assert.ok(rows.filter(r=>r.size).every(r=>r.size.document<=r.width),'horizontal document overflow');
 console.log(JSON.stringify(rows.map(r=>({stage:r.stage,width:r.width,violations:r.violations?.map(x=>x.id),size:r.size}))));
}finally{
 if(browser)await browser.close();
 server.kill('SIGTERM');await Promise.race([new Promise(r=>server.once('exit',r)),sleep(5000)]);
 fs.closeSync(log);save();
}
