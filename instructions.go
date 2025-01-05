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

type Command struct {
	name   string
	code   byte
	len    byte
	cycles byte
	flags  byte
	mask   byte
}

type Mode struct {
	string
	int
}

var (
	instructions = [][]string{
		0: {"BRK", "BPL", "JSR", "BMI", "RTI", "BVC", "RTS", "BVS", "", "BCC", "LDY", "BCS", "CPY", "BNE", "CPX", "BEQ"},
		1: {"ORA", "AND", "EOR", "ADC", "STA", "LDA", "CMP", "SBC"},
		2: {"", "", "", "", "", "LDX", "", ""},
		3: {"", "", "", "", "", "", "", ""},

		4: {"", "", "BIT", "", "", "", "", "", "STY", "STY", "LDY", "LDY", "CPY", "", "CPX", ""},
		5: {"ORA", "AND", "EOR", "ADC", "STA", "LDA", "CMP", "SBC"},
		6: {"ASL", "ROL", "LSR", "ROR", "STX", "LDX", "DEC", "INC"},
		7: {"", "", "", "", "", "", "", ""},

		8:  {"PHP", "CLC", "PLP", "SEC", "PHA", "CLI", "PLA", "SEI", "DEY", "TYA", "TAY", "CLV", "INY", "CLD", "INX", "SED"},
		9:  {"ORA", "AND", "EOR", "ADC", "STA", "LDA", "CMP", "SBC"},
		10: {"ASL", "", "ROL", "", "LSR", "", "ROR", "", "TXA", "TXS", "TAX", "TSX", "DEX", "", "NOP", ""},
		11: {"", "", "", "", "", "", "", ""},

		12: {"", "", "BIT", "", "JMP", "", "JMP", "", "STY", "", "LDY", "LDY", "CPY", "", "CPX", ""},
		13: {"ORA", "AND", "EOR", "ADC", "STA", "LDA", "CMP", "SBC"},
		14: {"ASL", "ROL", "LSR", "ROR", "STX", "LDX", "DEC", "INC"},
		15: {"", "", "", "", "", "", "", ""},
	}

	modes = []Mode{
		0:  {"A", 1},
		1:  {"abs", 3},
		2:  {"abs,X", 3},
		3:  {"abs,Y", 3},
		4:  {"#", 2},
		5:  {"impl", 1},
		6:  {"(ind)", 3},
		7:  {"(X,ind)", 3},
		8:  {"(ind),Y", 3},
		9:  {"rel", 2},
		10: {"zpg", 2},
		11: {"zpg,X", 2},
		12: {"zpg,Y", 2},
		13: {"ill", 1},
	}

	oper = [][]int{
		0:  {4, 9},
		1:  {7, 8},
		2:  {4, 13},
		3:  {13, 13},
		4:  {10, 11},
		5:  {10, 11},
		6:  {10, 11},
		7:  {13, 13},
		8:  {5, 5},
		9:  {4, 3},
		10: {0, 13},
		11: {13, 13},
		12: {1, 2},
		13: {1, 2},
		14: {1, 2},
		15: {13, 13},
	}

	illegalCodes = "80|[0-9B-F]2|.3|[013-7DF]4|.7|89|[1357DF]A|.B|[013579DF]C|9E|.F"
)
