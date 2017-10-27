// The MIT License (MIT)
//
// Copyright © 2017 Sven Agneessens <sven.agneessens@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package rn2483

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/bullettime/logger"
)

// Sleep puts the RN2483 chip to sleep for the specified number of milliseconds.
func Sleep(length uint32) bool {
	if length < 100 {
		logger.Warning.Println("sys sleep called with length lower than 100:", length)
		return false
	}

	err := serialWrite(fmt.Sprintf("sys sleep %v", length))
	if err != nil {
		logger.Warning.Println("sys sleep error:", err)
		return false
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) == invalidParameter {
		logger.Warning.Println("sys sleep error: invalid parameter")
		return false
	}

	return true
}

// Reset will reset and restart the RN2483 module.
func Reset() bool {
	err := serialWrite("sys reset")
	if err != nil {
		logger.Warning.Println("reset error:", err)
		return false
	}

	serialFlush()

	return true
}

// SaveByte allows the user to modify the EEPROM at the specified address
// with the specified data (one byte).
func SaveByte(address uint16, data uint8) bool {
	if address < 768 || address > 1023 {
		logger.Warning.Println("sys set nvm error: address out of range [768-1023]")
		return false
	}

	err := serialWrite(fmt.Sprintf("sys set nvm %X %X", address, data))
	if err != nil {
		logger.Warning.Println("sys set nvm error:", err)
		return false
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) == invalidParameter {
		logger.Warning.Println("sys set nvm error: invalid parameter")
		return false
	}

	return true
}

// ReadByte returns the data stored in the EEPROM at the specified address.
func ReadByte(address uint16) (byte, error) {
	if address < 768 || address > 1023 {
		logger.Warning.Println("sys get nvm error: address out of range [768-1023]")
		return 0, errors.New("address out of range [768-1023]")
	}

	err := serialWrite(fmt.Sprintf("sys get nvm %X", address))
	if err != nil {
		logger.Warning.Println("sys get nvm error:", err)
		return 0, err
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) == invalidParameter {
		logger.Warning.Println("sys get nvm error: invalid parameter")
		return 0, errors.New("invalid parameter")
	}

	value, err := strconv.ParseUint(string(sanitize(answer)), 16, 8)
	if err != nil {
		logger.Warning.Println("sys get nvm error:", err)
		return 0, err
	}

	return byte(value), nil
}

// Version returns the information related to the hardware platform,
// firmware version, release date and time stamp on firmware creation.
func Version() string {
	err := serialWrite("sys get ver")
	if err != nil {
		logger.Warning.Println("sys get ver error:", err)
		return ""
	}

	n, answer := serialRead()
	if n == 0 {
		logger.Warning.Println("sys get ver error: no answer")
		return ""
	}

	return string(sanitize(answer))
}

// Voltage will return the voltage measured on Vdd in millivolts
func Voltage() (uint16, error) {
	err := serialWrite("sys get vdd")
	if err != nil {
		logger.Warning.Println("sys get vdd error:", err)
		return 0, err
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) == invalidParameter {
		logger.Warning.Println("sys get vdd error: invalid parameter")
		return 0, errors.New("invalid parameter")
	}

	value, err := strconv.ParseUint(string(sanitize(answer)), 10, 16)
	if err != nil {
		logger.Warning.Println("sys get vdd error:", err)
		return 0, err
	}

	return uint16(value), nil
}

// HardwareID will return the HWEUI of the RN2483 module as a string.
// The HWEUI is actually an 8 bit hex string.
func HardwareID() string {
	err := serialWrite("sys get hweui")
	if err != nil {
		logger.Warning.Println("sys get hweui error:", err)
		return ""
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) == invalidParameter {
		logger.Warning.Println("sys get hweui error: invalid parameter")
		return ""
	}

	return string(sanitize(answer))
}
