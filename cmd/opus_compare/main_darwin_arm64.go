// Code generated for darwin/arm64 by 'ccgo celt_bands.o.go celt_celt.o.go celt_celt_lpc.o.go celt_cwrs.o.go celt_entcode.o.go celt_entdec.o.go celt_entenc.o.go celt_kiss_fft.o.go celt_laplace.o.go celt_mathops.o.go celt_mdct.o.go celt_modes.o.go celt_pitch.o.go celt_quant_bands.o.go celt_rate.o.go celt_vq.o.go opus_opus.o.go opus_opus_compare.o.go opus_opus_decoder.o.go opus_opus_encoder.o.go opus_opus_multistream.o.go opus_repacketizer.o.go silk_A2NLSF.o.go silk_CNG.o.go silk_HP_variable_cutoff.o.go silk_LPC_analysis_filter.o.go silk_LPC_inv_pred_gain.o.go silk_LP_variable_cutoff.o.go silk_NLSF2A.o.go silk_NLSF_VQ.o.go silk_NLSF_VQ_weights_laroia.o.go silk_NLSF_decode.o.go silk_NLSF_del_dec_quant.o.go silk_NLSF_encode.o.go silk_NLSF_stabilize.o.go silk_NLSF_unpack.o.go silk_NSQ.o.go silk_NSQ_del_dec.o.go silk_PLC.o.go silk_VAD.o.go silk_VQ_WMat_EC.o.go silk_ana_filt_bank_1.o.go silk_biquad_alt.o.go silk_bwexpander.o.go silk_bwexpander_32.o.go silk_check_control_input.o.go silk_code_signs.o.go silk_control_SNR.o.go silk_control_audio_bandwidth.o.go silk_control_codec.o.go silk_debug.o.go silk_dec_API.o.go silk_decode_core.o.go silk_decode_frame.o.go silk_decode_indices.o.go silk_decode_parameters.o.go silk_decode_pitch.o.go silk_decode_pulses.o.go silk_decoder_set_fs.o.go silk_enc_API.o.go silk_encode_indices.o.go silk_encode_pulses.o.go silk_float_LPC_analysis_filter_FLP.o.go silk_float_LPC_inv_pred_gain_FLP.o.go silk_float_LTP_analysis_filter_FLP.o.go silk_float_LTP_scale_ctrl_FLP.o.go silk_float_apply_sine_window_FLP.o.go silk_float_autocorrelation_FLP.o.go silk_float_burg_modified_FLP.o.go silk_float_bwexpander_FLP.o.go silk_float_corrMatrix_FLP.o.go silk_float_encode_frame_FLP.o.go silk_float_energy_FLP.o.go silk_float_find_LPC_FLP.o.go silk_float_find_LTP_FLP.o.go silk_float_find_pitch_lags_FLP.o.go silk_float_find_pred_coefs_FLP.o.go silk_float_inner_product_FLP.o.go silk_float_k2a_FLP.o.go silk_float_levinsondurbin_FLP.o.go silk_float_noise_shape_analysis_FLP.o.go silk_float_pitch_analysis_core_FLP.o.go silk_float_prefilter_FLP.o.go silk_float_process_gains_FLP.o.go silk_float_regularize_correlations_FLP.o.go silk_float_residual_energy_FLP.o.go silk_float_scale_copy_vector_FLP.o.go silk_float_scale_vector_FLP.o.go silk_float_schur_FLP.o.go silk_float_solve_LS_FLP.o.go silk_float_sort_FLP.o.go silk_float_warped_autocorrelation_FLP.o.go silk_float_wrappers_FLP.o.go silk_gain_quant.o.go silk_init_decoder.o.go silk_init_encoder.o.go silk_inner_prod_aligned.o.go silk_interpolate.o.go silk_lin2log.o.go silk_log2lin.o.go silk_pitch_est_tables.o.go silk_process_NLSFs.o.go silk_quant_LTP_gains.o.go silk_resampler.o.go silk_resampler_down2.o.go silk_resampler_down2_3.o.go silk_resampler_private_AR2.o.go silk_resampler_private_IIR_FIR.o.go silk_resampler_private_down_FIR.o.go silk_resampler_private_up2_HQ.o.go silk_resampler_rom.o.go silk_shell_coder.o.go silk_sigm_Q15.o.go silk_sort.o.go silk_stereo_LR_to_MS.o.go silk_stereo_MS_to_LR.o.go silk_stereo_decode_pred.o.go silk_stereo_encode_pred.o.go silk_stereo_find_predictor.o.go silk_stereo_quant_pred.o.go silk_sum_sqr_shift.o.go silk_table_LSF_cos.o.go silk_tables_LTP.o.go silk_tables_NLSF_CB_NB_MB.o.go silk_tables_NLSF_CB_WB.o.go silk_tables_gain.o.go silk_tables_other.o.go silk_tables_pitch_lag.o.go silk_tables_pulses_per_block.o.go', DO NOT EDIT.

//go:build darwin && arm64

package main

import (
	"reflect"
	"unsafe"

	"modernc.org/libc"
)

var _ reflect.Type
var _ unsafe.Pointer

const EXIT_FAILURE = 1
const EXIT_SUCCESS = 0

type size_t = uint64

const NBANDS = 21
const NFREQS = 240
const TEST_WIN_SIZE = 480
const TEST_WIN_STEP = 120

/* Security checking functions.  */
/*
 * Copyright (c) 2017, 2023 Apple Inc. All rights reserved.
 *
 * @APPLE_LICENSE_HEADER_START@
 *
 * This file contains Original Code and/or Modifications of Original Code
 * as defined in and that are subject to the Apple Public Source License
 * Version 2.0 (the 'License'). You may not use this file except in
 * compliance with the License. Please obtain a copy of the License at
 * http://www.opensource.apple.com/apsl/ and read it before using this
 * file.
 *
 * The Original Code and all software distributed under the License are
 * distributed on an 'AS IS' basis, WITHOUT WARRANTY OF ANY KIND, EITHER
 * EXPRESS OR IMPLIED, AND APPLE HEREBY DISCLAIMS ALL SUCH WARRANTIES,
 * INCLUDING WITHOUT LIMITATION, ANY WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE, QUIET ENJOYMENT OR NON-INFRINGEMENT.
 * Please see the License for the specific language governing rights and
 * limitations under the License.
 *
 * @APPLE_LICENSE_HEADER_END@
 */

