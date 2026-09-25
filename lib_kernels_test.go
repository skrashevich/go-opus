package opus

import (
	"math"
	"math/rand/v2"
	"testing"
	"unsafe"

	"modernc.org/libc"
)

// The reference functions below are verbatim scalar copies of the loops that
// the optimized kernels replaced. Results must match bit for bit, including
// whether multiply-adds are fused.

func randFloats(r *rand.Rand, n int, scale float32) []float32 {
	v := make([]float32, n)
	for i := range v {
		v[i] = (r.Float32()*2 - 1) * scale
	}
	return v
}

// ptr passes a buffer to the codec. Storing it in heapSink forces the backing
// array onto the heap, where a raw uintptr stays valid across stack growth.
func ptr(v []float32) uintptr {
	heapSink = v
	return uintptr(unsafe.Pointer(&v[0]))
}

func sameBits(a, b []float32) int {
	for i := range a {
		if math.Float32bits(a[i]) != math.Float32bits(b[i]) {
			return i
		}
	}
	return -1
}

func refPitchXcorrCoarse(x, y []float32, n, lags int, xcorr []float32) {
	for i := 0; i < lags; i++ {
		var sum, sum2, sum3, sum4 opus_val32
		j := 0
		for ; j+4 <= n; j += 4 {
			sum += x[j] * y[i+j]
			sum2 += x[j+1] * y[i+j+1]
			sum3 += x[j+2] * y[i+j+2]
			sum4 += x[j+3] * y[i+j+3]
		}
		sum += sum2 + sum3 + sum4
		for ; j < n; j++ {
			sum += x[j] * y[i+j]
		}
		if float32(-int32(1)) > sum {
			xcorr[i] = -1
		} else {
			xcorr[i] = sum
		}
	}
}

func refDotSerial(x, y []float32, n int) float32 {
	var sum1 opus_val32
	for j := 0; j < n; j++ {
		sum1 = sum1 + opus_val32(x[j]*y[j])
	}
	return sum1
}

func TestPitchXcorrKernels(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for iter := range 300 {
		n := 1 + r.IntN(300)
		lags := 1 + r.IntN(200)
		x := randFloats(r, n, 3000)
		y := randFloats(r, n+lags+4, 3000)
		want := make([]float32, lags)
		got := make([]float32, lags)
		refPitchXcorrCoarse(x, y, n, lags, want)
		pitchXcorrCoarse(ptr(x), ptr(y), int32(n), int32(lags), ptr(got))
		if i := sameBits(want, got); i >= 0 {
			t.Fatalf("coarse iter %d n=%d lag %d: got %v want %v", iter, n, i, got[i], want[i])
		}

		sel := []int32{}
		for k := 0; k < 1+r.IntN(4); k++ {
			sel = append(sel, int32(r.IntN(lags)))
		}
		fine := make([]float32, lags)
		pitchXcorrFine(ptr(x), ptr(y), int32(n), sel, ptr(fine))
		for _, lag := range sel {
			want := refDotSerial(x, y[lag:], n)
			if float32(-1) > want {
				want = -1
			}
			if math.Float32bits(fine[lag]) != math.Float32bits(want) {
				t.Fatalf("fine iter %d lag %d: got %v want %v", iter, lag, fine[lag], want)
			}
		}
	}
}

// ref_remove_doubling is remove_doubling as generated at 939ada7.
func ref_remove_doubling(tls *libc.TLS, x uintptr, maxperiod int32, minperiod int32, N int32, T0_ uintptr, prev_period int32, prev_gain opus_val16) (r opus_val16) {
	var T, T0, T1, T11, T1b, i, k, minperiod0, offset, v1 int32
	var best_xy, best_yy, xx, xy, yy, v2, v3 opus_val32
	var cont, g, g0, g1, pg, v5 opus_val16
	var xcorr [3]opus_val32
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = T, T0, T1, T11, T1b, best_xy, best_yy, cont, g, g0, g1, i, k, minperiod0, offset, pg, xcorr, xx, xy, yy, v1, v2, v3, v5
	minperiod0 = minperiod
	maxperiod = maxperiod / int32(2)
	minperiod = minperiod / int32(2)
	*(*int32)(unsafe.Pointer(T0_)) /= int32(2)
	prev_period = prev_period / int32(2)
	N = N / int32(2)
	x = x + uintptr(maxperiod)*4
	if *(*int32)(unsafe.Pointer(T0_)) >= maxperiod {
		*(*int32)(unsafe.Pointer(T0_)) = maxperiod - int32(1)
	}
	v1 = *(*int32)(unsafe.Pointer(T0_))
	T0 = v1
	T = v1
	v3 = float32(0)
	yy = v3
	v2 = v3
	xy = v2
	xx = v2
	i = 0
	for {
		if !(i < N) {
			break
		}
		xy = xy + opus_val32(*(*opus_val16)(unsafe.Pointer(x + uintptr(i)*4))**(*opus_val16)(unsafe.Pointer(x + uintptr(i-T0)*4)))
		xx = xx + opus_val32(*(*opus_val16)(unsafe.Pointer(x + uintptr(i)*4))**(*opus_val16)(unsafe.Pointer(x + uintptr(i)*4)))
		yy = yy + opus_val32(*(*opus_val16)(unsafe.Pointer(x + uintptr(i-T0)*4))**(*opus_val16)(unsafe.Pointer(x + uintptr(i-T0)*4)))
		goto _4
	_4:
		;
		i = i + 1
	}
	best_xy = xy
	best_yy = yy
	v5 = xy / float32(math.Sqrt(float64(float32(1)+opus_val32(xx*yy))))
	g0 = v5
	g = v5
	/* Look for any pitch at T/k */
	k = int32(2)
	for {
		if !(k <= int32(15)) {
			break
		}
		cont = float32(0)
		T1 = (int32(2)*T0 + k) / (int32(2) * k)
		if T1 < minperiod {
			break
		}
		/* Look for another strong correlation at T1b */
		if k == int32(2) {
			if T1+T0 > maxperiod {
				T1b = T0
			} else {
				T1b = T0 + T1
			}
		} else {
			T1b = (int32(2)*second_check[k]*T0 + k) / (int32(2) * k)
		}
		v2 = float32(0)
		yy = v2
		xy = v2
		i = 0
		for {
			if !(i < N) {
				break
			}
			xy = xy + opus_val32(*(*opus_val16)(unsafe.Pointer(x + uintptr(i)*4))**(*opus_val16)(unsafe.Pointer(x + uintptr(i-T1)*4)))
			yy = yy + opus_val32(*(*opus_val16)(unsafe.Pointer(x + uintptr(i-T1)*4))**(*opus_val16)(unsafe.Pointer(x + uintptr(i-T1)*4)))
			xy = xy + opus_val32(*(*opus_val16)(unsafe.Pointer(x + uintptr(i)*4))**(*opus_val16)(unsafe.Pointer(x + uintptr(i-T1b)*4)))
			yy = yy + opus_val32(*(*opus_val16)(unsafe.Pointer(x + uintptr(i-T1b)*4))**(*opus_val16)(unsafe.Pointer(x + uintptr(i-T1b)*4)))
			goto _8
		_8:
			;
			i = i + 1
		}
		g1 = xy / float32(math.Sqrt(float64(float32(1)+float32(float32(float32(libc.Float32FromFloat32(2)*xx)*libc.Float32FromFloat32(1))*yy))))
		if libc.Xabs(tls, T1-prev_period) <= int32(1) {
			cont = prev_gain
		} else {
			if libc.Xabs(tls, T1-prev_period) <= int32(2) && int32(5)*k*k < T0 {
				cont = float32(libc.Float32FromFloat32(0.5) * prev_gain)
			} else {
				cont = float32(0)
			}
		}
		if g1 > libc.Float32FromFloat32(0.3)+float32(libc.Float32FromFloat32(0.4)*g0)-cont {
			best_xy = xy
			best_yy = yy
			T = T1
			g = g1
		}
		goto _6
	_6:
		;
		k = k + 1
	}
	if best_yy <= best_xy {
		pg = libc.Float32FromFloat32(1)
	} else {
		pg = best_xy / (best_yy + float32(1))
	}
	k = 0
	for {
		if !(k < int32(3)) {
			break
		}
		T11 = T + k - int32(1)
		xy = float32(0)
		i = 0
		for {
			if !(i < N) {
				break
			}
			xy = xy + opus_val32(*(*opus_val16)(unsafe.Pointer(x + uintptr(i)*4))**(*opus_val16)(unsafe.Pointer(x + uintptr(i-T11)*4)))
			goto _10
		_10:
			;
			i = i + 1
		}
		xcorr[k] = xy
		goto _9
	_9:
		;
		k = k + 1
	}
	if xcorr[int32(2)]-xcorr[0] > float32(libc.Float32FromFloat32(0.7)*(xcorr[int32(1)]-xcorr[0])) {
		offset = int32(1)
	} else {
		if xcorr[0]-xcorr[int32(2)] > float32(libc.Float32FromFloat32(0.7)*(xcorr[int32(1)]-xcorr[int32(2)])) {
			offset = -int32(1)
		} else {
			offset = 0
		}
	}
	if pg > g {
		pg = g
	}
	*(*int32)(unsafe.Pointer(T0_)) = int32(2)*T + offset
	if *(*int32)(unsafe.Pointer(T0_)) < minperiod0 {
		*(*int32)(unsafe.Pointer(T0_)) = minperiod0
	}
	return pg
}

func TestRemoveDoubling(t *testing.T) {
	tls := libc.NewTLS()
	defer tls.Close()
	r := rand.New(rand.NewPCG(3, 4))
	const maxperiod, minperiod, n = 1024, 15, 960
	for iter := range 400 {
		x := make([]float32, maxperiod+n)
		period := 16 + r.Float64()*900
		noise := float32(r.IntN(4)) * 0.3
		for i := range x {
			x[i] = float32(1000*math.Sin(2*math.Pi*float64(i)/period)) + noise*1000*(r.Float32()*2-1)
		}
		t0 := int32(minperiod + r.IntN(maxperiod-minperiod))
		prevPeriod := int32(r.IntN(maxperiod))
		if r.IntN(2) == 0 {
			prevPeriod = t0/int32(2+r.IntN(4)) + int32(r.IntN(5)) - 2
		}
		prevGain := r.Float32()
		want, got := &cBuf[int32](1)[0], &cBuf[int32](1)[0]
		*want, *got = t0, t0
		wg := ref_remove_doubling(tls, ptr(x), maxperiod, minperiod, n, uintptr(unsafe.Pointer(want)), prevPeriod, prevGain)
		gg := remove_doubling(tls, ptr(x), maxperiod, minperiod, n, uintptr(unsafe.Pointer(got)), prevPeriod, prevGain)
		if *want != *got || math.Float32bits(wg) != math.Float32bits(gg) {
			t.Fatalf("iter %d: got (%d, %v) want (%d, %v)", iter, *got, gg, *want, wg)
		}
	}
}

// refcelt_fir is celt_fir as generated at 939ada7.
func refcelt_fir(tls *libc.TLS, x uintptr, num uintptr, y uintptr, N int32, ord int32, mem uintptr) {
	var i, j int32
	var sum opus_val32
	_, _, _ = i, j, sum
	i = 0
	for {
		if !(i < N) {
			break
		}
		sum = *(*opus_val16)(unsafe.Pointer(x + uintptr(i)*4))
		j = 0
		for {
			if !(j < ord) {
				break
			}
			sum = sum + opus_val32(*(*opus_val16)(unsafe.Pointer(num + uintptr(j)*4))**(*opus_val16)(unsafe.Pointer(mem + uintptr(j)*4)))
			goto _2
		_2:
			;
			j = j + 1
		}
		j = ord - int32(1)
		for {
			if !(j >= int32(1)) {
				break
			}
			*(*opus_val16)(unsafe.Pointer(mem + uintptr(j)*4)) = *(*opus_val16)(unsafe.Pointer(mem + uintptr(j-int32(1))*4))
			goto _3
		_3:
			;
			j = j - 1
		}
		*(*opus_val16)(unsafe.Pointer(mem)) = *(*opus_val16)(unsafe.Pointer(x + uintptr(i)*4))
		*(*opus_val16)(unsafe.Pointer(y + uintptr(i)*4)) = sum
		goto _1
	_1:
		;
		i = i + 1
	}
}

