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
	"strconv"
	"testing"
)

func TestMacResetWrongArgument(t *testing.T) {
	for i := uint16(0); i < maxUint16; i++ {
		if i != 433 && i != 868 && MacReset(i) != false {
			t.Errorf("MacReset(%v) returned true while the band <> 433 and 868", i)
			if testing.Short() {
				break
			}
		}
	}
}

func TestMacResetWriteError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return errors.New("Mock Write Error")
	}

	defer resetOriginals()

	if MacReset(433) != false {
		t.Errorf("MacReset(%v) returned true while the serial write failed", 433)
	}

	if MacReset(868) != false {
		t.Errorf("MacReset(%v) returned true while the serial write failed", 868)
	}
}

func TestMacResetReadError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		var b []byte
		return 0, b
	}

	defer resetOriginals()

	if MacReset(433) != false {
		t.Errorf("MacReset(%v) returned true while the serial read returned 0 bytes", 433)
	}

	if MacReset(868) != false {
		t.Errorf("MacReset(%v) returned true while the serial read returned 0 bytes", 868)
	}
}

func TestMacResetReadInvalidParam(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte(invalidParameter + "\r\n")
		return len(b), b
	}

	defer resetOriginals()

	if MacReset(433) != false {
		t.Errorf("MacReset(%v) returned true while the serial read returned invalid_param", 433)
	}

	if MacReset(868) != false {
		t.Errorf("MacReset(%v) returned true while the serial read returned invalid_param", 868)
	}
}

func TestMacResetSuccess(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte("ok\r\n")
		return len(b), b
	}

	defer resetOriginals()

	if MacReset(433) != true {
		t.Errorf("MacReset(%v) returned false while the serial read returned ok", 433)
	}

	if MacReset(868) != true {
		t.Errorf("MacReset(%v) returned false while the serial read returned ok", 868)
	}
}

func TestMacPauseWriteError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return errors.New("Mock Write Error")
	}

	defer resetOriginals()

	if MacPause() != 0 {
		t.Errorf("MacPause() returned non zero value while the serial write failed")
	}
}

func TestMacPauseReadError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		var b []byte
		return len(b), b
	}

	defer resetOriginals()

	if MacPause() != 0 {
		t.Errorf("MacPause() returned non zero value while the serial read failed")
	}
}

func TestMacPauseReadInvalidParamr(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte(invalidParameter + "\r\n")
		return len(b), b
	}

	defer resetOriginals()

	if MacPause() != 0 {
		t.Errorf("MacPause() returned non zero value while the serial read returned invalid_param")
	}
}

func TestMacPauseConversionError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte("nan\r\n")
		return len(b), b
	}

	defer resetOriginals()

	if MacPause() != 0 {
		t.Errorf("MacPause() returned non zero value while the conversion failed")
	}
}

func TestMacPauseSuccess(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte(strconv.Itoa(int(maxUint32)) + "\r\n")
		return len(b), b
	}

	defer resetOriginals()

	if MacPause() == 0 {
		t.Errorf("MacPause() returned zero value while the serial read succeeded")
	}
}

func TestMacResumeWriteError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return errors.New("Mock Write Error")
	}

	defer resetOriginals()

	if MacResume() != false {
		t.Errorf("MacResume() returned true value while the serial write failed")
	}
}

func TestMacResumeReadError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		var b []byte
		return len(b), b
	}

	defer resetOriginals()

	if MacResume() != false {
		t.Errorf("MacResume() returned true value while the serial read failed")
	}
}

func TestMacResumeSuccess(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte("ok\r\n")
		return len(b), b
	}

	defer resetOriginals()

	if MacResume() == false {
		t.Errorf("MacResume() returned false while the serial read succeeded")
	}
}

