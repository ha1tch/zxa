;-----------------------------
; Test constants file
;-----------------------------

; Memory addresses
SCREEN      EQU $4000
ATTR        EQU $5800
WORKSP      EQU $5B00

; System ports
ULA_PORT    EQU $FE
KEMPS_PORT  EQU $1F

; Color attributes
INK_BLACK   EQU 0
INK_BLUE    EQU 1
INK_RED     EQU 2
INK_MAG     EQU 3
INK_GREEN   EQU 4
INK_CYAN    EQU 5
INK_YELLOW  EQU 6
INK_WHITE   EQU 7

; Attribute bits
FLASH       EQU %10000000
BRIGHT      EQU %01000000
PAPER_MASK  EQU %00111000
INK_MASK    EQU %00000111