// ref_celt_autocorr is _celt_autocorr as generated at 939ada7.
func ref_celt_autocorr(tls *libc.TLS, x uintptr, ac uintptr, window uintptr, overlap int32, lag int32, n int32) {
	var d opus_val32
	var i int32
	var xx uintptr
	_, _, _ = d, i, xx
	_sp := _arenaSave()
	defer _arenaRestore(_sp)
	xx = _arenaAlloc(uint64(4) * uint64(n))
	i = 0
	for {
		if !(i < n) {
			break
		}
		*(*opus_val16)(unsafe.Pointer(xx + uintptr(i)*4)) = *(*opus_val16)(unsafe.Pointer(x + uintptr(i)*4))
		goto _1
	_1:
		;
		i = i + 1
	}
	i = 0
	for {
		if !(i < overlap) {
			break
		}
		*(*opus_val16)(unsafe.Pointer(xx + uintptr(i)*4)) = opus_val16(*(*opus_val16)(unsafe.Pointer(x + uintptr(i)*4)) * *(*opus_val16)(unsafe.Pointer(window + uintptr(i)*4)))
		*(*opus_val16)(unsafe.Pointer(xx + uintptr(n-i-int32(1))*4)) = opus_val16(*(*opus_val16)(unsafe.Pointer(x + uintptr(n-i-int32(1))*4)) * *(*opus_val16)(unsafe.Pointer(window + uintptr(i)*4)))
		goto _2
	_2:
		;
		i = i + 1
	}
	for lag >= 0 {
		i = lag
		d = float32(0)
		for {
			if !(i < n) {
				break
			}
			d = d + opus_val32(*(*opus_val16)(unsafe.Pointer(xx + uintptr(i)*4))**(*opus_val16)(unsafe.Pointer(xx + uintptr(i-lag)*4)))
			goto _3
		_3:
			;
			i = i + 1
		}
		*(*opus_val32)(unsafe.Pointer(ac + uintptr(lag)*4)) = d
		/*printf ("%f ", ac[lag]);*/
		lag = lag - 1
	}
	/*printf ("\n");*/
	*(*opus_val32)(unsafe.Pointer(ac)) += float32(10)
}

func TestCeltFir(t *testing.T) {
	tls := libc.NewTLS()
	defer tls.Close()
	r := rand.New(rand.NewPCG(5, 6))
	for iter := range 300 {
		ord := []int{1, 4, 24, 3}[iter%4]
		n := 1 + r.IntN(600)
		x := randFloats(r, n, 2000)
		num := randFloats(r, ord, 1)
		mem := randFloats(r, ord, 2000)
		inPlace := r.IntN(2) == 0
		wx, gx := append([]float32(nil), x...), append([]float32(nil), x...)
		wy, gy := wx, gx
		if !inPlace {
			wy, gy = make([]float32, n), make([]float32, n)
		}
		wm, gm := append([]float32(nil), mem...), append([]float32(nil), mem...)
		refcelt_fir(tls, ptr(wx), ptr(num), ptr(wy), int32(n), int32(ord), ptr(wm))
		celt_fir(tls, ptr(gx), ptr(num), ptr(gy), int32(n), int32(ord), ptr(gm))
		if i := sameBits(wy, gy); i >= 0 {
			t.Fatalf("iter %d ord %d n %d inPlace %v: y[%d] got %v want %v", iter, ord, n, inPlace, i, gy[i], wy[i])
		}
		if i := sameBits(wm, gm); i >= 0 {
			t.Fatalf("iter %d ord %d: mem[%d] differs", iter, ord, i)
		}
	}
}

func TestCeltAutocorr(t *testing.T) {
	tls := libc.NewTLS()
	defer tls.Close()
	r := rand.New(rand.NewPCG(7, 8))
	for iter := range 300 {
		n := 1 + r.IntN(1100)
		lag := r.IntN(25)
		overlap := 0
		var window uintptr
		win := randFloats(r, 130, 1)
		if r.IntN(2) == 0 {
			overlap = min(r.IntN(121), n/2)
			window = ptr(win)
		}
		x := randFloats(r, n, 3000)
		want := make([]float32, lag+1)
		got := make([]float32, lag+1)
		ref_celt_autocorr(tls, ptr(x), ptr(want), window, int32(overlap), int32(lag), int32(n))
		_celt_autocorr(tls, ptr(x), ptr(got), window, int32(overlap), int32(lag), int32(n))
		if i := sameBits(want, got); i >= 0 {
			t.Fatalf("iter %d n %d lag %d: ac[%d] got %v want %v", iter, n, lag, i, got[i], want[i])
		}
	}
}

func TestPreemphStereo(t *testing.T) {
	r := rand.New(rand.NewPCG(9, 10))
	coef := [4]opus_val16{0.8500061035, 0.2, 0.9, 1.0}
	for iter := range 200 {
		n := 1 + r.IntN(960)
		pcm := randFloats(r, 2*n, 2.5)
		if iter%5 == 0 {
			pcm[r.IntN(2*n)] = float32(math.NaN())
		}
		if iter%7 == 0 {
			clear(pcm)
		}
		clip := iter%2 == 0
		mem := [2]opus_val32{r.Float32() * 100, r.Float32() * 100}
		wmem := mem
		want := make([]float32, 2*n)
		wantNonzero := false
		for c := range 2 {
			m := wmem[c]
			for i := range n {
				x := opus_val16(pcm[2*i+c] * libc.Float32FromFloat32(32768))
				if !(x == x) {
					x = 0
				}
				if clip {
					if libc.Float32FromFloat32(65536) < x {
						x = 65536
					} else if -libc.Float32FromFloat32(65536) > x {
						x = -65536
					}
				}
				tmp1 := opus_val32(coef[2] * x)
				want[c*n+i] = tmp1 + m
				m = opus_val16(coef[1]*want[c*n+i]) - opus_val16(coef[0]*tmp1)
				if want[c*n+i] != 0 {
					wantNonzero = true
				}
			}
			wmem[c] = m
		}
		got := make([]float32, 2*n)
		nonzero := preemphStereo(ptr(pcm), ptr(got), ptr(got[n:]), int32(n), clip, &coef, &mem)
		if i := sameBits(want, got); i >= 0 || nonzero != wantNonzero || mem != wmem {
			t.Fatalf("iter %d: mismatch at %d (nonzero %v/%v, mem %v/%v)", iter, i, nonzero, wantNonzero, mem, wmem)
		}
	}
}

// refdeemphasis is deemphasis as generated at 939ada7.
func refdeemphasis(tls *libc.TLS, in uintptr, pcm uintptr, N int32, C int32, downsample int32, coef uintptr, mem uintptr) {
	var c, count, j, v1 int32
	var m, tmp celt_sig
	var x, y uintptr
	_, _, _, _, _, _, _, _ = c, count, j, m, tmp, x, y, v1
	count = 0
	c = 0
	for {
		m = *(*celt_sig)(unsafe.Pointer(mem + uintptr(c)*4))
		x = *(*uintptr)(unsafe.Pointer(in + uintptr(c)*8))
		y = pcm + uintptr(c)*4
		j = 0
		for {
			if !(j < N) {
				break
			}
			tmp = *(*celt_sig)(unsafe.Pointer(x)) + m
			m = opus_val16(*(*opus_val16)(unsafe.Pointer(coef))*tmp) - opus_val16(*(*opus_val16)(unsafe.Pointer(coef + 1*4))**(*celt_sig)(unsafe.Pointer(x)))
			tmp = opus_val16(*(*opus_val16)(unsafe.Pointer(coef + 3*4)) * tmp)
			x += 4
			/* Technically the store could be moved outside of the if because
			   the stores we don't want will just be overwritten */
			if count == 0 {
				*(*opus_val16)(unsafe.Pointer(y)) = opus_val16(SIG2WORD16(tls, tmp) * (float32(1) / libc.Float32FromFloat32(32768)))
			}
			count = count + 1
			v1 = count
			if v1 == downsample {
				y = y + uintptr(C)*4
				count = 0
			}
			goto _3
		_3:
			;
			j = j + 1
		}
		*(*celt_sig)(unsafe.Pointer(mem + uintptr(c)*4)) = m
		goto _2
	_2:
		;
		c = c + 1
		v1 = c
		if !(v1 < C) {
			break
		}
	}
}

func TestDeemphasis(t *testing.T) {
	tls := libc.NewTLS()
	defer tls.Close()
	r := rand.New(rand.NewPCG(11, 12))
	coef := cBuf[opus_val16](4)
	copy(coef, []opus_val16{0.8500061035, 0.2, 0.5, 0.9})
	for iter := range 200 {
		C := 1 + iter%2
		down := []int{1, 1, 2, 3}[r.IntN(4)]
		n := down * (1 + r.IntN(480))
		in := [2][]float32{randFloats(r, n, 30000), randFloats(r, n, 30000)}
		ins := cBuf[ptrslot](2)
		ins[0][0], ins[1][0] = ptr(in[0]), ptr(in[1])
		wmem, gmem := cBuf[float32](2), cBuf[float32](2)
		wmem[0], wmem[1] = r.Float32()*1000, r.Float32()*1000
		copy(gmem, wmem)
		want := make([]float32, C*n/down)
		got := make([]float32, C*n/down)
		refdeemphasis(tls, uintptr(unsafe.Pointer(&ins[0])), ptr(want), int32(n), int32(C), int32(down), ptr(coef), ptr(wmem))
		deemphasis(tls, uintptr(unsafe.Pointer(&ins[0])), ptr(got), int32(n), int32(C), int32(down), ptr(coef), ptr(gmem))
		if i := sameBits(want, got); i >= 0 || sameBits(wmem, gmem) >= 0 {
			t.Fatalf("iter %d C %d down %d: mismatch at %d, mem %v/%v", iter, C, down, i, gmem, wmem)
		}
	}
}

