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
	"strings"
	"time"
	"github.com/pkg/errors"
	"encoding/hex"
)

type receiveCallback func(port uint8, data []byte)

// MacReset will automatically reset the software LoRaWAN stack and initilize
// it with the parameters for the selected band.
func MacReset(band uint16) bool {
	if band != 433 && band != 868 {
		WARN.Println("mac reset error: invalid band selected (433 or 868)")
		return false
	}

	err := serialWrite(fmt.Sprintf("mac reset %v", band))
	if err != nil {
		WARN.Println("mac reset error:", err)
		return false
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) == invalidParameter {
		WARN.Println("mac reset error: invalid parameter")
		return false
	}

	//state.macPaused = false

	return true
}

// MacPause will pause the LoRaWAN stack functionality to allow transceiver (radio) configuration.
// The length is the time in milliseconds the stack will be paused, with a maximum of 4294967295
// (max of uint32), is returned as an uint32.
func MacPause() uint32 {
	err := serialWrite("mac pause")
	if err != nil {
		WARN.Println("mac pause error:", err)
		return 0
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) == invalidParameter {
		WARN.Println("mac pause error: invalid parameter")
		return 0
	}

	value, err := strconv.ParseUint(string(sanitize(answer)), 10, 32)
	if err != nil {
		WARN.Println("mac pause error:", err)
		return 0
	}

	//state.macPaused = true
	//state.macPausedEnd = time.Now().Add(time.Duration(value) * time.Millisecond)

	return uint32(value)
}

// MacResume will resume the LoRaWAN stack functionality, in order to continue normal
// functionality after being paused.
func MacResume() bool {
	err := serialWrite("mac resume")
	if err != nil {
		WARN.Println("mac resume error:", err)
		return false
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) == invalidParameter {
		WARN.Println("mac resume error: invalid parameter")
		return false
	}

	//state.macPaused = false

	return true
}

// The length is passed in milliseconds.
//func isMacPaused(length int) bool {
//	// An offset of 100 milliseconds is added to ensure we have enough
//	// time left over
//	d := time.Duration(length)*time.Millisecond + time.Duration(100)
//
//	if state.macPaused {
//		return time.Now().Add(d).Before(state.macPausedEnd)
//	}
//
//	return false
//}

// MacJoin will join the configured network with the given mode.
func MacJoin(mode string) bool {
	if mode != OTAA && mode != ABP {
		WARN.Println("mac join error: invalid mode (OTAA or ABP)")
		return false
	}

	err := serialWrite(fmt.Sprintf("mac join %s", mode))
	if err != nil {
		WARN.Println("mac join error:", err)
		return false
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) != "ok" {
		WARN.Println("mac join error:", string(sanitize(answer)))
		return false
	}

	tick := time.Tick(time.Second)
	timeout := time.After(time.Second * 15)

	for {
		select {
		case <- timeout:
			return false
		case <- tick:
			n, answer = serialRead()

			if n != 0 {
				if string(sanitize(answer)) != "accepted" {
					WARN.Println("mac join error:", string(sanitize(answer)))
					return false
				}

				return true
			}
		}
	}
}

// MacTX will transmit the given data on the given port. The transmission can
// either be confirmed (if the boolean is set), meaning that the server will
// response with an acknowledgement. If no acknowledgement is received, the
// message will be retransmitted by the number indicated by the MacSetRetx
// command. The port number has to be a value in the range of [1,223].
// The receiveCallback function passed is responsible to handle the received
// answer from the server. If no answers are expected, nil can be passed as the
// callback argument.
func MacTx(confirmed bool, port uint8, data []byte, callback receiveCallback) bool {
	if port < 1 || port > 223 {
		WARN.Printf("mac tx error: invalid port number (%v)", port)
		return false
	}

	if len(data) == 0 {
		WARN.Println("mac tx error: trying to send zero bytes")
		return false
	}

	uplinkType := UNCONFIRMED

	if confirmed {
		uplinkType = CONFIRMED
	}

	err := serialWrite(fmt.Sprintf("mac tx %s %v %X", uplinkType, port, data))
	if err != nil {
		WARN.Println("mac tx error:", err)
		return false
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) != "ok" {
		WARN.Println("mac tx error:", string(sanitize(answer)))
		return false
	}

	timeout := time.After(time.Second * 15)
	tick := time.Tick(time.Second)

	for {
		select {
		case <- timeout:
			WARN.Println("timed out")
			return false
		case <- tick:
			n, answer = serialRead()
			s := string(sanitize(answer))

			if n != 0 {
				if s == "mac_err" || s == "invalid_data_len" {
					WARN.Printf("mac tx error: %s", s)
					return false
				} else if s == "mac_tx_ok" {
					return true
				} else if strings.HasPrefix(s, "mac_rx") {
					if callback != nil {
						params := strings.Split(s, " ")

						port, err := strconv.ParseInt(params[1], 10, 8)
						if err != nil {
							WARN.Printf("mac_rx invalid port: %s", params[1])
							return true
						}

						decoded, err := hex.DecodeString(params[2])
						if err != nil {
							WARN.Printf("mac_rx invalid hex data: %s", params[2])
							return true
						}

						callback(uint8(port), []byte(decoded))
					}
					return true
				} else {
					return false
				}
			}
		}
	}
}

