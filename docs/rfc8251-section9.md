# RFC 8251 section 9: Hybrid folding

The codec applies the Hybrid folding update from RFC 8251 section 9 to the RFC
6716 core. The change extends folding data from the first Hybrid CELT band,
selects that data for the next band, and bounds the collapse-mask lookup to the
bands already decoded. It changes some Hybrid PCM reconstruction without changing
packet parsing, entropy final ranges, encoder packets, public APIs or dependencies.

This is section 9 only. It is not a claim of complete RFC 8251 conformance.

## Generated-source maintenance

The repository contains an optimized/transpiled lib.go; it does not contain a
complete end-to-end driver for regenerating every ccgo object and local kernel
optimization. Preserve section 9 in either of these source workflows:

1. Apply patches/rfc8251-section9.c.patch to the normative C source before
   transpilation, then reconcile the resulting quant_all_bands with this update.
2. When regenerating from the original RFC 6716 C source, run **go generate .**
   after the existing transpilation/runtime adaptation steps. The directive is
   kept in the handwritten export.go. scripts/apply_rfc8251.py applies exactly
   three bounded edits to quant_all_bands and is idempotent. It rejects mixed or
   unknown generated output before writing; inspect changed generator output
   rather than weakening that guard.

Checks:

    python3 scripts/apply_rfc8251.py --check
    python3 scripts/test_apply_rfc8251.py
    go generate .
    git diff --exit-code -- lib.go

The post-generation tests reconstruct the original anchors, reapply the change
and compare the result byte for byte with maintained lib.go; they also verify
idempotence, read-only checking and rejection without writes for partial/drifted
input. A complete regeneration of all prior hand optimizations is outside this
change's reproducibility claim.

## Independent regression fixtures

The frozen synthetic packets in testdata/rfc8251-hybrid.json cover mono/stereo,
fullband/superwideband, 10/20 ms frames, a short reported regression packet, and
stateful transient streams. CELT and SILK packets are compatibility controls.
There are no station recordings. Each stream has 16/24/48 kHz output expectations.

The expected integer and float PCM hashes and per-packet final ranges come from
**normative C RFC 6716 plus only section 9**, not from this Go implementation.
The tests compare exact hashes; the public float API has the same strict oracle.
The fixture also records the original C integer output hashes: CELT, SILK and
16 kHz output remain identical on these controls. An allocation test requires
zero steady-state public Hybrid float decoder allocations after warmup.

The original normative source is embedded in RFC 6716 Appendix A. The extracted
archive's SHA-1 is **86a927223e73d2476646a1b933fcd3fffb6ecc8c**. For local fixture
reproduction, obtain rfc6716.txt from the RFC Editor and extract its source:

    python3 - <<'PY'
    import base64, hashlib, pathlib, tarfile
    text = pathlib.Path('rfc6716.txt').read_text()
    data = base64.b64decode(''.join(x[6:].strip() for x in text.splitlines() if x.startswith('   ###')))
    assert hashlib.sha1(data).hexdigest() == '86a927223e73d2476646a1b933fcd3fffb6ecc8c'
    pathlib.Path('opus-rfc6716.tar.gz').write_bytes(data)
    with tarfile.open('opus-rfc6716.tar.gz') as tar:
        tar.extractall('reference', filter='data')
    PY

Keep an original reference copy and a second copy with just the supplied C patch.
Build both libraries in their own directories. Set CFLAGS to:

    -O2 -g -ffp-contract=off -DOPUS_BUILD -DUSE_ALLOCA -Drestrict= -Iinclude -Icelt -Isilk -Isilk/float -Isrc

Then use make -j6 lib. Link scripts/rfc8251_reference.c separately against each
libopus.a and the corresponding include/ directory with clang -O2
-ffp-contract=off and -lm. The driver takes rate, channels, int16|float, packet
file, output PCM, final-range file and expected samples/channel. Run:

    python3 scripts/update_rfc8251_fixtures.py \
      --reference /absolute/path/to/updated-c-driver \
      --legacy-reference /absolute/path/to/original-c-driver --check

Omit --check only when intentionally reviewing a fixture change. The script never
regenerates packet inputs and never uses Go output to produce expected PCM.
Compiler/platform numeric changes must be investigated rather than silently
refreshing the strict hashes.

Only two prior combined packet/PCM golden cases change: hybrid stereo 24k and
mix hybrid mono long. Their new hashes are independently recomputed with the
updated C decoder. The original C decoder reproduces all 12 original combined
hashes; all 12 encoder packet hashes are now checked separately. Multistream,
CELT/SILK and the other combined digests retain their existing values.

## Other RFC 8251 changes: scope inventory

This read-only source inspection identifies additional work; it does not validate
all edge cases or promise complete behavior/security equivalence:

| Section | Source observation | This branch |
|---|---|---|
| 3, SILK stereo reset | silk_InitDecoder still returns after channel initialization without the two extra reset statements | Unchanged |
| 4, packet padding | opus_packet_parse_impl still accumulates a separate int32 padding counter | Unchanged |
| 5, resampler buffer/copies | silk_resampler_private_IIR_FIR retains the original 488-int16 scratch declaration and 4-byte copy widths | Unchanged |
| 6, inverse gain overflow | LPC_inverse_pred_gain_QA retains the original arithmetic without the specified saturation and range check | Unchanged |
| 7, NLSF saturation | silk_NLSF_stabilize retains the original addition rather than the specified saturating add | Unchanged |
| 8, energy cap | log2Amp has no section-8 cap before exponentiation | Unchanged |
| 10, optional mono downmix phase behavior | quant_band retains the original inversion behavior; no new phase-control API is introduced | Unchanged |

Section 11 updates official vectors because sections 9/10 change decoder output.
Passing the targeted C regressions and the existing suite is not a substitute for
running the complete updated official vector set. No broader RFC update is
included automatically.

Sources: https://www.rfc-editor.org/rfc/rfc6716.html#appendix-A and
https://www.rfc-editor.org/rfc/rfc8251.html#section-9.
