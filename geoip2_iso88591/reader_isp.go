package geoip2_iso88591

import (
	"errors"
	"io/ioutil"
	"net"
	"strconv"
)

type ISPReader struct {
	*reader
}

func (r *ISPReader) Lookup(ip net.IP) (*ISP, error) {
	offset, err := r.getOffset(ip)
	if err != nil {
		return nil, err
	}
	dataType, size, offset, err := readControl(r.decoderBuffer, offset)
	if err != nil {
		return nil, err
	}
	result := &ISP{}
	switch dataType {
	case dataTypeMap:
		_, err = readISPMap(result, r.decoderBuffer, size, offset)
		if err != nil {
			return nil, err
		}
	case dataTypePointer:
		pointer, _, err := readPointer(r.decoderBuffer, size, offset)
		if err != nil {
			return nil, err
		}
		dataType, size, offset, err := readControl(r.decoderBuffer, pointer)
		if err != nil {
			return nil, err
		}
		if dataType != dataTypeMap {
			return nil, errors.New("invalid ISP pointer type: " + strconv.Itoa(int(dataType)))
		}
		_, err = readISPMap(result, r.decoderBuffer, size, offset)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid ISP type: " + strconv.Itoa(int(dataType)))
	}
	return result, nil
}

func NewISPReader(buffer []byte) (*ISPReader, error) {
	reader, err := newReader(buffer)
	if err != nil {
		return nil, err
	}
	if reader.metadata.DatabaseType != "GeoIP2-ISP" {
		return nil, errors.New("wrong MaxMind DB ISP type: " + reader.metadata.DatabaseType)
	}
	return &ISPReader{
		reader: reader,
	}, nil
}

func NewISPReaderFromFile(filename string) (*ISPReader, error) {
	buffer, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return NewISPReader(buffer)
}
