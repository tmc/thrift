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
	"strconv"
)

type Flusher interface {
	Flush() (err error)
}

// Encapsulates the I/O layer
type TTransport interface {
	IsOpen() bool
	Open() (err error)
	Close() (err error)
	Read(buf []byte) (n int, err error)
	ReadAll(buf []byte) (n int, err error)
	Write(buf []byte) (n int, err error)
	Flush() (err error)

	/**
	 * Is there more data to be read?
	 *
	 * @return True if the remote side is still alive and feeding us
	 */
	Peek() bool
}

/*
type TTransportBase struct {
}

func (p* TTransportBase) IsOpen() bool {
  return false;
};

func (p* TTransportBase) Peek() bool {
  return p.IsOpen();
}

func (p* TTransportBase) Open() error {
  return NewTTransportException(UNKNOWN, "Subclasses must implement TTransportBase.Open()");
}

func (p* TTransportBase) Close() error {
  return NewTTransportException(UNKNOWN, "Subclasses must implement TTransportBase.Close()");
}

func (p* TTransportBase) Read(buf []byte) (int, error) {
  return 0, NewTTransportExceptionDefaultString("Subclasses must implement TTransportBase.Read()");
}

func (p* TTransportBase) ReadAll(buf []byte) (n int, err error){
  ret := 0;
  size := len(buf);
  for (n < size) {
    ret, err = p.Read(buf[n:]);
    if ret <= 0 {
      if err != nil {
        err = NewTTransportExceptionDefaultString("Cannot read. Remote side has closed. Tried to read " + string(size) + " bytes, but only got " + string(n) + " bytes.");
      }
      return ret, err;
    }
    n += ret;
  }
  return n, err;
}

func (p* TTransportBase) Write(buf []byte) (int, error) {
  return 0, NewTTransportExceptionDefaultString("Subclasses must implement TTransportBase.Write()");
}

func (p* TTransportBase) Flush() error {
  return nil;
}
*/

// Reads all the bytes off the transport. Returns the umber of bytes read and an error
func ReadAllTransport(p TTransport, buf []byte) (n int, err error) {
	ret := 0
	size := len(buf)
	for n < size {
		ret, err = p.Read(buf[n:])
		if ret <= 0 {
			if err != nil {
				err = NewTTransportExceptionDefaultString("Cannot read. Remote side has closed. Tried to read " + strconv.Itoa(size) + " bytes, but only got " + strconv.Itoa(n) + " bytes.")
			}
			return ret, err
		}
		n += ret
	}
	return n, err
}
