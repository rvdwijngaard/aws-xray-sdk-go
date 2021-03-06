// Copyright 2017-2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may not use this file except in compliance with the License. A copy of the License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
package daemoncfg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDaemonEndpoints1(t *testing.T) { // default address set to udp and tcp
	udpAddr := "127.0.0.1:2000"
	tcpAddr := "127.0.0.1:2000"
	udpEndpt, _ := resolveUDPAddr(udpAddr)
	tcpEndpt, _ := resolveTCPAddr(tcpAddr)
	dEndpt := GetDaemonEndpoints()

	assert.Equal(t, dEndpt.UDPAddr, udpEndpt)
	assert.Equal(t, dEndpt.TCPAddr, tcpEndpt)
}

func TestGetDaemonEndpoints2(t *testing.T) { // default address set to udp and tcp
	udpAddr := "127.0.0.1:4000"
	tcpAddr := "127.0.0.1:5000"
	udpEndpt, _ := resolveUDPAddr(udpAddr)
	tcpEndpt, _ := resolveTCPAddr(tcpAddr)

	dAddr := "tcp:" + tcpAddr + " udp:" + udpAddr

	os.Setenv("AWS_XRAY_DAEMON_ADDRESS", dAddr) // env variable gets precedence over provided daemon addr
	defer os.Unsetenv("AWS_XRAY_DAEMON_ADDRESS")

	dEndpt := GetDaemonEndpoints()

	assert.Equal(t, dEndpt.UDPAddr, udpEndpt)
	assert.Equal(t, dEndpt.TCPAddr, tcpEndpt)
}

func TestGetDaemonEndpointsFromString1(t *testing.T) {
	udpAddr := "127.0.0.1:2000"
	tcpAddr := "127.0.0.1:2000"
	udpEndpt, _ := resolveUDPAddr(udpAddr)
	tcpEndpt, _ := resolveTCPAddr(tcpAddr)
	dAddr := udpAddr
	dEndpt, err := GetDaemonEndpointsFromString(dAddr)

	assert.Nil(t, err)
	assert.Equal(t, dEndpt.UDPAddr, udpEndpt)
	assert.Equal(t, dEndpt.TCPAddr, tcpEndpt)
}

func TestGetDaemonEndpointsFromString2(t *testing.T) {

	udpAddr := "127.0.0.1:2000"
	tcpAddr := "127.0.0.1:2000"

	dAddr := "127.0.0.1:2001"

	os.Setenv("AWS_XRAY_DAEMON_ADDRESS", udpAddr) // env variable gets precedence over provided daemon addr
	defer os.Unsetenv("AWS_XRAY_DAEMON_ADDRESS")

	udpEndpt, _ := resolveUDPAddr(udpAddr)
	tcpEndpt, _ := resolveTCPAddr(tcpAddr)

	dEndpt, err := GetDaemonEndpointsFromString(dAddr)

	assert.Nil(t, err)
	assert.Equal(t, dEndpt.UDPAddr, udpEndpt)
	assert.Equal(t, dEndpt.TCPAddr, tcpEndpt)
}

func TestGetDaemonEndpointsFromString3(t *testing.T) {
	udpAddr := "127.0.0.2:2001"
	tcpAddr := "127.0.0.1:2000"
	udpEndpt, _ := resolveUDPAddr(udpAddr)
	tcpEndpt, _ := resolveTCPAddr(tcpAddr)
	dAddr := "tcp:" + tcpAddr + " udp:" + udpAddr
	dEndpt, err := GetDaemonEndpointsFromString(dAddr)

	assert.Nil(t, err)
	assert.Equal(t, dEndpt.UDPAddr, udpEndpt)
	assert.Equal(t, dEndpt.TCPAddr, tcpEndpt)
}

func TestGetDaemonEndpointsFromString4(t *testing.T) {
	udpAddr := "127.0.0.2:2001"
	tcpAddr := "127.0.0.1:2000"
	udpEndpt, _ := resolveUDPAddr(udpAddr)
	tcpEndpt, _ := resolveTCPAddr(tcpAddr)
	dAddr := "udp:" + udpAddr + " tcp:" + tcpAddr
	dEndpt, err := GetDaemonEndpointsFromString(dAddr)

	assert.Nil(t, err)
	assert.Equal(t, dEndpt.UDPAddr, udpEndpt)
	assert.Equal(t, dEndpt.TCPAddr, tcpEndpt)
}

