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

import "time"

type myState struct {
	macPaused    bool
	macPausedEnd time.Time
}

const (
	maxUint8  = ^uint8(0)
	maxUint16 = ^uint16(0)
	maxUint32 = ^uint32(0)
	maxInt8   = int8(maxUint8 >> 1)
	maxInt16  = int16(maxUint16 >> 1)
	maxInt32  = int32(maxUint32 >> 1)
)

// The possible modulations
const (
	LoRa = "lora"
	FSK  = "fsk"
)

// The possible spreading factors
const (
	SF7  = "sf7"
	SF8  = "sf8"
	SF9  = "sf9"
	SF10 = "sf10"
	SF11 = "sf11"
	SF12 = "sf12"
)

// The possible coding rates
const (
	CR5 = "4/5"
	CR6 = "4/6"
	CR7 = "4/7"
	CR8 = "4/8"
)

// The possible bandwidths
const (
	BW1 = "125"
	BW2 = "250"
	BW3 = "500"
)

var (
	state            = new(myState)
	invalidParameter = "invalid_param"
	serialRead       = read
	serialWrite      = write
	serialFlush      = flush
)

var modulations = []string{
	LoRa,
	FSK,
}

// SFs is the mapping of the spreading factors
var SFs = map[uint8]string{
	7:  SF7,
	8:  SF8,
	9:  SF9,
	10: SF10,
	11: SF11,
	12: SF12,
}

// BWs is the mapping of the bandwidths
var BWs = map[uint16]string{
	125: BW1,
	250: BW2,
	500: BW3,
}

// CodingRates is the mapping of the coding rates
var CodingRates = map[uint8]string{
	5: CR5,
	6: CR6,
	7: CR7,
	8: CR8,
}

func sanitize(b []byte) []byte {
	l := len(b) - 2
	if l > 0 {
		return b[:l]
	}
	return b
}

func resetOriginals() {
	serialRead = read
	serialWrite = write
	serialFlush = flush
}

func stringInList(s string, list []string) bool {
	for _, l := range list {
		if l == s {
			return true
		}
	}
	return false
}