func ref_kf_bfly2(tls *libc.TLS, Fout uintptr, fstride size_t, st uintptr, m int32, N int32, mm int32) {
	var Fout2, Fout_beg, tw1 uintptr
	var i, j int32
	var t kiss_fft_cpx
	_, _, _, _, _, _ = Fout2, Fout_beg, i, j, t, tw1
	Fout_beg = Fout
	i = 0
	for {
		if !(i < N) {
			break
		}
		Fout = Fout_beg + uintptr(i*mm)*8
		Fout2 = Fout + uintptr(m)*8
		tw1 = (*kiss_fft_state)(unsafe.Pointer(st)).Ftwiddles[0]
		j = 0
		for {
			if !(j < m) {
				break
			}
			(*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fr = (*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fr
			(*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fi = (*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fi
			(*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fr = (*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fr
			(*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fi = (*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fi
			t.Fr = float32((*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fr*(*kiss_twiddle_cpx)(unsafe.Pointer(tw1)).Fr) - float32((*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fi*(*kiss_twiddle_cpx)(unsafe.Pointer(tw1)).Fi)
			t.Fi = float32((*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fr*(*kiss_twiddle_cpx)(unsafe.Pointer(tw1)).Fi) + float32((*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fi*(*kiss_twiddle_cpx)(unsafe.Pointer(tw1)).Fr)
			tw1 = tw1 + uintptr(fstride)*8
			(*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fr = (*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fr - t.Fr
			(*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fi = (*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fi - t.Fi
			(*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fr += t.Fr
			(*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fi += t.Fi
			Fout2 += 8
			Fout += 8
			goto _2
		_2:
			;
			j = j + 1
		}
		goto _1
	_1:
		;
		i = i + 1
	}
}

func ref_ki_bfly2(tls *libc.TLS, Fout uintptr, fstride size_t, st uintptr, m int32, N int32, mm int32) {
	var Fout2, Fout_beg, tw1 uintptr
	var i, j int32
	var t kiss_fft_cpx
	_, _, _, _, _, _ = Fout2, Fout_beg, i, j, t, tw1
	Fout_beg = Fout
	i = 0
	for {
		if !(i < N) {
			break
		}
		Fout = Fout_beg + uintptr(i*mm)*8
		Fout2 = Fout + uintptr(m)*8
		tw1 = (*kiss_fft_state)(unsafe.Pointer(st)).Ftwiddles[0]
		j = 0
		for {
			if !(j < m) {
				break
			}
			t.Fr = float32((*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fr*(*kiss_twiddle_cpx)(unsafe.Pointer(tw1)).Fr) + float32((*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fi*(*kiss_twiddle_cpx)(unsafe.Pointer(tw1)).Fi)
			t.Fi = float32((*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fi*(*kiss_twiddle_cpx)(unsafe.Pointer(tw1)).Fr) - float32((*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fr*(*kiss_twiddle_cpx)(unsafe.Pointer(tw1)).Fi)
			tw1 = tw1 + uintptr(fstride)*8
			(*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fr = (*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fr - t.Fr
			(*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fi = (*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fi - t.Fi
			(*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fr += t.Fr
			(*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fi += t.Fi
			Fout2 += 8
			Fout += 8
			goto _2
		_2:
			;
			j = j + 1
		}
		goto _1
	_1:
		;
		i = i + 1
	}
}

func ref_kf_bfly3(tls *libc.TLS, Fout uintptr, fstride size_t, st uintptr, m int32, N int32, mm int32) {
	var Fout_beg, tw1, tw2, v2 uintptr
	var epi3 kiss_twiddle_cpx
	var i int32
	var k, m2, v3 size_t
	var scratch [5]kiss_fft_cpx
	_, _, _, _, _, _, _, _, _, _ = Fout_beg, epi3, i, k, m2, scratch, tw1, tw2, v2, v3
	m2 = uint64(int32(2) * m)
	Fout_beg = Fout
	epi3 = *(*kiss_twiddle_cpx)(unsafe.Pointer((*kiss_fft_state)(unsafe.Pointer(st)).Ftwiddles[0] + uintptr(fstride*uint64(m))*8))
	i = 0
	for {
		if !(i < N) {
			break
		}
		Fout = Fout_beg + uintptr(i*mm)*8
		v2 = (*kiss_fft_state)(unsafe.Pointer(st)).Ftwiddles[0]
		tw2 = v2
		tw1 = v2
		k = uint64(m)
		for {
			scratch[int32(1)].Fr = float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fr*(*kiss_twiddle_cpx)(unsafe.Pointer(tw1)).Fr) - float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fi*(*kiss_twiddle_cpx)(unsafe.Pointer(tw1)).Fi)
			scratch[int32(1)].Fi = float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fr*(*kiss_twiddle_cpx)(unsafe.Pointer(tw1)).Fi) + float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fi*(*kiss_twiddle_cpx)(unsafe.Pointer(tw1)).Fr)
			scratch[int32(2)].Fr = float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fr*(*kiss_twiddle_cpx)(unsafe.Pointer(tw2)).Fr) - float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fi*(*kiss_twiddle_cpx)(unsafe.Pointer(tw2)).Fi)
			scratch[int32(2)].Fi = float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fr*(*kiss_twiddle_cpx)(unsafe.Pointer(tw2)).Fi) + float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fi*(*kiss_twiddle_cpx)(unsafe.Pointer(tw2)).Fr)
			scratch[int32(3)].Fr = scratch[int32(1)].Fr + scratch[int32(2)].Fr
			scratch[int32(3)].Fi = scratch[int32(1)].Fi + scratch[int32(2)].Fi
			scratch[0].Fr = scratch[int32(1)].Fr - scratch[int32(2)].Fr
			scratch[0].Fi = scratch[int32(1)].Fi - scratch[int32(2)].Fi
			tw1 = tw1 + uintptr(fstride)*8
			tw2 = tw2 + uintptr(fstride*uint64(2))*8
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fr = (*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fr - float32(scratch[int32(3)].Fr*libc.Float32FromFloat32(0.5))
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fi = (*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fi - float32(scratch[int32(3)].Fi*libc.Float32FromFloat32(0.5))
			scratch[0].Fr *= epi3.Fi
			scratch[0].Fi *= epi3.Fi
			(*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fr += scratch[int32(3)].Fr
			(*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fi += scratch[int32(3)].Fi
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fr = (*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fr + scratch[0].Fi
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fi = (*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fi - scratch[0].Fr
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fr -= scratch[0].Fi
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fi += scratch[0].Fr
			Fout += 8
			goto _4
		_4:
			;
			k = k - 1
			v3 = k
			if !(v3 != 0) {
				break
			}
		}
		goto _1
	_1:
		;
		i = i + 1
	}
}

func ref_ki_bfly3(tls *libc.TLS, Fout uintptr, fstride size_t, st uintptr, m int32, N int32, mm int32) {
	var Fout_beg, tw1, tw2, v2 uintptr
	var epi3 kiss_twiddle_cpx
	var i, k, v3 int32
	var m2 size_t
	var scratch [5]kiss_fft_cpx
	_, _, _, _, _, _, _, _, _, _ = Fout_beg, epi3, i, k, m2, scratch, tw1, tw2, v2, v3
	m2 = uint64(int32(2) * m)
	Fout_beg = Fout
	epi3 = *(*kiss_twiddle_cpx)(unsafe.Pointer((*kiss_fft_state)(unsafe.Pointer(st)).Ftwiddles[0] + uintptr(fstride*uint64(m))*8))
	i = 0
	for {
		if !(i < N) {
			break
		}
		Fout = Fout_beg + uintptr(i*mm)*8
		v2 = (*kiss_fft_state)(unsafe.Pointer(st)).Ftwiddles[0]
		tw2 = v2
		tw1 = v2
		k = m
		for {
			scratch[int32(1)].Fr = float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fr*(*kiss_twiddle_cpx)(unsafe.Pointer(tw1)).Fr) + float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fi*(*kiss_twiddle_cpx)(unsafe.Pointer(tw1)).Fi)
			scratch[int32(1)].Fi = float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fi*(*kiss_twiddle_cpx)(unsafe.Pointer(tw1)).Fr) - float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fr*(*kiss_twiddle_cpx)(unsafe.Pointer(tw1)).Fi)
			scratch[int32(2)].Fr = float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fr*(*kiss_twiddle_cpx)(unsafe.Pointer(tw2)).Fr) + float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fi*(*kiss_twiddle_cpx)(unsafe.Pointer(tw2)).Fi)
			scratch[int32(2)].Fi = float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fi*(*kiss_twiddle_cpx)(unsafe.Pointer(tw2)).Fr) - float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fr*(*kiss_twiddle_cpx)(unsafe.Pointer(tw2)).Fi)
			scratch[int32(3)].Fr = scratch[int32(1)].Fr + scratch[int32(2)].Fr
			scratch[int32(3)].Fi = scratch[int32(1)].Fi + scratch[int32(2)].Fi
			scratch[0].Fr = scratch[int32(1)].Fr - scratch[int32(2)].Fr
			scratch[0].Fi = scratch[int32(1)].Fi - scratch[int32(2)].Fi
			tw1 = tw1 + uintptr(fstride)*8
			tw2 = tw2 + uintptr(fstride*uint64(2))*8
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fr = (*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fr - float32(scratch[int32(3)].Fr*libc.Float32FromFloat32(0.5))
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fi = (*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fi - float32(scratch[int32(3)].Fi*libc.Float32FromFloat32(0.5))
			scratch[0].Fr *= -epi3.Fi
			scratch[0].Fi *= -epi3.Fi
			(*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fr += scratch[int32(3)].Fr
			(*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fi += scratch[int32(3)].Fi
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fr = (*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fr + scratch[0].Fi
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fi = (*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fi - scratch[0].Fr
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fr -= scratch[0].Fi
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fi += scratch[0].Fr
			Fout += 8
			goto _4
		_4:
			;
			k = k - 1
			v3 = k
			if !(v3 != 0) {
				break
			}
		}
		goto _1
	_1:
		;
		i = i + 1
	}
}

func ref_kf_bfly4(tls *libc.TLS, Fout uintptr, fstride size_t, st uintptr, m int32, N int32, mm int32) {
	var Fout_beg, tw1, tw2, tw3, v2, v3 uintptr
	var i, j int32
	var m2, m3 size_t
	var scratch [6]kiss_fft_cpx
	_, _, _, _, _, _, _, _, _, _, _ = Fout_beg, i, j, m2, m3, scratch, tw1, tw2, tw3, v2, v3
	m2 = uint64(int32(2) * m)
	m3 = uint64(int32(3) * m)
	Fout_beg = Fout
	i = 0
	for {
		if !(i < N) {
			break
		}
		Fout = Fout_beg + uintptr(i*mm)*8
		v3 = (*kiss_fft_state)(unsafe.Pointer(st)).Ftwiddles[0]
		tw1 = v3
		v2 = v3
		tw2 = v2
		tw3 = v2
		j = 0
		for {
			if !(j < m) {
				break
			}
			scratch[0].Fr = float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fr*(*kiss_twiddle_cpx)(unsafe.Pointer(tw1)).Fr) - float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fi*(*kiss_twiddle_cpx)(unsafe.Pointer(tw1)).Fi)
			scratch[0].Fi = float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fr*(*kiss_twiddle_cpx)(unsafe.Pointer(tw1)).Fi) + float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fi*(*kiss_twiddle_cpx)(unsafe.Pointer(tw1)).Fr)
			scratch[int32(1)].Fr = float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fr*(*kiss_twiddle_cpx)(unsafe.Pointer(tw2)).Fr) - float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fi*(*kiss_twiddle_cpx)(unsafe.Pointer(tw2)).Fi)
			scratch[int32(1)].Fi = float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fr*(*kiss_twiddle_cpx)(unsafe.Pointer(tw2)).Fi) + float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fi*(*kiss_twiddle_cpx)(unsafe.Pointer(tw2)).Fr)
			scratch[int32(2)].Fr = float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m3)*8))).Fr*(*kiss_twiddle_cpx)(unsafe.Pointer(tw3)).Fr) - float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m3)*8))).Fi*(*kiss_twiddle_cpx)(unsafe.Pointer(tw3)).Fi)
			scratch[int32(2)].Fi = float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m3)*8))).Fr*(*kiss_twiddle_cpx)(unsafe.Pointer(tw3)).Fi) + float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m3)*8))).Fi*(*kiss_twiddle_cpx)(unsafe.Pointer(tw3)).Fr)
			(*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fr = (*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fr
			(*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fi = (*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fi
			scratch[int32(5)].Fr = (*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fr - scratch[int32(1)].Fr
			scratch[int32(5)].Fi = (*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fi - scratch[int32(1)].Fi
			(*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fr += scratch[int32(1)].Fr
			(*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fi += scratch[int32(1)].Fi
			scratch[int32(3)].Fr = scratch[0].Fr + scratch[int32(2)].Fr
			scratch[int32(3)].Fi = scratch[0].Fi + scratch[int32(2)].Fi
			scratch[int32(4)].Fr = scratch[0].Fr - scratch[int32(2)].Fr
			scratch[int32(4)].Fi = scratch[0].Fi - scratch[int32(2)].Fi
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fr = (*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fr
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fi = (*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fi
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fr = (*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fr - scratch[int32(3)].Fr
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fi = (*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fi - scratch[int32(3)].Fi
			tw1 = tw1 + uintptr(fstride)*8
			tw2 = tw2 + uintptr(fstride*uint64(2))*8
			tw3 = tw3 + uintptr(fstride*uint64(3))*8
			(*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fr += scratch[int32(3)].Fr
			(*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fi += scratch[int32(3)].Fi
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fr = scratch[int32(5)].Fr + scratch[int32(4)].Fi
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fi = scratch[int32(5)].Fi - scratch[int32(4)].Fr
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m3)*8))).Fr = scratch[int32(5)].Fr - scratch[int32(4)].Fi
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m3)*8))).Fi = scratch[int32(5)].Fi + scratch[int32(4)].Fr
			Fout += 8
			goto _4
		_4:
			;
			j = j + 1
		}
		goto _1
	_1:
		;
		i = i + 1
	}
}

func ref_ki_bfly4(tls *libc.TLS, Fout uintptr, fstride size_t, st uintptr, m int32, N int32, mm int32) {
	var Fout_beg, tw1, tw2, tw3, v2, v3 uintptr
	var i, j int32
	var m2, m3 size_t
	var scratch [6]kiss_fft_cpx
	_, _, _, _, _, _, _, _, _, _, _ = Fout_beg, i, j, m2, m3, scratch, tw1, tw2, tw3, v2, v3
	m2 = uint64(int32(2) * m)
	m3 = uint64(int32(3) * m)
	Fout_beg = Fout
	i = 0
	for {
		if !(i < N) {
			break
		}
		Fout = Fout_beg + uintptr(i*mm)*8
		v3 = (*kiss_fft_state)(unsafe.Pointer(st)).Ftwiddles[0]
		tw1 = v3
		v2 = v3
		tw2 = v2
		tw3 = v2
		j = 0
		for {
			if !(j < m) {
				break
			}
			scratch[0].Fr = float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fr*(*kiss_twiddle_cpx)(unsafe.Pointer(tw1)).Fr) + float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fi*(*kiss_twiddle_cpx)(unsafe.Pointer(tw1)).Fi)
			scratch[0].Fi = float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fi*(*kiss_twiddle_cpx)(unsafe.Pointer(tw1)).Fr) - float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fr*(*kiss_twiddle_cpx)(unsafe.Pointer(tw1)).Fi)
			scratch[int32(1)].Fr = float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fr*(*kiss_twiddle_cpx)(unsafe.Pointer(tw2)).Fr) + float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fi*(*kiss_twiddle_cpx)(unsafe.Pointer(tw2)).Fi)
			scratch[int32(1)].Fi = float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fi*(*kiss_twiddle_cpx)(unsafe.Pointer(tw2)).Fr) - float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fr*(*kiss_twiddle_cpx)(unsafe.Pointer(tw2)).Fi)
			scratch[int32(2)].Fr = float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m3)*8))).Fr*(*kiss_twiddle_cpx)(unsafe.Pointer(tw3)).Fr) + float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m3)*8))).Fi*(*kiss_twiddle_cpx)(unsafe.Pointer(tw3)).Fi)
			scratch[int32(2)].Fi = float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m3)*8))).Fi*(*kiss_twiddle_cpx)(unsafe.Pointer(tw3)).Fr) - float32((*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m3)*8))).Fr*(*kiss_twiddle_cpx)(unsafe.Pointer(tw3)).Fi)
			scratch[int32(5)].Fr = (*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fr - scratch[int32(1)].Fr
			scratch[int32(5)].Fi = (*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fi - scratch[int32(1)].Fi
			(*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fr += scratch[int32(1)].Fr
			(*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fi += scratch[int32(1)].Fi
			scratch[int32(3)].Fr = scratch[0].Fr + scratch[int32(2)].Fr
			scratch[int32(3)].Fi = scratch[0].Fi + scratch[int32(2)].Fi
			scratch[int32(4)].Fr = scratch[0].Fr - scratch[int32(2)].Fr
			scratch[int32(4)].Fi = scratch[0].Fi - scratch[int32(2)].Fi
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fr = (*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fr - scratch[int32(3)].Fr
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m2)*8))).Fi = (*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fi - scratch[int32(3)].Fi
			tw1 = tw1 + uintptr(fstride)*8
			tw2 = tw2 + uintptr(fstride*uint64(2))*8
			tw3 = tw3 + uintptr(fstride*uint64(3))*8
			(*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fr += scratch[int32(3)].Fr
			(*kiss_fft_cpx)(unsafe.Pointer(Fout)).Fi += scratch[int32(3)].Fi
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fr = scratch[int32(5)].Fr - scratch[int32(4)].Fi
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m)*8))).Fi = scratch[int32(5)].Fi + scratch[int32(4)].Fr
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m3)*8))).Fr = scratch[int32(5)].Fr + scratch[int32(4)].Fi
			(*(*kiss_fft_cpx)(unsafe.Pointer(Fout + uintptr(m3)*8))).Fi = scratch[int32(5)].Fi - scratch[int32(4)].Fr
			Fout += 8
			goto _4
		_4:
			;
			j = j + 1
		}
		goto _1
	_1:
		;
		i = i + 1
	}
}

