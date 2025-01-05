# asm6502

This is an Assembler / Disassembler for the 6502 microprocessor.
At some point, I might try and add an emulator...

This is a work-in-progress. I'm doing this for the fun...
This was inspired from my days playing with an Apple 2+ back in the 80s,
and the desire to revive `call -151`

* 2025-01-04: Disassembler seems to be working.
```
Disassembling [Apple2_Plus.rom] at [d000]...
                          ORG $D000
D000:
D000:            MEM_1    EQU $00
D000:            MEM_2    EQU $01
D000:            MEM_3    EQU $02
D000:            MEM_4    EQU $03
D000:            MEM_5    EQU $04
D000:            MEM_6    EQU $05
D000:            MEM_7    EQU $06
D000:            MEM_8    EQU $07
D000:            MEM_9    EQU $08
D000:            MEM_10   EQU $09
D000:            MEM_11   EQU $0A
D000:            MEM_12   EQU $0B
D000:            MEM_13   EQU $0C
D000:            MEM_14   EQU $0D
D000:            MEM_15   EQU $0E
D000:            MEM_16   EQU $0F
D000:            MEM_17   EQU $10
(...)
D35E:  42                 ???
D35F:  52                 ???
D360:  45 41              EOR $41       ; MEM_64
D362:  4B                 ???
D363:  07                 ???
D364:  00                 BRK
D365:  BA        LBL_364  TSX
D366:  E8                 INX
D367:  E8                 INX
D368:  E8                 INX
D369:  E8                 INX
D36A:  BD 01 01  REL_13   LDA $0101,X   ; LBL_2
D36D:  C9 81              CMP #$81
D36F:  D0 21              BNE REL_16
D371:  A5 86              LDA $86       ; MEM_129
D373:  D0 0A              BNE REL_14
D375:  BD 02 01           LDA $0102,X   ; LBL_3
D378:  85 85              STA $85       ; MEM_128
D37A:  BD 03 01           LDA $0103,X   ; LBL_4
D37D:  85 86              STA $86       ; MEM_129
D37F:  DD 03 01  REL_14   CMP $0103,X   ; LBL_4
D382:  D0 07              BNE REL_15
D384:  A5 85              LDA $85       ; MEM_128
D386:  DD 02 01           CMP $0102,X   ; LBL_3
D389:  F0 07              BEQ REL_16
D38B:  8A        REL_15   TXA
D38C:  18                 CLC
D38D:  69 12              ADC #$12
D38F:  AA                 TAX
D390:  D0 D8              BNE REL_13
D392:  60        REL_16   RTS
D393:  20 E3 D3  LBL_365  JSR $D3E3     ; LBL_368
```

## Usage

### Disassembler

`./asm6502 -d [-l <load_address>] <hex_file>`

### Assembler

`./asm6502 [-o <hex_file>] <source>`

## instruction set
https://www.masswerk.at/6502/6502_instruction_set.html