// MacGetDeviceAddress will return the current end device address of the module.
// The address is represented as a 4-byte hexadecimal number and returned as a string.
// The default value of 00000000 will be returned in case of an error.
func MacGetDeviceAddress() string {
	err := serialWrite("mac get devaddr")
	if err != nil {
		WARN.Println("mac get devaddr error:", err)
		return "00000000"
	}

	n, answer := serialRead()
	if n != 0 {
		return string(sanitize(answer))
	}

	return "00000000"
}

// MacSetDeviceAddress will configure the module with a network device address.
// The address is a 4-byte hexadecimal value given as a string.
func MacSetDeviceAddress(address string) error {
	if len(address) != 8 {
		return errors.New("invalid address length")
	}

	err := serialWrite(fmt.Sprintf("mac set devaddr %s", address))
	if err != nil {
		return errors.Wrap(err, "could not set device address")
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) == invalidParameter {
		return errors.New("could not set device address: invalid parameter")
	}

	return nil
}

// MacGetDeviceEUI will return the current end device EUI of the module.
// The EUI is represented as a 8-byte hexadecimal number and returned as a string.
// The default value of 0000000000000000 will be returned in case of an error.
func MacGetDeviceEUI() string {
	err := serialWrite("mac get deveui")
	if err != nil {
		WARN.Println("mac get deveui error:", err)
		return "0000000000000000"
	}

	n, answer := serialRead()
	if n != 0 {
		return string(sanitize(answer))
	}

	return "0000000000000000"
}

// MacSetDeviceEUI will configure the module with a network device EUI.
// The EUI is a 8-byte hexadecimal value given as a string.
func MacSetDeviceEUI(eui string) error {
	if len(eui) != 16 {
		return errors.New("invalid eui length")
	}

	err := serialWrite(fmt.Sprintf("mac set deveui %s", eui))
	if err != nil {
		return errors.Wrap(err, "could not set device eui")
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) == invalidParameter {
		return errors.New("could not set device eui: invalid parameter")
	}

	return nil
}

// MacGetApplicationEUI will return the current configured application EUI.
// The EUI is represented as a 8-byte hexadecimal number and returned as a string.
// The default value of 0000000000000000 will be returned in case of an error.
func MacGetApplicationEUI() string {
	err := serialWrite("mac get appeui")
	if err != nil {
		WARN.Println("mac get appeui error:", err)
		return "0000000000000000"
	}

	n, answer := serialRead()
	if n != 0 {
		return string(sanitize(answer))
	}

	return "0000000000000000"
}

// MacSetApplicationEUI will configure the module with a network application EUI.
// The EUI is a 8-byte hexadecimal value given as a string.
func MacSetApplicationEUI(eui string) error {
	if len(eui) != 16 {
		return errors.New("invalid eui length")
	}

	err := serialWrite(fmt.Sprintf("mac set appeui %s", eui))
	if err != nil {
		return errors.Wrap(err, "could not set application eui")
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) == invalidParameter {
		return errors.New("could not set application eui: invalid parameter")
	}

	return nil
}

// MacSetNetworkSessionKey will configure the module with a network session key.
// The key is a 16-byte hexadecimal value given as a string.
func MacSetNetworkSessionKey(key string) error {
	if len(key) != 32 {
		return errors.New("invalid key length")
	}

	err := serialWrite(fmt.Sprintf("mac set nwkskey %s", key))
	if err != nil {
		return errors.Wrap(err, "could not set network session key")
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) == invalidParameter {
		return errors.New("could not set network session key: invalid parameter")
	}

	return nil
}