func ref_kf_bfly5(tls *libc.TLS, Fout uintptr, fstride size_t, st uintptr, m int32, N int32, mm int32) {
	var Fout0, Fout1, Fout2, Fout3, Fout4, Fout_beg, tw, twiddles uintptr
	var i, u int32
	var scratch [13]kiss_fft_cpx
	var ya, yb kiss_twiddle_cpx
	_, _, _, _, _, _, _, _, _, _, _, _, _ = Fout0, Fout1, Fout2, Fout3, Fout4, Fout_beg, i, scratch, tw, twiddles, u, ya, yb
	twiddles = (*kiss_fft_state)(unsafe.Pointer(st)).Ftwiddles[0]
	Fout_beg = Fout
	ya = *(*kiss_twiddle_cpx)(unsafe.Pointer(twiddles + uintptr(fstride*uint64(m))*8))
	yb = *(*kiss_twiddle_cpx)(unsafe.Pointer(twiddles + uintptr(fstride*uint64(2)*uint64(m))*8))
	tw = (*kiss_fft_state)(unsafe.Pointer(st)).Ftwiddles[0]
	i = 0
	for {
		if !(i < N) {
			break
		}
		Fout = Fout_beg + uintptr(i*mm)*8
		Fout0 = Fout
		Fout1 = Fout0 + uintptr(m)*8
		Fout2 = Fout0 + uintptr(int32(2)*m)*8
		Fout3 = Fout0 + uintptr(int32(3)*m)*8
		Fout4 = Fout0 + uintptr(int32(4)*m)*8
		u = 0
		for {
			if !(u < m) {
				break
			}
			scratch[0] = *(*kiss_fft_cpx)(unsafe.Pointer(Fout0))
			scratch[int32(1)].Fr = float32((*kiss_fft_cpx)(unsafe.Pointer(Fout1)).Fr*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(u)*fstride)*8))).Fr) - float32((*kiss_fft_cpx)(unsafe.Pointer(Fout1)).Fi*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(u)*fstride)*8))).Fi)
			scratch[int32(1)].Fi = float32((*kiss_fft_cpx)(unsafe.Pointer(Fout1)).Fr*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(u)*fstride)*8))).Fi) + float32((*kiss_fft_cpx)(unsafe.Pointer(Fout1)).Fi*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(u)*fstride)*8))).Fr)
			scratch[int32(2)].Fr = float32((*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fr*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(int32(2)*u)*fstride)*8))).Fr) - float32((*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fi*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(int32(2)*u)*fstride)*8))).Fi)
			scratch[int32(2)].Fi = float32((*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fr*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(int32(2)*u)*fstride)*8))).Fi) + float32((*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fi*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(int32(2)*u)*fstride)*8))).Fr)
			scratch[int32(3)].Fr = float32((*kiss_fft_cpx)(unsafe.Pointer(Fout3)).Fr*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(int32(3)*u)*fstride)*8))).Fr) - float32((*kiss_fft_cpx)(unsafe.Pointer(Fout3)).Fi*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(int32(3)*u)*fstride)*8))).Fi)
			scratch[int32(3)].Fi = float32((*kiss_fft_cpx)(unsafe.Pointer(Fout3)).Fr*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(int32(3)*u)*fstride)*8))).Fi) + float32((*kiss_fft_cpx)(unsafe.Pointer(Fout3)).Fi*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(int32(3)*u)*fstride)*8))).Fr)
			scratch[int32(4)].Fr = float32((*kiss_fft_cpx)(unsafe.Pointer(Fout4)).Fr*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(int32(4)*u)*fstride)*8))).Fr) - float32((*kiss_fft_cpx)(unsafe.Pointer(Fout4)).Fi*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(int32(4)*u)*fstride)*8))).Fi)
			scratch[int32(4)].Fi = float32((*kiss_fft_cpx)(unsafe.Pointer(Fout4)).Fr*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(int32(4)*u)*fstride)*8))).Fi) + float32((*kiss_fft_cpx)(unsafe.Pointer(Fout4)).Fi*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(int32(4)*u)*fstride)*8))).Fr)
			scratch[int32(7)].Fr = scratch[int32(1)].Fr + scratch[int32(4)].Fr
			scratch[int32(7)].Fi = scratch[int32(1)].Fi + scratch[int32(4)].Fi
			scratch[int32(10)].Fr = scratch[int32(1)].Fr - scratch[int32(4)].Fr
			scratch[int32(10)].Fi = scratch[int32(1)].Fi - scratch[int32(4)].Fi
			scratch[int32(8)].Fr = scratch[int32(2)].Fr + scratch[int32(3)].Fr
			scratch[int32(8)].Fi = scratch[int32(2)].Fi + scratch[int32(3)].Fi
			scratch[int32(9)].Fr = scratch[int32(2)].Fr - scratch[int32(3)].Fr
			scratch[int32(9)].Fi = scratch[int32(2)].Fi - scratch[int32(3)].Fi
			*(*float32)(unsafe.Pointer(Fout0)) += scratch[int32(7)].Fr + scratch[int32(8)].Fr
			*(*float32)(unsafe.Pointer(Fout0 + 4)) += scratch[int32(7)].Fi + scratch[int32(8)].Fi
			scratch[int32(5)].Fr = scratch[0].Fr + float32(scratch[int32(7)].Fr*ya.Fr) + float32(scratch[int32(8)].Fr*yb.Fr)
			scratch[int32(5)].Fi = scratch[0].Fi + float32(scratch[int32(7)].Fi*ya.Fr) + float32(scratch[int32(8)].Fi*yb.Fr)
			scratch[int32(6)].Fr = float32(scratch[int32(10)].Fi*ya.Fi) + float32(scratch[int32(9)].Fi*yb.Fi)
			scratch[int32(6)].Fi = -float32(scratch[int32(10)].Fr*ya.Fi) - float32(scratch[int32(9)].Fr*yb.Fi)
			(*kiss_fft_cpx)(unsafe.Pointer(Fout1)).Fr = scratch[int32(5)].Fr - scratch[int32(6)].Fr
			(*kiss_fft_cpx)(unsafe.Pointer(Fout1)).Fi = scratch[int32(5)].Fi - scratch[int32(6)].Fi
			(*kiss_fft_cpx)(unsafe.Pointer(Fout4)).Fr = scratch[int32(5)].Fr + scratch[int32(6)].Fr
			(*kiss_fft_cpx)(unsafe.Pointer(Fout4)).Fi = scratch[int32(5)].Fi + scratch[int32(6)].Fi
			scratch[int32(11)].Fr = scratch[0].Fr + float32(scratch[int32(7)].Fr*yb.Fr) + float32(scratch[int32(8)].Fr*ya.Fr)
			scratch[int32(11)].Fi = scratch[0].Fi + float32(scratch[int32(7)].Fi*yb.Fr) + float32(scratch[int32(8)].Fi*ya.Fr)
			scratch[int32(12)].Fr = -float32(scratch[int32(10)].Fi*yb.Fi) + float32(scratch[int32(9)].Fi*ya.Fi)
			scratch[int32(12)].Fi = float32(scratch[int32(10)].Fr*yb.Fi) - float32(scratch[int32(9)].Fr*ya.Fi)
			(*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fr = scratch[int32(11)].Fr + scratch[int32(12)].Fr
			(*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fi = scratch[int32(11)].Fi + scratch[int32(12)].Fi
			(*kiss_fft_cpx)(unsafe.Pointer(Fout3)).Fr = scratch[int32(11)].Fr - scratch[int32(12)].Fr
			(*kiss_fft_cpx)(unsafe.Pointer(Fout3)).Fi = scratch[int32(11)].Fi - scratch[int32(12)].Fi
			Fout0 += 8
			Fout1 += 8
			Fout2 += 8
			Fout3 += 8
			Fout4 += 8
			goto _2
		_2:
			;
			u = u + 1
		}
		goto _1
	_1:
		;
		i = i + 1
	}
}

