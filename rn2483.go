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
	"os"
	"time"

	"github.com/tarm/serial"
)

var (
	err    error
	n      int
	rn2483 *serial.Port
	config *serial.Config
)

func init() {
	// TODO make use of viper to get this from config
	config = &serial.Config{ReadTimeout: time.Millisecond * 100}
}

func read() (int, []byte) {
	defer func() {
		if r := recover(); r != nil {
			DEBUG.Println("Are you connected?")
			ERROR.Println(r)
			os.Exit(2)
		}
	}()

	var b []byte

	for {
		buf := make([]byte, 128)
		n, err = rn2483.Read(buf)
		if err != nil { // err will equal io.EOF
			break
		}
		b = append(b, buf[:n]...)
	}

	DEBUG.Printf("%v bytes read: %X", n, b)

	return len(b), b
}

func write(s string) error {
	defer func() {
		if r := recover(); r != nil {
			DEBUG.Println("Are you connected?")
			ERROR.Println(r)
			os.Exit(2)
		}
	}()

	b := append([]byte(s), []byte("\r\n")...)
	n, err = rn2483.Write(b)
	if err != nil {
		WARN.Println("RN2483 write error:", err)
		return err
	}
	DEBUG.Printf("%v bytes written: %X", n, b)
	return nil
}

func flush() {
	err = rn2483.Flush()
	if err != nil {
		WARN.Println("RN2483 flush error:", err)
	}
}

// Connect will connect to the serial device currently configured.
func Connect() {
	rn2483, err = serial.OpenPort(config)
	if err != nil {
		ERROR.Println(err)
	}
	flush()
	DEBUG.Println("RN2483 connected")
}

// Disconnect will disconnect the serial device that is currently connected.
// If no device is connected, it will recover from the panic.
func Disconnect() {
	defer func() {
		if r := recover(); r != nil {
			ERROR.Println("Recovered:", r)
		}
	}()

	err = rn2483.Close()
	if err != nil {
		ERROR.Println(err)
	}
	DEBUG.Println("RN2483 disconnected")
}

// SetName sets a new device name for the serial connection.
// Reconnect to the serial device required!
func SetName(name string) {
	config.Name = name
	DEBUG.Println("RN2483 serial device:", name)
}

// SetBaud sets a new baud rate for the serial connection.
// Reconnect to the serial device required!
func SetBaud(baud int) {
	config.Baud = baud
	DEBUG.Println("RN2483 baud rate:", baud)
}

// SetTimeout sets a new read timeout for the serial connection.
// Reconnect to the serial device required!
func SetTimeout(timeout time.Duration) {
	config.ReadTimeout = timeout
	DEBUG.Println("RN2483 read timeout:", timeout)
}
