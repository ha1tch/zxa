;-----------------------------
; Test assembly file
;-----------------------------

        INCLUDE "constants.inc"

        ORG $8000

; Test labels and constants
START       EQU $8000
SCREEN     EQU $4000
WORKSP     EQU $5B00

; Test data definitions
data1:  DEFB 1,2,3,4,5
data2:  DEFW $1234, $5678
space:  DEFS 16, $FF

; Test instructions
start:
        LD A,42
        LD B,10
loop:   DJNZ loop
        
        LD IX,data1
        LD (IX+0),B
        
        CALL subr
        RET

subr:   PUSH BC
        POP BC
        RET

; Test binary inclusion
sprite: INCBIN "sprite.bin"