func ref_ki_bfly5(tls *libc.TLS, Fout uintptr, fstride size_t, st uintptr, m int32, N int32, mm int32) {
	var Fout0, Fout1, Fout2, Fout3, Fout4, Fout_beg, tw, twiddles uintptr
	var i, u int32
	var scratch [13]kiss_fft_cpx
	var ya, yb kiss_twiddle_cpx
	_, _, _, _, _, _, _, _, _, _, _, _, _ = Fout0, Fout1, Fout2, Fout3, Fout4, Fout_beg, i, scratch, tw, twiddles, u, ya, yb
	twiddles = (*kiss_fft_state)(unsafe.Pointer(st)).Ftwiddles[0]
	Fout_beg = Fout
	ya = *(*kiss_twiddle_cpx)(unsafe.Pointer(twiddles + uintptr(fstride*uint64(m))*8))
	yb = *(*kiss_twiddle_cpx)(unsafe.Pointer(twiddles + uintptr(fstride*uint64(2)*uint64(m))*8))
	tw = (*kiss_fft_state)(unsafe.Pointer(st)).Ftwiddles[0]
	i = 0
	for {
		if !(i < N) {
			break
		}
		Fout = Fout_beg + uintptr(i*mm)*8
		Fout0 = Fout
		Fout1 = Fout0 + uintptr(m)*8
		Fout2 = Fout0 + uintptr(int32(2)*m)*8
		Fout3 = Fout0 + uintptr(int32(3)*m)*8
		Fout4 = Fout0 + uintptr(int32(4)*m)*8
		u = 0
		for {
			if !(u < m) {
				break
			}
			scratch[0] = *(*kiss_fft_cpx)(unsafe.Pointer(Fout0))
			scratch[int32(1)].Fr = float32((*kiss_fft_cpx)(unsafe.Pointer(Fout1)).Fr*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(u)*fstride)*8))).Fr) + float32((*kiss_fft_cpx)(unsafe.Pointer(Fout1)).Fi*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(u)*fstride)*8))).Fi)
			scratch[int32(1)].Fi = float32((*kiss_fft_cpx)(unsafe.Pointer(Fout1)).Fi*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(u)*fstride)*8))).Fr) - float32((*kiss_fft_cpx)(unsafe.Pointer(Fout1)).Fr*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(u)*fstride)*8))).Fi)
			scratch[int32(2)].Fr = float32((*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fr*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(int32(2)*u)*fstride)*8))).Fr) + float32((*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fi*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(int32(2)*u)*fstride)*8))).Fi)
			scratch[int32(2)].Fi = float32((*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fi*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(int32(2)*u)*fstride)*8))).Fr) - float32((*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fr*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(int32(2)*u)*fstride)*8))).Fi)
			scratch[int32(3)].Fr = float32((*kiss_fft_cpx)(unsafe.Pointer(Fout3)).Fr*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(int32(3)*u)*fstride)*8))).Fr) + float32((*kiss_fft_cpx)(unsafe.Pointer(Fout3)).Fi*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(int32(3)*u)*fstride)*8))).Fi)
			scratch[int32(3)].Fi = float32((*kiss_fft_cpx)(unsafe.Pointer(Fout3)).Fi*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(int32(3)*u)*fstride)*8))).Fr) - float32((*kiss_fft_cpx)(unsafe.Pointer(Fout3)).Fr*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(int32(3)*u)*fstride)*8))).Fi)
			scratch[int32(4)].Fr = float32((*kiss_fft_cpx)(unsafe.Pointer(Fout4)).Fr*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(int32(4)*u)*fstride)*8))).Fr) + float32((*kiss_fft_cpx)(unsafe.Pointer(Fout4)).Fi*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(int32(4)*u)*fstride)*8))).Fi)
			scratch[int32(4)].Fi = float32((*kiss_fft_cpx)(unsafe.Pointer(Fout4)).Fi*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(int32(4)*u)*fstride)*8))).Fr) - float32((*kiss_fft_cpx)(unsafe.Pointer(Fout4)).Fr*(*(*kiss_twiddle_cpx)(unsafe.Pointer(tw + uintptr(uint64(int32(4)*u)*fstride)*8))).Fi)
			scratch[int32(7)].Fr = scratch[int32(1)].Fr + scratch[int32(4)].Fr
			scratch[int32(7)].Fi = scratch[int32(1)].Fi + scratch[int32(4)].Fi
			scratch[int32(10)].Fr = scratch[int32(1)].Fr - scratch[int32(4)].Fr
			scratch[int32(10)].Fi = scratch[int32(1)].Fi - scratch[int32(4)].Fi
			scratch[int32(8)].Fr = scratch[int32(2)].Fr + scratch[int32(3)].Fr
			scratch[int32(8)].Fi = scratch[int32(2)].Fi + scratch[int32(3)].Fi
			scratch[int32(9)].Fr = scratch[int32(2)].Fr - scratch[int32(3)].Fr
			scratch[int32(9)].Fi = scratch[int32(2)].Fi - scratch[int32(3)].Fi
			*(*float32)(unsafe.Pointer(Fout0)) += scratch[int32(7)].Fr + scratch[int32(8)].Fr
			*(*float32)(unsafe.Pointer(Fout0 + 4)) += scratch[int32(7)].Fi + scratch[int32(8)].Fi
			scratch[int32(5)].Fr = scratch[0].Fr + float32(scratch[int32(7)].Fr*ya.Fr) + float32(scratch[int32(8)].Fr*yb.Fr)
			scratch[int32(5)].Fi = scratch[0].Fi + float32(scratch[int32(7)].Fi*ya.Fr) + float32(scratch[int32(8)].Fi*yb.Fr)
			scratch[int32(6)].Fr = -float32(scratch[int32(10)].Fi*ya.Fi) - float32(scratch[int32(9)].Fi*yb.Fi)
			scratch[int32(6)].Fi = float32(scratch[int32(10)].Fr*ya.Fi) + float32(scratch[int32(9)].Fr*yb.Fi)
			(*kiss_fft_cpx)(unsafe.Pointer(Fout1)).Fr = scratch[int32(5)].Fr - scratch[int32(6)].Fr
			(*kiss_fft_cpx)(unsafe.Pointer(Fout1)).Fi = scratch[int32(5)].Fi - scratch[int32(6)].Fi
			(*kiss_fft_cpx)(unsafe.Pointer(Fout4)).Fr = scratch[int32(5)].Fr + scratch[int32(6)].Fr
			(*kiss_fft_cpx)(unsafe.Pointer(Fout4)).Fi = scratch[int32(5)].Fi + scratch[int32(6)].Fi
			scratch[int32(11)].Fr = scratch[0].Fr + float32(scratch[int32(7)].Fr*yb.Fr) + float32(scratch[int32(8)].Fr*ya.Fr)
			scratch[int32(11)].Fi = scratch[0].Fi + float32(scratch[int32(7)].Fi*yb.Fr) + float32(scratch[int32(8)].Fi*ya.Fr)
			scratch[int32(12)].Fr = float32(scratch[int32(10)].Fi*yb.Fi) - float32(scratch[int32(9)].Fi*ya.Fi)
			scratch[int32(12)].Fi = -float32(scratch[int32(10)].Fr*yb.Fi) + float32(scratch[int32(9)].Fr*ya.Fi)
			(*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fr = scratch[int32(11)].Fr + scratch[int32(12)].Fr
			(*kiss_fft_cpx)(unsafe.Pointer(Fout2)).Fi = scratch[int32(11)].Fi + scratch[int32(12)].Fi
			(*kiss_fft_cpx)(unsafe.Pointer(Fout3)).Fr = scratch[int32(11)].Fr - scratch[int32(12)].Fr
			(*kiss_fft_cpx)(unsafe.Pointer(Fout3)).Fi = scratch[int32(11)].Fi - scratch[int32(12)].Fi
			Fout0 += 8
			Fout1 += 8
			Fout2 += 8
			Fout3 += 8
			Fout4 += 8
			goto _2
		_2:
			;
			u = u + 1
		}
		goto _1
	_1:
		;
		i = i + 1
	}
}

func TestFFTButterflies(t *testing.T) {
	tls := libc.NewTLS()
	defer tls.Close()
	r := rand.New(rand.NewPCG(13, 14))
	type bfly func(*libc.TLS, uintptr, size_t, uintptr, int32, int32, int32)
	cases := []struct {
		radix     int
		got, want bfly
	}{
		{2, kf_bfly2, ref_kf_bfly2}, {2, ki_bfly2, ref_ki_bfly2},
		{3, kf_bfly3, ref_kf_bfly3}, {3, ki_bfly3, ref_ki_bfly3},
		{4, kf_bfly4, ref_kf_bfly4}, {4, ki_bfly4, ref_ki_bfly4},
		{5, kf_bfly5, ref_kf_bfly5}, {5, ki_bfly5, ref_ki_bfly5},
	}
	st := cBuf[kiss_fft_state](1)
	for iter := range 400 {
		c := cases[iter%len(cases)]
		m := 1 + r.IntN(40)
		n := 1 + r.IntN(5)
		mm := c.radix * m
		fstride := 1 + r.IntN(4)
		tw := randFloats(r, 2*(fstride*c.radix*m+1), 1)
		st[0].Ftwiddles[0] = ptr(tw)
		data := randFloats(r, 2*n*mm, 1000)
		want := append([]float32(nil), data...)
		got := append([]float32(nil), data...)
		c.want(tls, ptr(want), size_t(fstride), uintptr(unsafe.Pointer(&st[0])), int32(m), int32(n), int32(mm))
		c.got(tls, ptr(got), size_t(fstride), uintptr(unsafe.Pointer(&st[0])), int32(m), int32(n), int32(mm))
		if i := sameBits(want, got); i >= 0 {
			t.Fatalf("iter %d radix %d m %d n %d: mismatch at %d: got %v want %v", iter, c.radix, m, n, i, got[i], want[i])
		}
	}
}

