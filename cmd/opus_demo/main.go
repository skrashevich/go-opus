// Code generated for darwin/arm64 by ccgo from opus_demo.c; the codec itself is
// the github.com/skrashevich/go-opus package. DO NOT EDIT.

package main

import (
	"os"
	"reflect"
	"runtime/pprof"
	"unsafe"

	"modernc.org/libc"

	opus "github.com/skrashevich/go-opus"
)

const EXIT_FAILURE = 1

const EXIT_SUCCESS = 0

const OPUS_APPLICATION_AUDIO = 2049

const OPUS_APPLICATION_RESTRICTED_LOWDELAY = 2051

const OPUS_APPLICATION_VOIP = 2048

const OPUS_BANDWIDTH_FULLBAND = 1105

const OPUS_BANDWIDTH_MEDIUMBAND = 1102

const OPUS_BANDWIDTH_NARROWBAND = 1101

const OPUS_BANDWIDTH_SUPERWIDEBAND = 1104

const OPUS_BANDWIDTH_WIDEBAND = 1103

const OPUS_GET_FINAL_RANGE_REQUEST = 4031

const OPUS_GET_LOOKAHEAD_REQUEST = 4027

const OPUS_OK = 0

const OPUS_SET_BANDWIDTH_REQUEST = 4008

const OPUS_SET_BITRATE_REQUEST = 4002

const OPUS_SET_COMPLEXITY_REQUEST = 4010

const OPUS_SET_DTX_REQUEST = 4016

const OPUS_SET_FORCE_CHANNELS_REQUEST = 4022

const OPUS_SET_INBAND_FEC_REQUEST = 4012

const OPUS_SET_PACKET_LOSS_PERC_REQUEST = 4014

const OPUS_SET_VBR_CONSTRAINT_REQUEST = 4020

const OPUS_SET_VBR_REQUEST = 4006

const SEEK_END = 2

const SEEK_SET = 0

type __predefined_ptrdiff_t = int64

type opus_int32 = int32

type opus_uint32 = uint32

const MODE_CELT_ONLY = 1002

const MODE_SILK_ONLY = 1000

const OPUS_SET_FORCE_MODE_REQUEST = 11002

const MAX_PACKET = 1500

func print_usage(tls *libc.TLS, argv uintptr) {
	bp := tls.Alloc(16)
	defer tls.Free(16)
	libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+165, libc.VaList(bp+8, *(*uintptr)(unsafe.Pointer(argv))))
	libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+279, libc.VaList(bp+8, *(*uintptr)(unsafe.Pointer(argv))))
	libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+359, 0)
	libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+401, 0)
	libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+411, 0)
	libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+481, 0)
	libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+559, 0)
	libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+634, 0)
	libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+718, 0)
	libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+820, 0)
	libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+885, 0)
	libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+954, 0)
	libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+1031, 0)
	libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+1078, 0)
	libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+1145, 0)
	libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+1185, 0)
}

/*
 * Copyright (c) 2000, 2007, 2010, 2023 Apple Inc. All rights reserved.
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
/*-
 * Copyright (c) 1990, 1993
 *	The Regents of the University of California.  All rights reserved.
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
 *	@(#)strings.h	8.1 (Berkeley) 6/2/93
 */

/*
 * Copyright (c) 2000, 2007, 2010, 2023 Apple Inc. All rights reserved.
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
/*-
 * Copyright (c) 1990, 1993
 *	The Regents of the University of California.  All rights reserved.
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
 *	@(#)strings.h	8.1 (Berkeley) 6/2/93
 */

/*
 * Copyright (c) 2000, 2007, 2010, 2023 Apple Inc. All rights reserved.
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
/*-
 * Copyright (c) 1990, 1993
 *	The Regents of the University of California.  All rights reserved.
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
 *	@(#)string.h	8.1 (Berkeley) 6/2/93
 */

/*
 * Copyright (c) 2023 Apple Inc. All rights reserved.
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
/*-
 * Copyright (c) 1990, 1993
 *	The Regents of the University of California.  All rights reserved.
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
 *	@(#)string.h	8.1 (Berkeley) 6/2/93
 */

func int_to_char(tls *libc.TLS, i opus_uint32, ch uintptr) {
	*(*uint8)(unsafe.Pointer(ch)) = uint8(i >> int32(24))
	*(*uint8)(unsafe.Pointer(ch + 1)) = uint8(i >> libc.Int32FromInt32(16) & uint32(0xFF))
	*(*uint8)(unsafe.Pointer(ch + 2)) = uint8(i >> libc.Int32FromInt32(8) & uint32(0xFF))
	*(*uint8)(unsafe.Pointer(ch + 3)) = uint8(i & uint32(0xFF))
}

func char_to_int(tls *libc.TLS, ch uintptr) (r opus_uint32) {
	return uint32(*(*uint8)(unsafe.Pointer(ch)))<<libc.Int32FromInt32(24) | uint32(*(*uint8)(unsafe.Pointer(ch + 1)))<<libc.Int32FromInt32(16) | uint32(*(*uint8)(unsafe.Pointer(ch + 2)))<<libc.Int32FromInt32(8) | uint32(*(*uint8)(unsafe.Pointer(ch + 3)))
}

func check_decoder_option(tls *libc.TLS, encode_only int32, opt uintptr) {
	bp := tls.Alloc(16)
	defer tls.Free(16)
	if encode_only != 0 {
		libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+1262, libc.VaList(bp+8, opt))
		libc.Xexit(tls, int32(EXIT_FAILURE))
	}
}

func check_encoder_option(tls *libc.TLS, decode_only int32, opt uintptr) {
	bp := tls.Alloc(16)
	defer tls.Free(16)
	if decode_only != 0 {
		libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+1294, libc.VaList(bp+8, opt))
		libc.Xexit(tls, int32(EXIT_FAILURE))
	}
}

