#!/usr/bin/env python3
"""Run the shipped systemd backup helper against disposable units only.
Usage (root on a systemd host): systemd-native.py BINARY NEW_OUTPUT
Never touches a pre-existing unit, workspace or backup catalog.
"""
import hashlib,json,os,secrets,shutil,socket,subprocess,sys,time
from pathlib import Path
from urllib.request import Request,urlopen
binary,out_s=sys.argv[1:];out=Path(out_s);out.mkdir(mode=0o700)
assert os.geteuid()==0,"isolated system unit creation requires root"
repo=Path(__file__).resolve().parents[2]
prefix="diffmind-validation-"+secrets.token_hex(6);root=Path("/run")/prefix;root.mkdir(mode=0o700)
app=prefix+"-app.service";backup=prefix+"-backup.service"
unit_paths=[Path("/run/systemd/system")/app,Path("/run/systemd/system")/backup]
assert not any(p.exists() for p in unit_paths)
runtime=Path("/run")/(prefix+"-backup")
secret=secrets.token_hex(24);rows=[]
def systemctl(*args):return subprocess.check_output(["systemctl",*args],text=True,stderr=subprocess.STDOUT).strip()
try:
 shutil.copy2(binary,root/"diffmind");shutil.copy2(repo/"scripts/backup-systemd.sh",root/"backup-systemd.sh")
 (root/"home").mkdir();(root/"backups").mkdir(mode=0o700)
 with socket.socket() as sock:sock.bind(("127.0.0.1",0));port=sock.getsockname()[1]
 unit_paths[0].write_text(f"[Unit]\nDescription=Disposable DiffMind validation service\n[Service]\nType=simple\nEnvironment=DIFFMIND_HOME={root}/home\nEnvironment=DIFFMIND_AUTH_TOKEN={secret}\nExecStart={root}/diffmind ui --no-spa-rebuild --host 127.0.0.1 --port {port} --refresh-interval 0 --refresh-on-start=false\nUMask=0077\n")
 env={"DIFFMIND_SERVICE":app,"DIFFMIND_BINARY":str(root/"diffmind"),"DIFFMIND_HOME":str(root/"home"),"DIFFMIND_BACKUP_DIRECTORY":str(root/"backups"),"DIFFMIND_BACKUP_KEEP_LAST":"1","DIFFMIND_BACKUP_LOCK":str(runtime/"maintenance.lock")}
 (root/"backup.env").write_text("".join(k+"="+v+"\n" for k,v in env.items()))
 unit=(repo/"deploy/systemd/diffmind-backup.service").read_text().replace("/etc/diffmind/backup.env",str(root/"backup.env")).replace("/usr/local/libexec/diffmind/backup-systemd.sh",str(root/"backup-systemd.sh")).replace("RuntimeDirectory=diffmind-backup","RuntimeDirectory="+runtime.name)
 unit_paths[1].write_text(unit)
 subprocess.run(["systemd-analyze","verify",str(unit_paths[0]),str(unit_paths[1])],check=True,capture_output=True)
 systemctl("daemon-reload");systemctl("start",app)
 for _ in range(100):
  try:
   with urlopen(f"http://127.0.0.1:{port}/healthz",timeout=1):break
  except OSError:time.sleep(.1)
 req=Request(f"http://127.0.0.1:{port}/api/projects",data=b'{"name":"Systemd fixture"}',headers={"Authorization":"Bearer "+secret,"Content-Type":"application/json"})
 with urlopen(req,timeout=10) as response:assert response.status<300
 for state in ["active","inactive"]:
  if state=="inactive":systemctl("stop",app)
  started=time.monotonic();systemctl("start",backup)
  assert systemctl("show","--property=Result","--value",backup)=="success"
  assert systemctl("show","--property=ActiveState","--value",app)==state
  archives=list((root/"backups").glob("snapshot-*/archive.tar.gz"));assert len(archives)==1,archives
  archive=archives[0];digest=hashlib.sha256(archive.read_bytes()).hexdigest()
  subprocess.run([str(root/"diffmind"),"backup","verify","--archive",str(archive),"--sha256",digest,"--json"],check=True,capture_output=True)
  assert archive.stat().st_mode&0o777==0o600
  assert not Path(env["DIFFMIND_BACKUP_LOCK"]+".restart").exists()
  rows.append({"initial_app_state":state,"final_app_state":state,"backup_result":"success","retained_archives":1,"private_archive":True,"sha256":digest,"integrity_verified":True,"restart_intent_cleared":True,"seconds":round(time.monotonic()-started,3)})
 (out/"results.json").write_text(json.dumps({"passed":True,"binary_sha256":hashlib.sha256(Path(binary).read_bytes()).hexdigest(),"rows":rows,"limits":"Disposable native systemd units with shipped helper and service protections. No production data or unit was used."},indent=2))
 print(json.dumps(rows),flush=True)
finally:
 for name in [backup,app]:subprocess.run(["systemctl","stop",name],capture_output=True)
 for p in unit_paths:
  if p.exists():p.unlink()
 subprocess.run(["systemctl","daemon-reload"],capture_output=True)
 # Delete only this invocation's exact private /run directories.
 for p in [root,runtime]:
  assert p.parent==Path("/run") and p.name.startswith(prefix)
  if p.exists():shutil.rmtree(p)