func TestGetDaemonEndpointsFromString5(t *testing.T) {
	udpAddr := "127.0.0.2:2001"
	tcpAddr := "127.0.0.1:2000"
	udpEndpt, _ := resolveUDPAddr(udpAddr)
	tcpEndpt, _ := resolveTCPAddr(tcpAddr)
	dAddr := "udp:" + udpAddr + " tcp:" + tcpAddr
	os.Setenv("AWS_XRAY_DAEMON_ADDRESS", dAddr) // env variable gets precedence over provided daemon addr
	defer os.Unsetenv("AWS_XRAY_DAEMON_ADDRESS")
	dEndpt, err := GetDaemonEndpointsFromString("tcp:127.0.0.5:2001 udp:127.0.0.5:2001")

	assert.Nil(t, err)
	assert.Equal(t, dEndpt.UDPAddr, udpEndpt)
	assert.Equal(t, dEndpt.TCPAddr, tcpEndpt)
}

func TestGetDaemonEndpointsFromStringInvalid1(t *testing.T) { // "udp:127.0.0.5:2001 udp:127.0.0.5:2001"
	udpAddr := "127.0.0.2:2001"
	tcpAddr := "127.0.0.1:2000"
	dAddr := "udp:" + udpAddr + " udp:" + tcpAddr
	dEndpt, err := GetDaemonEndpointsFromString(dAddr)

	assert.NotNil(t, err)
	assert.Nil(t, dEndpt)
}

func TestGetDaemonEndpointsFromStringInvalid2(t *testing.T) { // "tcp:127.0.0.5:2001 tcp:127.0.0.5:2001"
	udpAddr := "127.0.0.2:2001"
	tcpAddr := "127.0.0.1:2000"
	dAddr := "tcp:" + udpAddr + " tcp:" + tcpAddr
	dEndpt, err := GetDaemonEndpointsFromString(dAddr)

	assert.NotNil(t, err)
	assert.Nil(t, dEndpt)
}

func TestGetDaemonEndpointsFromStringInvalid3(t *testing.T) { // env variable set is invalid, string passed is valid
	udpAddr := "127.0.0.2:2001"
	tcpAddr := "127.0.0.1:2000"
	os.Setenv("AWS_XRAY_DAEMON_ADDRESS", "tcp:127.0.0.5:2001 tcp:127.0.0.5:2001") // invalid
	defer os.Unsetenv("AWS_XRAY_DAEMON_ADDRESS")
	dAddr := "udp:" + udpAddr + " tcp:" + tcpAddr
	dEndpt, err := GetDaemonEndpointsFromString(dAddr)

	assert.NotNil(t, err)
	assert.Nil(t, dEndpt)
}

func TestGetDaemonEndpointsFromStringInvalid4(t *testing.T) {
	udpAddr := "127.0.02:2001" // error in resolving address
	tcpAddr := "127.0.0.1:2000"

	dAddr := "udp:" + udpAddr + " tcp:" + tcpAddr
	dEndpt, err := GetDaemonEndpointsFromString(dAddr)

	assert.NotNil(t, err)
	assert.Nil(t, dEndpt)
}

func TestGetDaemonEndpointsFromStringInvalid5(t *testing.T) {
	udpAddr := "127.0.0.2:2001"
	tcpAddr := "127.0.a.1:2000" // error in resolving address

	dAddr := "udp:" + udpAddr + " tcp:" + tcpAddr
	dEndpt, err := GetDaemonEndpointsFromString(dAddr)

	assert.NotNil(t, err)
	assert.Nil(t, dEndpt)
}

func TestGetDaemonEndpointsFromStringInvalid6(t *testing.T) {
	udpAddr := "127.0.0.2:2001"
	dAddr := "udp:" + udpAddr // no tcp address present
	dEndpt, err := GetDaemonEndpointsFromString(dAddr)

	assert.NotNil(t, err)
	assert.Nil(t, dEndpt)
}

func TestGetDaemonEndpointsFromStringInvalid7(t *testing.T) {
	dAddr := ""
	dEndpt, err := GetDaemonEndpointsFromString(dAddr) // address passed is nil and env variable not set

	assert.Nil(t, err)
	assert.Nil(t, dEndpt)
}
