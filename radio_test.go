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
	"testing"
)

func TestRadioRxBlockingWriteError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return errors.New("Mock Write Error")
	}

	defer resetOriginals()

	if len(RadioRxBlocking(0)) > 0 {
		t.Errorf("RadioRxBlocking(%v) returned bytes while the serial write failed", 0)
	}
}

func TestRadioRxBlockingReadError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		var b []byte
		return 0, b
	}

	defer resetOriginals()

	if len(RadioRxBlocking(0)) > 0 {
		t.Errorf("RadioRxBlocking(%v) returned bytes while the serial read returned 0 bytes", 0)
	}
}

func TestRadioRxBlockingReadInvalidParam(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte(invalidParameter + "\r\n")
		return len(b), b
	}

	defer resetOriginals()

	if len(RadioRxBlocking(0)) > 0 {
		t.Errorf("RadioRxBlocking(%v) returned bytes while the serial read returned invalid_param", 0)
	}
}

func TestRadioRxBlockingReadBusy(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte("busy\r\n")
		return len(b), b
	}

	defer resetOriginals()

	if len(RadioRxBlocking(0)) > 0 {
		t.Errorf("RadioRxBlocking(%v) returned bytes while the serial read returned busy", 0)
	}
}

func TestRadioRxBlockingFailure(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte("radio_err\r\n")
		return len(b), b
	}

	defer resetOriginals()

	if len(RadioRxBlocking(0)) > 0 {
		t.Errorf("RadioRxBlocking(%v) returned bytes while the serial read returned radio_rx ...", 0)
	}
}

func TestRadioRxBlockingSuccess(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte("radio_rx  5376656E\r\n")
		return len(b), b
	}

	defer resetOriginals()

	if len(RadioRxBlocking(0)) == 0 {
		t.Errorf("RadioRxBlocking(%v) returned 0 bytes while the serial read returned radio_rx ...", 0)
	}
}

func TestRadioTxEmptyData(t *testing.T) {
	var data []byte
	if RadioTx(data) == true {
		t.Errorf("RadioTx(%v) returned true while the data is empty", data)
	}
}

func TestRadioTxWriteError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return errors.New("Mock Write Error")
	}

	defer resetOriginals()

	data := []byte("test")
	if RadioTx(data) == true {
		t.Errorf("RadioTx(%v) returned true while the serial write failed", data)
	}
}

func TestRadioTxReadError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		var b []byte
		return 0, b
	}

	defer resetOriginals()

	data := []byte("test")
	if RadioTx(data) == true {
		t.Errorf(`RadioTx(%v) returned true while the serial read returned
0 bytes or something else than ok`, data)
	}
}

func TestRadioTxTimedOut(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping time-out test in short mode")
	}

	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	run := false

	serialRead = func() (int, []byte) {
		var b = [][]byte{
			[]byte("ok\r\n"),
			[]byte("test\r\n"),
		}

		if run {
			return len(b[1]), b[1]
		}

		run = true
		return len(b[0]), b[0]
	}

	defer resetOriginals()

	data := []byte("test")
	if RadioTx(data) == true {
		t.Errorf(`RadioTx(%v) returned true while the serial read returned
ok and timed out afterwards`, data)
	}
}

func TestRadioTxRadioErr(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	run := false

	serialRead = func() (int, []byte) {
		var b = [][]byte{
			[]byte("ok\r\n"),
			[]byte("radio_err\r\n"),
		}

		if run {
			return len(b[1]), b[1]
		}

		run = true
		return len(b[0]), b[0]
	}

	defer resetOriginals()

	data := []byte("test")
	if RadioTx(data) == true {
		t.Errorf(`RadioTx(%v) returned true while the serial read returned
ok and radio_err afterwards`, data)
	}
}

func TestRadioTxOk(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	run := false

	serialRead = func() (int, []byte) {
		var b = [][]byte{
			[]byte("ok\r\n"),
			[]byte("radio_tx_ok\r\n"),
		}

		if run {
			return len(b[1]), b[1]
		}

		run = true
		return len(b[0]), b[0]
	}

	defer resetOriginals()

	data := []byte("test")
	if RadioTx(data) == false {
		t.Errorf(`RadioTx(%v) returned false while the serial read returned
ok and radio_tx_ok afterwards`, data)
	}
}

