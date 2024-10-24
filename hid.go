// usb - Self contained USB and HID library for Go
// Copyright 2019 The library Authors
//
// This library is free software: you can redistribute it and/or modify it under
// the terms of the GNU Lesser General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option) any
// later version.
//
// The library is distributed in the hope that it will be useful, but WITHOUT ANY
// WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR
// A PARTICULAR PURPOSE. See the GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License along
// with the library. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/binary"
	"errors"
	"log"

	"github.com/karalabe/usb"
)

const mhpVendorID = 0x12bf
const mhpProductID = 0xff03
const mhpMessageSize = 8

func hidSend(message int64) (err error) {
	// Enumerate all the HID devices matching the MHP ID
	hids, err := usb.EnumerateHid(mhpVendorID, mhpProductID)
	if err != nil {
		err = errors.New("mount hub pro not found")
		return
	}

	if len(hids) < 1 {
		err = errors.New("mount hub pro not found")
		return
	}

	if len(hids) > 1 {
		err = errors.New("only 1 Mount hub pro can be connected at the same time")
		return
	}

	myhid := hids[0]
	bs := make([]byte, mhpMessageSize)

	// note int64 is cast to uint64
	binary.LittleEndian.PutUint64(bs, uint64(message))

	mydevice, err := myhid.Open()
	if err != nil {
		return
	}
	defer mydevice.Close()
	mydevice.Write(bs)
	log.Printf("Command input: %x Little edian command sent: %x \n", message, bs)
	return
}