func ref_comb_filter(tls *libc.TLS, y uintptr, x uintptr, T0 int32, T1 int32, N int32, g0 opus_val16, g1 opus_val16, tapset0 int32, tapset1 int32, window uintptr, overlap int32) {
	var f, g00, g01, g02, g10, g11, g12 opus_val16
	var i int32
	_, _, _, _, _, _, _, _ = f, g00, g01, g02, g10, g11, g12, i
	g00 = opus_val16(g0 * *(*opus_val16)(unsafe.Pointer(uintptr(unsafe.Pointer(&gains)) + uintptr(tapset0)*12)))
	g01 = opus_val16(g0 * *(*opus_val16)(unsafe.Pointer(uintptr(unsafe.Pointer(&gains)) + uintptr(tapset0)*12 + 1*4)))
	g02 = opus_val16(g0 * *(*opus_val16)(unsafe.Pointer(uintptr(unsafe.Pointer(&gains)) + uintptr(tapset0)*12 + 2*4)))
	g10 = opus_val16(g1 * *(*opus_val16)(unsafe.Pointer(uintptr(unsafe.Pointer(&gains)) + uintptr(tapset1)*12)))
	g11 = opus_val16(g1 * *(*opus_val16)(unsafe.Pointer(uintptr(unsafe.Pointer(&gains)) + uintptr(tapset1)*12 + 1*4)))
	g12 = opus_val16(g1 * *(*opus_val16)(unsafe.Pointer(uintptr(unsafe.Pointer(&gains)) + uintptr(tapset1)*12 + 2*4)))
	i = 0
	for {
		if !(i < overlap) {
			break
		}
		f = opus_val16(*(*opus_val16)(unsafe.Pointer(window + uintptr(i)*4)) * *(*opus_val16)(unsafe.Pointer(window + uintptr(i)*4)))
		*(*opus_val32)(unsafe.Pointer(y + uintptr(i)*4)) = *(*opus_val32)(unsafe.Pointer(x + uintptr(i)*4)) + float32(float32((libc.Float32FromFloat32(1)-f)*g00)**(*opus_val32)(unsafe.Pointer(x + uintptr(i-T0)*4))) + float32(float32((libc.Float32FromFloat32(1)-f)*g01)**(*opus_val32)(unsafe.Pointer(x + uintptr(i-T0-int32(1))*4))) + float32(float32((libc.Float32FromFloat32(1)-f)*g01)**(*opus_val32)(unsafe.Pointer(x + uintptr(i-T0+int32(1))*4))) + float32(float32((libc.Float32FromFloat32(1)-f)*g02)**(*opus_val32)(unsafe.Pointer(x + uintptr(i-T0-int32(2))*4))) + float32(float32((libc.Float32FromFloat32(1)-f)*g02)**(*opus_val32)(unsafe.Pointer(x + uintptr(i-T0+int32(2))*4))) + opus_val16(opus_val16(f*g10)**(*opus_val32)(unsafe.Pointer(x + uintptr(i-T1)*4))) + opus_val16(opus_val16(f*g11)**(*opus_val32)(unsafe.Pointer(x + uintptr(i-T1-int32(1))*4))) + opus_val16(opus_val16(f*g11)**(*opus_val32)(unsafe.Pointer(x + uintptr(i-T1+int32(1))*4))) + opus_val16(opus_val16(f*g12)**(*opus_val32)(unsafe.Pointer(x + uintptr(i-T1-int32(2))*4))) + opus_val16(opus_val16(f*g12)**(*opus_val32)(unsafe.Pointer(x + uintptr(i-T1+int32(2))*4)))
		goto _1
	_1:
		;
		i = i + 1
	}
	i = overlap
	for {
		if !(i < N) {
			break
		}
		*(*opus_val32)(unsafe.Pointer(y + uintptr(i)*4)) = *(*opus_val32)(unsafe.Pointer(x + uintptr(i)*4)) + opus_val16(g10**(*opus_val32)(unsafe.Pointer(x + uintptr(i-T1)*4))) + opus_val16(g11**(*opus_val32)(unsafe.Pointer(x + uintptr(i-T1-int32(1))*4))) + opus_val16(g11**(*opus_val32)(unsafe.Pointer(x + uintptr(i-T1+int32(1))*4))) + opus_val16(g12**(*opus_val32)(unsafe.Pointer(x + uintptr(i-T1-int32(2))*4))) + opus_val16(g12**(*opus_val32)(unsafe.Pointer(x + uintptr(i-T1+int32(2))*4)))
		goto _2
	_2:
		;
		i = i + 1
	}
}

func TestCombFilter(t *testing.T) {
	tls := libc.NewTLS()
	defer tls.Close()
	r := rand.New(rand.NewPCG(15, 16))
	const hist = 1030
	for iter := range 500 {
		n := 1 + r.IntN(960)
		overlap := min(r.IntN(121), n)
		t0, t1 := int32(r.IntN(1024)), int32(r.IntN(1024))
		if iter%5 == 0 {
			t0, t1 = int32(r.IntN(7)), int32(r.IntN(7))
		}
		g0, g1 := r.Float32()*0.8, r.Float32()*0.8
		tap0, tap1 := int32(r.IntN(3)), int32(r.IntN(3))
		win := randFloats(r, overlap+1, 1)
		base := randFloats(r, hist+n+4, 1000) // taps may read up to x[n+1]
		inPlace := r.IntN(2) == 0
		wx, gx := append([]float32(nil), base...), append([]float32(nil), base...)
		wy, gy := wx, gx
		if !inPlace {
			wy, gy = make([]float32, hist+n+4), make([]float32, hist+n+4)
		}
		ref_comb_filter(tls, ptr(wy[hist:]), ptr(wx[hist:]), t0, t1, int32(n), g0, g1, tap0, tap1, ptr(win), int32(overlap))
		comb_filter(tls, ptr(gy[hist:]), ptr(gx[hist:]), t0, t1, int32(n), g0, g1, tap0, tap1, ptr(win), int32(overlap))
		if i := sameBits(wy, gy); i >= 0 {
			t.Fatalf("iter %d T0 %d T1 %d n %d overlap %d inPlace %v: mismatch at %d", iter, t0, t1, n, overlap, inPlace, i-hist)
		}
	}
}

func ref_clt_mdct_backward(tls *libc.TLS, l uintptr, in uintptr, out uintptr, window uintptr, overlap int32, shift int32, stride int32) {
	var N, N2, N4, i int32
	var f, f2, fp, fp1, fp11, fp2, fp21, t, t1, wp1, wp11, wp2, wp21, xp1, xp11, xp2, xp21, yp, yp1, yp11, yp2, v2 uintptr
	var im, re, sine, x1, x2, yi, yi1, yr, yr1 float32
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = N, N2, N4, f, f2, fp, fp1, fp11, fp2, fp21, i, im, re, sine, t, t1, wp1, wp11, wp2, wp21, x1, x2, xp1, xp11, xp2, xp21, yi, yi1, yp, yp1, yp11, yp2, yr, yr1, v2
	_sp := _arenaSave()
	defer _arenaRestore(_sp)
	N = (*mdct_lookup)(unsafe.Pointer(l)).Fn
	N = N >> shift
	N2 = N >> int32(1)
	N4 = N >> int32(2)
	f = _arenaAlloc(uint64(4) * uint64(N2))
	f2 = _arenaAlloc(uint64(4) * uint64(N2))
	/* sin(x) ~= x here */
	sine = float32(float32(float32(2)*libc.Float32FromFloat32(3.141592653))*libc.Float32FromFloat32(0.125)) / float32(N)
	/* Pre-rotate */
	/* Temp pointers to make it really clear to the compiler what we're doing */
	xp1 = in
	xp2 = in + uintptr(stride*(N2-int32(1)))*4
	yp = f2
	t = (*mdct_lookup)(unsafe.Pointer(l)).Ftrig[0]
	i = 0
	for {
		if !(i < N4) {
			break
		}
		yr = -float32(*(*float32)(unsafe.Pointer(xp2))**(*float32)(unsafe.Pointer(t + uintptr(i<<shift)*4))) + float32(*(*float32)(unsafe.Pointer(xp1))**(*float32)(unsafe.Pointer(t + uintptr((N4-i)<<shift)*4)))
		yi = -float32(*(*float32)(unsafe.Pointer(xp2))**(*float32)(unsafe.Pointer(t + uintptr((N4-i)<<shift)*4))) - float32(*(*float32)(unsafe.Pointer(xp1))**(*float32)(unsafe.Pointer(t + uintptr(i<<shift)*4)))
		/* works because the cos is nearly one */
		v2 = yp
		yp += 4
		*(*float32)(unsafe.Pointer(v2)) = yr - float32(yi*sine)
		v2 = yp
		yp += 4
		*(*float32)(unsafe.Pointer(v2)) = yi + float32(yr*sine)
		xp1 = xp1 + uintptr(int32(2)*stride)*4
		xp2 = xp2 - uintptr(int32(2)*stride)*4
		goto _1
	_1:
		;
		i = i + 1
	}
	/* Inverse N/4 complex FFT. This one should *not* downscale even in fixed-point */
	opus_ifft(tls, *(*uintptr)(unsafe.Pointer(l + 8 + uintptr(shift)*8)), f2, f)
	/* Post-rotate */
	fp = f
	t1 = (*mdct_lookup)(unsafe.Pointer(l)).Ftrig[0]
	i = 0
	for {
		if !(i < N4) {
			break
		}
		re = *(*float32)(unsafe.Pointer(fp))
		im = *(*float32)(unsafe.Pointer(fp + 1*4))
		/* We'd scale up by 2 here, but instead it's done when mixing the windows */
		yr1 = float32(re**(*float32)(unsafe.Pointer(t1 + uintptr(i<<shift)*4))) - float32(im**(*float32)(unsafe.Pointer(t1 + uintptr((N4-i)<<shift)*4)))
		yi1 = float32(im**(*float32)(unsafe.Pointer(t1 + uintptr(i<<shift)*4))) + float32(re**(*float32)(unsafe.Pointer(t1 + uintptr((N4-i)<<shift)*4)))
		/* works because the cos is nearly one */
		v2 = fp
		fp += 4
		*(*float32)(unsafe.Pointer(v2)) = yr1 - float32(yi1*sine)
		v2 = fp
		fp += 4
		*(*float32)(unsafe.Pointer(v2)) = yi1 + float32(yr1*sine)
		goto _4
	_4:
		;
		i = i + 1
	}
	/* De-shuffle the components for the middle of the window only */
	fp1 = f
	fp2 = f + uintptr(N2)*4 - uintptr(1)*4
	yp1 = f2
	i = 0
	for {
		if !(i < N4) {
			break
		}
		v2 = yp1
		yp1 += 4
		*(*float32)(unsafe.Pointer(v2)) = -*(*float32)(unsafe.Pointer(fp1))
		v2 = yp1
		yp1 += 4
		*(*float32)(unsafe.Pointer(v2)) = *(*float32)(unsafe.Pointer(fp2))
		fp1 = fp1 + uintptr(2)*4
		fp2 = fp2 - uintptr(2)*4
		goto _7
	_7:
		;
		i = i + 1
	}
	out = out - uintptr((N2-overlap)>>int32(1))*4
	/* Mirror on both sides for TDAC */
	fp11 = f2 + uintptr(N4)*4 - uintptr(1)*4
	xp11 = out + uintptr(N2)*4 - uintptr(1)*4
	yp11 = out + uintptr(N4)*4 - uintptr(overlap/int32(2))*4
	wp1 = window
	wp2 = window + uintptr(overlap)*4 - uintptr(1)*4
	i = 0
	for {
		if !(i < N4-overlap/int32(2)) {
			break
		}
		*(*float32)(unsafe.Pointer(xp11)) = *(*float32)(unsafe.Pointer(fp11))
		xp11 -= 4
		fp11 -= 4
		goto _10
	_10:
		;
		i = i + 1
	}
	for {
		if !(i < N4) {
			break
		}
		v2 = fp11
		fp11 -= 4
		x1 = *(*float32)(unsafe.Pointer(v2))
		v2 = yp11
		yp11 += 4
		*(*float32)(unsafe.Pointer(v2)) += -opus_val16(*(*opus_val16)(unsafe.Pointer(wp1)) * x1)
		v2 = xp11
		xp11 -= 4
		*(*float32)(unsafe.Pointer(v2)) += opus_val16(*(*opus_val16)(unsafe.Pointer(wp2)) * x1)
		wp1 += 4
		wp2 -= 4
		goto _11
	_11:
		;
		i = i + 1
	}
	fp21 = f2 + uintptr(N4)*4
	xp21 = out + uintptr(N2)*4
	yp2 = out + uintptr(N)*4 - uintptr(1)*4 - uintptr(N4-overlap/int32(2))*4
	wp11 = window
	wp21 = window + uintptr(overlap)*4 - uintptr(1)*4
	i = 0
	for {
		if !(i < N4-overlap/int32(2)) {
			break
		}
		*(*float32)(unsafe.Pointer(xp21)) = *(*float32)(unsafe.Pointer(fp21))
		xp21 += 4
		fp21 += 4
		goto _15
	_15:
		;
		i = i + 1
	}
	for {
		if !(i < N4) {
			break
		}
		v2 = fp21
		fp21 += 4
		x2 = *(*float32)(unsafe.Pointer(v2))
		v2 = yp2
		yp2 -= 4
		*(*float32)(unsafe.Pointer(v2)) = opus_val16(*(*opus_val16)(unsafe.Pointer(wp11)) * x2)
		v2 = xp21
		xp21 += 4
		*(*float32)(unsafe.Pointer(v2)) = opus_val16(*(*opus_val16)(unsafe.Pointer(wp21)) * x2)
		wp11 += 4
		wp21 -= 4
		goto _16
	_16:
		;
		i = i + 1
	}
}

