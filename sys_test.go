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

func TestSleepWrongArgument(t *testing.T) {
	for i := uint32(0); i < 100; i++ {
		if Sleep(i) != false {
			t.Errorf("Sleep(%v) returned true while the length < 100", i)
			if testing.Short() {
				break
			}
		}
	}
}

func TestSleepWriteError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return errors.New("Mock Write Error")
	}

	defer resetOriginals()

	for i := uint32(100); i <= 1000100; i += 100000 {
		if Sleep(i) != false {
			t.Errorf("Sleep(%v) returned true while the serial write failed", i)
			if testing.Short() {
				if testing.Short() {
					break
				}
			}
		}
	}
}

func TestSleepReadError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		var b []byte
		return 0, b
	}

	defer resetOriginals()

	for i := uint32(100); i <= 1000100; i += 100000 {
		if Sleep(i) != false {
			t.Errorf("Sleep(%v) returned true while the serial read returned 0 bytes", i)
			if testing.Short() {
				break
			}
		}
	}
}

func TestSleepReadInvalidParam(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte(invalidParameter + "\r\n")
		return len(b), b
	}

	defer resetOriginals()

	for i := uint32(100); i <= 1000100; i += 100000 {
		if Sleep(i) != false {
			t.Errorf("Sleep(%v) returned true while the serial read returned invalid_param", i)
			if testing.Short() {
				break
			}
		}
	}
}

func TestSleepSuccess(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte("ok\r\n")
		return len(b), b
	}

	defer resetOriginals()

	for i := uint32(100); i <= 1000100; i += 100000 {
		if Sleep(i) != true {
			t.Errorf("Sleep(%v) returned false while the serial read returned ok", i)
			if testing.Short() {
				break
			}
		}
	}
}

func TestResetFail(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return errors.New("Mock Write Error")
	}

	defer resetOriginals()

	if Reset() != false {
		t.Errorf("Reset() returned true while the serial write failed")
	}
}

func TestResetSuccess(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialFlush = func() {}

	defer resetOriginals()

	if Reset() != true {
		t.Errorf("Reset() returned false while the serial write succeeded")
	}
}

func TestSaveByteWrongArgument(t *testing.T) {
	for i := uint16(0); i < 768; i++ {
		for j := uint8(0); j < maxUint8; j++ {
			if SaveByte(i, j) != false {
				t.Errorf("SaveByte(%v, %v) returned true while the address < 768", i, j)
				if testing.Short() {
					break
				}
			}
		}
	}

	for i := uint16(1024); i < 2048; i += 100 {
		for j := uint8(0); j < maxUint8; j++ {
			if SaveByte(i, j) != false {
				t.Errorf("SaveByte(%v, %v) returned true while the address > 1023", i, j)
				if testing.Short() {
					break
				}
			}
		}
	}
}

func TestSaveByteWriteError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return errors.New("Mock Write Error")
	}

	defer resetOriginals()

	for i := uint16(768); i <= 1023; i++ {
		for j := uint8(0); j < maxUint8; j++ {
			if SaveByte(i, j) != false {
				t.Errorf("SaveByte(%v, %v) returned true while the serial write failed", i, j)
				if testing.Short() {
					break
				}
			}
		}
	}
}

func TestSaveByteReadError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		var b []byte
		return 0, b
	}

	defer resetOriginals()

	for i := uint16(768); i <= 1023; i++ {
		for j := uint8(0); j < maxUint8; j++ {
			if SaveByte(i, j) != false {
				t.Errorf("SaveByte(%v, %v) returned true while the serial read returned 0 bytes", i, j)
				if testing.Short() {
					break
				}
			}
		}
	}
}

func TestSaveByteReadInvalidParam(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte(invalidParameter + "\r\n")
		return len(b), b
	}

	defer resetOriginals()

	for i := uint16(768); i <= 1023; i++ {
		for j := uint8(0); j < maxUint8; j++ {
			if SaveByte(i, j) != false {
				t.Errorf("SaveByte(%v, %v) returned true while the serial read returned invalid_param", i, j)
				if testing.Short() {
					break
				}
			}
		}
	}
}

func TestSaveByteSuccess(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte("ok\r\n")
		return len(b), b
	}

	defer resetOriginals()

	for i := uint16(768); i <= 1023; i++ {
		for j := uint8(0); j < maxUint8; j++ {
			if SaveByte(i, j) != true {
				t.Errorf("SaveByte(%v, %v) returned false while the serial read returned ok", i, j)
				if testing.Short() {
					break
				}
			}
		}
	}
}

func TestReadByteWrongArgument(t *testing.T) {
	for i := uint16(0); i < 768; i++ {
		if b, err := ReadByte(i); b != 0 && err == nil {
			t.Errorf("ReadByte(%v) returned no error while the address < 768", i)
			if testing.Short() {
				break
			}
		}
	}

	for i := uint16(1024); i < 2048; i += 100 {
		if b, err := ReadByte(i); b != 0 && err == nil {
			t.Errorf("ReadByte(%v) returned no error while the address > 1023", i)
			if testing.Short() {
				break
			}
		}
	}
}

