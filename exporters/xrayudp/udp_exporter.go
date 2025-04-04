// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package xrayudp

import (
	"encoding/base64"
	"fmt"
	"net"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel"
)

const (
	DefaultEndpoint = "127.0.0.1:2000"
	ProtocolHeader  = "{\"format\":\"json\",\"version\":1}\n"
)

type Conn interface {
	Write([]byte) (int, error)
	Close() error
}

type UdpExporter struct {
	endpoint string
	host     string
	port     int
	conn     Conn
}

func NewUdpExporter(endpoint string) (*UdpExporter, error) {
	if endpoint == "" {
		endpoint = DefaultEndpoint
	}

	exporter := &UdpExporter{
		endpoint: endpoint,
	}

	var err error
	exporter.host, exporter.port, err = parseEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", exporter.host, exporter.port))
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return nil, err
	}
	exporter.conn = conn

	return exporter, nil
}

func (e *UdpExporter) SendData(data []byte, signalFormatPrefix string) error {
	base64EncodedString := base64.StdEncoding.EncodeToString(data)
	message := fmt.Sprintf("%s%s%s", ProtocolHeader, signalFormatPrefix, base64EncodedString)

	_, err := e.conn.Write([]byte(message))
	if err != nil {
		otel.Handle(fmt.Errorf("error sending UDP data: %w", err))
		return err
	}

	return nil
}

func (e *UdpExporter) Shutdown() error {
	return e.conn.Close()
}

func parseEndpoint(endpoint string) (string, int, error) {
	parts := strings.Split(endpoint, ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid endpoint: %s", endpoint)
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, fmt.Errorf("invalid port in endpoint: %s", endpoint)
	}

	return parts[0], port, nil
}
