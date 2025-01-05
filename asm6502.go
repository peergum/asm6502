/*
   asm6502,
   Copyright (C) 2024  Phil Hilger

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

// Package asm6502 does ...
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
)

var (
	help   bool
	disasm bool
	emul   bool
	start  string
	load   string
)

func init() {
	flag.BoolVar(&help, "h", false, "show help")
	flag.BoolVar(&disasm, "d", false, "disassemble source")
	flag.BoolVar(&emul, "e", false, "emulate source")
	flag.StringVar(&start, "s", "0000", "start address")
	flag.StringVar(&load, "l", "0000", "load address")
}

func main() {
	flag.Parse()
	srcName := flag.Arg(0)

	if help {
		usage()
		return
	}

	var sAddr, lAddr int
	if _, err := fmt.Sscanf(start, "%x", &sAddr); err != nil {
		log.Fatalf("Error parsing start address: %v", err)
	}
	if _, err := fmt.Sscanf(load, "%x", &lAddr); err != nil {
		log.Fatalf("Error parsing load address: %v", err)
	}
	fmt.Println(load, lAddr)

	switch {
	case disasm:
		disassemble(srcName, lAddr)
	case emul:
		//emulator(srcName, lAddr, sAddr)
	default:
		assemble(srcName)
	}
}

func usage() {
	fmt.Printf("\n6502 Assembler/Disassembler\n(c) Copyright 2024, Phil Hilger\n\n")
	fmt.Printf("Usage:" + path.Base(os.Args[0]) + " [-d][-o dest] file\n")
	fmt.Println()
}