func TestReadByteWriteError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return errors.New("Mock Write Error")
	}

	defer resetOriginals()

	for i := uint16(768); i <= 1023; i++ {
		if b, err := ReadByte(i); b != 0 && err == nil {
			t.Errorf("ReadByte(%v) returned no error while the serial write failed", i)
			if testing.Short() {
				break
			}
		}
	}
}

func TestReadByteReadError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		var b []byte
		return 0, b
	}

	defer resetOriginals()

	for i := uint16(768); i <= 1023; i++ {
		if b, err := ReadByte(i); b != 0 && err == nil {
			t.Errorf("ReadByte(%v) returned no error while the serial read returned 0 bytes", i)
			if testing.Short() {
				break
			}
		}
	}
}

func TestReadByteReadInvalidParam(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte(invalidParameter + "\r\n")
		return len(b), b
	}

	defer resetOriginals()

	for i := uint16(768); i <= 1023; i++ {
		if b, err := ReadByte(i); b != 0 && err == nil {
			t.Errorf("ReadByte(%v) returned no error while the serial read returned invalid_param", i)
			if testing.Short() {
				break
			}
		}
	}
}

func TestReadByteConversionError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte("ok\r\n")
		return len(b), b
	}

	defer resetOriginals()

	for i := uint16(768); i <= 1023; i++ {
		if b, err := ReadByte(i); b != 0 && err == nil {
			t.Errorf("ReadByte(%v) returned no error while the serial read returned something not convertable", i)
			if testing.Short() {
				break
			}
		}
	}
}

func TestReadByteSuccess(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte("FF\r\n")
		return len(b), b
	}

	defer resetOriginals()

	for i := uint16(768); i <= 1023; i++ {
		if b, err := ReadByte(i); b == 0 && err != nil {
			t.Errorf("ReadByte(%v) returned an error while the serial read returned byte", i)
			if testing.Short() {
				break
			}
		}
	}
}

func TestVersionWriteError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return errors.New("Mock Write Error")
	}

	defer resetOriginals()

	if Version() != "" {
		t.Errorf("Version() returned non empty string while the serial write failed")
	}
}

func TestVersionReadError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		var b []byte
		return len(b), b
	}

	defer resetOriginals()

	if Version() != "" {
		t.Errorf("Version() returned non empty string while the serial read failed")
	}
}

func TestVersionSuccess(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte("Version test\r\n")
		return len(b), b
	}

	defer resetOriginals()

	if Version() == "" {
		t.Errorf("Version() returned empty string while the serial write and read succeeded")
	}
}

func TestVoltageWriteError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return errors.New("Mock Write Error")
	}

	defer resetOriginals()

	if value, err := Voltage(); value != 0 && err == nil {
		t.Errorf("Voltage() returned no error while the serial write failed")
	}
}

func TestVoltageReadError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		var b []byte
		return len(b), b
	}

	defer resetOriginals()

	if value, err := Voltage(); value != 0 && err == nil {
		t.Errorf("Voltage() returned no error while the serial read failed")
	}
}

func TestVoltageReadInvalidParamr(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte(invalidParameter + "\r\n")
		return len(b), b
	}

	defer resetOriginals()

	if value, err := Voltage(); value != 0 && err == nil {
		t.Errorf("Voltage() returned no error while the serial read returned invalid_param")
	}
}

func TestVoltageConversionError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte("nan\r\n")
		return len(b), b
	}

	defer resetOriginals()

	if value, err := Voltage(); value != 0 && err == nil {
		t.Errorf("Voltage() returned no error while the conversion failed")
	}
}

func TestVoltageSuccess(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte("3000\r\n")
		return len(b), b
	}

	defer resetOriginals()

	if value, err := Voltage(); value == 0 && err != nil {
		t.Errorf("Voltage() returned an error while the serial read succeeded")
	}
}

func TestHardwareIDWriteError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return errors.New("Mock Write Error")
	}

	defer resetOriginals()

	if HardwareID() != "" {
		t.Errorf("HardwareID() returned non empty string while the serial write failed")
	}
}

func TestHardwareIDReadError(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		var b []byte
		return len(b), b
	}

	defer resetOriginals()

	if HardwareID() != "" {
		t.Errorf("HardwareID() returned non empty string while the serial read failed")
	}
}

func TestHardwareIDSuccess(t *testing.T) {
	serialWrite = func(s string) error {
		t.Logf("String written to serial: %v", s)
		return nil
	}

	serialRead = func() (int, []byte) {
		b := []byte("HardwareID test\r\n")
		return len(b), b
	}

	defer resetOriginals()

	if HardwareID() == "" {
		t.Errorf("HardwareID() returned empty string while the serial write and read succeeded")
	}
}
