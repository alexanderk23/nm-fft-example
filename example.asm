		org $8000

FFT_DST_ADDR	= $4000
FFT_BASE_PORT	= $20fb

		call init

main		halt

		ld hl, FFT_DST_ADDR
		ld bc, FFT_BASE_PORT
		inir

		call draw
		jp main

init		xor a
		out ($fe), a

		ld hl, $5800
		ld de, $5801
		ld bc, 767
		ld (hl), l
		ldir

		ld hl, $4000
		ld de, $4001
		ld (hl), %01111110
		ld bc, 6143
		ldir

		ld hl, $4000
		ld b, 192 / 4
.l1		ld a, 32
		ld d, l
.l2		ld (hl), c
		inc l
		dec a
		jr nz, .l2
		ld l, d
		call down4
		djnz .l1
		ei
		ret

down4		ld a, h
		add 4
		ld h, a
		and $07
		ret nz
		ld a, l
		sub $e0
		ld l, a
		sbc a
		and $f8
		add h
		ld h, a
		ret

draw		ld ix, FFT_DST_ADDR
		ld hl, $5b00-32
		ld de, -32

.l0		push hl

		ld a, (ix)
		.4 rra
		and 15
		ld c, a
		jr z, .empty

		ld b, a
		cp 12
		ld a, 005o
		jr c, $+4
		or %1000000

.l1		ld (hl), a
		add hl, de
		djnz .l1

.empty		ld a, 16
		sub c
		ld b, a
		xor a

.l2		ld (hl), a
		add hl, de
		djnz .l2

		pop hl

		inc ixl
		inc l
		jp nz, .l0
		ret