func TestMdctBackward(t *testing.T) {
	tls := libc.NewTLS()
	defer tls.Close()
	r := rand.New(rand.NewPCG(17, 18))
	mode := uintptr(unsafe.Pointer(&mode48000_960_120))
	m := &mode48000_960_120
	for iter := range 200 {
		shift := int32(r.IntN(4))
		stride := int32(1) << r.IntN(4)
		n := int32(1920) >> shift
		in := randFloats(r, int(n/2*stride), 3000)
		pad := int(n) // output spans out-(N/2-overlap)/2 .. out+N
		base := randFloats(r, int(n)+2*pad, 1000)
		want := append([]float32(nil), base...)
		got := append([]float32(nil), base...)
		ref_clt_mdct_backward(tls, mode+80, ptr(in), ptr(want[pad:]), m.Fwindow, m.Foverlap, shift, stride)
		clt_mdct_backward(tls, mode+80, ptr(in), ptr(got[pad:]), m.Fwindow, m.Foverlap, shift, stride)
		if i := sameBits(want, got); i >= 0 {
			t.Fatalf("iter %d shift %d stride %d: mismatch at %d", iter, shift, stride, i)
		}
	}
}

func ref_deinterleave_hadamard(tls *libc.TLS, X uintptr, N0 int32, stride int32, hadamard int32) {
	var N, i, j int32
	var ordery, tmp uintptr
	_, _, _, _, _ = N, i, j, ordery, tmp
	_sp := _arenaSave()
	defer _arenaRestore(_sp)
	N = N0 * stride
	tmp = _arenaAlloc(uint64(4) * uint64(N))
	if hadamard != 0 {
		ordery = uintptr(unsafe.Pointer(&ordery_table)) + uintptr(stride)*4 - uintptr(2)*4
		i = 0
		for {
			if !(i < stride) {
				break
			}
			j = 0
			for {
				if !(j < N0) {
					break
				}
				*(*celt_norm)(unsafe.Pointer(tmp + uintptr(*(*int32)(unsafe.Pointer(ordery + uintptr(i)*4))*N0+j)*4)) = *(*celt_norm)(unsafe.Pointer(X + uintptr(j*stride+i)*4))
				goto _2
			_2:
				;
				j = j + 1
			}
			goto _1
		_1:
			;
			i = i + 1
		}
	} else {
		i = 0
		for {
			if !(i < stride) {
				break
			}
			j = 0
			for {
				if !(j < N0) {
					break
				}
				*(*celt_norm)(unsafe.Pointer(tmp + uintptr(i*N0+j)*4)) = *(*celt_norm)(unsafe.Pointer(X + uintptr(j*stride+i)*4))
				goto _4
			_4:
				;
				j = j + 1
			}
			goto _3
		_3:
			;
			i = i + 1
		}
	}
	j = 0
	for {
		if !(j < N) {
			break
		}
		*(*celt_norm)(unsafe.Pointer(X + uintptr(j)*4)) = *(*celt_norm)(unsafe.Pointer(tmp + uintptr(j)*4))
		goto _5
	_5:
		;
		j = j + 1
	}
}

func ref_interleave_hadamard(tls *libc.TLS, X uintptr, N0 int32, stride int32, hadamard int32) {
	var N, i, j int32
	var ordery, tmp uintptr
	_, _, _, _, _ = N, i, j, ordery, tmp
	_sp := _arenaSave()
	defer _arenaRestore(_sp)
	N = N0 * stride
	tmp = _arenaAlloc(uint64(4) * uint64(N))
	if hadamard != 0 {
		ordery = uintptr(unsafe.Pointer(&ordery_table)) + uintptr(stride)*4 - uintptr(2)*4
		i = 0
		for {
			if !(i < stride) {
				break
			}
			j = 0
			for {
				if !(j < N0) {
					break
				}
				*(*celt_norm)(unsafe.Pointer(tmp + uintptr(j*stride+i)*4)) = *(*celt_norm)(unsafe.Pointer(X + uintptr(*(*int32)(unsafe.Pointer(ordery + uintptr(i)*4))*N0+j)*4))
				goto _2
			_2:
				;
				j = j + 1
			}
			goto _1
		_1:
			;
			i = i + 1
		}
	} else {
		i = 0
		for {
			if !(i < stride) {
				break
			}
			j = 0
			for {
				if !(j < N0) {
					break
				}
				*(*celt_norm)(unsafe.Pointer(tmp + uintptr(j*stride+i)*4)) = *(*celt_norm)(unsafe.Pointer(X + uintptr(i*N0+j)*4))
				goto _4
			_4:
				;
				j = j + 1
			}
			goto _3
		_3:
			;
			i = i + 1
		}
	}
	j = 0
	for {
		if !(j < N) {
			break
		}
		*(*celt_norm)(unsafe.Pointer(X + uintptr(j)*4)) = *(*celt_norm)(unsafe.Pointer(tmp + uintptr(j)*4))
		goto _5
	_5:
		;
		j = j + 1
	}
}

func TestHadamardInterleave(t *testing.T) {
	tls := libc.NewTLS()
	defer tls.Close()
	r := rand.New(rand.NewPCG(19, 20))
	for iter := range 200 {
		stride := int32(2) << r.IntN(4)
		n0 := int32(1 + r.IntN(40))
		had := int32(r.IntN(2))
		x := randFloats(r, int(n0*stride), 1)
		for k, pair := range [][2]func(*libc.TLS, uintptr, int32, int32, int32){
			{ref_deinterleave_hadamard, deinterleave_hadamard},
			{ref_interleave_hadamard, interleave_hadamard},
		} {
			want := append([]float32(nil), x...)
			got := append([]float32(nil), x...)
			pair[0](tls, ptr(want), n0, stride, had)
			pair[1](tls, ptr(got), n0, stride, had)
			if i := sameBits(want, got); i >= 0 {
				t.Fatalf("iter %d fn %d stride %d n0 %d hadamard %d: mismatch at %d", iter, k, stride, n0, had, i)
			}
		}
	}
}

func ref_uprev(tls *libc.TLS, _ui uintptr, _n uint32, _ui0 opus_uint32) {
	var j, v1 uint32
	var ui1 opus_uint32
	_, _, _ = j, ui1, v1
	/*This do-while will overrun the array if we don't have storage for at least
	  2 values.*/
	j = uint32(1)
	for {
		ui1 = *(*opus_uint32)(unsafe.Pointer(_ui + uintptr(j)*4)) - *(*opus_uint32)(unsafe.Pointer(_ui + uintptr(j-uint32(1))*4)) - _ui0
		*(*opus_uint32)(unsafe.Pointer(_ui + uintptr(j-uint32(1))*4)) = _ui0
		_ui0 = ui1
		goto _2
	_2:
		;
		j = j + 1
		v1 = j
		if !(v1 < _n) {
			break
		}
	}
	*(*opus_uint32)(unsafe.Pointer(_ui + uintptr(j-uint32(1))*4)) = _ui0
}

func ref_unext(tls *libc.TLS, _ui uintptr, _len uint32, _ui0 opus_uint32) {
	var j, v1 uint32
	var ui1 opus_uint32
	_, _, _ = j, ui1, v1
	/*This do-while will overrun the array if we don't have storage for at least
	  2 values.*/
	j = uint32(1)
	for {
		ui1 = *(*opus_uint32)(unsafe.Pointer(_ui + uintptr(j)*4)) + *(*opus_uint32)(unsafe.Pointer(_ui + uintptr(j-uint32(1))*4)) + _ui0
		*(*opus_uint32)(unsafe.Pointer(_ui + uintptr(j-uint32(1))*4)) = _ui0
		_ui0 = ui1
		goto _2
	_2:
		;
		j = j + 1
		v1 = j
		if !(v1 < _len) {
			break
		}
	}
	*(*opus_uint32)(unsafe.Pointer(_ui + uintptr(j-uint32(1))*4)) = _ui0
}

func ref_ncwrs_urow(tls *libc.TLS, _n uint32, _k uint32, _u uintptr) (r opus_uint32) {
	var k, len1, v2 uint32
	var n2m1, um1, um2, v1, v5 opus_uint32
	_, _, _, _, _, _, _, _ = k, len1, n2m1, um1, um2, v1, v2, v5
	len1 = _k + uint32(2)
	/*We require storage at least 3 values (e.g., _k>0).*/
	*(*opus_uint32)(unsafe.Pointer(_u)) = uint32(0)
	v1 = uint32(1)
	um2 = v1
	*(*opus_uint32)(unsafe.Pointer(_u + 1*4)) = v1
	/*_k>52 doesn't work in the false branch due to the limits of INV_TABLE,
	  but _k isn't tested here because k<=52 for n=7*/
	if _n <= uint32(6) {
		/*If _n==0, _u[0] should be 1 and the rest should be 0.*/
		/*If _n==1, _u[i] should be 1 for i>1.*/
		/*If _k==0, the following do-while loop will overflow the buffer.*/
		k = uint32(2)
		for {
			*(*opus_uint32)(unsafe.Pointer(_u + uintptr(k)*4)) = k<<int32(1) - uint32(1)
			goto _3
		_3:
			;
			k = k + 1
			v2 = k
			if !(v2 < len1) {
				break
			}
		}
		k = uint32(2)
		for {
			if !(k < _n) {
				break
			}
			ref_unext(tls, _u+uintptr(1)*4, _k+uint32(1), uint32(1))
			goto _4
		_4:
			;
			k = k + 1
		}
	} else {
		v5 = _n<<int32(1) - uint32(1)
		um1 = v5
		v1 = v5
		n2m1 = v1
		*(*opus_uint32)(unsafe.Pointer(_u + 2*4)) = v1
		k = uint32(3)
		for {
			if !(k < len1) {
				break
			}
			/*U(N,K) = ((2*N-1)*U(N,K-1)-U(N,K-2))/(K-1) + U(N,K-2)*/
			v1 = imusdiv32even(tls, n2m1, um1, um2, int32(k-uint32(1))) + um2
			um2 = v1
			*(*opus_uint32)(unsafe.Pointer(_u + uintptr(k)*4)) = v1
			k = k + 1
			v2 = k
			if v2 >= len1 {
				break
			}
			v1 = imusdiv32odd(tls, n2m1, um2, um1, int32((k-uint32(1))>>int32(1))) + um1
			um1 = v1
			*(*opus_uint32)(unsafe.Pointer(_u + uintptr(k)*4)) = v1
			goto _7
		_7:
			;
			k = k + 1
		}
	}
	return *(*opus_uint32)(unsafe.Pointer(_u + uintptr(_k)*4)) + *(*opus_uint32)(unsafe.Pointer(_u + uintptr(_k+uint32(1))*4))
}

func ref_cwrsi(tls *libc.TLS, _n int32, _k int32, _i opus_uint32, _y uintptr, _u uintptr) {
	var j, s, yj, v1 int32
	var p opus_uint32
	_, _, _, _, _ = j, p, s, yj, v1
	j = 0
	for {
		p = *(*opus_uint32)(unsafe.Pointer(_u + uintptr(_k+int32(1))*4))
		s = -libc.BoolInt32(_i >= p)
		_i = _i - p&uint32(s)
		yj = _k
		p = *(*opus_uint32)(unsafe.Pointer(_u + uintptr(_k)*4))
		for p > _i {
			_k = _k - 1
			v1 = _k
			p = *(*opus_uint32)(unsafe.Pointer(_u + uintptr(v1)*4))
		}
		_i = _i - p
		yj = yj - _k
		*(*int32)(unsafe.Pointer(_y + uintptr(j)*4)) = yj + s ^ s
		ref_uprev(tls, _u, uint32(_k+int32(2)), uint32(0))
		goto _2
	_2:
		;
		j = j + 1
		v1 = j
		if !(v1 < _n) {
			break
		}
	}
}

