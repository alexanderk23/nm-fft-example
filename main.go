package main

import (
	"encoding/json"
	"flag"
	"log"
	"math"
	"math/cmplx"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/recws-org/recws"
	"gonum.org/v1/gonum/dsp/fourier"
)

type PortValues struct {
	Port   uint16
	Values []uint16
}

func (r *PortValues) MarshalJSON() ([]byte, error) {
	arr := []any{r.Port, r.Values}
	return json.Marshal(arr)
}

type PortsMessage struct {
	Ports []PortValues `json:"ports"`
}

var addr = flag.String("s", "wss://nm.alexanderk.ru/nyukomatic/", "bonzomatic server url")
var sampleRate = flag.Int("r", 48000, "sample rate")

func initAudio(buffer []float32) (*portaudio.Stream, error) {
	err := portaudio.Initialize()
	if err != nil {
		return nil, err
	}

	stream, err := portaudio.OpenDefaultStream(
		1, 0, float64(*sampleRate), len(buffer), &buffer)
	if err != nil {
		return nil, err
	}

	return stream, nil
}

type Analyzer struct {
	buffer *[]float32
	data   []float64
	fft    *fourier.FFT
	Bands  []uint16
}

func NewAnalyzer(buffer []float32) *Analyzer {
	size := len(buffer)
	return &Analyzer{
		fft:    fourier.NewFFT(size),
		buffer: &buffer,
		data:   make([]float64, size),
		Bands:  make([]uint16, 32),
	}
}

func (a *Analyzer) Process() {
	for i, v := range *a.buffer {
		a.data[i] = float64(v)
	}

	coeff := a.fft.Coefficients(nil, a.data)
	magnitude := make([]float64, len(coeff)/2)
	for i := range magnitude {
		magnitude[i] = cmplx.Abs(coeff[i])
	}
	bandSize := len(magnitude) / len(a.Bands)

	for i := range a.Bands {
		start := i * bandSize
		end := start + bandSize
		if i == 31 {
			end = len(magnitude)
		}

		maxVal := 0.0
		for j := start; j < end; j++ {
			if magnitude[j] > maxVal {
				maxVal = magnitude[j]
			}
		}

		if maxVal > 0 {
			dbVal := 20 * math.Log10(maxVal)
			normalized := min(max((dbVal+10)*8, 0), 255)
			a.Bands[i] = uint16(normalized)
		}
	}
}

func readLoop(conn *recws.RecConn) {
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			log.Printf("ws read error: %v", err)
			time.Sleep(5 * time.Second)
		}
	}
}

func main() {
	flag.Parse()

	conn := recws.RecConn{KeepAliveTimeout: 15 * time.Second}
	conn.Dial(*addr, nil)
	defer conn.Close()
	go readLoop(&conn)

	buffer := make([]float32, 2048)
	analyzer := NewAnalyzer(buffer)

	stream, err := initAudio(buffer)
	if err != nil {
		log.Fatalf("error initializing audio: %v", err)
		return
	}
	defer stream.Close()

	if err := stream.Start(); err != nil {
		log.Fatalf("error starting audio stream: %v", err)
		return
	}
	defer stream.Stop()

	for {
		if err := stream.Read(); err != nil {
			log.Printf("audio stream read error: %v", err)
			continue
		}

		if !conn.IsConnected() {
			continue
		}

		analyzer.Process()
		portsMessage := PortsMessage{
			Ports: []PortValues{{0x20fb, analyzer.Bands}},
		}

		if err := conn.WriteJSON(portsMessage); err != nil {
			log.Printf("ws write error: %v", err)
		}
	}
}