//func TestIsMacPausedFalse(t *testing.T) {
//	state.macPaused = false
//	length := 0
//
//	if isMacPaused(length) == true {
//		t.Errorf("isMacPaused(%v) returned true while the state is set to false", length)
//	}
//}
//
//func TestIsMacPausedFalseBecauseOffset(t *testing.T) {
//	state.macPaused = true
//	state.macPausedEnd = time.Now().Add(time.Duration(100) * time.Millisecond)
//	length := 100
//
//	if isMacPaused(length) == true {
//		t.Errorf("isMacPaused(%v) returned true wile the end time + offset is too close", length)
//	}
//}
//
//func TestIsMacPausedTrue(t *testing.T) {
//	state.macPaused = true
//	state.macPausedEnd = time.Now().Add(time.Duration(200) * time.Millisecond)
//	length := 100
//
//	if isMacPaused(length) == false {
//		t.Errorf("isMacPaused(%v) returned false wile the state is set to true", length)
//	}
//}

func TestMacJoinInvalidMode(t *testing.T) {
	mode := "test"
	if MacJoin(mode) == true {
		t.Errorf("MacJoin(%s) returned true with invalid mode", mode)
	}
}

func TestMacJoinAccepted(t *testing.T) {
	counter := 0

	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		var b = [][]byte{
			[]byte("ok\r\n"),
			[]byte("accepted\r\n"),
		}

		data := b[counter]

		counter = (counter + 1) % len(b)

		return len(data), data
	}

	defer resetOriginals()

	mode := OTAA

	if MacJoin(mode) == false {
		t.Errorf("MacTX(%v) returned false", mode)
	}

	//if state.joined != true {
	//	t.Error("Invalid joined state")
	//}
}

func TestMacJoinDenied(t *testing.T) {
	counter := 0

	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		var b = [][]byte{
			[]byte("ok\r\n"),
			[]byte("denied\r\n"),
		}

		data := b[counter]

		counter = (counter + 1) % len(b)

		return len(data), data
	}

	defer resetOriginals()

	mode := ABP

	if MacJoin(mode) == true {
		t.Errorf("MacTX(%v) returned true", mode)
	}

	//if state.joined != false {
	//	t.Error("Invalid joined state")
	//}
}

func TestMacTxInvalidPort(t *testing.T) {
	uplinkType := false
	ports := []uint8{0, 224}
	data := []byte("test")

	if MacTx(uplinkType, ports[0], data, nil) == true {
		t.Errorf("MacTX(%v, %v, %v, nil) returned true with invalid port", uplinkType, ports[0], data)
	}

	if MacTx(uplinkType, ports[1], data, nil) == true {
		t.Errorf("MacTX(%v, %v, %v, nil) returned true with invalid port", uplinkType, ports[1], data)
	}
}

func TestMacTxOk(t *testing.T) {
	counter := 0

	serialWrite = func(s string) error {
		t.Logf("String written to serial: %s", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		var b = [][]byte{
			[]byte("ok\r\n"),
			[]byte("mac_tx_ok\r\n"),
		}

		data := b[counter]

		counter = (counter + 1) % len(b)

		return len(data), data
	}

	defer resetOriginals()

	uplinkType := true
	port := uint8(1)
	data := []byte("test")

	if MacTx(uplinkType, port, data, nil) == false {
		t.Errorf("MacTX(%v, %v, %v, nil) returned false", uplinkType, port, data)
	}
}

func TestMacTxOkRx(t *testing.T) {
	counter := 0

	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		var b = [][]byte{
			[]byte("ok\r\n"),
			[]byte("mac_rx 1 7265636569766564\r\n"),
		}

		data := b[counter]

		counter = (counter + 1) % len(b)

		return len(data), data
	}

	defer resetOriginals()

	uplinkType := true
	port := uint8(1)
	data := []byte("test")
	callback := func(port uint8, data []byte) {
		t.Logf("Received message on port %v: %s", port, string(data))
	}

	if MacTx(uplinkType, port, data, callback) == false {
		t.Errorf("MacTX(%v, %v, %v, callback) returned false", uplinkType, port, data)
	}
}