// MacSetApplicationSessionKey will configure the module with an application session key.
// The key is a 16-byte hexadecimal value given as a string.
func MacSetApplicationSessionKey(key string) error {
	if len(key) != 32 {
		return errors.New("invalid key length")
	}

	err := serialWrite(fmt.Sprintf("mac set appskey %s", key))
	if err != nil {
		return errors.Wrap(err, "could not set application session key")
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) == invalidParameter {
		return errors.New("could not set application session key: invalid parameter")
	}

	return nil
}

// MacSetApplicationKey will configure the module with an application key.
// The key is a 16-byte hexadecimal value given as a string.
func MacSetApplicationKey(key string) error {
	if len(key) != 32 {
		return errors.New("invalid key length")
	}

	err := serialWrite(fmt.Sprintf("mac set appkey %s", key))
	if err != nil {
		return errors.Wrap(err, "could not set application key")
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) == invalidParameter {
		return errors.New("could not set application key: invalid parameter")
	}

	return nil
}

// MacGetDataRate will return the current data rate.
// The data rate is a number in the range of [0-5],
// with 0 = SF12BW125 and 5 = SF7BW125.
func MacGetDataRate() uint8 {
	err := serialWrite("mac get dr")
	if err != nil {
		WARN.Println("mac get dr error:", err)
		return 0
	}

	n, answer := serialRead()
	if n != 0 {
		dr, err := strconv.ParseUint(string(sanitize(answer)), 10, 8)
		if err != nil {
			WARN.Println("mac get dr uint conversion error:", err)
			return 0
		}

		return uint8(dr)
	}

	return 0
}

// MacSetDataRate will configure the data rate for the next transmission.
// The data rate has to be in the range of [0-5],
// with 0 = SF12BW125 and 5 = SF7BW125.
func MacSetDataRate(dr uint8) error {
	if dr > 5 {
		return errors.New("invalid data rate")
	}

	err := serialWrite(fmt.Sprintf("mac set dr %v", dr))
	if err != nil {
		return errors.Wrap(err, "could not set data rate")
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) == invalidParameter {
		return errors.New("could not set data rate: invalid parameter")
	}

	return nil
}

// MacGetPowerIndex will return the current power index.
// The power index is a number in the range of [0-5],
// with 0 = 20 dBm (if available), 1 = 14 dBm, 2 = 11 dBm,
// 3 = 8 dBm, 4 = 5dBm and 5 = 2 dBm.
func MacGetPowerIndex() uint8 {
	err := serialWrite("mac get pwridx")
	if err != nil {
		WARN.Println("mac get pwridx error:", err)
		return 1
	}

	n, answer := serialRead()
	if n != 0 {
		pwr, err := strconv.ParseUint(string(sanitize(answer)), 10, 8)
		if err != nil {
			WARN.Println("mac get pwridx uint conversion error:", err)
			return 1
		}

		return uint8(pwr)
	}

	return 1
}

// MacSetPowerIndex will configure the power index for the next transmission.
// The index has to be in the range of [1-5] for 868 MHz and [0-5] for 433 MHz.
func MacSetPowerIndex(index uint8) error {
	if index > 5 {
		return errors.New("invalid power index")
	}

	err := serialWrite(fmt.Sprintf("mac set pwridx %v", index))
	if err != nil {
		return errors.Wrap(err, "could not set power index")
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) == invalidParameter {
		return errors.New("could not set power index: invalid parameter")
	}

	return nil
}

// MacGetADR will return the state of the adpative data rate mechanism.
func MacGetADR() bool {
	err := serialWrite("mac get adr")
	if err != nil {
		WARN.Println("mac get adr error:", err)
		return false
	}

	n, answer := serialRead()
	if n == 0 {
		WARN.Println("mac get adr error: no answer")
		return false
	}

	if string(sanitize(answer)) == "on" {
		return true
	}

	return false
}

// MacSetADR will set the adaptive data rate.
func MacSetADR(adr bool) error {
	var state = "off"

	if adr {
		state = "on"
	}

	err := serialWrite(fmt.Sprintf("mac set adr %s", state))
	if err != nil {
		return errors.Wrap(err, "could not set adaptive data rate")
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) == invalidParameter {
		return errors.New("could not set adaptive data rate: invalid parameter")
	}

	return nil
}

// MacSetLinkCheck will set the time interval for the link check process to be triggered.
func MacSetLinkCheck(interval uint16) error {
	err := serialWrite(fmt.Sprintf("mac set linkchk %v", interval))
	if err != nil {
		return errors.Wrap(err, "could not set link check")
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) == invalidParameter {
		return errors.New("could not set link check: invalid parameter")
	}

	return nil
}

//TODO: implement mac get status