func main1(tls *libc.TLS, argc int32, argv uintptr) (r int32) {
	bp := tls.Alloc(80)
	defer tls.Free(80)
	var application, args, bandwidth, channels, complexity, curr_mode, curr_mode_count, curr_read, cvbr, decode_only, delayed_celt, encode_only, forcechannels, frame_size, i, i1, k, lost, lost_prev, max_frame_size, max_payload_bytes, mode_switch_time, nb_modes_in_list, newsize, output_samples, packet_loss_perc, random_fec, random_framesize, size, stop, sweep_bps, sweep_max, sweep_min, toggle, use_dtx, use_inbandfec, use_vbr int32
	var bandwidth_string, dec, enc, fbytes, fin, fout, in, inFile, mode_list, out, outFile, v3 uintptr
	var bitrate_bps, count, count_act, s, sampling_rate opus_int32
	var bits, bits2, bits_act, bits_max, nrg, v6 float64
	var data [2]uintptr
	var len1 [2]int32
	var s1 int16
	var _ /* ch at bp+20 */ [4]uint8
	var _ /* dec_final_range at bp+16 */ opus_uint32
	var _ /* enc_final_range at bp+8 */ [2]opus_uint32
	var _ /* err at bp+0 */ int32
	var _ /* int_field at bp+24 */ [4]uint8
	var _ /* skip at bp+4 */ int32
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = application, args, bandwidth, bandwidth_string, bitrate_bps, bits, bits2, bits_act, bits_max, channels, complexity, count, count_act, curr_mode, curr_mode_count, curr_read, cvbr, data, dec, decode_only, delayed_celt, enc, encode_only, fbytes, fin, forcechannels, fout, frame_size, i, i1, in, inFile, k, len1, lost, lost_prev, max_frame_size, max_payload_bytes, mode_list, mode_switch_time, nb_modes_in_list, newsize, nrg, out, outFile, output_samples, packet_loss_perc, random_fec, random_framesize, s, s1, sampling_rate, size, stop, sweep_bps, sweep_max, sweep_min, toggle, use_dtx, use_inbandfec, use_vbr, v3, v6
	enc = libc.UintptrFromInt32(0)
	dec = libc.UintptrFromInt32(0)
	bitrate_bps = 0
	cvbr = 0
	count = 0
	count_act = 0
	*(*int32)(unsafe.Pointer(bp + 4)) = 0
	stop = 0
	application = int32(OPUS_APPLICATION_AUDIO)
	bits = float64(0)
	bits_max = float64(0)
	bits_act = float64(0)
	bits2 = float64(0)
	bandwidth = -int32(1)
	lost = 0
	lost_prev = int32(1)
	toggle = 0
	encode_only = 0
	decode_only = 0
	max_frame_size = libc.Int32FromInt32(960) * libc.Int32FromInt32(6)
	curr_read = 0
	sweep_bps = 0
	random_framesize = 0
	newsize = 0
	delayed_celt = 0
	sweep_max = 0
	sweep_min = 0
	random_fec = 0
	mode_list = libc.UintptrFromInt32(0)
	nb_modes_in_list = 0
	curr_mode = 0
	curr_mode_count = 0
	mode_switch_time = int32(48000)
	if argc < int32(5) {
		print_usage(tls, argv)
		return int32(EXIT_FAILURE)
	}
	libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+1326, libc.VaList(bp+40, opus.GetVersionString(tls)))
	args = int32(1)
	if libc.Xstrcmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1330) == 0 {
		encode_only = int32(1)
		args = args + 1
	} else {
		if libc.Xstrcmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1333) == 0 {
			decode_only = int32(1)
			args = args + 1
		}
	}
	if !(decode_only != 0) && argc < int32(7) {
		print_usage(tls, argv)
		return int32(EXIT_FAILURE)
	}
	if !(decode_only != 0) {
		if libc.Xstrcmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1336) == 0 {
			application = int32(OPUS_APPLICATION_VOIP)
		} else {
			if libc.Xstrcmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1341) == 0 {
				application = int32(OPUS_APPLICATION_RESTRICTED_LOWDELAY)
			} else {
				if libc.Xstrcmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1361) != 0 {
					libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+1367, libc.VaList(bp+40, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize))))
					print_usage(tls, argv)
					return int32(EXIT_FAILURE)
				}
			}
		}
		args = args + 1
	}
	sampling_rate = int32(libc.Xatol(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize))))
	args = args + 1
	channels = libc.Xatoi(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)))
	args = args + 1
	if !(decode_only != 0) {
		bitrate_bps = int32(libc.Xatol(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize))))
		args = args + 1
	}
	if sampling_rate != int32(8000) && sampling_rate != int32(12000) && sampling_rate != int32(16000) && sampling_rate != int32(24000) && sampling_rate != int32(48000) {
		libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+1392, 0)
		return int32(EXIT_FAILURE)
	}
	frame_size = sampling_rate / int32(50)
	/* defaults: */
	use_vbr = int32(1)
	bandwidth = -int32(1000)
	max_payload_bytes = int32(MAX_PACKET)
	complexity = int32(10)
	use_inbandfec = 0
	forcechannels = -int32(1000)
	use_dtx = 0
	packet_loss_perc = 0
	max_frame_size = libc.Int32FromInt32(960) * libc.Int32FromInt32(6)
	curr_read = 0
	for args < argc-int32(2) {
		/* process command line options */
		if libc.Xstrcasecmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1459) == 0 {
			check_encoder_option(tls, decode_only, __ccgo_ts+1459)
			use_vbr = 0
			args = args + 1
		} else {
			if libc.Xstrcasecmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1464) == 0 {
				check_encoder_option(tls, decode_only, __ccgo_ts+1464)
				if libc.Xstrcmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args+int32(1))*ptrSize)), __ccgo_ts+1475) == 0 {
					bandwidth = int32(OPUS_BANDWIDTH_NARROWBAND)
				} else {
					if libc.Xstrcmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args+int32(1))*ptrSize)), __ccgo_ts+1478) == 0 {
						bandwidth = int32(OPUS_BANDWIDTH_MEDIUMBAND)
					} else {
						if libc.Xstrcmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args+int32(1))*ptrSize)), __ccgo_ts+1481) == 0 {
							bandwidth = int32(OPUS_BANDWIDTH_WIDEBAND)
						} else {
							if libc.Xstrcmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args+int32(1))*ptrSize)), __ccgo_ts+1484) == 0 {
								bandwidth = int32(OPUS_BANDWIDTH_SUPERWIDEBAND)
							} else {
								if libc.Xstrcmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args+int32(1))*ptrSize)), __ccgo_ts+1488) == 0 {
									bandwidth = int32(OPUS_BANDWIDTH_FULLBAND)
								} else {
									libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+1491, libc.VaList(bp+40, *(*uintptr)(unsafe.Pointer(argv + uintptr(args+int32(1))*ptrSize))))
									return int32(EXIT_FAILURE)
								}
							}
						}
					}
				}
				args = args + int32(2)
			} else {
				if libc.Xstrcasecmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1549) == 0 {
					check_encoder_option(tls, decode_only, __ccgo_ts+1549)
					if libc.Xstrcmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args+int32(1))*ptrSize)), __ccgo_ts+1560) == 0 {
						frame_size = sampling_rate / int32(400)
					} else {
						if libc.Xstrcmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args+int32(1))*ptrSize)), __ccgo_ts+1564) == 0 {
							frame_size = sampling_rate / int32(200)
						} else {
							if libc.Xstrcmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args+int32(1))*ptrSize)), __ccgo_ts+1566) == 0 {
								frame_size = sampling_rate / int32(100)
							} else {
								if libc.Xstrcmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args+int32(1))*ptrSize)), __ccgo_ts+1569) == 0 {
									frame_size = sampling_rate / int32(50)
								} else {
									if libc.Xstrcmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args+int32(1))*ptrSize)), __ccgo_ts+1572) == 0 {
										frame_size = sampling_rate / int32(25)
									} else {
										if libc.Xstrcmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args+int32(1))*ptrSize)), __ccgo_ts+1575) == 0 {
											frame_size = int32(3) * sampling_rate / int32(50)
										} else {
											libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+1578, libc.VaList(bp+40, *(*uintptr)(unsafe.Pointer(argv + uintptr(args+int32(1))*ptrSize))))
											return int32(EXIT_FAILURE)
										}
									}
								}
							}
						}
					}
					args = args + int32(2)
				} else {
					if libc.Xstrcasecmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1648) == 0 {
						check_encoder_option(tls, decode_only, __ccgo_ts+1648)
						max_payload_bytes = libc.Xatoi(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args+int32(1))*ptrSize)))
						args = args + int32(2)
					} else {
						if libc.Xstrcasecmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1661) == 0 {
							check_encoder_option(tls, decode_only, __ccgo_ts+1661)
							complexity = libc.Xatoi(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args+int32(1))*ptrSize)))
							args = args + int32(2)
						} else {
							if libc.Xstrcasecmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1673) == 0 {
								use_inbandfec = int32(1)
								args = args + 1
							} else {
								if libc.Xstrcasecmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1684) == 0 {
									check_encoder_option(tls, decode_only, __ccgo_ts+1684)
									forcechannels = int32(1)
									args = args + 1
								} else {
									if libc.Xstrcasecmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1695) == 0 {
										check_encoder_option(tls, decode_only, __ccgo_ts+1695)
										cvbr = int32(1)
										args = args + 1
									} else {
										if libc.Xstrcasecmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1701) == 0 {
											check_encoder_option(tls, decode_only, __ccgo_ts+1701)
											use_dtx = int32(1)
											args = args + 1
										} else {
											if libc.Xstrcasecmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1706) == 0 {
												check_decoder_option(tls, encode_only, __ccgo_ts+1706)
												packet_loss_perc = libc.Xatoi(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args+int32(1))*ptrSize)))
												args = args + int32(2)
											} else {
												if libc.Xstrcasecmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1712) == 0 {
													check_encoder_option(tls, decode_only, __ccgo_ts+1712)
													sweep_bps = libc.Xatoi(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args+int32(1))*ptrSize)))
													args = args + int32(2)
												} else {
													if libc.Xstrcasecmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1719) == 0 {
														check_encoder_option(tls, decode_only, __ccgo_ts+1719)
														random_framesize = int32(1)
														args = args + 1
													} else {
														if libc.Xstrcasecmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1737) == 0 {
															check_encoder_option(tls, decode_only, __ccgo_ts+1737)
															sweep_max = libc.Xatoi(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args+int32(1))*ptrSize)))
															args = args + int32(2)
														} else {
															if libc.Xstrcasecmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1748) == 0 {
																check_encoder_option(tls, decode_only, __ccgo_ts+1748)
																random_fec = int32(1)
																args = args + 1
															} else {
																if libc.Xstrcasecmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1760) == 0 {
																	check_encoder_option(tls, decode_only, __ccgo_ts+1760)
																	mode_list = uintptr(unsafe.Pointer(&silk8_test))
																	nb_modes_in_list = int32(8)
																	args = args + 1
																} else {
																	if libc.Xstrcasecmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1773) == 0 {
																		check_encoder_option(tls, decode_only, __ccgo_ts+1773)
																		mode_list = uintptr(unsafe.Pointer(&silk12_test))
																		nb_modes_in_list = int32(8)
																		args = args + 1
																	} else {
																		if libc.Xstrcasecmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1787) == 0 {
																			check_encoder_option(tls, decode_only, __ccgo_ts+1787)
																			mode_list = uintptr(unsafe.Pointer(&silk16_test))
																			nb_modes_in_list = int32(8)
																			args = args + 1
																		} else {
																			if libc.Xstrcasecmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1801) == 0 {
																				check_encoder_option(tls, decode_only, __ccgo_ts+1801)
																				mode_list = uintptr(unsafe.Pointer(&hybrid24_test))
																				nb_modes_in_list = int32(4)
																				args = args + 1
																			} else {
																				if libc.Xstrcasecmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1817) == 0 {
																					check_encoder_option(tls, decode_only, __ccgo_ts+1817)
																					mode_list = uintptr(unsafe.Pointer(&hybrid48_test))
																					nb_modes_in_list = int32(4)
																					args = args + 1
																				} else {
																					if libc.Xstrcasecmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1833) == 0 {
																						check_encoder_option(tls, decode_only, __ccgo_ts+1833)
																						mode_list = uintptr(unsafe.Pointer(&celt_test))
																						nb_modes_in_list = int32(32)
																						args = args + 1
																					} else {
																						if libc.Xstrcasecmp(tls, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize)), __ccgo_ts+1844) == 0 {
																							check_encoder_option(tls, decode_only, __ccgo_ts+1844)
																							mode_list = uintptr(unsafe.Pointer(&celt_hq_test))
																							nb_modes_in_list = int32(4)
																							args = args + 1
																						} else {
																							libc.Xprintf(tls, __ccgo_ts+1858, libc.VaList(bp+40, *(*uintptr)(unsafe.Pointer(argv + uintptr(args)*ptrSize))))
																							print_usage(tls, argv)
																							return int32(EXIT_FAILURE)
																						}
																					}
																				}
																			}
																		}
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	if sweep_max != 0 {
		sweep_min = bitrate_bps
	}
	if max_payload_bytes < 0 || max_payload_bytes > int32(MAX_PACKET) {
		libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+1892, libc.VaList(bp+40, int32(MAX_PACKET)))
		return int32(EXIT_FAILURE)
	}
	inFile = *(*uintptr)(unsafe.Pointer(argv + uintptr(argc-int32(2))*ptrSize))
	fin = libc.Xfopen(tls, inFile, __ccgo_ts+1936)
	if !(fin != 0) {
		libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+1939, libc.VaList(bp+40, *(*uintptr)(unsafe.Pointer(argv + uintptr(argc-int32(2))*ptrSize))))
		return int32(EXIT_FAILURE)
	}
	if mode_list != 0 {
		libc.Xfseek(tls, fin, 0, int32(SEEK_END))
		size = int32(libc.Xftell(tls, fin))
		libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+1969, libc.VaList(bp+40, size))
		libc.Xfseek(tls, fin, 0, SEEK_SET)
		mode_switch_time = libc.Int32FromUint64(libc.Uint64FromInt32(size) / uint64(2) / libc.Uint64FromInt32(channels) / libc.Uint64FromInt32(nb_modes_in_list))
		libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+1992, libc.VaList(bp+40, mode_switch_time))
	}
	outFile = *(*uintptr)(unsafe.Pointer(argv + uintptr(argc-int32(1))*ptrSize))
	fout = libc.Xfopen(tls, outFile, __ccgo_ts+2025)
	if !(fout != 0) {
		libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+2029, libc.VaList(bp+40, *(*uintptr)(unsafe.Pointer(argv + uintptr(argc-int32(1))*ptrSize))))
		libc.Xfclose(tls, fin)
		return int32(EXIT_FAILURE)
	}
	if !(decode_only != 0) {
		enc = opus.EncoderCreate(tls, sampling_rate, channels, application, bp)
		if *(*int32)(unsafe.Pointer(bp)) != OPUS_OK {
			libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+2060, libc.VaList(bp+40, opus.Strerror(tls, *(*int32)(unsafe.Pointer(bp)))))
			libc.Xfclose(tls, fin)
			libc.Xfclose(tls, fout)
			return int32(EXIT_FAILURE)
		}
		_ = bitrate_bps == libc.Int32FromInt32(0)
		opus.EncoderCtl(tls, enc, int32(OPUS_SET_BITRATE_REQUEST), libc.VaList(bp+40, bitrate_bps))
		_ = bandwidth == libc.Int32FromInt32(0)
		opus.EncoderCtl(tls, enc, int32(OPUS_SET_BANDWIDTH_REQUEST), libc.VaList(bp+40, bandwidth))
		_ = use_vbr == libc.Int32FromInt32(0)
		opus.EncoderCtl(tls, enc, int32(OPUS_SET_VBR_REQUEST), libc.VaList(bp+40, use_vbr))
		_ = cvbr == libc.Int32FromInt32(0)
		opus.EncoderCtl(tls, enc, int32(OPUS_SET_VBR_CONSTRAINT_REQUEST), libc.VaList(bp+40, cvbr))
		_ = complexity == libc.Int32FromInt32(0)
		opus.EncoderCtl(tls, enc, int32(OPUS_SET_COMPLEXITY_REQUEST), libc.VaList(bp+40, complexity))
		_ = use_inbandfec == libc.Int32FromInt32(0)
		opus.EncoderCtl(tls, enc, int32(OPUS_SET_INBAND_FEC_REQUEST), libc.VaList(bp+40, use_inbandfec))
		_ = forcechannels == libc.Int32FromInt32(0)
		opus.EncoderCtl(tls, enc, int32(OPUS_SET_FORCE_CHANNELS_REQUEST), libc.VaList(bp+40, forcechannels))
		_ = use_dtx == libc.Int32FromInt32(0)
		opus.EncoderCtl(tls, enc, int32(OPUS_SET_DTX_REQUEST), libc.VaList(bp+40, use_dtx))
		_ = packet_loss_perc == libc.Int32FromInt32(0)
		opus.EncoderCtl(tls, enc, int32(OPUS_SET_PACKET_LOSS_PERC_REQUEST), libc.VaList(bp+40, packet_loss_perc))
		opus.EncoderCtl(tls, enc, int32(OPUS_GET_LOOKAHEAD_REQUEST), libc.VaList(bp+40, bp+4+uintptr((__predefined_ptrdiff_t(bp+4)-int64(bp+4))/4)*4))
	}
	if !(encode_only != 0) {
		dec = opus.DecoderCreate(tls, sampling_rate, channels, bp)
		if *(*int32)(unsafe.Pointer(bp)) != OPUS_OK {
			libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+2087, libc.VaList(bp+40, opus.Strerror(tls, *(*int32)(unsafe.Pointer(bp)))))
			libc.Xfclose(tls, fin)
			libc.Xfclose(tls, fout)
			return int32(EXIT_FAILURE)
		}
	}
	switch bandwidth {
	case int32(OPUS_BANDWIDTH_NARROWBAND):
		bandwidth_string = __ccgo_ts + 2114
	case int32(OPUS_BANDWIDTH_MEDIUMBAND):
		bandwidth_string = __ccgo_ts + 2125
	case int32(OPUS_BANDWIDTH_WIDEBAND):
		bandwidth_string = __ccgo_ts + 2136
	case int32(OPUS_BANDWIDTH_SUPERWIDEBAND):
		bandwidth_string = __ccgo_ts + 2145
	case int32(OPUS_BANDWIDTH_FULLBAND):
		bandwidth_string = __ccgo_ts + 2159
	case -int32(1000):
		bandwidth_string = __ccgo_ts + 2168
	default:
		bandwidth_string = __ccgo_ts + 2173
		break
	}
	if decode_only != 0 {
		libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+2181, libc.VaList(bp+40, int64(sampling_rate), channels))
	} else {
		libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+2224, libc.VaList(bp+40, int64(sampling_rate), float64(float64(bitrate_bps)*float64(0.001)), bandwidth_string, frame_size))
	}
	in = xmalloc(tls, libc.Uint64FromInt32(max_frame_size*channels)*uint64(2))
	out = xmalloc(tls, libc.Uint64FromInt32(max_frame_size*channels)*uint64(2))
	fbytes = xmalloc(tls, libc.Uint64FromInt32(max_frame_size*channels)*uint64(2))
	data[0] = xcalloc(tls, libc.Uint64FromInt32(max_payload_bytes), uint64(1))
	if use_inbandfec != 0 {
		data[int32(1)] = xcalloc(tls, libc.Uint64FromInt32(max_payload_bytes), uint64(1))
	}
	for !(stop != 0) {
		if delayed_celt != 0 {
			frame_size = newsize
			delayed_celt = 0
		} else {
			if random_framesize != 0 && libc.Xrand(tls)%int32(20) == 0 {
				newsize = libc.Xrand(tls) % int32(6)
				switch newsize {
				case 0:
					newsize = sampling_rate / int32(400)
				case int32(1):
					newsize = sampling_rate / int32(200)
				case int32(2):
					newsize = sampling_rate / int32(100)
				case int32(3):
					newsize = sampling_rate / int32(50)
				case int32(4):
					newsize = sampling_rate / int32(25)
				case int32(5):
					newsize = int32(3) * sampling_rate / int32(50)
					break
				}
				for newsize < sampling_rate/int32(25) && float64(bitrate_bps)-libc.Xfabs(tls, float64(sweep_bps)) <= float64(libc.Int32FromInt32(3)*libc.Int32FromInt32(12)*sampling_rate/newsize) {
					newsize = newsize * int32(2)
				}
				if newsize < sampling_rate/int32(100) && frame_size >= sampling_rate/int32(100) {
					_ = libc.Int32FromInt32(MODE_CELT_ONLY) == libc.Int32FromInt32(0)
					opus.EncoderCtl(tls, enc, int32(OPUS_SET_FORCE_MODE_REQUEST), libc.VaList(bp+40, libc.Int32FromInt32(MODE_CELT_ONLY)))
					delayed_celt = int32(1)
				} else {
					frame_size = newsize
				}
			}
		}
		if random_fec != 0 && libc.Xrand(tls)%int32(30) == 0 {
			_ = libc.BoolInt32(libc.Xrand(tls)%libc.Int32FromInt32(4) == libc.Int32FromInt32(0)) == libc.Int32FromInt32(0)
			opus.EncoderCtl(tls, enc, int32(OPUS_SET_INBAND_FEC_REQUEST), libc.VaList(bp+40, libc.BoolInt32(libc.Xrand(tls)%libc.Int32FromInt32(4) == libc.Int32FromInt32(0))))
		}
		if decode_only != 0 {
			*(*int32)(unsafe.Pointer(bp)) = libc.Int32FromUint64(xfread(tls, bp+20, uint64(1), uint64(4), fin))
			if libc.Xfeof(tls, fin) != 0 {
				break
			}
			len1[toggle] = libc.Int32FromUint32(char_to_int(tls, bp+20))
			if len1[toggle] > max_payload_bytes || len1[toggle] < 0 {
				libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+2294, libc.VaList(bp+40, len1[toggle]))
				break
			}
			*(*int32)(unsafe.Pointer(bp)) = libc.Int32FromUint64(xfread(tls, bp+20, uint64(1), uint64(4), fin))
			(*(*[2]opus_uint32)(unsafe.Pointer(bp + 8)))[toggle] = char_to_int(tls, bp+20)
			*(*int32)(unsafe.Pointer(bp)) = libc.Int32FromUint64(xfread(tls, data[toggle], uint64(1), libc.Uint64FromInt32(len1[toggle]), fin))
			if *(*int32)(unsafe.Pointer(bp)) < len1[toggle] {
				libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+2322, libc.VaList(bp+40, len1[toggle], *(*int32)(unsafe.Pointer(bp))))
				break
			}
		} else {
			if mode_list != libc.UintptrFromInt32(0) {
				_ = *(*int32)(unsafe.Pointer(mode_list + uintptr(curr_mode)*16 + 1*4)) == libc.Int32FromInt32(0)
				opus.EncoderCtl(tls, enc, int32(OPUS_SET_BANDWIDTH_REQUEST), libc.VaList(bp+40, *(*int32)(unsafe.Pointer(mode_list + uintptr(curr_mode)*16 + 1*4))))
				_ = *(*int32)(unsafe.Pointer(mode_list + uintptr(curr_mode)*16)) == libc.Int32FromInt32(0)
				opus.EncoderCtl(tls, enc, int32(OPUS_SET_FORCE_MODE_REQUEST), libc.VaList(bp+40, *(*int32)(unsafe.Pointer(mode_list + uintptr(curr_mode)*16))))
				_ = *(*int32)(unsafe.Pointer(mode_list + uintptr(curr_mode)*16 + 3*4)) == libc.Int32FromInt32(0)
				opus.EncoderCtl(tls, enc, int32(OPUS_SET_FORCE_CHANNELS_REQUEST), libc.VaList(bp+40, *(*int32)(unsafe.Pointer(mode_list + uintptr(curr_mode)*16 + 3*4))))
				frame_size = *(*int32)(unsafe.Pointer(mode_list + uintptr(curr_mode)*16 + 2*4))
			}
			*(*int32)(unsafe.Pointer(bp)) = libc.Int32FromUint64(xfread(tls, fbytes, uint64(2)*libc.Uint64FromInt32(channels), libc.Uint64FromInt32(frame_size), fin))
			curr_read = *(*int32)(unsafe.Pointer(bp))
			i = 0
			for {
				if !(i < curr_read*channels) {
					break
				}
				s = libc.Int32FromUint8(*(*uint8)(unsafe.Pointer(fbytes + uintptr(int32(2)*i+int32(1)))))<<int32(8) | libc.Int32FromUint8(*(*uint8)(unsafe.Pointer(fbytes + uintptr(int32(2)*i))))
				s = s&int32(0xFFFF) ^ int32(0x8000) - int32(0x8000)
				*(*int16)(unsafe.Pointer(in + uintptr(i)*2)) = int16(s)
				goto _1
			_1:
				;
				i = i + 1
			}
			if curr_read < frame_size {
				i = curr_read * channels
				for {
					if !(i < frame_size*channels) {
						break
					}
					*(*int16)(unsafe.Pointer(in + uintptr(i)*2)) = 0
					goto _2
				_2:
					;
					i = i + 1
				}
				stop = int32(1)
			}
			len1[toggle] = opus.Encode(tls, enc, in, frame_size, data[toggle], max_payload_bytes)
			if sweep_bps != 0 {
				bitrate_bps = bitrate_bps + sweep_bps
				if sweep_max != 0 {
					if bitrate_bps > sweep_max {
						sweep_bps = -sweep_bps
					} else {
						if bitrate_bps < sweep_min {
							sweep_bps = -sweep_bps
						}
					}
				}
				/* safety */
				if bitrate_bps < int32(1000) {
					bitrate_bps = int32(1000)
				}
				_ = bitrate_bps == libc.Int32FromInt32(0)
				opus.EncoderCtl(tls, enc, int32(OPUS_SET_BITRATE_REQUEST), libc.VaList(bp+40, bitrate_bps))
			}
			opus.EncoderCtl(tls, enc, int32(OPUS_GET_FINAL_RANGE_REQUEST), libc.VaList(bp+40, bp+8+uintptr(toggle)*4+uintptr((__predefined_ptrdiff_t(bp+8+uintptr(toggle)*4)-int64(bp+8+uintptr(toggle)*4))/4)*4))
			if len1[toggle] < 0 {
				libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+2367, libc.VaList(bp+40, len1[toggle]))
				libc.Xfclose(tls, fin)
				libc.Xfclose(tls, fout)
				return int32(EXIT_FAILURE)
			}
			curr_mode_count = curr_mode_count + frame_size
			if curr_mode_count > mode_switch_time && curr_mode < nb_modes_in_list-int32(1) {
				curr_mode = curr_mode + 1
				curr_mode_count = 0
			}
		}
		if encode_only != 0 {
			int_to_char(tls, libc.Uint32FromInt32(len1[toggle]), bp+24)
			if xfwrite(tls, bp+24, uint64(1), uint64(4), fout) != uint64(4) {
				libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+2394, 0)
				return int32(EXIT_FAILURE)
			}
			int_to_char(tls, (*(*[2]opus_uint32)(unsafe.Pointer(bp + 8)))[toggle], bp+24)
			if xfwrite(tls, bp+24, uint64(1), uint64(4), fout) != uint64(4) {
				libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+2394, 0)
				return int32(EXIT_FAILURE)
			}
			if xfwrite(tls, data[toggle], uint64(1), libc.Uint64FromInt32(len1[toggle]), fout) != uint64(libc.Uint32FromInt32(len1[toggle])) {
				libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+2394, 0)
				return int32(EXIT_FAILURE)
			}
		} else {
			lost = libc.BoolInt32(len1[toggle] == 0 || packet_loss_perc > 0 && libc.Xrand(tls)%int32(100) < packet_loss_perc)
			if count >= use_inbandfec {
				/* delay by one packet when using in-band FEC */
				if use_inbandfec != 0 {
					if lost_prev != 0 {
						/* attempt to decode with in-band FEC from next packet */
						if lost != 0 {
							v3 = libc.UintptrFromInt32(0)
						} else {
							v3 = data[toggle]
						}
						output_samples = opus.Decode(tls, dec, v3, len1[toggle], out, max_frame_size, int32(1))
					} else {
						/* regular decode */
						output_samples = opus.Decode(tls, dec, data[int32(1)-toggle], len1[int32(1)-toggle], out, max_frame_size, 0)
					}
				} else {
					if lost != 0 {
						v3 = libc.UintptrFromInt32(0)
					} else {
						v3 = data[toggle]
					}
					output_samples = opus.Decode(tls, dec, v3, len1[toggle], out, max_frame_size, 0)
				}
				if output_samples > 0 {
					if output_samples > *(*int32)(unsafe.Pointer(bp + 4)) {
						i1 = 0
						for {
							if !(i1 < (output_samples-*(*int32)(unsafe.Pointer(bp + 4)))*channels) {
								break
							}
							s1 = *(*int16)(unsafe.Pointer(out + uintptr(i1+*(*int32)(unsafe.Pointer(bp + 4))*channels)*2))
							*(*uint8)(unsafe.Pointer(fbytes + uintptr(int32(2)*i1))) = libc.Uint8FromInt32(int32(s1) & int32(0xFF))
							*(*uint8)(unsafe.Pointer(fbytes + uintptr(int32(2)*i1+int32(1)))) = libc.Uint8FromInt32(int32(s1) >> int32(8) & int32(0xFF))
							goto _5
						_5:
							;
							i1 = i1 + 1
						}
						if xfwrite(tls, fbytes, uint64(2)*libc.Uint64FromInt32(channels), libc.Uint64FromInt32(output_samples-*(*int32)(unsafe.Pointer(bp + 4))), fout) != uint64(libc.Uint32FromInt32(output_samples-*(*int32)(unsafe.Pointer(bp + 4)))) {
							libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+2394, 0)
							return int32(EXIT_FAILURE)
						}
					}
					if output_samples < *(*int32)(unsafe.Pointer(bp + 4)) {
						*(*int32)(unsafe.Pointer(bp + 4)) = *(*int32)(unsafe.Pointer(bp + 4)) - output_samples
					} else {
						*(*int32)(unsafe.Pointer(bp + 4)) = 0
					}
				} else {
					libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+2410, libc.VaList(bp+40, opus.Strerror(tls, output_samples)))
				}
			}
		}
		if !(encode_only != 0) {
			opus.DecoderCtl(tls, dec, int32(OPUS_GET_FINAL_RANGE_REQUEST), libc.VaList(bp+40, bp+16+uintptr((__predefined_ptrdiff_t(bp+16)-int64(bp+16))/4)*4))
		}
		/* compare final range encoder rng values of encoder and decoder */
		if (*(*[2]opus_uint32)(unsafe.Pointer(bp + 8)))[toggle^use_inbandfec] != uint32(0) && !(encode_only != 0) && !(lost != 0) && !(lost_prev != 0) && *(*opus_uint32)(unsafe.Pointer(bp + 16)) != (*(*[2]opus_uint32)(unsafe.Pointer(bp + 8)))[toggle^use_inbandfec] {
			libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+2436, libc.VaList(bp+40, int64(count), uint64((*(*[2]opus_uint32)(unsafe.Pointer(bp + 8)))[toggle^use_inbandfec]), uint64(*(*opus_uint32)(unsafe.Pointer(bp + 16)))))
			libc.Xfclose(tls, fin)
			libc.Xfclose(tls, fout)
			return int32(EXIT_FAILURE)
		}
		lost_prev = lost
		/* count bits */
		bits = bits + float64(len1[toggle]*int32(8))
		if float64(len1[toggle]*int32(8)) > bits_max {
			v6 = float64(len1[toggle] * int32(8))
		} else {
			v6 = bits_max
		}
		bits_max = v6
		if count >= use_inbandfec {
			nrg = float64(0)
			if !(decode_only != 0) {
				k = 0
				for {
					if !(k < frame_size*channels) {
						break
					}
					nrg = nrg + float64(float64(*(*int16)(unsafe.Pointer(in + uintptr(k)*2)))*float64(*(*int16)(unsafe.Pointer(in + uintptr(k)*2))))
					goto _7
				_7:
					;
					k = k + 1
				}
			}
			if nrg/float64(frame_size*channels) > float64(100000) {
				bits_act = bits_act + float64(len1[toggle]*int32(8))
				count_act = count_act + 1
			}
			/* Variance */
			bits2 = bits2 + float64(len1[toggle]*len1[toggle]*int32(64))
		}
		count = count + 1
		toggle = (toggle + use_inbandfec) & int32(1)
	}
	libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+2530, libc.VaList(bp+40, float64(float64(float64(0.001)*bits)*float64(sampling_rate))/float64(float64(frame_size)*float64(count))))
	libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+2571, libc.VaList(bp+40, float64(float64(float64(0.001)*bits_max)*float64(sampling_rate))/float64(frame_size)))
	if !(decode_only != 0) {
		libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+2613, libc.VaList(bp+40, float64(float64(float64(0.001)*bits_act)*float64(sampling_rate))/float64(float64(frame_size)*float64(count_act))))
	}
	libc.Xfprintf(tls, libc.Xstderr, __ccgo_ts+2654, libc.VaList(bp+40, float64(float64(float64(0.001)*libc.Xsqrt(tls, bits2/float64(count)-float64(bits*bits)/float64(float64(count)*float64(count))))*float64(sampling_rate))/float64(frame_size)))
	/* Close any files to which intermediate results were stored */
	opus.EncoderDestroy(tls, enc)
	opus.DecoderDestroy(tls, dec)
	libc.Xfree(tls, data[0])
	if use_inbandfec != 0 {
		libc.Xfree(tls, data[int32(1)])
	}
	libc.Xfclose(tls, fin)
	libc.Xfclose(tls, fout)
	libc.Xfree(tls, in)
	libc.Xfree(tls, out)
	libc.Xfree(tls, fbytes)
	return EXIT_SUCCESS
}

