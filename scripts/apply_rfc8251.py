#!/usr/bin/env python3
"""Reapply only RFC 8251 section 9 after regenerating the RFC 6716 Go core.

Idempotent. Each old/new anchor must occur exactly once; mixed or unknown output
is rejected before any write. --check validates the maintained generated file.
"""
import argparse
import pathlib

root = pathlib.Path(__file__).resolve().parents[1]
old_condition = '\t\tif resynth != 0 && M*int32(*(*opus_int16)(unsafe.Pointer(eBands + uintptr(i)*2)))-N >= M*int32(*(*opus_int16)(unsafe.Pointer(eBands + uintptr(start)*2))) && (update_lowband != 0 || lowband_offset == 0) {'
new_condition = '\t\tif resynth != 0 && (M*int32(*(*opus_int16)(unsafe.Pointer(eBands + uintptr(i)*2)))-N >= M*int32(*(*opus_int16)(unsafe.Pointer(eBands + uintptr(start)*2))) || i == start+1) && (update_lowband != 0 || lowband_offset == 0) {'
old_copy = '\t\ttf_change = *(*int32)(unsafe.Pointer(tf_res + uintptr(i)*4))'
new_copy = '''		// RFC 8251 section 9: extend the first Hybrid band's folding data.
		// n2 == n1 in CELT-only mode, so that path copies no data.
		if i == start+1 {
			n1 := M * (int32(*(*opus_int16)(unsafe.Pointer(eBands + uintptr(start+1)*2))) - int32(*(*opus_int16)(unsafe.Pointer(eBands + uintptr(start)*2))))
			n2 := M * (int32(*(*opus_int16)(unsafe.Pointer(eBands + uintptr(start+2)*2))) - int32(*(*opus_int16)(unsafe.Pointer(eBands + uintptr(start+1)*2))))
			offset := M * int32(*(*opus_int16)(unsafe.Pointer(eBands + uintptr(start)*2)))
			crtXmemcpy(tls, norm+uintptr(offset+n1)*4, norm+uintptr(offset+2*n1-n2)*4, crtTsize_t((n2-n1)*4))
			if C == 2 {
				crtXmemcpy(tls, norm2+uintptr(offset+n1)*4, norm2+uintptr(offset+2*n1-n2)*4, crtTsize_t((n2-n1)*4))
			}
		}
		tf_change = *(*int32)(unsafe.Pointer(tf_res + uintptr(i)*4))'''
old_end = '\t\t\t\tif !(M*int32(*(*opus_int16)(unsafe.Pointer(eBands + uintptr(v1)*2))) < effective_lowband+N) {'
new_end = '\t\t\t\tif !(v1 < i && M*int32(*(*opus_int16)(unsafe.Pointer(eBands + uintptr(v1)*2))) < effective_lowband+N) {'


def transform(text):
    anchors = [(old_condition, new_condition), (old_copy, new_copy), (old_end, new_end)]
    # The unchanged tf_change line is also the last line of the new block.
    applied = [text.count(new) == 1 for _, new in anchors]
    if all(applied):
        if text.count(old_condition) or text.count(old_end):
            raise ValueError('mixed RFC 6716/RFC 8251 output')
        return text
    if any(applied) or any(text.count(old) != 1 for old, _ in anchors):
        raise ValueError('unknown or partially patched quant_all_bands; inspect regeneration')
    for old, new in anchors:
        text = text.replace(old, new, 1)
    return text


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--file', type=pathlib.Path, default=root / 'lib.go')
    parser.add_argument('--check', action='store_true')
    args = parser.parse_args()
    before = args.file.read_text()
    try:
        after = transform(before)
    except ValueError as error:
        raise SystemExit(str(error)) from error
    if args.check and before != after:
        raise SystemExit('RFC 8251 section 9 is not applied')
    if before != after:
        args.file.write_text(after)
    print('RFC 8251 section 9 verified')


if __name__ == '__main__':
    main()
