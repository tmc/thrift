/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package thrift

import (
	"bytes"
	"net"
	"time"
)

type TSocket struct {
	writeBuffer *bytes.Buffer
	conn        net.Conn
	addr        net.Addr
	nsecTimeout int64
}

// Creates a Socket from an existing net.Conn
func NewTSocketConn(connection net.Conn) (*TSocket, TTransportException) {
	return NewTSocketConnTimeout(connection, 0)
}

// Creates a Socket from an existing net.Conn
func NewTSocketConnTimeout(connection net.Conn, nsecTimeout int64) (*TSocket, TTransportException) {
	address := connection.RemoteAddr()
	if address == nil {
		address = connection.LocalAddr()
	}
	p := &TSocket{conn: connection, addr: address, nsecTimeout: nsecTimeout, writeBuffer: bytes.NewBuffer(make([]byte, 0, 4096))}
	return p, nil
}

// Creates a Socket from a net.Addr
func NewTSocketAddr(address net.Addr) *TSocket {
	return NewTSocket(address, 0)
}

// Creates a Socket from a net.Addr
func NewTSocket(address net.Addr, nsecTimeout int64) *TSocket {
	sock := &TSocket{addr: address, nsecTimeout: nsecTimeout, writeBuffer: bytes.NewBuffer(make([]byte, 0, 4096))}
	return sock
}

// Sets the socket timeout
func (p *TSocket) SetTimeout(nsecTimeout int64) error {
	p.nsecTimeout = nsecTimeout
	return nil
}

func (p *TSocket) pushDeadline(read, write bool) {
	var t time.Time
	if p.nsecTimeout > 0 {
		t = time.Now().Add(time.Duration(p.nsecTimeout))
	}
	if read && write {
		p.conn.SetDeadline(t)
	} else if read {
		p.conn.SetReadDeadline(t)
	} else if write {
		p.conn.SetWriteDeadline(t)
	}
}

func (p *TSocket) Conn() net.Conn {
	return p.conn
}

func (p *TSocket) IsOpen() bool {
	if p.conn == nil {
		return false
	}
	return true
}

// Connects the socket, creating a new socket object if necessary.
func (p *TSocket) Open() error {
	if p.IsOpen() {
		return NewTTransportException(ALREADY_OPEN, "Socket already connected.")
	}
	if p.addr == nil {
		return NewTTransportException(NOT_OPEN, "Cannot open nil address.")
	}
	if len(p.addr.Network()) == 0 {
		return NewTTransportException(NOT_OPEN, "Cannot open bad network name.")
	}
	if len(p.addr.String()) == 0 {
		return NewTTransportException(NOT_OPEN, "Cannot open bad address.")
	}
	var err error
	if p.nsecTimeout > 0 {
		if p.conn, err = net.DialTimeout(p.addr.Network(), p.addr.String(), time.Duration(p.nsecTimeout)); err != nil {
			return NewTTransportException(NOT_OPEN, err.Error())
		}
	} else {
		if p.conn, err = net.Dial(p.addr.Network(), p.addr.String()); err != nil {
			return NewTTransportException(NOT_OPEN, err.Error())
		}
	}
	return nil
}

// Closes the socket.
func (p *TSocket) Close() error {
	// Close the socket
	if p.conn != nil {
		err := p.conn.Close()
		if err != nil {
			return err
		}
		p.conn = nil
	}
	return nil
}

func (p *TSocket) Read(buf []byte) (int, error) {
	if !p.IsOpen() {
		return 0, NewTTransportException(NOT_OPEN, "Connection not open")
	}
	p.pushDeadline(true, false)
	n, err := p.conn.Read(buf)
	return n, NewTTransportExceptionFromError(err)
}

func (p *TSocket) Write(buf []byte) (int, error) {
	if !p.IsOpen() {
		return 0, NewTTransportException(NOT_OPEN, "Connection not open")
	}
	p.pushDeadline(false, true)
	p.writeBuffer.Write(buf)
	return len(buf), nil
}

func (p *TSocket) Peek() bool {
	return p.IsOpen()
}

func (p *TSocket) Flush() error {
	if !p.IsOpen() {
		return NewTTransportException(NOT_OPEN, "Connection not open")
	}
	_, err := p.writeBuffer.WriteTo(p.conn)
	return NewTTransportExceptionFromError(err)
}

func (p *TSocket) Interrupt() error {
	if !p.IsOpen() {
		return nil
	}
	// TODO(pomack) fix Interrupt as this is probably wrong
	return p.conn.Close()
}