func main() {
	if p := os.Getenv("CPUPROFILE"); p != "" {
		f, _ := os.Create(p)
		pprof.StartCPUProfile(f)
		// Register atexit to flush profile since libc.Start calls os.Exit
		libc.AtExit(func() {
			pprof.StopCPUProfile()
			f.Close()
		})
	}
	libc.Start(main1)
}

var celt_hq_test = [4][4]int32{
	0: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_FULLBAND),
		2: int32(960),
		3: int32(2),
	},
	1: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_FULLBAND),
		2: int32(480),
		3: int32(2),
	},
	2: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_FULLBAND),
		2: int32(240),
		3: int32(2),
	},
	3: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_FULLBAND),
		2: int32(120),
		3: int32(2),
	},
}

var celt_test = [32][4]int32{
	0: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_FULLBAND),
		2: int32(960),
		3: int32(1),
	},
	1: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_SUPERWIDEBAND),
		2: int32(960),
		3: int32(1),
	},
	2: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_WIDEBAND),
		2: int32(960),
		3: int32(1),
	},
	3: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_NARROWBAND),
		2: int32(960),
		3: int32(1),
	},
	4: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_FULLBAND),
		2: int32(480),
		3: int32(1),
	},
	5: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_SUPERWIDEBAND),
		2: int32(480),
		3: int32(1),
	},
	6: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_WIDEBAND),
		2: int32(480),
		3: int32(1),
	},
	7: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_NARROWBAND),
		2: int32(480),
		3: int32(1),
	},
	8: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_FULLBAND),
		2: int32(240),
		3: int32(1),
	},
	9: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_SUPERWIDEBAND),
		2: int32(240),
		3: int32(1),
	},
	10: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_WIDEBAND),
		2: int32(240),
		3: int32(1),
	},
	11: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_NARROWBAND),
		2: int32(240),
		3: int32(1),
	},
	12: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_FULLBAND),
		2: int32(120),
		3: int32(1),
	},
	13: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_SUPERWIDEBAND),
		2: int32(120),
		3: int32(1),
	},
	14: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_WIDEBAND),
		2: int32(120),
		3: int32(1),
	},
	15: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_NARROWBAND),
		2: int32(120),
		3: int32(1),
	},
	16: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_FULLBAND),
		2: int32(960),
		3: int32(2),
	},
	17: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_SUPERWIDEBAND),
		2: int32(960),
		3: int32(2),
	},
	18: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_WIDEBAND),
		2: int32(960),
		3: int32(2),
	},
	19: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_NARROWBAND),
		2: int32(960),
		3: int32(2),
	},
	20: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_FULLBAND),
		2: int32(480),
		3: int32(2),
	},
	21: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_SUPERWIDEBAND),
		2: int32(480),
		3: int32(2),
	},
	22: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_WIDEBAND),
		2: int32(480),
		3: int32(2),
	},
	23: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_NARROWBAND),
		2: int32(480),
		3: int32(2),
	},
	24: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_FULLBAND),
		2: int32(240),
		3: int32(2),
	},
	25: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_SUPERWIDEBAND),
		2: int32(240),
		3: int32(2),
	},
	26: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_WIDEBAND),
		2: int32(240),
		3: int32(2),
	},
	27: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_NARROWBAND),
		2: int32(240),
		3: int32(2),
	},
	28: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_FULLBAND),
		2: int32(120),
		3: int32(2),
	},
	29: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_SUPERWIDEBAND),
		2: int32(120),
		3: int32(2),
	},
	30: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_WIDEBAND),
		2: int32(120),
		3: int32(2),
	},
	31: {
		0: int32(MODE_CELT_ONLY),
		1: int32(OPUS_BANDWIDTH_NARROWBAND),
		2: int32(120),
		3: int32(2),
	},
}

