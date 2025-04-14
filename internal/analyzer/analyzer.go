package analyzer

import (
	"math"
	"math/cmplx"

	"gonum.org/v1/gonum/dsp/fourier"
	"gonum.org/v1/gonum/dsp/window"
)

type Analyzer struct {
	fft        *fourier.FFT
	coeff      []complex128
	data       []float64
	sampleRate float64
	bands      []uint16
}

func NewAnalyzer(bufferLen, sampleRate int) *Analyzer {
	return &Analyzer{
		fft:        fourier.NewFFT(bufferLen),
		coeff:      make([]complex128, bufferLen/2+1),
		data:       make([]float64, bufferLen),
		sampleRate: float64(sampleRate),
		bands:      make([]uint16, 3),
	}
}

func (a *Analyzer) FreqIndex(f float64) int {
	step := float64(len(a.data)+1) / a.sampleRate
	return max(1, min(int(math.Round(step*f)), len(a.data)/2-1))
}

func (a *Analyzer) PeakMagnitude(fs, fe float64) float64 {
	start := a.FreqIndex(fs)
	end := a.FreqIndex(fe)
	maxVal := 0.0

	for i := start; i < end; i++ {
		magnitude := cmplx.Abs(a.coeff[i])
		if magnitude > maxVal {
			maxVal = magnitude
		}
	}

	return maxVal
}

func (a *Analyzer) BarHeight(fs, fe float64) uint16 {
	magnitude := a.PeakMagnitude(fs, fe)
	if magnitude > 0 {
		dbVal := 20 * math.Log10(magnitude*magnitude)
		scaled := min(max((dbVal+10)*2.5, 0), 255)
		return uint16(scaled)
	}
	return 0
}

func (a *Analyzer) Process(buffer []float32) []uint16 {
	for i, v := range buffer {
		a.data[i] = float64(v)
	}

	window.Hamming(a.data)
	a.fft.Coefficients(a.coeff, a.data)

	a.bands[0] = a.BarHeight(20.0, 250.0)
	a.bands[1] = a.BarHeight(250.0, 4000.0)
	a.bands[2] = a.BarHeight(4000.0, 20000.0)

	return a.bands
}
