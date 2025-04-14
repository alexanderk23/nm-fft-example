# nm-fft-example

### Usage
```sh
go build
./nm-fft-example -s wss://nm.alexanderk.ru/nyukomatic/
```

### Description
This tool continuously processes audio in real time at a fixed frame rate
(default: **50 FPS**, configurable via `-f`):

1. **Captures audio** from the default input device at `44.1 kHz` (configurable via `-r`).
2. **Performs FFT analysis** on the streaming audio data, calculating peak magnitudes
for **three** defined frequency bands (bass, midrange and treble).
3. **Normalizes** these values to a `0..255` range and **sends** them to
[nyukomatic](https://github.com/alexanderk23/nyukomatic/) via WebSocket
using a [modified Bonzomatic server](https://github.com/alexanderk23/BonzomaticServer)
as a relay.

### Accessing Spectrum Data in Z80 Assembly
The FFT results are mapped to three ports:

| **Port** | **Frequency Range** | **Band**     |
|----------|---------------------|--------------|
| `$03FB`  | 20–250 Hz           | Bass         |
| `$02FB`  | 250–4000 Hz         | Midrange     |
| `$01FB`  | 4000–20000 Hz       | Treble       |

### Example Code
See [example.asm](example.asm) for a usage example.