var hybrid24_test = [4][4]int32{
	0: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_SUPERWIDEBAND),
		2: int32(960),
		3: int32(1),
	},
	1: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_SUPERWIDEBAND),
		2: int32(480),
		3: int32(1),
	},
	2: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_SUPERWIDEBAND),
		2: int32(960),
		3: int32(2),
	},
	3: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_SUPERWIDEBAND),
		2: int32(480),
		3: int32(2),
	},
}

var hybrid48_test = [4][4]int32{
	0: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_FULLBAND),
		2: int32(960),
		3: int32(1),
	},
	1: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_FULLBAND),
		2: int32(480),
		3: int32(1),
	},
	2: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_FULLBAND),
		2: int32(960),
		3: int32(2),
	},
	3: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_FULLBAND),
		2: int32(480),
		3: int32(2),
	},
}

var silk12_test = [8][4]int32{
	0: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_MEDIUMBAND),
		2: libc.Int32FromInt32(960) * libc.Int32FromInt32(3),
		3: int32(1),
	},
	1: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_MEDIUMBAND),
		2: libc.Int32FromInt32(960) * libc.Int32FromInt32(2),
		3: int32(1),
	},
	2: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_MEDIUMBAND),
		2: int32(960),
		3: int32(1),
	},
	3: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_MEDIUMBAND),
		2: int32(480),
		3: int32(1),
	},
	4: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_MEDIUMBAND),
		2: libc.Int32FromInt32(960) * libc.Int32FromInt32(3),
		3: int32(2),
	},
	5: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_MEDIUMBAND),
		2: libc.Int32FromInt32(960) * libc.Int32FromInt32(2),
		3: int32(2),
	},
	6: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_MEDIUMBAND),
		2: int32(960),
		3: int32(2),
	},
	7: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_MEDIUMBAND),
		2: int32(480),
		3: int32(2),
	},
}