/*
 * Copyright (c) 2000-2018 Apple Inc. All rights reserved.
 *
 * @APPLE_OSREFERENCE_LICENSE_HEADER_START@
 *
 * This file contains Original Code and/or Modifications of Original Code
 * as defined in and that are subject to the Apple Public Source License
 * Version 2.0 (the 'License'). You may not use this file except in
 * compliance with the License. The rights granted to you under the License
 * may not be used to create, or enable the creation or redistribution of,
 * unlawful or unlicensed copies of an Apple operating system, or to
 * circumvent, violate, or enable the circumvention or violation of, any
 * terms of an Apple operating system software license agreement.
 *
 * Please obtain a copy of the License at
 * http://www.opensource.apple.com/apsl/ and read it before using this file.
 *
 * The Original Code and all software distributed under the License are
 * distributed on an 'AS IS' basis, WITHOUT WARRANTY OF ANY KIND, EITHER
 * EXPRESS OR IMPLIED, AND APPLE HEREBY DISCLAIMS ALL SUCH WARRANTIES,
 * INCLUDING WITHOUT LIMITATION, ANY WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE, QUIET ENJOYMENT OR NON-INFRINGEMENT.
 * Please see the License for the specific language governing rights and
 * limitations under the License.
 *
 * @APPLE_OSREFERENCE_LICENSE_HEADER_END@
 */
/* Copyright 1995 NeXT Computer, Inc. All rights reserved. */
/*
 * Copyright (c) 1991, 1993
 *	The Regents of the University of California.  All rights reserved.
 *
 * This code is derived from software contributed to Berkeley by
 * Berkeley Software Design, Inc.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 * 3. All advertising materials mentioning features or use of this software
 *    must display the following acknowledgement:
 *	This product includes software developed by the University of
 *	California, Berkeley and its contributors.
 * 4. Neither the name of the University nor the names of its contributors
 *    may be used to endorse or promote products derived from this software
 *    without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE REGENTS AND CONTRIBUTORS ``AS IS'' AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED.  IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
 * OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
 * LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
 * OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
 * SUCH DAMAGE.
 *
 *	@(#)cdefs.h	8.8 (Berkeley) 1/9/95
 */

/*
 * Copyright (c) 2007-2016 by Apple Inc.. All rights reserved.
 *
 * @APPLE_LICENSE_HEADER_START@
 *
 * This file contains Original Code and/or Modifications of Original Code
 * as defined in and that are subject to the Apple Public Source License
 * Version 2.0 (the 'License'). You may not use this file except in
 * compliance with the License. Please obtain a copy of the License at
 * http://www.opensource.apple.com/apsl/ and read it before using this
 * file.
 *
 * The Original Code and all software distributed under the License are
 * distributed on an 'AS IS' basis, WITHOUT WARRANTY OF ANY KIND, EITHER
 * EXPRESS OR IMPLIED, AND APPLE HEREBY DISCLAIMS ALL SUCH WARRANTIES,
 * INCLUDING WITHOUT LIMITATION, ANY WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE, QUIET ENJOYMENT OR NON-INFRINGEMENT.
 * Please see the License for the specific language governing rights and
 * limitations under the License.
 *
 * @APPLE_LICENSE_HEADER_END@
 */

// This is explicitly outside the header guard
/*
 * Copyright (c) 2007, 2008 Apple Inc. All rights reserved.
 *
 * @APPLE_LICENSE_HEADER_START@
 *
 * This file contains Original Code and/or Modifications of Original Code
 * as defined in and that are subject to the Apple Public Source License
 * Version 2.0 (the 'License'). You may not use this file except in
 * compliance with the License. Please obtain a copy of the License at
 * http://www.opensource.apple.com/apsl/ and read it before using this
 * file.
 *
 * The Original Code and all software distributed under the License are
 * distributed on an 'AS IS' basis, WITHOUT WARRANTY OF ANY KIND, EITHER
 * EXPRESS OR IMPLIED, AND APPLE HEREBY DISCLAIMS ALL SUCH WARRANTIES,
 * INCLUDING WITHOUT LIMITATION, ANY WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE, QUIET ENJOYMENT OR NON-INFRINGEMENT.
 * Please see the License for the specific language governing rights and
 * limitations under the License.
 *
 * @APPLE_LICENSE_HEADER_END@
 */

/* bcopy and bzero */

/* Removed in Issue 7 */

/* Security checking functions.  */
/*
 * Copyright (c) 2007,2017,2023 Apple Inc. All rights reserved.
 *
 * @APPLE_LICENSE_HEADER_START@
 *
 * This file contains Original Code and/or Modifications of Original Code
 * as defined in and that are subject to the Apple Public Source License
 * Version 2.0 (the 'License'). You may not use this file except in
 * compliance with the License. Please obtain a copy of the License at
 * http://www.opensource.apple.com/apsl/ and read it before using this
 * file.
 *
 * The Original Code and all software distributed under the License are
 * distributed on an 'AS IS' basis, WITHOUT WARRANTY OF ANY KIND, EITHER
 * EXPRESS OR IMPLIED, AND APPLE HEREBY DISCLAIMS ALL SUCH WARRANTIES,
 * INCLUDING WITHOUT LIMITATION, ANY WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE, QUIET ENJOYMENT OR NON-INFRINGEMENT.
 * Please see the License for the specific language governing rights and
 * limitations under the License.
 *
 * @APPLE_LICENSE_HEADER_END@
 */

