;-----------------------------
; Test Sprites - demonstrating numeric formats
;-----------------------------

        ORG 0            ; Start at 0 since it's just data

; Simple arrow sprite (16x16)
arrow:  
        ; Top part - using % binary format
        DEFB %00000110, %00000000
        DEFB %00001111, %00000000
        ; Middle part - using 0b binary format
        DEFB 0b00011111, 0b10000000
        DEFB 0b00111111, 0b11000000
        ; Bottom part - using $ hex format
        DEFB $3F, $E0
        DEFB $1F, $80
        ; Final part - using 0x hex format
        DEFB 0x1F, 0x80
        DEFB 0x0F, 0x00

; Ball using different formats
ball:   
        DEFB %00111100      ; Binary %
        DEFB 0b01111110     ; Binary 0b
        DEFB $FF            ; Hex $
        DEFB 0xFF           ; Hex 0x
        DEFB 255            ; Decimal
        DEFB 0b11111111     ; Binary again
        DEFB $7E            ; Hex again
        DEFB 0x3C           ; Hex again