func ref_icwrs(tls *libc.TLS, _n int32, _k int32, _nc uintptr, _y uintptr, _u uintptr) (r opus_uint32) {
	bp := tls.Alloc(16)
	defer tls.Free(16)
	var i opus_uint32
	var j, v2 int32
	var _ /* k at bp+0 */ int32
	_, _, _ = i, j, v2
	/*We can't unroll the first two iterations of the loop unless _n>=2.*/
	*(*opus_uint32)(unsafe.Pointer(_u)) = uint32(0)
	*(*int32)(unsafe.Pointer(bp)) = int32(1)
	for {
		if !(*(*int32)(unsafe.Pointer(bp)) <= _k+int32(1)) {
			break
		}
		*(*opus_uint32)(unsafe.Pointer(_u + uintptr(*(*int32)(unsafe.Pointer(bp)))*4)) = uint32(*(*int32)(unsafe.Pointer(bp))<<int32(1) - int32(1))
		goto _1
	_1:
		;
		*(*int32)(unsafe.Pointer(bp)) = *(*int32)(unsafe.Pointer(bp)) + 1
	}
	i = icwrs1(tls, _y+uintptr(_n)*4-uintptr(1)*4, bp)
	j = _n - int32(2)
	i = i + *(*opus_uint32)(unsafe.Pointer(_u + uintptr(*(*int32)(unsafe.Pointer(bp)))*4))
	*(*int32)(unsafe.Pointer(bp)) = *(*int32)(unsafe.Pointer(bp)) + libc.Xabs(tls, *(*int32)(unsafe.Pointer(_y + uintptr(j)*4)))
	if *(*int32)(unsafe.Pointer(_y + uintptr(j)*4)) < 0 {
		i = i + *(*opus_uint32)(unsafe.Pointer(_u + uintptr(*(*int32)(unsafe.Pointer(bp))+int32(1))*4))
	}
	for {
		v2 = j
		j = j - 1
		if !(v2 > 0) {
			break
		}
		ref_unext(tls, _u, uint32(_k+int32(2)), uint32(0))
		i = i + *(*opus_uint32)(unsafe.Pointer(_u + uintptr(*(*int32)(unsafe.Pointer(bp)))*4))
		*(*int32)(unsafe.Pointer(bp)) = *(*int32)(unsafe.Pointer(bp)) + libc.Xabs(tls, *(*int32)(unsafe.Pointer(_y + uintptr(j)*4)))
		if *(*int32)(unsafe.Pointer(_y + uintptr(j)*4)) < 0 {
			i = i + *(*opus_uint32)(unsafe.Pointer(_u + uintptr(*(*int32)(unsafe.Pointer(bp))+int32(1))*4))
		}
	}
	*(*opus_uint32)(unsafe.Pointer(_nc)) = *(*opus_uint32)(unsafe.Pointer(_u + uintptr(*(*int32)(unsafe.Pointer(bp)))*4)) + *(*opus_uint32)(unsafe.Pointer(_u + uintptr(*(*int32)(unsafe.Pointer(bp))+int32(1))*4))
	return i
}

func TestCWRSTable(t *testing.T) {
	tls := libc.NewTLS()
	defer tls.Close()
	r := rand.New(rand.NewPCG(21, 22))
	u := cBuf[uint32](cwrsCols + 2)
	y := cBuf[int32](cwrsMaxN)
	y2 := cBuf[int32](cwrsMaxN)
	nc := cBuf[uint32](1)
	checked := 0
	for n := int32(5); n <= cwrsMaxN; n++ {
		// Exact U(n, 0..K+1) in uint64 to find the codebooks that fit in 32 bits.
		exact := make([][]uint64, n+1)
		for m := range exact {
			exact[m] = make([]uint64, cwrsCols)
		}
		exact[0][0] = 1
		for m := 1; m <= int(n); m++ {
			for k := 1; k < cwrsCols; k++ {
				exact[m][k] = min(exact[m-1][k]+exact[m][k-1]+exact[m-1][k-1], 1<<40)
			}
		}
		for k := int32(1); k+1 < cwrsCols; k++ {
			v := exact[n][k] + exact[n][k+1]
			if v >= 1<<32 {
				break
			}
			want := ref_ncwrs_urow(tls, uint32(n), uint32(k), ptr32(u))
			row := cwrsRow(n)
			if got := row[k] + row[k+1]; got != want || uint64(want) != v {
				t.Fatalf("V(%d,%d): table %d, ncwrs_urow %d, exact %d", n, k, got, want, v)
			}
			for s := range 6 {
				idx := uint32(r.Uint64N(v))
				if s == 0 {
					idx = 0
				} else if s == 1 {
					idx = uint32(v - 1)
				}
				ref_ncwrs_urow(tls, uint32(n), uint32(k), ptr32(u))
				ref_cwrsi(tls, n, k, idx, ptr32i(y), ptr32(u))
				cwrsiTable(n, k, idx, ptr32i(y2))
				for j := range n {
					if y[j] != y2[j] {
						t.Fatalf("cwrsi(%d,%d,%d): y[%d] got %d want %d", n, k, idx, j, y2[j], y[j])
					}
				}
				wi := ref_icwrs(tls, n, k, ptr32(nc), ptr32i(y), ptr32(u))
				gi, gnc := icwrsTable(n, ptr32i(y))
				if wi != gi || nc[0] != gnc || wi != idx {
					t.Fatalf("icwrs(%d,%d): got (%d,%d) want (%d,%d), index %d", n, k, gi, gnc, wi, nc[0], idx)
				}
				checked++
			}
		}
	}
	t.Logf("checked %d codewords", checked)
}

func ptr32(v []uint32) uintptr { return uintptr(unsafe.Pointer(&v[0])) }
func ptr32i(v []int32) uintptr { return uintptr(unsafe.Pointer(&v[0])) }

func ref_transient_analysis(tls *libc.TLS, in uintptr, len1 int32, C int32, overlap int32) (r int32) {
	var N, block, conseq, i, is_transient, j, j1 int32
	var bins, tmp uintptr
	var max_abs, t1, t2, t3, v7, v8, v9 opus_val16
	var mem0, mem1, x, y opus_val32
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = N, bins, block, conseq, i, is_transient, j, j1, max_abs, mem0, mem1, t1, t2, t3, tmp, x, y, v7, v8, v9
	_sp := _arenaSave()
	defer _arenaRestore(_sp)
	mem0 = float32(0)
	mem1 = float32(0)
	is_transient = 0
	tmp = _arenaAlloc(uint64(4) * uint64(len1))
	block = overlap / int32(2)
	N = len1 / block
	bins = _arenaAlloc(uint64(4) * uint64(N))
	if C == int32(1) {
		i = 0
		for {
			if !(i < len1) {
				break
			}
			*(*opus_val16)(unsafe.Pointer(tmp + uintptr(i)*4)) = *(*opus_val32)(unsafe.Pointer(in + uintptr(i)*4))
			goto _1
		_1:
			;
			i = i + 1
		}
	} else {
		i = 0
		for {
			if !(i < len1) {
				break
			}
			*(*opus_val16)(unsafe.Pointer(tmp + uintptr(i)*4)) = *(*opus_val32)(unsafe.Pointer(in + uintptr(i)*4)) + *(*opus_val32)(unsafe.Pointer(in + uintptr(i+len1)*4))
			goto _2
		_2:
			;
			i = i + 1
		}
	}
	/* High-pass filter: (1 - 2*z^-1 + z^-2) / (1 - z^-1 + .5*z^-2) */
	i = 0
	for {
		if !(i < len1) {
			break
		}
		x = *(*opus_val16)(unsafe.Pointer(tmp + uintptr(i)*4))
		y = mem0 + x
		mem0 = mem1 + y - opus_val32(float32(2)*x)
		mem1 = x - opus_val32(libc.Float32FromFloat32(0.5)*y)
		*(*opus_val16)(unsafe.Pointer(tmp + uintptr(i)*4)) = y
		goto _3
	_3:
		;
		i = i + 1
	}
	/* First few samples are bad because we don't propagate the memory */
	i = 0
	for {
		if !(i < int32(12)) {
			break
		}
		*(*opus_val16)(unsafe.Pointer(tmp + uintptr(i)*4)) = float32(0)
		goto _4
	_4:
		;
		i = i + 1
	}
	i = 0
	for {
		if !(i < N) {
			break
		}
		max_abs = float32(0)
		j = 0
		for {
			if !(j < block) {
				break
			}
			if *(*opus_val16)(unsafe.Pointer(tmp + uintptr(i*block+j)*4)) < float32(0) {
				v8 = -*(*opus_val16)(unsafe.Pointer(tmp + uintptr(i*block+j)*4))
			} else {
				v8 = *(*opus_val16)(unsafe.Pointer(tmp + uintptr(i*block+j)*4))
			}
			if max_abs > v8 {
				v7 = max_abs
			} else {
				if *(*opus_val16)(unsafe.Pointer(tmp + uintptr(i*block+j)*4)) < float32(0) {
					v9 = -*(*opus_val16)(unsafe.Pointer(tmp + uintptr(i*block+j)*4))
				} else {
					v9 = *(*opus_val16)(unsafe.Pointer(tmp + uintptr(i*block+j)*4))
				}
				v7 = v9
			}
			max_abs = v7
			goto _6
		_6:
			;
			j = j + 1
		}
		*(*opus_val16)(unsafe.Pointer(bins + uintptr(i)*4)) = max_abs
		goto _5
	_5:
		;
		i = i + 1
	}
	i = 0
	for {
		if !(i < N) {
			break
		}
		conseq = 0
		t1 = float32(libc.Float32FromFloat32(0.15) * *(*opus_val16)(unsafe.Pointer(bins + uintptr(i)*4)))
		t2 = float32(libc.Float32FromFloat32(0.4) * *(*opus_val16)(unsafe.Pointer(bins + uintptr(i)*4)))
		t3 = float32(libc.Float32FromFloat32(0.15) * *(*opus_val16)(unsafe.Pointer(bins + uintptr(i)*4)))
		j1 = 0
		for {
			if !(j1 < i) {
				break
			}
			if *(*opus_val16)(unsafe.Pointer(bins + uintptr(j1)*4)) < t1 {
				conseq = conseq + 1
			}
			if *(*opus_val16)(unsafe.Pointer(bins + uintptr(j1)*4)) < t2 {
				conseq = conseq + 1
			} else {
				conseq = 0
			}
			goto _11
		_11:
			;
			j1 = j1 + 1
		}
		if conseq >= int32(3) {
			is_transient = int32(1)
		}
		conseq = 0
		j1 = i + int32(1)
		for {
			if !(j1 < N) {
				break
			}
			if *(*opus_val16)(unsafe.Pointer(bins + uintptr(j1)*4)) < t3 {
				conseq = conseq + 1
			} else {
				conseq = 0
			}
			goto _12
		_12:
			;
			j1 = j1 + 1
		}
		if conseq >= int32(7) {
			is_transient = int32(1)
		}
		goto _10
	_10:
		;
		i = i + 1
	}
	return is_transient
}

func TestTransientAnalysis(t *testing.T) {
	tls := libc.NewTLS()
	defer tls.Close()
	r := rand.New(rand.NewPCG(23, 24))
	transients := 0
	for iter := range 400 {
		C := int32(1 + iter%2)
		overlap := int32(120)
		len1 := int32(120 + 120*r.IntN(16))
		in := randFloats(r, int(C*len1), 100)
		// Occasional bursts make some frames transient.
		for k := 0; k < r.IntN(3); k++ {
			at := r.IntN(len(in) - 30)
			for j := range 30 {
				in[at+j] *= 300
			}
		}
		if iter%50 == 0 {
			in[r.IntN(len(in))] = float32(math.NaN())
		}
		want := ref_transient_analysis(tls, ptr(in), len1, C, overlap)
		got := transient_analysis(tls, ptr(in), len1, C, overlap)
		if want != got {
			t.Fatalf("iter %d C %d len %d: got %d want %d", iter, C, len1, got, want)
		}
		transients += int(want)
	}
	if transients == 0 || transients == 400 {
		t.Fatalf("test signal never/always transient (%d)", transients)
	}
}
