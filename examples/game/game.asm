;-----------------------------
; Simple Game Example
;-----------------------------

        INCLUDE "constants.inc"

        ORG $8000

start:
        ; Initialize
        CALL init
        
main_loop:
        ; Wait for frame
        HALT
        
        ; Handle input
        CALL read_keys
        
        ; Update sprite position
        CALL move_sprite
        
        ; Draw sprite
        CALL draw_sprite
        
        ; Check for exit (SPACE)
        LD A,(key_space)
        CP 1
        JR NZ,main_loop
        
        RET

init:
        ; Set up initial position
        LD A,100
        LD (sprite_x),A
        LD A,100
        LD (sprite_y),A
        
        ; Clear screen
        LD A,7
        LD (23693),A
        CALL 3503
        RET

read_keys:
        ; Read Q key (up)
        LD BC,$FBFE
        IN A,(C)
        BIT 0,A
        LD A,0
        JR NZ,not_up
        LD A,1
not_up: LD (key_up),A

        ; Similar code for other keys...
        RET

move_sprite:
        ; Move based on keys
        LD A,(key_up)
        CP 1
        JR NZ,not_move_up
        LD A,(sprite_y)
        DEC A
        LD (sprite_y),A
not_move_up:
        ; Similar code for other directions...
        RET

draw_sprite:
        ; Load sprite data and draw it
        LD IX,sprite_data
        LD A,(sprite_x)
        LD B,A
        LD A,(sprite_y)
        LD C,A
        CALL draw_16x16
        RET

; Variables
sprite_x:    DEFB 0
sprite_y:    DEFB 0
key_up:      DEFB 0
key_down:    DEFB 0
key_left:    DEFB 0
key_right:   DEFB 0
key_space:   DEFB 0

; Include sprite data
sprite_data:
        INCBIN "sprite.bin"