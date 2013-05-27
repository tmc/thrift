package main

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

import (
	"fmt"
	"git.apache.org/thrift/lib/go/thrift"
	"tutorial"
)

func runServer(transportFactory thrift.TTransportFactory, protocolFactory thrift.TProtocolFactory) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:9090")
	if err != nil {
		fmt.Println("Error resolving address: ", err)
		return
	}

	handler := NewCalculatorHandler()
	processor := tutorial.NewCalculatorProcessor(handler)
	transport, err := thrift.NewTServerSocketAddr(addr)
	if err != nil {
		fmt.Println("Error creating server socket: ", err)
		return
	}

	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)

	fmt.Println("Starting the simple server... on ", addr)
	if err = server.Serve(); err != nil {
		fmt.Println("Error during simple server: ", err)
	}
}
