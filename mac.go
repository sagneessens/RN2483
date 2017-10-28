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
	"fmt"
	"strconv"
	"time"

	"github.com/bullettime/logger"
)

// MacReset will automatically reset the software LoRaWAN stack and initilize
// it with the parameters for the selected band.
func MacReset(band uint16) bool {
	if band != 433 && band != 868 {
		logger.Warning.Println("mac reset error: invalid band selected (433 or 868)")
		return false
	}

	err := serialWrite(fmt.Sprintf("mac reset %v", band))
	if err != nil {
		logger.Warning.Println("mac reset error:", err)
		return false
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) == invalidParameter {
		logger.Warning.Println("mac reset error: invalid parameter")
		return false
	}

	state.macPaused = false

	return true
}

// MacPause will pause the LoRaWAN stack functionality to allow transceiver (radio) configuration.
// The length is the time in milliseconds the stack will be paused, with a maximum of 4294967295
// (max of uint32), is returned as an uint32.
func MacPause() uint32 {
	err := serialWrite("mac pause")
	if err != nil {
		logger.Warning.Println("mac pause error:", err)
		return 0
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) == invalidParameter {
		logger.Warning.Println("mac pause error: invalid parameter")
		return 0
	}

	value, err := strconv.ParseUint(string(sanitize(answer)), 10, 32)
	if err != nil {
		logger.Warning.Println("mac pause error:", err)
		return 0
	}

	state.macPaused = true
	state.macPausedEnd = time.Now().Add(time.Duration(value) * time.Millisecond)

	return uint32(value)
}

// MacResume will resume the LoRaWAN stack functionality, in order to continue normal
// functionality after being paused.
func MacResume() bool {
	err := serialWrite("mac resume")
	if err != nil {
		logger.Warning.Println("mac resume error:", err)
		return false
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) == invalidParameter {
		logger.Warning.Println("mac resume error: invalid parameter")
		return false
	}

	state.macPaused = false

	return true
}

// The length is passed in milliseconds.
func isMacPaused(length int) bool {
	// An offset of 100 milliseconds is added to ensure we have enough
	// time left over
	d := time.Duration(length)*time.Millisecond + time.Duration(100)

	if state.macPaused {
		return time.Now().Add(d).Before(state.macPausedEnd)
	}

	return false
}
