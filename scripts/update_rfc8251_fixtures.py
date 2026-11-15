#!/usr/bin/env python3
"""Refresh frozen packet decode expectations from an independently built C oracle.
See docs/rfc8251-section9.md. No packet generation or Go codec is used here.
"""
import argparse
import base64
import hashlib
import json
import pathlib
import struct
import subprocess
import tempfile

root = pathlib.Path(__file__).resolve().parents[1]
parser = argparse.ArgumentParser(description=__doc__)
parser.add_argument('--reference', type=pathlib.Path, required=True, help='C RFC 6716 + section 9 decoder driver')
parser.add_argument('--legacy-reference', type=pathlib.Path, required=True, help='original C RFC 6716 driver')
parser.add_argument('--check', action='store_true', help='check without changing the fixture')
args = parser.parse_args()
fixture = root / 'testdata/rfc8251-hybrid.json'
data = json.loads(fixture.read_text())
data['oracle'] = {'rfc6716_archive_sha1': '86a927223e73d2476646a1b933fcd3fffb6ecc8c', 'patch_sha256': hashlib.sha256((root / 'patches/rfc8251-section9.c.patch').read_bytes()).hexdigest(), 'flags': '-O2 -ffp-contract=off; float C; section 9 only'}
with tempfile.TemporaryDirectory(prefix='rfc8251-fixtures-') as directory:
    temp = pathlib.Path(directory)
    for case in data['cases']:
        packets = [base64.b64decode(p, validate=True) for p in case['packets']]
        packed = b''.join(struct.pack('<I', len(p)) + p for p in packets)
        (temp / 'packets').write_bytes(packed)
        outputs = []
        for rate in (16000, 24000, 48000):
            frame = case['frame'] * rate // case['rate']
            expected = {'rate': rate, 'frame': frame}
            for kind in ('int16', 'float'):
                blobs = []
                for reference in (args.reference, args.legacy_reference):
                    subprocess.run([str(reference.resolve()), str(rate), str(case['channels']), kind, str(temp / 'packets'), str(temp / 'pcm'), str(temp / 'ranges'), str(frame)], check=True, capture_output=True)
                    pcm = (temp / 'pcm').read_bytes()
                    width = 2 if kind == 'int16' else 4
                    assert len(pcm) == len(packets) * frame * case['channels'] * width
                    blobs.append(pcm)
                    ranges = [int(r) for r in (temp / 'ranges').read_text().split()]
                    assert len(ranges) == len(packets)
                    if reference == args.reference:
                        expected['final_ranges'] = ranges
                expected[kind + '_sha256'] = hashlib.sha256(blobs[0]).hexdigest()
                expected['legacy_' + kind + '_sha256'] = hashlib.sha256(blobs[1]).hexdigest()
            outputs.append(expected)
        case['outputs'] = outputs
encoded = json.dumps(data, indent=2) + '\n'
if args.check:
    if fixture.read_text() != encoded:
        raise SystemExit('C fixture output differs; inspect the oracle configuration and diff')
else:
    fixture.write_text(encoded)
print('C oracle fixture expectations verified' if args.check else 'C oracle fixture expectations refreshed')
