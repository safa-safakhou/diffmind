import fs from 'node:fs';
import path from 'node:path';
import net from 'node:net';
import assert from 'node:assert/strict';
import {spawn} from 'node:child_process';
import {randomBytes} from 'node:crypto';
import {pathToFileURL} from 'node:url';

const [binary, home, output] = process.argv.slice(2);
assert.ok(binary && home && output, 'BINARY EXISTING_PRIVATE_HOME NEW_EVIDENCE_DIRECTORY');
fs.mkdirSync(output, {mode:0o700});
const {chromium}=await import(pathToFileURL(process.env.DIFFMIND_PLAYWRIGHT_MODULE).href);
const {default:AxeBuilder}=await import(pathToFileURL(process.env.DIFFMIND_AXE_MODULE).href);
const serverSocket=net.createServer();
await new Promise(r=>serverSocket.listen(0,'127.0.0.1',r));
const port=serverSocket.address().port;await new Promise(r=>serverSocket.close(r));
const secret=randomBytes(32).toString('hex');
const env=Object.fromEntries(Object.entries(process.env).filter(([k])=>!k.startsWith('DIFFMIND_')&&!['GITHUB_TOKEN','GH_TOKEN'].includes(k)));
const log=fs.openSync(path.join(output,'private-server.log'),'w',0o600);
const server=spawn(binary,['ui','--no-spa-rebuild','--host','127.0.0.1','--port',String(port),'--refresh-interval','0','--refresh-on-start=false','--project-access','scoped'],{env:{...env,DIFFMIND_HOME:home,DIFFMIND_BINARY:binary,DIFFMIND_TRUSTED_PROXY_SECRET:secret},stdio:['ignore',log,log]});
const url='http://127.0.0.1:'+port;
const headers={'X-DiffMind-Proxy-Secret':secret,'X-DiffMind-User':'audit-admin','X-DiffMind-Role':'admin'};
const rows=[], errors=[];let browser;
const sleep=ms=>new Promise(r=>setTimeout(r,ms));
const save=()=>fs.writeFileSync(path.join(output,'browser-results.json'),JSON.stringify({rows,errors},null,2));
try{
 for(let i=0;i<100;i++){try{if((await fetch(url+'/healthz')).ok)break;}catch{}await sleep(100);}
 const projects=await (await fetch(url+'/api/projects',{headers})).json();
 assert.ok(projects.projects.length);
 const pid=projects.projects[0].id;
 browser=await chromium.launch({headless:true,executablePath:process.env.DIFFMIND_CHROMIUM});
 const context=await browser.newContext({extraHTTPHeaders:headers});
 const page=await context.newPage();
 page.on('pageerror',e=>errors.push(e.message));
 const checkGraphSpace=async()=>{
  const layout=await page.evaluate(()=>{
   const box=s=>{const r=document.querySelector(s).getBoundingClientRect();return {top:r.top,bottom:r.bottom,left:r.left,right:r.right,width:r.width,height:r.height};};
   const board=box('.workspace-board'),toolbar=box('.graph-mode-toolbar'),header=box('.workspace-topbar'),readiness=box('.workspace-alerts');
   const controls=[...document.querySelectorAll('.workspace-topbar button,.workspace-topbar h1,.graph-mode-toolbar input,.graph-mode-toolbar select,.graph-mode-toolbar button')].map(e=>{const r=e.getBoundingClientRect();return {left:r.left,right:r.right,top:r.top,bottom:r.bottom};});
   return {board,toolbar,header,readiness,controls,viewport:innerWidth,height:innerHeight};
  });
  assert.ok(layout.readiness.height>=100,'readiness must remain readable above the graph');
  assert.ok(layout.board.height>=320,'graph needs usable height');
  assert.ok(layout.board.bottom-layout.toolbar.bottom>=180,'toolbar must leave visible graph space');
  assert.ok(layout.controls.every(r=>r.left>=0&&r.right<=layout.viewport+1),'workspace controls must fit horizontally: '+JSON.stringify(layout));
  assert.ok(layout.controls.every(r=>r.top>=-1),'workspace controls must not clip above the page: '+JSON.stringify(layout));
  return layout;
 };
 const routes=[['projects','/'],['workspace','/projects/'+pid],['operations','/projects/'+pid+'/operations'],['access','/projects/'+pid+'/access'],['pull-requests','/projects/'+pid+'/pull-requests'],['compare','/projects/'+pid+'/compare']];
 for(const width of [1440,768,390]){
  await page.setViewportSize({width,height:1000});
  for(const [name,route] of routes){
   const started=performance.now();
   await page.goto(url+'/#'+route);
   await page.waitForLoadState('networkidle');
   if(name==='workspace')await page.locator('.compact-service-name').first().waitFor();
   await page.waitForTimeout(350);
   const loadMs=Math.round(performance.now()-started);
   const audit=await new AxeBuilder({page}).withTags(['wcag2a','wcag2aa','wcag21a','wcag21aa']).analyze();
   if(name==='workspace')await page.evaluate(()=>{scrollTo(0,0);document.querySelector('.workspace').scrollTop=0;});
   const layout=name==='workspace'?await checkGraphSpace():undefined;
   const size=await page.evaluate(()=>({viewport:innerWidth,document:document.documentElement.scrollWidth}));
   rows.push({name,width,load_ms:loadMs,size,layout,violations:audit.violations.map(v=>({id:v.id,impact:v.impact,description:v.description,nodes:v.nodes.map(n=>({target:n.target,html:n.html,summary:n.failureSummary}))}))});
   await page.screenshot({path:path.join(output,name+'-'+width+'.png'),fullPage:true});
   save();console.log(name,width,'violations',audit.violations.map(v=>v.id).join(','),'overflow',size.document-width);
  }
 }
 // A short phone must retain a usable, scroll-reachable graph and inspector.
 await page.setViewportSize({width:390,height:640});
 await page.goto(url+'/#/projects/'+pid);await page.locator('.compact-service-name').first().waitFor();
 await page.waitForLoadState('networkidle');
 await page.evaluate(()=>{scrollTo(0,0);document.querySelector('.workspace').scrollTop=0;});
 const phoneLayout=await checkGraphSpace();
 await page.getByLabel('Graph team',{exact:true}).selectOption('checkout');
 const search=page.getByLabel('Search service',{exact:true});
 await search.fill('checkout-api');await search.press('Enter');
 const phoneNode=page.locator('.service-system[data-select-id="checkout-api"]');
 await phoneNode.waitFor();await phoneNode.focus();await page.keyboard.press('Enter');
 assert.match(await page.locator('.workspace-right').innerText(),/checkout-api/);
 await page.locator('.workspace-right').scrollIntoViewIfNeeded();
 await page.screenshot({path:path.join(output,'phone-graph-selection.png'),fullPage:true});
 rows.push({name:'short-phone-graph',layout:phoneLayout,selected:true});
 // Real keyboard-only dialog entry, trap, dismissal and focus restoration.
 await page.setViewportSize({width:1440,height:1000});
 await page.goto(url+'/#/projects/'+pid);await page.getByLabel('Workspace readiness').waitFor();
 const opener=page.getByRole('button',{name:'Knowledge packs',exact:true});
 await opener.focus();await page.keyboard.press('Enter');
 const dialog=page.getByRole('dialog');await dialog.waitFor();
 let trapped=true;
 for(let i=0;i<15;i++){await page.keyboard.press(i%3===0?'Shift+Tab':'Tab');trapped&&=await dialog.evaluate(d=>d.contains(document.activeElement));}
 await page.keyboard.press('Escape');await dialog.waitFor({state:'hidden'});
 const restored=await opener.evaluate(e=>e===document.activeElement);
 rows.push({name:'keyboard-dialog',trapped,restored});assert.ok(trapped&&restored);
 await page.getByLabel('Graph team',{exact:true}).selectOption('checkout');
 await page.getByLabel('Search service',{exact:true}).fill('checkout-api');
 await page.getByLabel('Search service',{exact:true}).press('Enter');
 const node=page.locator('.service-system[data-select-id="checkout-api"]');
 await node.waitFor();await node.focus();await page.keyboard.press('Enter');
 await page.screenshot({path:path.join(output,'keyboard-graph-selection.png'),fullPage:true});
 rows.push({name:'keyboard-graph-selection',selected:await node.count(),details:await page.locator('.workspace-right').innerText()});
 assert.match(rows.at(-1).details,/checkout-api/);
 await page.getByRole('button',{name:'Full detail',exact:true}).click();
 await page.waitForLoadState('networkidle');await page.waitForTimeout(350);
 const detailed=await new AxeBuilder({page}).withTags(['wcag2a','wcag2aa','wcag21a','wcag21aa']).analyze();
 rows.push({name:'full-detail',violations:detailed.violations.map(v=>({id:v.id,impact:v.impact,nodes:v.nodes.map(n=>({target:n.target,html:n.html,summary:n.failureSummary}))}))});
 assert.deepEqual(errors,[]);save();
 assert.equal(rows.reduce((n,r)=>n+(r.violations?.length||0),0),0,'accessibility violations');

}finally{
 if(browser)await browser.close();
 const closed=new Promise(r=>server.once('exit',r));server.kill('SIGTERM');await Promise.race([closed,sleep(5000)]);
 fs.closeSync(log);save();
}
