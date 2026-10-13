# Changelog

## v0.1.0

- Add the public pion/opus decoder and encoder API, including PCM slice methods,
  encoder options and controls, packet loss concealment, and Ogg reader support.
- Support direct imports of `github.com/skrashevich/go-opus` and a single remote
  `replace github.com/pion/opus => github.com/skrashevich/go-opus v0.1.0` directive.
- Keep the private C runtime in the root package so replacement consumers do not
  need a second module identity or a local adapter.
- Retain the existing low-level TLS, allocation and codec API. No CGo or external
  module dependencies are required. The minimum Go version is 1.26.1.

Compatibility covers the public API of pion/opus v0.1.0 and upstream commit
`44637de087b3f506170ab738ed96ac1be338006b`. Codec output is not promised to be
bit-identical to Pion's implementation.