func TestRadioGetModulationWriteError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return errors.New("Mock Write Error")
	}

	defer resetOriginals()

	x := RadioGetModulation()
	if x != "" {
		t.Errorf("RadioGetModulation() returned non empty string (%v) while serial write failed", x)
	}
}

func TestRadioGetModulationReadError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		var b []byte
		return 0, b
	}

	defer resetOriginals()

	x := RadioGetModulation()
	if x != "" {
		t.Errorf("RadioGetModulation() returned non empty string (%v) while serial read failed", x)
	}
}

func TestRadioGetModulationSuccess(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte("lora")
		return len(b), b
	}

	defer resetOriginals()

	if RadioGetModulation() == "" {
		t.Error("RadioGetModulation() returned empty string when it should succeed")
	}
}

func TestRadioSetModulationArgumentError(t *testing.T) {
	mod := "nonse"
	if RadioSetModulation(mod) == true {
		t.Errorf("RadioSetModulation(%v) returned true while the modulation doesn't exist", mod)
	}
}

func TestRadioSetModulationWriteError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return errors.New("Mock Write Error")
	}

	defer resetOriginals()

	mod := LoRa
	if RadioSetModulation(mod) == true {
		t.Errorf("RadioSetModulation(%v) returned true while the serial write failed", mod)
	}
}

func TestRadioModulationReadError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		var b []byte
		return 0, b
	}

	defer resetOriginals()

	mod := LoRa
	if RadioSetModulation(mod) == true {
		t.Errorf(`RadioSetModulation(%v) returned true while the serial read returned
0 bytes or something else than ok`, mod)
	}
}

func TestRadioModulationSuccess(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte("ok\r\n")
		return len(b), b
	}

	defer resetOriginals()

	mod := LoRa
	if RadioSetModulation(mod) == false {
		t.Errorf(`RadioSetModulation(%v) returned false while the serial read returned ok`, mod)
	}
}

func TestRadioGetFrequencyWriteError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return errors.New("Mock Write Error")
	}

	defer resetOriginals()

	if RadioGetFrequency() != 0 {
		t.Error("RadioGetFrequency() returned non zero while serial write failed")
	}
}

func TestRadioGetFrequencyReadError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		var b []byte
		return 0, b
	}

	defer resetOriginals()

	if RadioGetFrequency() != 0 {
		t.Error("RadioGetFrequency() returned non zero while serial read failed")
	}
}

func TestRadioGetFrequencyConversionError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte("nan\r\n")
		return len(b), b
	}

	defer resetOriginals()

	if RadioGetFrequency() != 0 {
		t.Error("RadioGetFrequency() returned non zero while conversion failed")
	}
}

func TestRadioGetFrequencySuccess(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte("863000000\r\n")
		return len(b), b
	}

	defer resetOriginals()

	if RadioGetFrequency() != uint32(863000000) {
		t.Error("RadioGetFrequency() returned wrong value while it should succeed")
	}
}

func TestRadioSetFrequencyArgumentError(t *testing.T) {
	freq := uint32(100000000)
	if RadioSetFrequency(freq) != false {
		t.Errorf("RadioSetFrequency(%v) returned true with the wrong frequency", freq)
	}
}

func TestRadioSetFrequencyWriteError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return errors.New("Mock Write Error")
	}

	defer resetOriginals()

	freq := uint32(863000000)
	if RadioSetFrequency(freq) != false {
		t.Errorf("RadioSetFrequency(%v) returned true while serial write failed", freq)
	}
}

func TestRadioSetFrequencyReadError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		var b []byte
		return 0, b
	}

	defer resetOriginals()

	freq := uint32(863000000)
	if RadioSetFrequency(freq) != false {
		t.Errorf("RadioSetFrequency(%v) returned true while serial read failed", freq)
	}
}

func TestRadioSetFrequencySuccess(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte("ok\r\n")
		return len(b), b
	}

	defer resetOriginals()

	freq := uint32(863000000)
	if RadioSetFrequency(freq) != true {
		t.Errorf("RadioSetFrequency(%v) returned false while it should succeed", freq)
	}
}

func TestRadioGetPowerWriteError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return errors.New("Mock Write Error")
	}

	defer resetOriginals()

	if RadioGetPower() != -15 {
		t.Error("RadioGetPower() returned value other than -15 while serial write failed")
	}
}

func TestRadioGetPowerReadError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		var b []byte
		return 0, b
	}

	defer resetOriginals()

	if RadioGetPower() != -15 {
		t.Error("RadioGetPower() returned value other than -15 while serial read failed")
	}
}

