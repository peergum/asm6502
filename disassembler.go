/*
   main,
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

// Package main does ...
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"sort"
)

type Labels map[string]string

var labels = make(Labels)

type Line struct {
	addr     int
	codes    string
	dis      string
	value    int
	comment1 string
	comment2 string
}

func disassemble(srcName string, lAddr int) {

	fmt.Printf("Disassembling [%s] at [%04x]...\n", srcName, lAddr)

	var src *os.File
	var err error
	if srcName == "" {
		src = os.Stdin
	} else {
		src, err = os.Open(srcName)
	}
	scanner := bufio.NewReader(src)

	var b [3]byte
	addr := lAddr

	var lines []Line
	var keys []int
	var rels []int
	//addrCnt := 0
	//memCnt := 0
disLoop:
	for {
		b[0], err = scanner.ReadByte()
		if err != nil {
			break
		}

		hi := b[0] & 0xf0 >> 4
		lo := b[0] & 0x0f

		var code string
		if len(instructions[lo]) == 8 {
			code = instructions[lo][hi>>1]
		} else {
			code = instructions[lo][hi]
		}
		mode := oper[lo][hi%2]
		instr := modes[mode]

		// special cases
		hex := fmt.Sprintf("%02X", b[0])
		if ill, err := regexp.MatchString(illegalCodes, hex); err == nil && ill {
			code = ""
			instr = modes[13]
		} else if lo == 0 && hi < 8 && hi%2 == 0 {
			if hex == "20" {
				instr = modes[1]
			} else {
				instr = modes[5]
			}
		} else if lo == 10 && hi >= 8 {
			instr = modes[5]
		} else if hex == "6C" {
			instr = modes[6]
		}

		iType := instr.string
		iLen := instr.int

		value := 0
		valueSet := false
		for i := range iLen - 1 {
			v, err := scanner.ReadByte()
			if err != nil {
				break disLoop
			}
			value += int(v) << (8 * i)
			b[1+i] = v
			valueSet = true
		}

		var dis string
		var comment1 string
		var comment2 string

		sCodes := ""
		for i := range 3 {
			if i < iLen {
				sCodes += fmt.Sprintf(" %02X", b[i])
				//if b[i]&0x7f < 32 {
				//	comment2 += "."
				//} else if b[i]&0x80 == 0 {
				//	comment2 += string(b[i])
				//} else {
				//	comment2 += "[" + string(b[i]&0x7f) + "]"
				//}
			} else {
				sCodes += "   "
			}
		}

		switch iType {
		case "A":
			dis = fmt.Sprintf("%s A", code)
		case "#":
			dis = fmt.Sprintf("%s #$%02X", code, value)
			comment1 = fmt.Sprintf("%02X", value)
		case "abs":
			dis = fmt.Sprintf("%s $%04X", code, value)
			comment1 = fmt.Sprintf("%04X", value)
		case "abs,X":
			dis = fmt.Sprintf("%s $%04X,X", code, value)
			comment1 = fmt.Sprintf("%04X", value)
		case "abs,Y":
			dis = fmt.Sprintf("%s $%04X,Y", code, value)
			comment1 = fmt.Sprintf("%04X", value)
		case "zpg":
			dis = fmt.Sprintf("%s $%02X", code, value)
			comment1 = fmt.Sprintf("%02X", value)
		case "zpg,X":
			dis = fmt.Sprintf("%s $%02X,X", code, value)
			comment1 = fmt.Sprintf("%02X", value)
		case "zpg,Y":
			dis = fmt.Sprintf("%s $%02X,Y", code, value)
			comment1 = fmt.Sprintf("%02X", value)
		case "(ind)":
			dis = fmt.Sprintf("%s ($%04X)", code, value)
			comment1 = fmt.Sprintf("%04X", value)
		case "(ind,X)":
			dis = fmt.Sprintf("%s ($%04X,X)", code, value)
			comment1 = fmt.Sprintf("%04X", value)
		case "(ind),Y":
			dis = fmt.Sprintf("%s ($%04X),Y", code, value)
			comment1 = fmt.Sprintf("%04X", value)
		case "rel":
			v := int8(value)
			rel := uint16(int16(addr+iLen) + int16(v))
			value = int(rel)
			dis = fmt.Sprintf("%s $%04X", code, value)
			comment1 = fmt.Sprintf("%04X", value)
		case "ill":
			dis = "???"
		default:
			dis = code
		}

		lines = append(lines, Line{addr, sCodes, dis, value, comment1, comment2})
		if code == "RTS" || code == "RTI" || code == "JMP" {
			lines = append(lines, Line{addr, "", "", 0, "", ""})
		}
		if valueSet {
			if iType == "rel" {
				i := sort.SearchInts(rels, value)
				if i == len(rels) || rels[i] != value {
					rels = slices.Insert(rels, i, value)
				}
			} else {
				i := sort.SearchInts(keys, value)
				if i == len(keys) || keys[i] != value {
					keys = slices.Insert(keys, i, value)
				}
			}
		}
		addr += iLen
	}

	fmt.Printf("%25s ORG $%04X\n", "", lAddr)
	fmt.Printf("%04X:\n", lAddr)

	sort.Ints(keys)

	// memory and labels
	var labels = make(map[string]string)
	memCnt := 1
	addrCnt := 1
	for _, addr := range keys {
		if addr < 256 {
			label := fmt.Sprintf("MEM_%d", memCnt)
			fmt.Printf("%04X: %-10s %-8s EQU $%02X\n", lAddr, "", label, addr)
			labels[fmt.Sprintf("%02X", addr)] = label
			memCnt++
		} else {
			label := fmt.Sprintf("LBL_%d", addrCnt)
			fmt.Printf("%04X: %-10s %-8s EQU $%04X\n", lAddr, "", label, addr)
			labels[fmt.Sprintf("%04X", addr)] = label
			addrCnt++
		}
	}

	// simple relative labels
	addrCnt = 1
	for _, addr := range rels {
		label := fmt.Sprintf("REL_%d", addrCnt)
		labels[fmt.Sprintf("%04X", addr)] = label
		addrCnt++
	}

	fmt.Printf("%04X:\n", lAddr)

	for _, line := range lines {
		if line.dis == "" {
			fmt.Println()
			continue
		}
		re1 := regexp.MustCompile(`^([0-9A-F][0-9A-F][0-9A-F][0-9A-F])$`)
		re2 := regexp.MustCompile(`^([0-9A-F][0-9A-F])$`)
		comment1 := re1.ReplaceAllString(line.comment1, labels[fmt.Sprintf("%04X", line.value)])
		comment1 = re2.ReplaceAllString(comment1, labels[fmt.Sprintf("%02X", line.value)])
		label := labels[fmt.Sprintf("%04X", line.addr)]
		fmt.Printf("%04X: %-10s %-8s %-15s ; %-15s   %-3s\n", line.addr, line.codes, label, line.dis, comment1, line.comment2)
	}
}

func getLabel(s string) string {
	log.Printf("%s %s", s, labels[s])
	return labels[s]
}
