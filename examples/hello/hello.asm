;-----------------------------
; Hello World for ZX Spectrum
;-----------------------------

        INCLUDE "constants.inc"

        ORG $8000            ; Start at 32768

start:
        ; Clear screen
        LD A,7               ; White paper, black ink
        LD (23693),A        ; Set system colors
        CALL 3503           ; ROM routine: Clear screen

        ; Print message
        LD DE,message       ; Point DE to message
        LD BC,msg_len       ; Length of message
        CALL print          ; Print it

        RET                 ; Return to BASIC

; Print routine
print:  LD A,(DE)          ; Get character
        CP $FF             ; End marker?
        RET Z              ; Yes, done
        RST 16             ; Print character
        INC DE             ; Next character
        JR print           ; Loop

message:
        DEFB "Hello, ZX Spectrum World!"
        DEFB 13            ; New line
        DEFB $FF           ; End marker
msg_len EQU $ - message