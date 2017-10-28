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

// RadioRxBlocking will open the receiver.
// The window size is the number of symbols for LoRa modulation and the
// time in milliseconds for FSK modulation. In order to enable continuous
// reception, the window size should be 0. Don't forget to set the radio
// watchdog timer time-out! This function will return a valid packet that has
// been received, or an empty array of bytes when the receiver was busy or it
// timed out without receiving a valid packet. This function is blocking, which
// means if you enabled continous reception, it will block the program until a
// valid packet has been received or until a time out occured.
func RadioRxBlocking(window uint16) []byte {
	var b []byte

	// TODO Should get wdt to get the length
	// if !isMacPaused(length)

	err := serialWrite(fmt.Sprintf("radio rx %v", window))
	if err != nil {
		logger.Warning.Println("radio rx error:", err)
		return b
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) == invalidParameter ||
		string(sanitize(answer)) == "busy" {
		logger.Warning.Println("radio rx error: busy or invalid parameter")
		return b
	}

	for {
		n, answer := serialRead()
		if n != 0 && string(sanitize(answer)) == "radio_err" {
			return b
		}

		if n != 0 && string(answer[:8]) == "radio_rx" {
			return answer[10:]
		}
	}
}

// RadioTx will transmit the given data. The data has to have a length > 0
// but has to be smaller than 255 if LoRa modulation is active or smaller
// than 64 if FSK modulation is active. It will return a boolean, true if
// the transmit was succesful, false is there was an error. For more info
// about the error, the user can check the log file.
func RadioTx(data []byte) bool {
	//TODO check modulation to get maximum bytes allowed: 255 LoRa and 64 FSK
	if len(data) == 0 {
		logger.Warning.Println("radio tx error: trying to send zero bytes")
		return false
	}

	// TODO check air time to check isMacPaused

	err := serialWrite(fmt.Sprintf("radio tx %X", data))
	if err != nil {
		logger.Warning.Println("radio tx error:", err)
		return false
	}

	n, firstAnswer := serialRead()
	if n == 0 || string(sanitize(firstAnswer)) != "ok" {
		logger.Warning.Println("radio tx error:", string(sanitize(firstAnswer)))
		return false
	}

	timeout := time.After(time.Second * 5)

	for {
		select {
		case <-timeout:
			return false
		default:
			n, answer := serialRead()
			if n != 0 && string(sanitize(answer)) == "radio_err" {
				logger.Warning.Println("radio tx error: radio_err")
				return false
			}

			if n != 0 && string(sanitize(answer)) == "radio_tx_ok" {
				return true
			}
		}
	}
}

// RadioGetModulation reads back the current mode of operation of the module.
// It returns an empty string if something went wrong.
func RadioGetModulation() string {
	err := serialWrite("radio get mod")
	if err != nil {
		logger.Warning.Println("radio get mod error:", err)
		return ""
	}

	n, answer := serialRead()
	if n == 0 {
		logger.Warning.Println("radio get mod error: no answer")
		return ""
	}

	return string(sanitize(answer))
}

// RadioSetModulation changes the modulation method being used by the module.
// The modulations are available as constants in the package.
// The function will return true when the change is accepted by the module.
// When the change isn't accepted or the modulation is wrong, it will return false.
func RadioSetModulation(mod string) bool {
	if !stringInList(mod, modulations) {
		logger.Warning.Println("radio set mod error: invalid modulation")
		return false
	}

	err := serialWrite(fmt.Sprintf("radio set mod %s", mod))
	if err != nil {
		logger.Warning.Println("radio set mod error:", err)
		return false
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) != "ok" {
		logger.Warning.Println("radio set mod:", string(sanitize(answer)))
		return false
	}

	return true
}

