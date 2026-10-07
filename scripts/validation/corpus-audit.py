#!/usr/bin/env python3
"""Run extraction and Protocol validation against pinned local source snapshots.

Input is JSONL with repository, commit, path and status=downloaded (or empty).
No code is downloaded, installed or executed from the source repositories.
Extraction success is not a claim of complete architecture coverage.
"""
import argparse
import collections
import concurrent.futures
import json
from pathlib import Path
import subprocess
import time


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--inventory', type=Path, required=True)
    parser.add_argument('--binary', type=Path, required=True)
    parser.add_argument('--out', type=Path, required=True)
    parser.add_argument('--workers', type=int, default=2)
    parser.add_argument('--timeout', type=int, default=600)
    args = parser.parse_args()
    if args.workers < 1 or args.workers > 8:
        parser.error('workers must be between 1 and 8')
    binary = args.binary.resolve()
    repos = [json.loads(line) for line in args.inventory.read_text().splitlines() if line.strip()]
    args.out.mkdir(parents=True, exist_ok=True)

    def run(repo):
        record = {key: repo.get(key) for key in ('repository', 'commit', 'archived', 'language')}
        if repo['status'] != 'downloaded':
            return {**record, 'status': repo['status']}
        name = repo['repository'].split('/')[-1]
        if name in ('', '.', '..') or '\\' in name:
            return {**record, 'status': 'failed', 'error': 'Invalid repository name'}
        base = args.out / name
        base.mkdir(parents=True, exist_ok=True)
        started = time.monotonic()
        try:
            with (base / 'extractor.log').open('w') as log:
                result = subprocess.run([str(binary), 'run', '--repo', repo['path'], '--out', str(base), '--workers', '3'], stdout=log, stderr=subprocess.STDOUT, timeout=args.timeout)
            record['exit_code'] = result.returncode
            manifests = sorted(base.glob('*/run_manifest.json'))
            if result.returncode or not manifests:
                raise RuntimeError(f'Extractor exit {result.returncode}; manifest available={bool(manifests)}')
            manifest_path = manifests[-1]
            manifest = json.loads(manifest_path.read_text())
            with (base / 'validation.log').open('w') as log:
                result = subprocess.run([str(binary), 'validate', '--out', str(base), '--run', manifest['run_id']], stdout=log, stderr=subprocess.STDOUT, timeout=args.timeout)
            if result.returncode:
                raise RuntimeError(f'Protocol validation exit {result.returncode}')
            metrics = manifest.get('repo_metrics') or {}
            inventory = metrics.get('dependency_inventory') or {}
            record.update(
                status='completed', counts=manifest['counts'], warnings=manifest.get('warnings') or [],
                run_dir=str(manifest_path.parent), detector_revision=metrics.get('detector_revision'),
                dependencies=len(inventory.get('dependencies') or []),
                version_limits=inventory.get('limitations') or [],
                detector_coverage=metrics.get('detector_coverage') or [],
                source_languages=metrics.get('languages') or [],
            )
        except Exception as error:
            record.update(status='failed', error=str(error))
        record['seconds'] = round(time.monotonic() - started, 2)
        return record

    records = []
    with (args.out / 'results.jsonl').open('w') as log, concurrent.futures.ThreadPoolExecutor(max_workers=args.workers) as pool:
        for index, record in enumerate(pool.map(run, repos), 1):
            records.append(record)
            log.write(json.dumps(record) + '\n')
            log.flush()
            print(f"{index}/{len(repos)} {record['repository']} {record['status']} {record.get('counts', '')}", flush=True)
    summary = {
        'repositories': len(records),
        'active': sum(not record.get('archived') for record in records),
        'archived': sum(bool(record.get('archived')) for record in records),
        'statuses': dict(collections.Counter(record['status'] for record in records)),
        'coverage_states': dict(collections.Counter(item['status'] for record in records for item in record.get('detector_coverage', []))),
        'repos_with_version_limits': sum(bool(record.get('version_limits')) for record in records),
        'repos_without_extracted_evidence': sum(record.get('status') == 'completed' and not any(record.get('counts', {}).get(key, 0) for key in ('exposures', 'dependencies', 'connections')) for record in records),
        'failures': [record for record in records if record['status'] == 'failed'],
        'scope': 'Pinned source snapshots; deterministic extraction and Protocol validation. No runtime execution or completeness oracle.',
    }
    (args.out / 'summary.json').write_text(json.dumps(summary, indent=2))
    print(json.dumps(summary, indent=2), flush=True)
    return 1 if summary['failures'] else 0


if __name__ == '__main__':
    raise SystemExit(main())
