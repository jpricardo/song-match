package audioutil

import (
	"bytes"
	"fmt"
	"math"
	"song-match-backend/domain"
	"sort"

	"github.com/go-audio/wav"
	"github.com/madelynnblue/go-dsp/fft"
)

func HannWindow(signal []float64) {
	for i := range signal {
		hannValue := 0.5 * (1.0 - math.Cos(2.0*math.Pi*float64(i)/float64(len(signal)-1)))
		signal[i] *= hannValue
	}
}

type peakData struct {
	bin       int
	magnitude float64
}

// Define our frequency bands in Hz
var frequencyBands = [5][2]float64{
	{40, 250},     // Bass (Kick drums, Bassline)
	{250, 500},    // Low-Mids (Cello, Low guitars, Deep vocals)
	{500, 2000},   // Mids (Vocals, Lead guitars, Synths)
	{2000, 4000},  // High-Mids (Snare attack, Cymbals)
	{4000, 16000}, // Treble (Air, Hi-hats)
}

// Update FindPeaks to accept sampleRate to calculate Hz correctly
func FindPeaks(spectrum []float64, sampleRate int) []int {
	peaks := []int{}

	// The full FFT size is twice the spectrum length
	fftSize := len(spectrum) * 2
	hzPerBin := float64(sampleRate) / float64(fftSize)

	// Iterate through each of the 5 frequency bands
	for _, band := range frequencyBands {
		var localPeaks []peakData

		// Convert the Hz targets into FFT bin indices
		minBin := int(band[0] / hzPerBin)
		maxBin := int(band[1] / hzPerBin)

		// Ensure we don't go out of bounds of our array
		if minBin < 1 {
			minBin = 1
		}
		if maxBin > len(spectrum)-2 {
			maxBin = len(spectrum) - 2
		}

		// 1. Find all local maxima within THIS specific band
		for i := minBin; i <= maxBin; i++ {
			if spectrum[i] > spectrum[i-1] && spectrum[i] > spectrum[i+1] {
				// Noise gate: Ignore total silence
				if spectrum[i] > 0.05 {
					localPeaks = append(localPeaks, peakData{
						bin:       i,
						magnitude: spectrum[i],
					})
				}
			}
		}

		// 2. Sort the peaks in THIS band descending by magnitude
		sort.Slice(localPeaks, func(i, j int) bool {
			return localPeaks[i].magnitude > localPeaks[j].magnitude
		})

		// 3. Keep only the Top 2 loudest peaks from this band
		maxPeaksForBand := 2
		if len(localPeaks) < maxPeaksForBand {
			maxPeaksForBand = len(localPeaks)
		}

		for i := 0; i < maxPeaksForBand; i++ {
			peaks = append(peaks, localPeaks[i].bin)
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
			Peaks:     FindPeaks(spectrum, sampleRate),
		}

		fingerprints = append(fingerprints, fingerprint)
	}

	fmt.Printf("Extracted %d fingerprints at %d Hz\n", len(fingerprints), sampleRate)

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

	samples := make([]float64, n)
	for i := 0; i < n; i++ {
		samples[i] = float64(buf.Data[i]) / 32768.0
	}

	return samples, int(decoder.Format().SampleRate), nil
}

func computeFFT(signal []float64) []float64 {
	fftSize := 1
	for fftSize < len(signal) {
		fftSize *= 2
	}

	padded := make([]complex128, fftSize)
	for i, s := range signal {
		padded[i] = complex(s, 0)
	}

	result := fft.FFT(padded)

	spectrum := make([]float64, len(result)/2)
	for i := range spectrum {
		spectrum[i] = math.Abs(real(result[i]))
	}

	return spectrum
}

type AudioHash struct {
	HashValue string
	Time      float64
}

func GenerateHashes(fingerprints []domain.TrackFingerprint) []AudioHash {
	var hashes []AudioHash

	// How many future fingerprints to pair with the anchor
	targetZone := 5

	for i := range fingerprints {
		anchor := fingerprints[i]

		end := i + targetZone
		if end > len(fingerprints) {
			end = len(fingerprints)
		}

		// Create pairs between the anchor and the targets
		for j := i + 1; j < end; j++ {
			target := fingerprints[j]

			// Time difference in milliseconds (the 3rd part of our hash)
			dt := int((target.Timestamp - anchor.Timestamp) * 1000)

			// Cross-match every peak in the anchor with every peak in the target
			for _, f1 := range anchor.Peaks {
				for _, f2 := range target.Peaks {

					// The cryptographic hash: "f1|f2|dt"
					hashVal := fmt.Sprintf("%d|%d|%d", f1, f2, dt)

					hashes = append(hashes, AudioHash{
						HashValue: hashVal,
						Time:      anchor.Timestamp,
					})
				}
			}
		}
	}

	return hashes
}
