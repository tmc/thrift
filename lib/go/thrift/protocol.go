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

const (
	VERSION_MASK = 0xffff0000
	VERSION_1    = 0x80010000
)

type TProtocol interface {
	WriteMessageBegin(name string, typeId TMessageType, seqid int32) TProtocolException
	WriteMessageEnd() TProtocolException
	WriteStructBegin(name string) TProtocolException
	WriteStructEnd() TProtocolException
	WriteFieldBegin(name string, typeId TType, id int16) TProtocolException
	WriteFieldEnd() TProtocolException
	WriteFieldStop() TProtocolException
	WriteMapBegin(keyType TType, valueType TType, size int) TProtocolException
	WriteMapEnd() TProtocolException
	WriteListBegin(elemType TType, size int) TProtocolException
	WriteListEnd() TProtocolException
	WriteSetBegin(elemType TType, size int) TProtocolException
	WriteSetEnd() TProtocolException
	WriteBool(value bool) TProtocolException
	WriteByte(value byte) TProtocolException
	WriteI16(value int16) TProtocolException
	WriteI32(value int32) TProtocolException
	WriteI64(value int64) TProtocolException
	WriteDouble(value float64) TProtocolException
	WriteString(value string) TProtocolException
	WriteBinary(value []byte) TProtocolException

	ReadMessageBegin() (name string, typeId TMessageType, seqid int32, err TProtocolException)
	ReadMessageEnd() TProtocolException
	ReadStructBegin() (name string, err TProtocolException)
	ReadStructEnd() TProtocolException
	ReadFieldBegin() (name string, typeId TType, id int16, err TProtocolException)
	ReadFieldEnd() TProtocolException
	ReadMapBegin() (keyType TType, valueType TType, size int, err TProtocolException)
	ReadMapEnd() TProtocolException
	ReadListBegin() (elemType TType, size int, err TProtocolException)
	ReadListEnd() TProtocolException
	ReadSetBegin() (elemType TType, size int, err TProtocolException)
	ReadSetEnd() TProtocolException
	ReadBool() (value bool, err TProtocolException)
	ReadByte() (value byte, err TProtocolException)
	ReadI16() (value int16, err TProtocolException)
	ReadI32() (value int32, err TProtocolException)
	ReadI64() (value int64, err TProtocolException)
	ReadDouble() (value float64, err TProtocolException)
	ReadString() (value string, err TProtocolException)
	ReadBinary() (value []byte, err TProtocolException)

	Skip(fieldType TType) (err TProtocolException)
	Flush() (err TProtocolException)

	Transport() TTransport
}

// The maximum recursive depth the skip() function will traverse
var MaxSkipDepth = 1<<31 - 1

// Skips over the next data element from the provided input TProtocol object.
func SkipDefaultDepth(prot TProtocol, typeId TType) (err TProtocolException) {
	return Skip(prot, typeId, MaxSkipDepth)
}

// Skips over the next data element from the provided input TProtocol object.
func Skip(self TProtocol, fieldType TType, maxDepth int) (err TProtocolException) {
	switch fieldType {
	case STOP:
		return
	case BOOL:
		_, err = self.ReadBool()
		return
	case BYTE:
		_, err = self.ReadByte()
		return
	case I16:
		_, err = self.ReadI16()
		return
	case I32:
		_, err = self.ReadI32()
		return
	case I64:
		_, err = self.ReadI64()
		return
	case DOUBLE:
		_, err = self.ReadDouble()
		return
	case STRING:
		_, err = self.ReadString()
		return
	case STRUCT:
		if _, err = self.ReadStructBegin(); err != nil {
			return err
		}
		for {
			_, typeId, _, _ := self.ReadFieldBegin()
			if typeId == STOP {
				break
			}
			Skip(self, typeId, maxDepth-1)
			self.ReadFieldEnd()
		}
		return self.ReadStructEnd()
	case MAP:
		keyType, valueType, size, err := self.ReadMapBegin()
		if err != nil {
			return err
		}
		for i := 0; i < size; i++ {
			Skip(self, keyType, maxDepth-1)
			self.Skip(valueType)
		}
		return self.ReadMapEnd()
	case SET:
		elemType, size, err := self.ReadSetBegin()
		if err != nil {
			return err
		}
		for i := 0; i < size; i++ {
			Skip(self, elemType, maxDepth-1)
		}
		return self.ReadSetEnd()
	case LIST:
		elemType, size, err := self.ReadListBegin()
		if err != nil {
			return err
		}
		for i := 0; i < size; i++ {
			Skip(self, elemType, maxDepth-1)
		}
		return self.ReadListEnd()
	}
	return nil
}