// MacGetChannelFrequency will return the frequency on the requested channelID.
// This frequency is returned in Hz.
// The channelID has to be in the range of [0-15].
func MacGetChannelFrequency(channelID uint8) uint32 {
	if channelID > 15 {
		WARN.Println("mac get ch freq error: invalid channel")
		return 0
	}

	err := serialWrite(fmt.Sprintf("mac get ch freq %v", channelID))
	if err != nil {
		WARN.Println("mac get ch freq error:", err)
		return 0
	}

	n, answer := serialRead()
	if n != 0 {
		value, err := strconv.ParseUint(string(sanitize(answer)), 10, 32)
		if err != nil {
			WARN.Println("mac get ch freq uint conversion error:", err)
			return 0
		}

		return uint32(value)
	}

	return 0
}

// MacSetChannelFrequency will set the frequency on the given channel id.
// The default channels (0-2) cannot be modified.
// The applicable range for the channel id is [3-15].
// The frequency has to be given in Hz.
func MacSetChannelFrequency(channelID uint8, frequency uint32) error {
	if channelID < 3 || channelID > 15 {
		return errors.New("invalid channel id")
	}

	if frequency < 433050000 || (frequency > 434790000 && frequency < 863000000) || frequency > 870000000 {
		return errors.New("invalid frequency")
	}

	err := serialWrite(fmt.Sprintf("mac set ch freq %v %v", channelID, frequency))
	if err != nil {
		return errors.Wrap(err, "could not set channel frequency")
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) == invalidParameter {
		return errors.New("could not set channel frequency: invalid parameter")
	}

	return nil
}

// MacGetChannelDutyCycle will return the duty cycle on the requested channelID.
// The duty cycle will be returned as a percentage.
// The channelID has to be in the range of [0-15].
func MacGetChannelDutyCycle(channelID uint8) float32 {
	if channelID > 15 {
		WARN.Println("mac get ch dcycle error: invalid channel")
		return 0
	}

	err := serialWrite(fmt.Sprintf("mac get ch dcycle %v", channelID))
	if err != nil {
		WARN.Println("mac get ch dcycle error:", err)
		return 0
	}

	n, answer := serialRead()
	if n != 0 {
		value, err := strconv.ParseUint(string(sanitize(answer)), 10, 16)
		if err != nil {
			WARN.Println("mac get ch dcycle uint conversion error:", err)
			return 0
		}

		return 100 / float32(value+1)
	}

	return 0
}

// MacSetChannelDutyCycle will set the duty cycle used on the given channel id.
// The applicable range for the channel id is [0-15].
// The duty cycle can be given as a percentage.
func MacSetChannelDutyCycle(channelID uint8, dcycle float32) error {
	if channelID > 15 {
		return errors.New("invalid channel id")
	}

	value := uint64((100 / dcycle) - 1)

	if value > uint64(^uint16(0)) {
		value = uint64(^uint16(0))
	}

	err := serialWrite(fmt.Sprintf("mac set ch dcycle %v %v", channelID, uint16(value)))
	if err != nil {
		return errors.Wrap(err, "could not set channel duty cycle")
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) == invalidParameter {
		return errors.New("could not set channel duty cycle: invalid parameter")
	}

	return nil
}

// MacGetChannelStatus will return if the given channelID is currently enabled for use.
// The channelID has to be in the range of [0-15].
func MacGetChannelStatus(channelID uint8) bool {
	if channelID > 15 {
		WARN.Println("mac get ch status error: invalid channel")
		return false
	}

	err := serialWrite(fmt.Sprintf("mac get ch status %v", channelID))
	if err != nil {
		WARN.Println("mac get ch status error:", err)
		return false
	}

	n, answer := serialRead()
	if n == 0 {
		WARN.Println("mac get ch status error: no answer")
		return false
	}

	if string(sanitize(answer)) == "on" {
		return true
	}

	return false
}

// MacSetChannelStatus will set the operation on the given channel id.
// The applicable range for the channel id is [0-15].
func MacSetChannelStatus(channelID uint8, status bool) error {
	var state = "off"

	if status {
		state = "on"
	}

	if channelID > 15 {
		return errors.New("invalid channel id")
	}

	err := serialWrite(fmt.Sprintf("mac set ch dcycle %v %s", channelID, state))
	if err != nil {
		return errors.Wrap(err, "could not set channel status")
	}

	n, answer := serialRead()
	if n == 0 || string(sanitize(answer)) == invalidParameter {
		return errors.New("could not set channel status: invalid parameter")
	}

	return nil
}