/*
 * Copyright (c) 2000-2018 Apple Inc. All rights reserved.
 *
 * @APPLE_OSREFERENCE_LICENSE_HEADER_START@
 *
 * This file contains Original Code and/or Modifications of Original Code
 * as defined in and that are subject to the Apple Public Source License
 * Version 2.0 (the 'License'). You may not use this file except in
 * compliance with the License. The rights granted to you under the License
 * may not be used to create, or enable the creation or redistribution of,
 * unlawful or unlicensed copies of an Apple operating system, or to
 * circumvent, violate, or enable the circumvention or violation of, any
 * terms of an Apple operating system software license agreement.
 *
 * Please obtain a copy of the License at
 * http://www.opensource.apple.com/apsl/ and read it before using this file.
 *
 * The Original Code and all software distributed under the License are
 * distributed on an 'AS IS' basis, WITHOUT WARRANTY OF ANY KIND, EITHER
 * EXPRESS OR IMPLIED, AND APPLE HEREBY DISCLAIMS ALL SUCH WARRANTIES,
 * INCLUDING WITHOUT LIMITATION, ANY WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE, QUIET ENJOYMENT OR NON-INFRINGEMENT.
 * Please see the License for the specific language governing rights and
 * limitations under the License.
 *
 * @APPLE_OSREFERENCE_LICENSE_HEADER_END@
 */
/* Copyright 1995 NeXT Computer, Inc. All rights reserved. */
/*
 * Copyright (c) 1991, 1993
 *	The Regents of the University of California.  All rights reserved.
 *
 * This code is derived from software contributed to Berkeley by
 * Berkeley Software Design, Inc.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 * 3. All advertising materials mentioning features or use of this software
 *    must display the following acknowledgement:
 *	This product includes software developed by the University of
 *	California, Berkeley and its contributors.
 * 4. Neither the name of the University nor the names of its contributors
 *    may be used to endorse or promote products derived from this software
 *    without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE REGENTS AND CONTRIBUTORS ``AS IS'' AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED.  IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
 * OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
 * LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
 * OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
 * SUCH DAMAGE.
 *
 *	@(#)cdefs.h	8.8 (Berkeley) 1/9/95
 */

/*
 * Copyright (c) 2007-2016 by Apple Inc.. All rights reserved.
 *
 * @APPLE_LICENSE_HEADER_START@
 *
 * This file contains Original Code and/or Modifications of Original Code
 * as defined in and that are subject to the Apple Public Source License
 * Version 2.0 (the 'License'). You may not use this file except in
 * compliance with the License. Please obtain a copy of the License at
 * http://www.opensource.apple.com/apsl/ and read it before using this
 * file.
 *
 * The Original Code and all software distributed under the License are
 * distributed on an 'AS IS' basis, WITHOUT WARRANTY OF ANY KIND, EITHER
 * EXPRESS OR IMPLIED, AND APPLE HEREBY DISCLAIMS ALL SUCH WARRANTIES,
 * INCLUDING WITHOUT LIMITATION, ANY WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE, QUIET ENJOYMENT OR NON-INFRINGEMENT.
 * Please see the License for the specific language governing rights and
 * limitations under the License.
 *
 * @APPLE_LICENSE_HEADER_END@
 */

// This is explicitly outside the header guard
/*
 * Copyright (c) 2007, 2008 Apple Inc. All rights reserved.
 *
 * @APPLE_LICENSE_HEADER_START@
 *
 * This file contains Original Code and/or Modifications of Original Code
 * as defined in and that are subject to the Apple Public Source License
 * Version 2.0 (the 'License'). You may not use this file except in
 * compliance with the License. Please obtain a copy of the License at
 * http://www.opensource.apple.com/apsl/ and read it before using this
 * file.
 *
 * The Original Code and all software distributed under the License are
 * distributed on an 'AS IS' basis, WITHOUT WARRANTY OF ANY KIND, EITHER
 * EXPRESS OR IMPLIED, AND APPLE HEREBY DISCLAIMS ALL SUCH WARRANTIES,
 * INCLUDING WITHOUT LIMITATION, ANY WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE, QUIET ENJOYMENT OR NON-INFRINGEMENT.
 * Please see the License for the specific language governing rights and
 * limitations under the License.
 *
 * @APPLE_LICENSE_HEADER_END@
 */

/* __is_gcc(gcc_major, gcc_minor)
* Special values:
* 10.0 means "test should always fail when __has_builtin isn't supported"
  (because gcc got __has_builtin in version 10.0); this is used for builtins
  that gcc did not support yet at the time __has_builtin was introduced, so
  there is no point checking the compiler version.
* 0.0 means that we did not research when gcc started supporting this builtin,
  but it's believed to have been the case at least since gcc 4.0, which came
  out in 2005. (Hello from 2025! What year is it now? Can't believe we're still
  using C!)
*/

/* memccpy, memcpy, mempcpy, memmove, memset, strcpy, strlcpy, stpcpy,
   strncpy, stpncpy, strcat, strlcat, and strncat */

/* The use of .../__VA_ARGS__ is load-bearing. If the macros take fixed
 * arguments, they are unable to themselves accept macros that expand to
 * multiple arguments, like this:
 *  #define memcpy(a, b, c) ...
 *  #define FOO(data) get_bytes(data), get_length(data)
 *  memcpy(bar, FOO(d));
 * This will fail because the preprocessor only sees two arguments on the first
 * expansion of memcpy, when 3 are required.
 * This is also required to support syntaxes that embed commas. The preprocessor
 * recognizes parentheses for the isolation of arguments but not brackets. This
 * expands to 3 arguments:
 *  strcpy(destination, [NSString stringWithFormat:@"%i", 4].UTF8String);
 *         ^            ^                                 ^
 *         |destination |                                 |
 *                      |[NSString stringWithFormat:@"%i" |
 *							  |4].UTF8String
 * This expands to 4 arguments:
 *  memcpy(destination, (uint8_t[]) { 1, 2 }, 2);
 *         ^            ^                ^    ^
 * To work correctly under these hostile circumstances, chk_func macros
 * need to expand to a bare identifier (like #define memcpy_chk_func __memcpy)
 * or to a macro that also takes variadic arguments.
 */