func TestRadioGetPowerConversionError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte("nan\r\n")
		return len(b), b
	}

	defer resetOriginals()

	if RadioGetPower() != -15 {
		t.Error("RadioGetPower() returned value other than -15 while conversion failed")
	}
}

func TestRadioGetPowerSuccess(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte("14\r\n")
		return len(b), b
	}

	defer resetOriginals()

	if RadioGetPower() != 14 {
		t.Error("RadioGetPower() returned value other than 14 which it should")
	}
}

func TestRadioSetPowerArgumentError(t *testing.T) {
	for i := int8(-125); i < -3; i++ {
		pwr := i
		if RadioSetPower(pwr) != false {
			t.Errorf("RadioSetPower(%v) returned true with the wrong power", pwr)
		}
	}

	for i := int8(16); i < maxInt8; i++ {
		pwr := i
		if RadioSetPower(pwr) != false {
			t.Errorf("RadioSetPower(%v) returned true with the wrong power", pwr)
		}
	}
}

func TestRadioSetPowerWriteError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return errors.New("Mock Write Error")
	}

	defer resetOriginals()

	pwr := int8(14)
	if RadioSetPower(pwr) != false {
		t.Errorf("RadioSetPower(%v) returned true while serial write failed", pwr)
	}
}

func TestRadioSetPowerReadError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		var b []byte
		return 0, b
	}

	defer resetOriginals()

	pwr := int8(14)
	if RadioSetPower(pwr) != false {
		t.Errorf("RadioSetPower(%v) returned true while serial read failed", pwr)
	}
}

func TestRadioSetPowerSuccess(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte("ok\r\n")
		return len(b), b
	}

	defer resetOriginals()

	pwr := int8(14)
	if RadioSetPower(pwr) != true {
		t.Errorf("RadioSetPower(%v) returned false while it should succeed", pwr)
	}
}

func TestRadioGetSFWriteError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return errors.New("Mock Write Error")
	}

	defer resetOriginals()

	if RadioGetSpreadingFactor() != 0 {
		t.Error("RadioGetSpreadingFactor() returned non zero value while serial write failed")
	}
}

func TestRadioGetSFReadError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		var b []byte
		return 0, b
	}

	defer resetOriginals()

	if RadioGetSpreadingFactor() != 0 {
		t.Error("RadioGetSpreadingFactor() returned non zero value while serial read failed")
	}
}

func TestRadioGetSFConversionError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte("nan\r\n")
		return len(b), b
	}

	defer resetOriginals()

	if RadioGetSpreadingFactor() != 0 {
		t.Error("RadioGetSpreadingFactor() returned non zero value while conversion failed")
	}
}

func TestRadioGetSFSuccess(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte("SF12\r\n")
		return len(b), b
	}

	defer resetOriginals()

	if RadioGetSpreadingFactor() == 0 {
		t.Error("RadioGetSpreadingFactor() returned zero value while it should succeed")
	}
}

func TestRadioSetSFArgumentError(t *testing.T) {
	for i := uint8(0); i < 7; i++ {
		sf := i
		if RadioSetSpreadingFactor(sf) != false {
			t.Errorf("RadioSetSpreadingFactor(%v) returned true with the wrong sf", sf)
		}
	}

	for i := uint8(13); i < maxUint8; i++ {
		sf := i
		if RadioSetSpreadingFactor(sf) != false {
			t.Errorf("RadioSetSpreadingFactor(%v) returned true with the wrong sf", sf)
		}
	}
}

func TestRadioSetSFWriteError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return errors.New("Mock Write Error")
	}

	defer resetOriginals()

	sf := uint8(7)
	if RadioSetSpreadingFactor(sf) != false {
		t.Errorf("RadioSetSpreadingFactor(%v) returned true while serial write failed", sf)
	}
}

func TestRadioSetSFReadError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		var b []byte
		return 0, b
	}

	defer resetOriginals()

	sf := uint8(7)
	if RadioSetSpreadingFactor(sf) != false {
		t.Errorf("RadioSetSpreadingFactor(%v) returned true while serial read failed", sf)
	}
}

func TestRadioSetSFSuccess(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte("ok\r\n")
		return len(b), b
	}

	defer resetOriginals()

	sf := uint8(7)
	if RadioSetSpreadingFactor(sf) != true {
		t.Errorf("RadioSetSpreadingFactor(%v) returned false while it should succeed", sf)
	}
}
