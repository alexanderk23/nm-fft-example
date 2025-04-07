# nm-fft-example

```sh
go build
./nm-fft-example -s wss://nm.alexanderk.ru/nyukomatic/ -r 48000
```
This tool samples the default audio input device, performs an FFT and sends
32 bytes of spectrum analyzer band heights (0..255) to the Bonzomatic server.
These can be obtained from Z80 assembly code by reading ports `$20FB`..`$01FB`.

See [example.asm](example.asm) for an example of usage.