func check_alloc(tls *libc.TLS, _ptr uintptr) (r uintptr) {
	if _ptr == libc.UintptrFromInt32(0) {
		libc.Xfprintf(tls, libc.X__stderrp, __ccgo_ts+165, 0)
		libc.Xexit(tls, int32(EXIT_FAILURE))
	}
	return _ptr
}

func opus_malloc(tls *libc.TLS, _size size_t) (r uintptr) {
	return check_alloc(tls, libc.Xmalloc(tls, _size))
}

func opus_realloc(tls *libc.TLS, _ptr uintptr, _size size_t) (r uintptr) {
	return check_alloc(tls, libc.Xrealloc(tls, _ptr, _size))
}

func read_pcm16(tls *libc.TLS, _samples uintptr, _fin uintptr, _nchannels int32) (r size_t) {
	bp := tls.Alloc(1024)
	defer tls.Free(1024)
	var ci, s int32
	var csamples, nread, nsamples, xi, v1 size_t
	var samples uintptr
	var _ /* buf at bp+0 */ [1024]uint8
	_, _, _, _, _, _, _, _ = ci, csamples, nread, nsamples, s, samples, xi, v1
	samples = libc.UintptrFromInt32(0)
	v1 = libc.Uint64FromInt32(0)
	csamples = v1
	nsamples = v1
	for {
		nread = libc.Xfread(tls, bp, libc.Uint64FromInt32(int32(2)*_nchannels), libc.Uint64FromInt32(int32(1024)/(int32(2)*_nchannels)), _fin)
		if nread <= uint64(0) {
			break
		}
		if nsamples+nread > csamples {
			for cond := true; cond; cond = nsamples+nread > csamples {
				csamples = csamples<<int32(1) | uint64(1)
			}
			samples = opus_realloc(tls, samples, libc.Uint64FromInt32(_nchannels)*csamples*uint64(4))
		}
		xi = uint64(0)
		for {
			if !(xi < nread) {
				break
			}
			ci = 0
			for {
				if !(ci < _nchannels) {
					break
				}
				s = libc.Int32FromUint8((*(*[1024]uint8)(unsafe.Pointer(bp)))[uint64(2)*(xi*libc.Uint64FromInt32(_nchannels)+libc.Uint64FromInt32(ci))+uint64(1)])<<int32(8) | libc.Int32FromUint8((*(*[1024]uint8)(unsafe.Pointer(bp)))[uint64(2)*(xi*libc.Uint64FromInt32(_nchannels)+libc.Uint64FromInt32(ci))])
				s = s&int32(0xFFFF) ^ int32(0x8000) - int32(0x8000)
				*(*float32)(unsafe.Pointer(samples + uintptr((nsamples+xi)*libc.Uint64FromInt32(_nchannels)+libc.Uint64FromInt32(ci))*4)) = float32(s)
				goto _4
			_4:
				;
				ci = ci + 1
			}
			goto _3
		_3:
			;
			xi = xi + 1
		}
		nsamples = nsamples + nread
		goto _2
	_2:
	}
	*(*uintptr)(unsafe.Pointer(_samples)) = opus_realloc(tls, samples, libc.Uint64FromInt32(_nchannels)*nsamples*uint64(4))
	return nsamples
}