var silk16_test = [8][4]int32{
	0: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_WIDEBAND),
		2: libc.Int32FromInt32(960) * libc.Int32FromInt32(3),
		3: int32(1),
	},
	1: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_WIDEBAND),
		2: libc.Int32FromInt32(960) * libc.Int32FromInt32(2),
		3: int32(1),
	},
	2: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_WIDEBAND),
		2: int32(960),
		3: int32(1),
	},
	3: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_WIDEBAND),
		2: int32(480),
		3: int32(1),
	},
	4: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_WIDEBAND),
		2: libc.Int32FromInt32(960) * libc.Int32FromInt32(3),
		3: int32(2),
	},
	5: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_WIDEBAND),
		2: libc.Int32FromInt32(960) * libc.Int32FromInt32(2),
		3: int32(2),
	},
	6: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_WIDEBAND),
		2: int32(960),
		3: int32(2),
	},
	7: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_WIDEBAND),
		2: int32(480),
		3: int32(2),
	},
}

var silk8_test = [8][4]int32{
	0: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_NARROWBAND),
		2: libc.Int32FromInt32(960) * libc.Int32FromInt32(3),
		3: int32(1),
	},
	1: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_NARROWBAND),
		2: libc.Int32FromInt32(960) * libc.Int32FromInt32(2),
		3: int32(1),
	},
	2: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_NARROWBAND),
		2: int32(960),
		3: int32(1),
	},
	3: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_NARROWBAND),
		2: int32(480),
		3: int32(1),
	},
	4: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_NARROWBAND),
		2: libc.Int32FromInt32(960) * libc.Int32FromInt32(3),
		3: int32(2),
	},
	5: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_NARROWBAND),
		2: libc.Int32FromInt32(960) * libc.Int32FromInt32(2),
		3: int32(2),
	},
	6: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_NARROWBAND),
		2: int32(960),
		3: int32(2),
	},
	7: {
		0: int32(MODE_SILK_ONLY),
		1: int32(OPUS_BANDWIDTH_NARROWBAND),
		2: int32(480),
		3: int32(2),
	},
}

