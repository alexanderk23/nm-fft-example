package main

import (
	"flag"
	"log"
	"time"

	"github.com/alexanderk23/nm-fft-example/internal/analyzer"
	"github.com/alexanderk23/nm-fft-example/internal/client"
	"github.com/gorilla/websocket"
	"github.com/hophiphip/portaudio"
)

var addr = flag.String("s", "wss://nm.alexanderk.ru/nyukomatic/", "bonzomatic server url")
var sampleRate = flag.Int("r", 44100, "sample rate")
var fps = flag.Float64("f", 50.0, "fps")

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

func audioStreamPump(bufferLen int, c chan []float32) {
	buffer := make([]float32, bufferLen)

	stream, err := initAudio(buffer)
	if err != nil {
		log.Fatalf("error initializing audio: %v", err)
		return
	}
	defer func() {
		stream.Close()
		portaudio.Terminate()
	}()

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
		select {
		case c <- buffer[:]:
		default:
		}
	}
}

func run(c chan []float32, conn *websocket.Conn, analyzerInstance *analyzer.Analyzer) {
	defer conn.Close()

	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				log.Printf("ws read error: %v", err)
				return
			}
		}
	}()

	for {
		select {
		case buffer := <-c:
			values := analyzerInstance.Process(buffer)
			port := (uint16(len(values)) << 8) | 0x00fb

			portsMessage := client.PortsMessage{
				Ports: []client.PortValues{{
					Port:   port,
					Values: values,
				}},
			}

			if err := conn.WriteJSON(portsMessage); err != nil {
				log.Printf("ws write error: %v", err)
				return
			}
		}
	}
}

func main() {
	flag.Parse()

	bufferLen := int(float64(*sampleRate) / *fps)
	analyzerInstance := analyzer.NewAnalyzer(bufferLen, *sampleRate)
	log.Printf("sample rate: %dHz fps: %.1f buffer size: %db", *sampleRate, *fps, bufferLen)

	c := make(chan []float32, 2)
	go audioStreamPump(bufferLen, c)

	for {
		conn, _, err := websocket.DefaultDialer.Dial(*addr, nil)
		if err != nil {
			log.Println("dial:", err)
			time.Sleep(3 * time.Second)
		} else {
			log.Printf("dial: connected to %s", *addr)
			run(c, conn, analyzerInstance)
		}
	}
}
