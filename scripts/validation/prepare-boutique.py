#!/usr/bin/env python3
"""Prepare an absent private trial root; clone and copy public source, never run it."""
import hashlib
import json
from pathlib import Path
import re
import shutil
import subprocess
import sys

PIN = "b9a978db9e01f4ad3dca9494a22cb9edc17548fe"
NAMES = ["frontend", "cartservice", "productcatalogservice", "currencyservice",
         "shippingservice", "recommendationservice", "adservice", "checkoutservice"]
VARIABLES = {"PRODUCT_CATALOG", "CURRENCY", "CART", "RECOMMENDATION",
             "SHIPPING", "CHECKOUT", "AD"}
if len(sys.argv) != 2:
    raise SystemExit("usage: python3 prepare-boutique.py NEW_PRIVATE_TRIAL_ROOT")
root = Path(sys.argv[1]).resolve()
root.mkdir(mode=0o700)  # Must be absent; never overwrite an earlier trial.
(root / "evidence").mkdir()
source = root / "boutique"
subprocess.run(["git", "clone", "--filter=blob:none", "--no-checkout",
                "https://github.com/GoogleCloudPlatform/microservices-demo.git",
                str(source)], check=True)
subprocess.run(["git", "-C", str(source), "checkout", "--detach", PIN], check=True)
labels, mapping = [], []
for name in NAMES:
    destination = root / "repositories" / name
    shutil.copytree(source / "src" / name, destination)
    upstream = source / "kubernetes-manifests" / (name + ".yaml")
    (destination / "deployment.yaml").write_bytes(upstream.read_bytes())
    mapping.append({"repository": name, "deployment_source": str(upstream.relative_to(source)),
                    "sha256": hashlib.sha256(upstream.read_bytes()).hexdigest()})
    lines = upstream.read_text().splitlines()
    for index, line in enumerate(lines):
        variable = re.fullmatch(r"\s*- name: ([A-Z_]+)_SERVICE_ADDR\s*", line)
        if variable and variable.group(1) in VARIABLES:
            value = re.fullmatch(r'\s*value: "([a-z0-9-]+):([0-9]+)"\s*', lines[index+1])
            if not value:
                raise ValueError("Unreviewed declaration syntax")
            labels.append({"from": name, "to": value.group(1), "type": "rpc",
                           "source": str(upstream.relative_to(source)), "line": index+1,
                           "basis": "literal configured address, declared evidence only"})
(root / "evidence" / "source-labels.json").write_text(json.dumps({
    "revision": PIN,
    "label_author": "implementation-agent source-first review before candidate output; no independent blind reviewer",
    "scope": "eight src directories plus exact corresponding deployment files copied without execution",
    "mappings": mapping, "expected_declared_edges": labels,
    "excluded": "Other variable names, comments, placeholders, indirect/runtime calls; source-extracted extras remain unlabeled"
}, indent=2) + "\n")
print(root)