func band_energy(tls *libc.TLS, _out uintptr, _ps uintptr, _bands uintptr, _nbands int32, _in uintptr, _nchannels int32, _nframes size_t, _window_sz int32, _step int32, _downsample int32) {
	bp := tls.Alloc(16)
	defer tls.Free(16)
	var bi, ci, ps_sz, ti, xj, xk, v8 int32
	var c, s, window, x uintptr
	var im, re, v11 float32
	var xi size_t
	var _ /* p at bp+0 */ [2]float32
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = bi, c, ci, im, ps_sz, re, s, ti, window, x, xi, xj, xk, v11, v8
	window = opus_malloc(tls, libc.Uint64FromInt32((int32(3)+_nchannels)*_window_sz)*uint64(4))
	c = window + uintptr(_window_sz)*4
	s = c + uintptr(_window_sz)*4
	x = s + uintptr(_window_sz)*4
	ps_sz = _window_sz / int32(2)
	xj = 0
	for {
		if !(xj < _window_sz) {
			break
		}
		*(*float32)(unsafe.Pointer(window + uintptr(xj)*4)) = libc.Float32FromFloat32(0.5) - float32(libc.Float32FromFloat32(0.5)*float32(libc.Xcos(tls, float64(float32(libc.Float32FromInt32(2)*libc.Float32FromFloat32(3.14159265))/float32(_window_sz-libc.Int32FromInt32(1))*float32(xj)))))
		goto _1
	_1:
		;
		xj = xj + 1
	}
	xj = 0
	for {
		if !(xj < _window_sz) {
			break
		}
		*(*float32)(unsafe.Pointer(c + uintptr(xj)*4)) = float32(libc.Xcos(tls, float64(float32(libc.Float32FromInt32(2)*libc.Float32FromFloat32(3.14159265))/float32(_window_sz)*float32(xj))))
		goto _2
	_2:
		;
		xj = xj + 1
	}
	xj = 0
	for {
		if !(xj < _window_sz) {
			break
		}
		*(*float32)(unsafe.Pointer(s + uintptr(xj)*4)) = float32(libc.Xsin(tls, float64(float32(libc.Float32FromInt32(2)*libc.Float32FromFloat32(3.14159265))/float32(_window_sz)*float32(xj))))
		goto _3
	_3:
		;
		xj = xj + 1
	}
	xi = uint64(0)
	for {
		if !(xi < _nframes) {
			break
		}
		ci = 0
		for {
			if !(ci < _nchannels) {
				break
			}
			xk = 0
			for {
				if !(xk < _window_sz) {
					break
				}
				*(*float32)(unsafe.Pointer(x + uintptr(ci*_window_sz+xk)*4)) = float32(*(*float32)(unsafe.Pointer(window + uintptr(xk)*4)) * *(*float32)(unsafe.Pointer(_in + uintptr((xi*libc.Uint64FromInt32(_step)+libc.Uint64FromInt32(xk))*libc.Uint64FromInt32(_nchannels)+libc.Uint64FromInt32(ci))*4)))
				goto _6
			_6:
				;
				xk = xk + 1
			}
			goto _5
		_5:
			;
			ci = ci + 1
		}
		v8 = libc.Int32FromInt32(0)
		xj = v8
		bi = v8
		for {
			if !(bi < _nbands) {
				break
			}
			*(*[2]float32)(unsafe.Pointer(bp)) = [2]float32{}
			for {
				if !(xj < *(*int32)(unsafe.Pointer(_bands + uintptr(bi+int32(1))*4))) {
					break
				}
				ci = 0
				for {
					if !(ci < _nchannels) {
						break
					}
					ti = 0
					v11 = libc.Float32FromInt32(0)
					im = v11
					re = v11
					xk = 0
					for {
						if !(xk < _window_sz) {
							break
						}
						re = re + float32(*(*float32)(unsafe.Pointer(c + uintptr(ti)*4))**(*float32)(unsafe.Pointer(x + uintptr(ci*_window_sz+xk)*4)))
						im = im - float32(*(*float32)(unsafe.Pointer(s + uintptr(ti)*4))**(*float32)(unsafe.Pointer(x + uintptr(ci*_window_sz+xk)*4)))
						ti = ti + xj
						if ti >= _window_sz {
							ti = ti - _window_sz
						}
						goto _12
					_12:
						;
						xk = xk + 1
					}
					re = re * float32(_downsample)
					im = im * float32(_downsample)
					*(*float32)(unsafe.Pointer(_ps + uintptr((xi*libc.Uint64FromInt32(ps_sz)+libc.Uint64FromInt32(xj))*libc.Uint64FromInt32(_nchannels)+libc.Uint64FromInt32(ci))*4)) = float32(re*re) + float32(im*im) + libc.Float32FromInt32(100000)
					*(*float32)(unsafe.Pointer(bp + uintptr(ci)*4)) += *(*float32)(unsafe.Pointer(_ps + uintptr((xi*libc.Uint64FromInt32(ps_sz)+libc.Uint64FromInt32(xj))*libc.Uint64FromInt32(_nchannels)+libc.Uint64FromInt32(ci))*4))
					goto _10
				_10:
					;
					ci = ci + 1
				}
				goto _9
			_9:
				;
				xj = xj + 1
			}
			if _out != 0 {
				*(*float32)(unsafe.Pointer(_out + uintptr((xi*libc.Uint64FromInt32(_nbands)+libc.Uint64FromInt32(bi))*libc.Uint64FromInt32(_nchannels))*4)) = (*(*[2]float32)(unsafe.Pointer(bp)))[0] / float32(*(*int32)(unsafe.Pointer(_bands + uintptr(bi+int32(1))*4))-*(*int32)(unsafe.Pointer(_bands + uintptr(bi)*4)))
				if _nchannels == int32(2) {
					*(*float32)(unsafe.Pointer(_out + uintptr((xi*libc.Uint64FromInt32(_nbands)+libc.Uint64FromInt32(bi))*libc.Uint64FromInt32(_nchannels)+uint64(1))*4)) = (*(*[2]float32)(unsafe.Pointer(bp)))[int32(1)] / float32(*(*int32)(unsafe.Pointer(_bands + uintptr(bi+int32(1))*4))-*(*int32)(unsafe.Pointer(_bands + uintptr(bi)*4)))
				}
			}
			goto _7
		_7:
			;
			bi = bi + 1
		}
		goto _4
	_4:
		;
		xi = xi + 1
	}
	libc.Xfree(tls, window)
}

// C documentation
//
//	/*Bands on which we compute the pseudo-NMR (Bark-derived
//	  CELT bands).*/
var BANDS = [22]int32{
	1:  int32(2),
	2:  int32(4),
	3:  int32(6),
	4:  int32(8),
	5:  int32(10),
	6:  int32(12),
	7:  int32(14),
	8:  int32(16),
	9:  int32(20),
	10: int32(24),
	11: int32(28),
	12: int32(32),
	13: int32(40),
	14: int32(48),
	15: int32(56),
	16: int32(68),
	17: int32(80),
	18: int32(96),
	19: int32(120),
	20: int32(156),
	21: int32(200),
}

