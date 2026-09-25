package opus

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

// layoutTypes lists the codec structs whose memory layout the transpiled code
// relies on through hardcoded 64-bit offsets and sizes.
var layoutTypes = []struct {
	name string
	typ  reflect.Type
}{
	{"OpusCustomMode", reflect.TypeFor[OpusCustomMode]()},
	{"ec_ctx", reflect.TypeFor[ec_ctx]()},
	{"ec_enc", reflect.TypeFor[ec_enc]()},
	{"ec_dec", reflect.TypeFor[ec_dec]()},
	{"kiss_fft_cpx", reflect.TypeFor[kiss_fft_cpx]()},
	{"kiss_twiddle_cpx", reflect.TypeFor[kiss_twiddle_cpx]()},
	{"kiss_fft_state", reflect.TypeFor[kiss_fft_state]()},
	{"mdct_lookup", reflect.TypeFor[mdct_lookup]()},
	{"PulseCache", reflect.TypeFor[PulseCache]()},
	{"OpusCustomEncoder", reflect.TypeFor[OpusCustomEncoder]()},
	{"OpusCustomDecoder", reflect.TypeFor[OpusCustomDecoder]()},
	{"OpusRepacketizer", reflect.TypeFor[OpusRepacketizer]()},
	{"OpusDecoder", reflect.TypeFor[OpusDecoder]()},
	{"silk_EncControlStruct", reflect.TypeFor[silk_EncControlStruct]()},
	{"silk_DecControlStruct", reflect.TypeFor[silk_DecControlStruct]()},
	{"silk_TOC_struct", reflect.TypeFor[silk_TOC_struct]()},
	{"silk_resampler_state_struct", reflect.TypeFor[silk_resampler_state_struct]()},
	{"silk_nsq_state", reflect.TypeFor[silk_nsq_state]()},
	{"silk_VAD_state", reflect.TypeFor[silk_VAD_state]()},
	{"silk_LP_state", reflect.TypeFor[silk_LP_state]()},
	{"silk_NLSF_CB_struct", reflect.TypeFor[silk_NLSF_CB_struct]()},
	{"stereo_enc_state", reflect.TypeFor[stereo_enc_state]()},
	{"stereo_dec_state", reflect.TypeFor[stereo_dec_state]()},
	{"SideInfoIndices", reflect.TypeFor[SideInfoIndices]()},
	{"silk_encoder_state", reflect.TypeFor[silk_encoder_state]()},
	{"silk_PLC_struct", reflect.TypeFor[silk_PLC_struct]()},
	{"silk_CNG_struct", reflect.TypeFor[silk_CNG_struct]()},
	{"silk_decoder_state", reflect.TypeFor[silk_decoder_state]()},
	{"silk_decoder_control", reflect.TypeFor[silk_decoder_control]()},
	{"OpusEncoder", reflect.TypeFor[OpusEncoder]()},
	{"silk_shape_state_FLP", reflect.TypeFor[silk_shape_state_FLP]()},
	{"silk_prefilter_state_FLP", reflect.TypeFor[silk_prefilter_state_FLP]()},
	{"silk_encoder_state_FLP", reflect.TypeFor[silk_encoder_state_FLP]()},
	{"silk_encoder_control_FLP", reflect.TypeFor[silk_encoder_control_FLP]()},
	{"silk_encoder", reflect.TypeFor[silk_encoder]()},
	{"OpusMSEncoder", reflect.TypeFor[OpusMSEncoder]()},
	{"OpusMSDecoder", reflect.TypeFor[OpusMSDecoder]()},
	{"ChannelLayout", reflect.TypeFor[ChannelLayout]()},
	{"NSQ_del_dec_struct", reflect.TypeFor[NSQ_del_dec_struct]()},
	{"NSQ_sample_struct", reflect.TypeFor[NSQ_sample_struct]()},
	{"silk_decoder", reflect.TypeFor[silk_decoder]()},
}

func dumpLayout(b *strings.Builder, t reflect.Type, prefix string, base uintptr) {
	if t == reflect.TypeFor[ptrslot]() {
		// One pointer in 8 bytes; its element count depends on the pointer width.
		fmt.Fprintf(b, "%s off=%d size=%d\n", prefix, base, t.Size())
		return
	}
	switch t.Kind() {
	case reflect.Struct:
		for i := range t.NumField() {
			f := t.Field(i)
			if f.Name == "_" {
				continue
			}
			dumpLayout(b, f.Type, prefix+"."+f.Name, base+f.Offset)
		}
	case reflect.Array:
		fmt.Fprintf(b, "%s off=%d size=%d stride=%d\n", prefix, base, t.Size(), t.Elem().Size())
		if t.Len() > 0 {
			dumpLayout(b, t.Elem(), prefix+"[0]", base)
		}
	default:
		fmt.Fprintf(b, "%s off=%d\n", prefix, base)
	}
}

// TestStructLayout checks that every codec struct keeps the layout of the
// 64-bit platform lib.go was generated for: the transpiled code addresses
// fields through hardcoded offsets, so 32-bit targets must match it exactly.
// Regenerate the golden file on a 64-bit host with UPDATE_LAYOUT=1.
func TestStructLayout(t *testing.T) {
	var b strings.Builder
	for _, lt := range layoutTypes {
		fmt.Fprintf(&b, "%s size=%d\n", lt.name, lt.typ.Size())
		dumpLayout(&b, lt.typ, lt.name, 0)
	}
	if os.Getenv("UPDATE_LAYOUT") != "" {
		if err := os.WriteFile("testdata/layout64.txt", []byte(b.String()), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile("testdata/layout64.txt")
	if err != nil {
		t.Fatal(err)
	}
	got := strings.Split(b.String(), "\n")
	for i, w := range strings.Split(string(want), "\n") {
		if i >= len(got) || got[i] != w {
			g := "<missing>"
			if i < len(got) {
				g = got[i]
			}
			t.Fatalf("layout differs at line %d:\n got  %s\n want %s", i+1, g, w)
		}
	}
}
