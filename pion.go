// Copyright (c) 2026 Sergey Krashevich.
// SPDX-License-Identifier: Apache-2.0

package opus

import "errors"

// Bandwidth identifies the encoded audio bandwidth, independently of the PCM output rate.
type Bandwidth byte

const (
	BandwidthAuto Bandwidth = iota
	BandwidthNarrowband
	BandwidthMediumband
	BandwidthWideband
	BandwidthSuperwideband
	BandwidthFullband
)

func (b Bandwidth) String() string {
	names := [...]string{"Auto", "Narrowband", "Mediumband", "Wideband", "Superwideband", "Fullband"}
	if int(b) < len(names) {
		return names[b]
	}
	return "Invalid Bandwidth"
}

// SampleRate returns the effective sample rate, or zero for Auto or an invalid bandwidth.
func (b Bandwidth) SampleRate() int {
	rates := [...]int{0, 8000, 12000, 16000, 24000, 48000}
	if int(b) < len(rates) {
		return rates[b]
	}
	return 0
}

// Configuration is the TOC configuration number. Retained for pion/opus v0.1.0 callers.
type Configuration byte

// Application selects an encoder tuning profile.
type Application int

const (
	ApplicationVoIP               Application = OPUS_APPLICATION_VOIP
	ApplicationAudio              Application = OPUS_APPLICATION_AUDIO
	ApplicationRestrictedLowDelay Application = OPUS_APPLICATION_RESTRICTED_LOWDELAY
)

var (
	errInvalidSampleRate                = errors.New("invalid sample rate")
	errInvalidChannelCount              = errors.New("invalid channel count")
	errOutBufferTooSmall                = errors.New("out isn't large enough")
	errMalformedPacket                  = errors.New("malformed packet")
	errTooShortForTableOfContentsHeader = errors.New("packet is too short to contain table of contents header")
	errInvalidPLCFrameSize              = errors.New("PLC output has an unsupported duration for the previous packet mode")
	errInvalidInputLength               = errors.New("invalid input length")
	errInvalidFrameSize                 = errors.New("invalid frame size")
	errInvalidBandwidth                 = errors.New("invalid bandwidth")
	errInvalidApplication               = errors.New("invalid application")
	errBitrateOutOfRange                = errors.New("bitrate out of range")
	errInvalidComplexity                = errors.New("invalid complexity")
	errInvalidLossRate                  = errors.New("loss rate must be 0-100")
)

func validOutput(rate, channels int) error {
	switch rate {
	case 8000, 12000, 16000, 24000, 48000:
	default:
		return errInvalidSampleRate
	}
	if channels != 1 && channels != 2 {
		return errInvalidChannelCount
	}
	return nil
}

func packetBandwidth(toc byte) Bandwidth {
	c := toc >> 3
	switch {
	case c < 4:
		return BandwidthNarrowband
	case c < 8:
		return BandwidthMediumband
	case c < 12:
		return BandwidthWideband
	case c < 14:
		return BandwidthSuperwideband
	case c < 16:
		return BandwidthFullband
	case c < 20:
		return BandwidthNarrowband
	case c < 24:
		return BandwidthWideband
	case c < 28:
		return BandwidthSuperwideband
	default:
		return BandwidthFullband
	}
}
