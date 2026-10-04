#!/usr/bin/env python3
"""Exercise the installed model host as a DiffMind client; no task delegation.

Usage: installed-host.py SPEC_JSON NEW_OUTPUT_DIRECTORY
A spec supplies binary, mode (mcp/agent/http), home or url, work, prompt and
optional token_env and approved_tools (explicit test-scope authorization). The existing Codex model/effort settings are preserved;
unrelated MCP configuration is not loaded. Raw private transcripts may contain
tool data: review/redact them before retaining evidence in Git.
"""
import hashlib,json,os,signal,subprocess,sys,time,tomllib
from pathlib import Path

spec=json.loads(Path(sys.argv[1]).read_text())
out=Path(sys.argv[2]).resolve();out.mkdir(mode=0o700)
work=Path(spec["work"]).resolve();work.mkdir(parents=True,exist_ok=True)
profile=Path(os.environ.get("CODEX_HOME",str(Path.home()/".codex")))/"config.toml"
settings=tomllib.loads(profile.read_text()) if profile.exists() else {}
args=["codex","exec","--ignore-user-config","--ephemeral","--skip-git-repo-check",
      "--sandbox","read-only","--disable","apps","--json","--color","never","-C",str(work),
      "-c",'approval_policy="never"']
for key in ("model","model_reasoning_effort"):
 if settings.get(key):args+=["-c",key+"="+json.dumps(settings[key])]
if spec["mode"]=="http":
 args+=["-c","mcp_servers.diffmind.url="+json.dumps(spec["url"])]
 if spec.get("token_env"):
  args+=["-c","mcp_servers.diffmind.bearer_token_env_var="+json.dumps(spec["token_env"])]
else:
 assert spec["mode"] in ("mcp","agent")
 args+=["-c","mcp_servers.diffmind.command="+json.dumps(str(Path(spec["binary"]).resolve())),
        "-c","mcp_servers.diffmind.args="+json.dumps([spec["mode"]]),
        "-c","mcp_servers.diffmind.env={ DIFFMIND_HOME = "+json.dumps(str(Path(spec["home"]).resolve()))+" }"]
approved_tools=spec.get("approved_tools",[])
assert all(name in ("manage_workspace","agent_runtime") for name in approved_tools),"unsupported test approval"
for name in approved_tools:
 args+=["-c",f'mcp_servers.diffmind.tools.{name}.approval_mode="approve"']
args+=["-o",str(out/"answer.txt"),"-"]
env={k:v for k,v in os.environ.items() if not k.startswith("DIFFMIND_") and k not in ("GITHUB_TOKEN","GH_TOKEN")}
if spec.get("token_env"):
 assert spec["token_env"] in os.environ,"required test identity is unavailable"
 env[spec["token_env"]]=os.environ[spec["token_env"]]
start=time.monotonic();timed_out=False
with (out/"private-events.jsonl").open("w") as events,(out/"private-stderr.txt").open("w") as err:
 process=subprocess.Popen(args,stdin=subprocess.PIPE,stdout=events,stderr=err,text=True,env=env,start_new_session=True)
 try:process.communicate(spec["prompt"],timeout=spec.get("timeout_seconds",300))
 except subprocess.TimeoutExpired:
  timed_out=True;os.killpg(process.pid,signal.SIGTERM)
  try:process.communicate(timeout=10)
  except subprocess.TimeoutExpired:os.killpg(process.pid,signal.SIGKILL);process.communicate()
records=[]
for line in (out/"private-events.jsonl").read_text().splitlines():
 try:records.append(json.loads(line))
 except json.JSONDecodeError:pass
items=[e["item"] for e in records if e.get("type")=="item.completed"]
calls=[{"server":x.get("server"),"tool":x["tool"],"status":x.get("status"),"error":x.get("error"),"arguments":x.get("arguments")} for x in items if x.get("type")=="mcp_tool_call"]
result={"id":spec["id"],"mode":spec["mode"],"host_version":subprocess.check_output(["codex","--version"],text=True).strip(),
        "model":settings.get("model","host default"),"effort":settings.get("model_reasoning_effort","host default"),
        "exit_code":process.returncode,"timed_out":timed_out,"elapsed_seconds":round(time.monotonic()-start,3),
        "prompt":spec["prompt"],"tools":calls,
        "shell_commands":[x.get("command") for x in items if x.get("type")=="command_execution"],
        "answer_exists":(out/"answer.txt").exists(),"unrelated_mcp_config_loaded":False,"approved_tools":approved_tools,"apps_feature_disabled":True,
        "binary_sha256":hashlib.sha256(Path(spec["binary"]).read_bytes()).hexdigest()}
(out/"result.json").write_text(json.dumps(result,indent=2)+"\n")
print(json.dumps({k:result[k] for k in ("id","exit_code","timed_out","elapsed_seconds","answer_exists")}),flush=True)
sys.exit(1 if timed_out else process.returncode)
