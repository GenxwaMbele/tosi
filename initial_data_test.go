/*
 Copyright 2013-2021 Daniele Pala <pala.daniele@gmail.com>

 This file is part of tosi.

 tosi is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 tosi is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with tosi. If not, see <http://www.gnu.org/licenses/>.

*/

package github.com/GenxwaMbele/tosi
import (
	"bytes"
	"strconv"
	"testing"
	"time"
)

const (
	// each test uses different ports for servers,
	// in order to avoid possible conflicts.
	initTest1Port = 8107
	initTest2Port = 8108
	initTest3Port = 8109
)

// Test 1
// test initial data write with 5 bytes. No error should occur.
func TestWrite5bytesIn(t *testing.T) {
	// start a server
	go tosiServerRead5bytesIn(t, initTest1Port)
	// wait for server to come up
	time.Sleep(time.Millisecond)
	remAddr := "127.0.0.1:" + strconv.Itoa(initTest1Port) + ":100"
	tosiAddr, err := ResolveTOSIAddr("tosi", remAddr)
	checkErrorIn(err, t)
	// try to connect with initial data
	data := []byte{0x01, 0xff, 0x66, 0x93, 0x20}
	opt := DialOpt{Data: data}
	conn, err := DialOptTOSI("tosi", nil, tosiAddr, opt)
	checkErrorIn(err, t)
	if conn.UseExpedited == true {
		t.Log("Expedited service available but not requested")
		t.FailNow()
	}
	buf := make([]byte, 100)
	read, err := conn.ReadTOSI(buf)
	checkErrorIn(err, t)
	if read.N != 5 {
		t.Log("Wrong data size")
		t.FailNow()
	}
	if !bytes.Equal(buf[:5], data) {
		t.Log("Wrong data values")
		t.FailNow()
	}
	if read.Expedited == true {
		t.Log("expedited data received")
		t.FailNow()
	}
	// close connection
	err = conn.Close()
	checkErrorIn(err, t)
}

// Test 2
// test initial data write with 35 bytes. No error should occur,
// but only 32 bytes should be transferred.
func TestWrite35bytesIn(t *testing.T) {
	// start a server
	go tosiServerRead5bytesIn(t, initTest2Port)
	// wait for server to come up
	time.Sleep(time.Millisecond)
	remAddr := "127.0.0.1:" + strconv.Itoa(initTest2Port) + ":100"
	tosiAddr, err := ResolveTOSIAddr("tosi", remAddr)
	checkErrorIn(err, t)
	// try to connect with initial data
	data := make([]byte, 35)
	opt := DialOpt{Data: data}
	conn, err := DialOptTOSI("tosi", nil, tosiAddr, opt)
	checkErrorIn(err, t)
	if conn.UseExpedited == true {
		t.Log("Expedited service available but not requested")
		t.FailNow()
	}
	buf := make([]byte, 100)
	read, err := conn.ReadTOSI(buf)
	checkErrorIn(err, t)
	if read.N != 32 {
		t.Log("Wrong data size")
		t.FailNow()
	}
	if !bytes.Equal(buf[:32], data[:32]) {
		t.Log("Wrong data values")
		t.FailNow()
	}
	if read.Expedited == true {
		t.Log("expedited data received")
		t.FailNow()
	}
	// close connection
	err = conn.Close()
	checkErrorIn(err, t)
}

// Test 3
// test initial data write with 5 bytes. The server should read this data with
// a normal ReadTOSI call. No error should occur.
func TestWrite5bytes(t *testing.T) {
	// start a server
	go tosiServerRead5bytes(t, initTest3Port)
	// wait for server to come up
	time.Sleep(time.Millisecond)
	remAddr := "127.0.0.1:" + strconv.Itoa(initTest3Port) + ":100"
	tosiAddr, err := ResolveTOSIAddr("tosi", remAddr)
	checkErrorIn(err, t)
	// try to connect with initial data
	data := []byte{0x01, 0xff, 0x66, 0x93, 0x20}
	opt := DialOpt{Data: data}
	conn, err := DialOptTOSI("tosi", nil, tosiAddr, opt)
	checkErrorIn(err, t)
	if conn.UseExpedited == true {
		t.Log("Expedited service available but not requested")
		t.FailNow()
	}
	// close connection
	err = conn.Close()
	checkErrorIn(err, t)
}

// a tosi server reading 5 bytes of initial data. No fault is expected.
func tosiServerRead5bytesIn(t *testing.T, port int) {
	locAddr := "127.0.0.1:" + strconv.Itoa(port) + ":100"
	tosiAddr, err := ResolveTOSIAddr("tosi", locAddr)
	checkErrorIn(err, t)
	listener, err := ListenTOSI("tosi", tosiAddr)
	checkErrorIn(err, t)
	// listen for connections
	conn, err := listener.AcceptTOSI(func(b []byte) []byte { return b })
	checkErrorIn(err, t)
	if conn.UseExpedited == true {
		t.Log("Expedited service available but not requested")
		t.FailNow()
	}
	err = listener.Close()
	checkErrorIn(err, t)
}

// a tosi server reading 5 bytes. No fault is expected.
func tosiServerRead5bytes(t *testing.T, port int) {
	locAddr := "127.0.0.1:" + strconv.Itoa(port) + ":100"
	tosiAddr, err := ResolveTOSIAddr("tosi", locAddr)
	checkErrorIn(err, t)
	listener, err := ListenTOSI("tosi", tosiAddr)
	checkErrorIn(err, t)
	// listen for connections
	conn, err := listener.Accept()
	checkErrorIn(err, t)
	buf := make([]byte, 100)
	read, err := conn.(*TOSIConn).ReadTOSI(buf)
	checkErrorIn(err, t)
	if read.N != 5 {
		t.Log("Wrong data size")
		t.FailNow()
	}
	if !bytes.Equal(buf[:5], []byte{0x01, 0xff, 0x66, 0x93, 0x20}) {
		t.Log("Wrong data values")
		t.FailNow()
	}
	if read.Expedited == true {
		t.Log("Expedited data received")
		t.FailNow()
	}
	if read.EndOfTSDU == false {
		t.Log("Wrong EndOfTSDU indication")
		t.FailNow()
	}
	err = listener.Close()
	checkErrorIn(err, t)
}

// check for unexpected errors
func checkErrorIn(err error, t *testing.T) {
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}
}

// check for expected errors
func checkWantedErrorIn(err error, t *testing.T) {
	if err == nil {
		t.FailNow()
	}
}
