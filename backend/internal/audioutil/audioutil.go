package audioutil

import (
	"bytes"
	"fmt"
	"math"
	"song-match-backend/domain"

	"github.com/go-audio/wav"
	"github.com/madelynnblue/go-dsp/fft"
)

func HannWindow(signal []float64) {
	for i := range signal {
		hannValue := 0.5 * (1.0 - math.Cos(2.0*math.Pi*float64(i)/float64(len(signal)-1)))
		signal[i] *= hannValue
	}
}

func FindPeaks(spectrum []float64) []int {
	peaks := []int{}

	for i := 1; i < len(spectrum)-1; i++ {
		if spectrum[i] > spectrum[i-1] && spectrum[i] > spectrum[i+1] {
			peaks = append(peaks, i)
		}
	}
	return peaks
}

func ExtractFingerprints(samples []float64, sampleRate int) ([]domain.TrackFingerprint, error) {
	windowSize := 8192
	hopSize := 4096

	var fingerprints []domain.TrackFingerprint

	// Slide window over audio
	for i := 0; i+windowSize < len(samples); i += hopSize {
		window := make([]float64, windowSize)
		copy(window, samples[i:i+windowSize])

		HannWindow(window)
		spectrum := computeFFT(window)

		fingerprint := domain.TrackFingerprint{
			Timestamp: float64(i) / float64(sampleRate),
			Peaks:     FindPeaks(spectrum),
		}

		fingerprints = append(fingerprints, fingerprint)
	}

	return fingerprints, nil
}

func DecodeAudio(data []byte) ([]float64, int, error) {
	if len(data) < 4 || string(data[:4]) != "RIFF" {
		return nil, 0, fmt.Errorf("invalid WAV file: missing RIFF header (got %q)", string(data[:4]))
	}

	reader := bytes.NewReader(data)
	decoder := wav.NewDecoder(reader)

	if !decoder.IsValidFile() {
		return nil, 0, fmt.Errorf("invalid WAV file")
	}

	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, 0, err
	}

	n := len(buf.Data)

	// Convert to float64 and mix channels if stereo
	samples := make([]float64, n)
	for i := 0; i < n; i++ {
		samples[i] = float64(buf.Data[i]) / 32768.0
	}

	return samples, int(decoder.Format().SampleRate), nil
}

func computeFFT(signal []float64) []float64 {
	// Pad to next power of 2 for FFT efficiency
	fftSize := 1
	for fftSize < len(signal) {
		fftSize *= 2
	}

	padded := make([]complex128, fftSize)
	for i, s := range signal {
		padded[i] = complex(s, 0)
	}

	result := fft.FFT(padded)

	// Get magnitude spectrum
	spectrum := make([]float64, len(result)/2)
	for i := range spectrum {
		spectrum[i] = math.Abs(real(result[i]))
	}

	return spectrum
}