// RadioGetFrequency returns the current operation frequency of the module.
// If there was an error, the function will return 0.
func RadioGetFrequency() uint32 {
	err := serialWrite("radio get freq")
	if err != nil {
		logger.Warning.Println("radio get freq error:", err)
		return 0
	}

	n, answer := serialRead()
	if n == 0 {
		logger.Warning.Println("radio get freq error: no answer")
		return 0
	}

	value, err := strconv.ParseUint(string(sanitize(answer)), 10, 32)
	if err != nil {
		logger.Warning.Println("radio get freq error:", err)
		return 0
	}

	return uint32(value)
}

// RadioSetFrequency changes the communication frequency of the radio transceiver.
// It will only accept frequencies between [433050000, 434790000] and [863000000, 870000000].
// The function will return true when the frequency changed and false when an error occured.
func RadioSetFrequency(freq uint32) bool {
	if (freq < 433050000 || freq > 434790000) && (freq < 863000000 || freq > 870000000) {
		logger.Warning.Println("radio set freq error: invalid frequency", freq)
		return false
	}

	err := serialWrite(fmt.Sprintf("radio set freq %v", freq))
	if err != nil {
		logger.Warning.Println("radio set freq error:", err)
		return false
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) != "ok" {
		logger.Warning.Println("radio set freq error: invalid parameter")
		return false
	}

	return true
}

// RadioGetPower reads back the current power level setting used in operation.
// The function will return an int8 value, which will be between [-3, 15].
// If an error occured, it will return -15.
func RadioGetPower() int8 {
	err := serialWrite("radio get pwr")
	if err != nil {
		logger.Warning.Println("radio get pwr error:", err)
		return -15
	}

	n, answer := serialRead()
	if n == 0 {
		logger.Warning.Println("radio get pwr error: no answer")
		return -15
	}

	value, err := strconv.ParseInt(string(sanitize(answer)), 10, 8)
	if err != nil {
		logger.Warning.Println("radio get pwr error:", err)
		return -15
	}

	return int8(value)
}

// RadioSetPower changes the transceiver output power.
// The output power has to be passed as an int8 value between [-3, 15].
// The function will return true if the change succeeeded, or false when
// an error occured.
func RadioSetPower(pwr int8) bool {
	if pwr < -3 || pwr > 15 {
		logger.Warning.Println("radio set pwr error: invalid power", pwr)
		return false
	}

	err := serialWrite(fmt.Sprintf("radio set pwr %v", pwr))
	if err != nil {
		logger.Warning.Println("radio set pwr error:", err)
		return false
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) != "ok" {
		logger.Warning.Println("radio set pwr error: invalid parameter")
		return false
	}

	return true
}

// RadioGetSpreadingFactor reads back the current spreading factor
// being used by the transceiver.
// It will return an int8 between [7, 12].
// If an error occured, it will return 0.
func RadioGetSpreadingFactor() uint8 {
	err := serialWrite("radio get sf")
	if err != nil {
		logger.Warning.Println("radio get sf error:", err)
		return 0
	}

	n, answer := serialRead()
	if n == 0 {
		logger.Warning.Println("radio get sf error: no answer")
		return 0
	}

	value, err := strconv.ParseUint(string(sanitize(answer[2:])), 10, 8)
	if err != nil {
		logger.Warning.Println("radio get sf error:", err)
		return 0
	}

	return uint8(value)
}

// RadioSetSpreadingFactor sets the spreading factor used during transmission.
// The spreading factor has to be passed as an uint8 between [7, 12].
// The function will return true if the command succeeded.
// If an error occured, it will return false.
func RadioSetSpreadingFactor(sf uint8) bool {
	if sf < 7 || sf > 12 {
		logger.Warning.Println("radio set sf error: invalid spreading factor", sf)
		return false
	}

	err := serialWrite(fmt.Sprintf("radio set sf %v", sfs[sf]))
	if err != nil {
		logger.Warning.Println("radio set sf error:", err)
		return false
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) != "ok" {
		logger.Warning.Println("radio set sf error: invalid parameter")
		return false
	}

	return true
}
