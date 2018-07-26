# RN2483
[![Build Status](https://travis-ci.org/BulletTime/RN2483.svg?branch=master)](https://travis-ci.org/BulletTime/RN2483)
[![GoDoc](https://godoc.org/github.com/BulletTime/RN2483?status.svg)](https://godoc.org/github.com/BulletTime/RN2483)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/BulletTime/RN2483/blob/master/LICENSE)

This library gives users the ability to use the RN2483 module in their application. For more info about the RN2483 module commands, see <http://ww1.microchip.com/downloads/en/DeviceDoc/40001784F.pdf>.

As this is a work in progress, not all commands are supported yet. For more info about what is supported, check the go docs reference.

## Usage
```
package main

import (
  "fmt"
  "time"

  "github.com/bullettime/rn2483"
)

func main() {
  // Setup the serial information
  rn2483.SetName("/dev/cu.usbmodem14111")
  rn2483.setBaud(57600)
  rn2483.SetTimeout(time.Millisecond * 500)

  // Connect the RN2483 via serial
  rn2483.Connect()
  // Make sure the app closes the connection at the end the free the resource
  defer rn2483.Disconnect()

  // Do you stuff
  fmt.Println(rn2483.Version())
}
```