var __ccgo_ts = (*reflect.StringHeader)(unsafe.Pointer(&__ccgo_ts1)).Data

var __ccgo_ts1 = "success\x00invalid argument\x00buffer too small\x00internal error\x00corrupted stream\x00request not implemented\x00invalid state\x00memory allocation failed\x00unknown error\x00libopus 1.0.0\x00Usage: %s [-e] <application> <sampling rate (Hz)> <channels (1/2)> <bits per second>  [options] <input> <output>\n\x00       %s -d <sampling rate (Hz)> <channels (1/2)> [options] <input> <output>\n\n\x00mode: voip | audio | restricted-lowdelay\n\x00options:\n\x00-e                   : only runs the encoder (output the bit-stream)\n\x00-d                   : only runs the decoder (reads the bit-stream as input)\n\x00-cbr                 : enable constant bitrate; default: variable bitrate\n\x00-cvbr                : enable constrained variable bitrate; default: unconstrained\n\x00-bandwidth <NB|MB|WB|SWB|FB> : audio bandwidth (from narrowband to fullband); default: sampling rate\n\x00-framesize <2.5|5|10|20|40|60> : frame size in ms; default: 20 \n\x00-max_payload <bytes> : maximum payload size in bytes, default: 1024\n\x00-complexity <comp>   : complexity, 0 (lowest) ... 10 (highest); default: 10\n\x00-inbandfec           : enable SILK inband FEC\n\x00-forcemono           : force mono encoding, even for stereo input\n\x00-dtx                 : enable SILK DTX\n\x00-loss <perc>         : simulate packet loss, in percent (0-100); default: 0\n\x00option %s is only for decoding\n\x00option %s is only for encoding\n\x00%s\n\x00-e\x00-d\x00voip\x00restricted-lowdelay\x00audio\x00unknown application: %s\n\x00Supported sampling rates are 8000, 12000, 16000, 24000 and 48000.\n\x00-cbr\x00-bandwidth\x00NB\x00MB\x00WB\x00SWB\x00FB\x00Unknown bandwidth %s. Supported are NB, MB, WB, SWB, FB.\n\x00-framesize\x002.5\x005\x0010\x0020\x0040\x0060\x00Unsupported frame size: %s ms. Supported are 2.5, 5, 10, 20, 40, 60.\n\x00-max_payload\x00-complexity\x00-inbandfec\x00-forcemono\x00-cvbr\x00-dtx\x00-loss\x00-sweep\x00-random_framesize\x00-sweep_max\x00-random_fec\x00-silk8k_test\x00-silk12k_test\x00-silk16k_test\x00-hybrid24k_test\x00-hybrid48k_test\x00-celt_test\x00-celt_hq_test\x00Error: unrecognized setting: %s\n\n\x00max_payload_bytes must be between 0 and %d\n\x00rb\x00Could not open input file %s\n\x00File size is %d bytes\n\x00Switching mode every %d samples\n\x00wb+\x00Could not open output file %s\n\x00Cannot create encoder: %s\n\x00Cannot create decoder: %s\n\x00narrowband\x00mediumband\x00wideband\x00superwideband\x00fullband\x00auto\x00unknown\x00Decoding with %ld Hz output (%d channels)\n\x00Encoding %ld Hz input at %.3f kb/s in %s mode with %d-sample frames.\n\x00Invalid payload length: %d\n\x00Ran out of input, expecting %d bytes got %d\n\x00opus_encode() returned %d\n\x00Error writing.\n\x00error decoding frame: %s\n\x00Error: Range coder state mismatch between encoder and decoder in frame %ld: 0x%8lx vs 0x%8lx\n\x00average bitrate:             %7.3f kb/s\n\x00maximum bitrate:             %7.3f bkp/s\n\x00active bitrate:              %7.3f kb/s\n\x00bitrate standard deviation:  %7.3f kb/s\n\x00"
