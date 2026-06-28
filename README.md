# ⚠️ zxa is superseded by zenas

**zxa is no longer developed. Use [zenas](https://github.com/ha1tch/zenas) instead.**

zxa was an early attempt at a Go-based Z80 cross-assembler for the ZX Spectrum. Its successor, **[zenas](https://github.com/ha1tch/zenas)**, is a complete, working assembler that delivers everything zxa set out to do — and is verified against real-world code.

## Why zenas

- **It actually assembles.** zenas builds a full, real-world Z80 operating system kernel **byte-for-byte identically** to the long-established pasmo assembler — every opcode and every symbol address matching exactly. Correctness is the design priority, not an afterthought.
- **Complete Z80 instruction set.** The entire documented instruction set is covered and continuously verified against a reference assembler, including the undocumented IX/IY half-register operations (IXH/IXL/IYH/IYL) with proper rejection of the illegal combinations. A reproducible coverage checker keeps this honest.
- **A real toolchain, not a skeleton.** Conditional assembly (`IF`/`IFDEF`/`ELSE`/`ENDIF`) with command-line `--define` build flags, file inclusion with forward references across boundaries, case-sensitive symbols matching the conventions of pasmo and sjasmplus, pasmo-compatible symbol-file output, and symbol arithmetic in operands (`LD HL,base+offset`, `vtable+N*3`).
- **Same spirit, finished.** Pure Go, zero dependencies, Apache-2.0, clean CLI. The goals zxa described — full Z80 support, directives, multiple output formats, modular includes — are all present and working in zenas.

## Migrating

zenas uses conventional, pasmo-style Z80 syntax, so most source written for zxa (or assumed by it) will assemble directly. Point your build at zenas:

```
zenas assemble yourprogram.asm yourprogram.bin
```

See the [zenas README](https://github.com/ha1tch/zenas) for the full directive and option reference.

