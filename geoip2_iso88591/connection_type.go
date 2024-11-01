package geoip2_iso88591

import "errors"

func readConnectionTypeMap(result *ConnectionType, buffer []byte, mapSize, offset uint) (uint, error) {
	var key []byte
	var err error
	for i := uint(0); i < mapSize; i++ {
		key, offset, err = readMapKey(buffer, offset)
		if err != nil {
			return 0, err
		}
		switch bytesToKeyString(key) {
		case "connection_type":
			result.ConnectionType, offset, err = readString(buffer, offset)
			if err != nil {
				return 0, err
			}
		default:
			return 0, errors.New("unknown connectionType key: " + string(key))
		}
	}
	return offset, nil
}
