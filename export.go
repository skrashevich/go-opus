package opus

import "modernc.org/libc"

// Exported entry points of the libopus C API (opus.h) for callers outside
// this package, such as cmd/opus_demo. They follow the transpiled calling
// convention: every call takes a *libc.TLS, pointers are uintptr values into
// C-managed or heap memory, and ctl requests take a libc.VaList.

func EncoderCreate(tls *libc.TLS, fs int32, channels int32, application int32, errPtr uintptr) uintptr {
	return opus_encoder_create(tls, fs, channels, application, errPtr)
}

func Encode(tls *libc.TLS, st uintptr, pcm uintptr, frameSize int32, data uintptr, maxDataBytes int32) int32 {
	return opus_encode(tls, st, pcm, frameSize, data, maxDataBytes)
}

func EncoderCtl(tls *libc.TLS, st uintptr, request int32, va uintptr) int32 {
	return opus_encoder_ctl(tls, st, request, va)
}

func EncoderDestroy(tls *libc.TLS, st uintptr) {
	opus_encoder_destroy(tls, st)
}

func DecoderCreate(tls *libc.TLS, fs int32, channels int32, errPtr uintptr) uintptr {
	return opus_decoder_create(tls, fs, channels, errPtr)
}

func Decode(tls *libc.TLS, st uintptr, data uintptr, length int32, pcm uintptr, frameSize int32, decodeFEC int32) int32 {
	return opus_decode(tls, st, data, length, pcm, frameSize, decodeFEC)
}

func DecoderCtl(tls *libc.TLS, st uintptr, request int32, va uintptr) int32 {
	return opus_decoder_ctl(tls, st, request, va)
}

func DecoderDestroy(tls *libc.TLS, st uintptr) {
	opus_decoder_destroy(tls, st)
}

// Strerror returns a pointer to a NUL-terminated message for an error code.
func Strerror(tls *libc.TLS, code int32) uintptr {
	return opus_strerror(tls, code)
}

// GetVersionString returns a pointer to the NUL-terminated version string.
func GetVersionString(tls *libc.TLS) uintptr {
	return opus_get_version_string(tls)
}