func main1(tls *libc.TLS, _argc int32, _argv uintptr) (r1 int32) {
	bp := tls.Alloc(48)
	defer tls.Free(48)
	var Eb, Ef, err float64
	var Q, im, l, r, re, xtmp, xtmp2, ytmp, ytmp2 float32
	var X, Y, fin1, fin2, xb uintptr
	var bi, ci, downsample, max_compare, nchannels, xj, ybands, yfreqs, v6 int32
	var nframes, xi, xlength, ylength size_t
	var rate uint32
	var _ /* x at bp+0 */ uintptr
	var _ /* y at bp+8 */ uintptr
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = Eb, Ef, Q, X, Y, bi, ci, downsample, err, fin1, fin2, im, l, max_compare, nchannels, nframes, r, rate, re, xb, xi, xj, xlength, xtmp, xtmp2, ybands, yfreqs, ylength, ytmp, ytmp2, v6
	if _argc < int32(3) || _argc > int32(6) {
		libc.Xfprintf(tls, libc.X__stderrp, __ccgo_ts+181, libc.VaList(bp+24, *(*uintptr)(unsafe.Pointer(_argv))))
		return int32(EXIT_FAILURE)
	}
	nchannels = int32(1)
	if libc.Xstrcmp(tls, *(*uintptr)(unsafe.Pointer(_argv + 1*8)), __ccgo_ts+230) == 0 {
		nchannels = int32(2)
		_argv += 8
	}
	rate = uint32(48000)
	ybands = int32(NBANDS)
	yfreqs = int32(NFREQS)
	downsample = int32(1)
	if libc.Xstrcmp(tls, *(*uintptr)(unsafe.Pointer(_argv + 1*8)), __ccgo_ts+233) == 0 {
		rate = libc.Uint32FromInt32(libc.Xatoi(tls, *(*uintptr)(unsafe.Pointer(_argv + 2*8))))
		if rate != uint32(8000) && rate != uint32(12000) && rate != uint32(16000) && rate != uint32(24000) && rate != uint32(48000) {
			libc.Xfprintf(tls, libc.X__stderrp, __ccgo_ts+236, 0)
			return int32(EXIT_FAILURE)
		}
		downsample = libc.Int32FromUint32(uint32(48000) / rate)
		switch rate {
		case uint32(8000):
			ybands = int32(13)
		case uint32(12000):
			ybands = int32(15)
		case uint32(16000):
			ybands = int32(17)
		case uint32(24000):
			ybands = int32(19)
			break
		}
		yfreqs = int32(NFREQS) / downsample
		_argv = _argv + uintptr(2)*8
	}
	fin1 = libc.Xfopen(tls, *(*uintptr)(unsafe.Pointer(_argv + 1*8)), __ccgo_ts+295)
	if fin1 == libc.UintptrFromInt32(0) {
		libc.Xfprintf(tls, libc.X__stderrp, __ccgo_ts+298, libc.VaList(bp+24, *(*uintptr)(unsafe.Pointer(_argv + 1*8))))
		return int32(EXIT_FAILURE)
	}
	fin2 = libc.Xfopen(tls, *(*uintptr)(unsafe.Pointer(_argv + 2*8)), __ccgo_ts+295)
	if fin2 == libc.UintptrFromInt32(0) {
		libc.Xfprintf(tls, libc.X__stderrp, __ccgo_ts+298, libc.VaList(bp+24, *(*uintptr)(unsafe.Pointer(_argv + 2*8))))
		libc.Xfclose(tls, fin1)
		return int32(EXIT_FAILURE)
	}
	/*Read in the data and allocate scratch space.*/
	xlength = read_pcm16(tls, bp, fin1, int32(2))
	if nchannels == int32(1) {
		xi = uint64(0)
		for {
			if !(xi < xlength) {
				break
			}
			*(*float32)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(bp)) + uintptr(xi)*4)) = float32(float64(0.5) * float64(*(*float32)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(bp)) + uintptr(uint64(2)*xi)*4))+*(*float32)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(bp)) + uintptr(uint64(2)*xi+uint64(1))*4))))
			goto _1
		_1:
			;
			xi = xi + 1
		}
	}
	libc.Xfclose(tls, fin1)
	ylength = read_pcm16(tls, bp+8, fin2, nchannels)
	libc.Xfclose(tls, fin2)
	if xlength != ylength*libc.Uint64FromInt32(downsample) {
		libc.Xfprintf(tls, libc.X__stderrp, __ccgo_ts+319, libc.VaList(bp+24, xlength, ylength*libc.Uint64FromInt32(downsample)))
		return int32(EXIT_FAILURE)
	}
	if xlength < libc.Uint64FromInt32(libc.Int32FromInt32(TEST_WIN_SIZE)) {
		libc.Xfprintf(tls, libc.X__stderrp, __ccgo_ts+359, libc.VaList(bp+24, xlength, int32(TEST_WIN_SIZE)))
		return int32(EXIT_FAILURE)
	}
	nframes = (xlength - libc.Uint64FromInt32(libc.Int32FromInt32(TEST_WIN_SIZE)) + libc.Uint64FromInt32(libc.Int32FromInt32(TEST_WIN_STEP))) / libc.Uint64FromInt32(libc.Int32FromInt32(TEST_WIN_STEP))
	xb = opus_malloc(tls, nframes*libc.Uint64FromInt32(libc.Int32FromInt32(NBANDS))*libc.Uint64FromInt32(nchannels)*uint64(4))
	X = opus_malloc(tls, nframes*libc.Uint64FromInt32(libc.Int32FromInt32(NFREQS))*libc.Uint64FromInt32(nchannels)*uint64(4))
	Y = opus_malloc(tls, nframes*libc.Uint64FromInt32(yfreqs)*libc.Uint64FromInt32(nchannels)*uint64(4))
	/*Compute the per-band spectral energy of the original signal
	  and the error.*/
	band_energy(tls, xb, X, uintptr(unsafe.Pointer(&BANDS)), int32(NBANDS), *(*uintptr)(unsafe.Pointer(bp)), nchannels, nframes, int32(TEST_WIN_SIZE), int32(TEST_WIN_STEP), int32(1))
	libc.Xfree(tls, *(*uintptr)(unsafe.Pointer(bp)))
	band_energy(tls, libc.UintptrFromInt32(0), Y, uintptr(unsafe.Pointer(&BANDS)), ybands, *(*uintptr)(unsafe.Pointer(bp + 8)), nchannels, nframes, int32(TEST_WIN_SIZE)/downsample, int32(TEST_WIN_STEP)/downsample, downsample)
	libc.Xfree(tls, *(*uintptr)(unsafe.Pointer(bp + 8)))
	xi = uint64(0)
	for {
		if !(xi < nframes) {
			break
		}
		/*Frequency masking (low to high): 10 dB/Bark slope.*/
		bi = int32(1)
		for {
			if !(bi < int32(NBANDS)) {
				break
			}
			ci = 0
			for {
				if !(ci < nchannels) {
					break
				}
				*(*float32)(unsafe.Pointer(xb + uintptr((xi*libc.Uint64FromInt32(libc.Int32FromInt32(NBANDS))+libc.Uint64FromInt32(bi))*libc.Uint64FromInt32(nchannels)+libc.Uint64FromInt32(ci))*4)) += float32(libc.Float32FromFloat32(0.1) * *(*float32)(unsafe.Pointer(xb + uintptr((xi*libc.Uint64FromInt32(libc.Int32FromInt32(NBANDS))+libc.Uint64FromInt32(bi)-uint64(1))*libc.Uint64FromInt32(nchannels)+libc.Uint64FromInt32(ci))*4)))
				goto _4
			_4:
				;
				ci = ci + 1
			}
			goto _3
		_3:
			;
			bi = bi + 1
		}
		/*Frequency masking (high to low): 15 dB/Bark slope.*/
		bi = libc.Int32FromInt32(NBANDS) - libc.Int32FromInt32(1)
		for {
			v6 = bi
			bi = bi - 1
			if !(v6 > 0) {
				break
			}
			ci = 0
			for {
				if !(ci < nchannels) {
					break
				}
				*(*float32)(unsafe.Pointer(xb + uintptr((xi*libc.Uint64FromInt32(libc.Int32FromInt32(NBANDS))+libc.Uint64FromInt32(bi))*libc.Uint64FromInt32(nchannels)+libc.Uint64FromInt32(ci))*4)) += float32(libc.Float32FromFloat32(0.03) * *(*float32)(unsafe.Pointer(xb + uintptr((xi*libc.Uint64FromInt32(libc.Int32FromInt32(NBANDS))+libc.Uint64FromInt32(bi)+uint64(1))*libc.Uint64FromInt32(nchannels)+libc.Uint64FromInt32(ci))*4)))
				goto _7
			_7:
				;
				ci = ci + 1
			}
			goto _5
		_5:
		}
		if xi > uint64(0) {
			/*Temporal masking: -3 dB/2.5ms slope.*/
			bi = 0
			for {
				if !(bi < int32(NBANDS)) {
					break
				}
				ci = 0
				for {
					if !(ci < nchannels) {
						break
					}
					*(*float32)(unsafe.Pointer(xb + uintptr((xi*libc.Uint64FromInt32(libc.Int32FromInt32(NBANDS))+libc.Uint64FromInt32(bi))*libc.Uint64FromInt32(nchannels)+libc.Uint64FromInt32(ci))*4)) += float32(libc.Float32FromFloat32(0.5) * *(*float32)(unsafe.Pointer(xb + uintptr(((xi-uint64(1))*libc.Uint64FromInt32(libc.Int32FromInt32(NBANDS))+libc.Uint64FromInt32(bi))*libc.Uint64FromInt32(nchannels)+libc.Uint64FromInt32(ci))*4)))
					goto _9
				_9:
					;
					ci = ci + 1
				}
				goto _8
			_8:
				;
				bi = bi + 1
			}
		}
		/* Allowing some cross-talk */
		if nchannels == int32(2) {
			bi = 0
			for {
				if !(bi < int32(NBANDS)) {
					break
				}
				l = *(*float32)(unsafe.Pointer(xb + uintptr((xi*libc.Uint64FromInt32(libc.Int32FromInt32(NBANDS))+libc.Uint64FromInt32(bi))*libc.Uint64FromInt32(nchannels)+uint64(0))*4))
				r = *(*float32)(unsafe.Pointer(xb + uintptr((xi*libc.Uint64FromInt32(libc.Int32FromInt32(NBANDS))+libc.Uint64FromInt32(bi))*libc.Uint64FromInt32(nchannels)+uint64(1))*4))
				*(*float32)(unsafe.Pointer(xb + uintptr((xi*libc.Uint64FromInt32(libc.Int32FromInt32(NBANDS))+libc.Uint64FromInt32(bi))*libc.Uint64FromInt32(nchannels)+uint64(0))*4)) += float32(libc.Float32FromFloat32(0.01) * r)
				*(*float32)(unsafe.Pointer(xb + uintptr((xi*libc.Uint64FromInt32(libc.Int32FromInt32(NBANDS))+libc.Uint64FromInt32(bi))*libc.Uint64FromInt32(nchannels)+uint64(1))*4)) += float32(libc.Float32FromFloat32(0.01) * l)
				goto _10
			_10:
				;
				bi = bi + 1
			}
		}
		/* Apply masking */
		bi = 0
		for {
			if !(bi < ybands) {
				break
			}
			xj = BANDS[bi]
			for {
				if !(xj < BANDS[bi+int32(1)]) {
					break
				}
				ci = 0
				for {
					if !(ci < nchannels) {
						break
					}
					*(*float32)(unsafe.Pointer(X + uintptr((xi*libc.Uint64FromInt32(libc.Int32FromInt32(NFREQS))+libc.Uint64FromInt32(xj))*libc.Uint64FromInt32(nchannels)+libc.Uint64FromInt32(ci))*4)) += float32(libc.Float32FromFloat32(0.1) * *(*float32)(unsafe.Pointer(xb + uintptr((xi*libc.Uint64FromInt32(libc.Int32FromInt32(NBANDS))+libc.Uint64FromInt32(bi))*libc.Uint64FromInt32(nchannels)+libc.Uint64FromInt32(ci))*4)))
					*(*float32)(unsafe.Pointer(Y + uintptr((xi*libc.Uint64FromInt32(yfreqs)+libc.Uint64FromInt32(xj))*libc.Uint64FromInt32(nchannels)+libc.Uint64FromInt32(ci))*4)) += float32(libc.Float32FromFloat32(0.1) * *(*float32)(unsafe.Pointer(xb + uintptr((xi*libc.Uint64FromInt32(libc.Int32FromInt32(NBANDS))+libc.Uint64FromInt32(bi))*libc.Uint64FromInt32(nchannels)+libc.Uint64FromInt32(ci))*4)))
					goto _13
				_13:
					;
					ci = ci + 1
				}
				goto _12
			_12:
				;
				xj = xj + 1
			}
			goto _11
		_11:
			;
			bi = bi + 1
		}
		goto _2
	_2:
		;
		xi = xi + 1
	}
	/* Average of consecutive frames to make comparison slightly less sensitive */
	bi = 0
	for {
		if !(bi < ybands) {
			break
		}
		xj = BANDS[bi]
		for {
			if !(xj < BANDS[bi+int32(1)]) {
				break
			}
			ci = 0
			for {
				if !(ci < nchannels) {
					break
				}
				xtmp = *(*float32)(unsafe.Pointer(X + uintptr(xj*nchannels+ci)*4))
				ytmp = *(*float32)(unsafe.Pointer(Y + uintptr(xj*nchannels+ci)*4))
				xi = uint64(1)
				for {
					if !(xi < nframes) {
						break
					}
					xtmp2 = *(*float32)(unsafe.Pointer(X + uintptr((xi*libc.Uint64FromInt32(libc.Int32FromInt32(NFREQS))+libc.Uint64FromInt32(xj))*libc.Uint64FromInt32(nchannels)+libc.Uint64FromInt32(ci))*4))
					ytmp2 = *(*float32)(unsafe.Pointer(Y + uintptr((xi*libc.Uint64FromInt32(yfreqs)+libc.Uint64FromInt32(xj))*libc.Uint64FromInt32(nchannels)+libc.Uint64FromInt32(ci))*4))
					*(*float32)(unsafe.Pointer(X + uintptr((xi*libc.Uint64FromInt32(libc.Int32FromInt32(NFREQS))+libc.Uint64FromInt32(xj))*libc.Uint64FromInt32(nchannels)+libc.Uint64FromInt32(ci))*4)) += xtmp
					*(*float32)(unsafe.Pointer(Y + uintptr((xi*libc.Uint64FromInt32(yfreqs)+libc.Uint64FromInt32(xj))*libc.Uint64FromInt32(nchannels)+libc.Uint64FromInt32(ci))*4)) += ytmp
					xtmp = xtmp2
					ytmp = ytmp2
					goto _17
				_17:
					;
					xi = xi + 1
				}
				goto _16
			_16:
				;
				ci = ci + 1
			}
			goto _15
		_15:
			;
			xj = xj + 1
		}
		goto _14
	_14:
		;
		bi = bi + 1
	}
	/*If working at a lower sampling rate, don't take into account the last
	   300 Hz to allow for different transition bands.
	  For 12 kHz, we don't skip anything, because the last band already skips
	   400 Hz.*/
	if rate == uint32(48000) {
		max_compare = BANDS[int32(NBANDS)]
	} else {
		if rate == uint32(12000) {
			max_compare = BANDS[ybands]
		} else {
			max_compare = BANDS[ybands] - int32(3)
		}
	}
	err = libc.Float64FromInt32(0)
	xi = uint64(0)
	for {
		if !(xi < nframes) {
			break
		}
		Ef = libc.Float64FromInt32(0)
		bi = 0
		for {
			if !(bi < ybands) {
				break
			}
			Eb = libc.Float64FromInt32(0)
			xj = BANDS[bi]
			for {
				if !(xj < BANDS[bi+int32(1)] && xj < max_compare) {
					break
				}
				ci = 0
				for {
					if !(ci < nchannels) {
						break
					}
					re = *(*float32)(unsafe.Pointer(Y + uintptr((xi*libc.Uint64FromInt32(yfreqs)+libc.Uint64FromInt32(xj))*libc.Uint64FromInt32(nchannels)+libc.Uint64FromInt32(ci))*4)) / *(*float32)(unsafe.Pointer(X + uintptr((xi*libc.Uint64FromInt32(libc.Int32FromInt32(NFREQS))+libc.Uint64FromInt32(xj))*libc.Uint64FromInt32(nchannels)+libc.Uint64FromInt32(ci))*4))
					im = float32(float64(re) - libc.Xlog(tls, float64(re)) - libc.Float64FromInt32(1))
					/*Make comparison less sensitive around the SILK/CELT cross-over to
					  allow for mode freedom in the filters.*/
					if xj >= int32(79) && xj <= int32(81) {
						im = im * libc.Float32FromFloat32(0.1)
					}
					if xj == int32(80) {
						im = im * libc.Float32FromFloat32(0.1)
					}
					Eb = Eb + float64(im)
					goto _21
				_21:
					;
					ci = ci + 1
				}
				goto _20
			_20:
				;
				xj = xj + 1
			}
			Eb = Eb / float64((BANDS[bi+int32(1)]-BANDS[bi])*nchannels)
			Ef = Ef + float64(Eb*Eb)
			goto _19
		_19:
			;
			bi = bi + 1
		}
		/*Using a fixed normalization value means we're willing to accept slightly
		  lower quality for lower sampling rates.*/
		Ef = Ef / float64(libc.Int32FromInt32(NBANDS))
		Ef = Ef * Ef
		err = err + float64(Ef*Ef)
		goto _18
	_18:
		;
		xi = xi + 1
	}
	err = libc.Xpow(tls, err/float64(nframes), libc.Float64FromFloat64(1)/libc.Float64FromInt32(16))
	Q = float32(libc.Float64FromInt32(100) * (libc.Float64FromInt32(1) - float64(float64(0.5)*libc.Xlog(tls, libc.Float64FromInt32(1)+err))/libc.Xlog(tls, float64(1.13))))
	if Q < libc.Float32FromInt32(0) {
		libc.Xfprintf(tls, libc.X__stderrp, __ccgo_ts+395, 0)
		libc.Xfprintf(tls, libc.X__stderrp, __ccgo_ts+414, libc.VaList(bp+24, err))
		return int32(EXIT_FAILURE)
	} else {
		libc.Xfprintf(tls, libc.X__stderrp, __ccgo_ts+445, 0)
		libc.Xfprintf(tls, libc.X__stderrp, __ccgo_ts+465, libc.VaList(bp+24, float64(Q), err))
		return EXIT_SUCCESS
	}
	return r1
}

func main() {
	libc.Start(main1)
}

var __ccgo_ts = (*reflect.StringHeader)(unsafe.Pointer(&__ccgo_ts1)).Data

var __ccgo_ts1 = "success\x00invalid argument\x00buffer too small\x00internal error\x00corrupted stream\x00request not implemented\x00invalid state\x00memory allocation failed\x00unknown error\x00libopus 1.0.0\x00Out of memory.\n\x00Usage: %s [-s] [-r rate2] <file1.sw> <file2.sw>\n\x00-s\x00-r\x00Sampling rate must be 8000, 12000, 16000, 24000, or 48000\n\x00rb\x00Error opening '%s'.\n\x00Sample counts do not match (%lu!=%lu).\n\x00Insufficient sample data (%lu<%i).\n\x00Test vector FAILS\n\x00Internal weighted error is %f\n\x00Test vector PASSES\n\x00Opus quality metric: %.1f %% (internal weighted error is %f)\n\x00"
