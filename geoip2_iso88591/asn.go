package geoip2_iso88591

import "errors"

func readASNMap(result *ASN, buffer []byte, mapSize, offset uint) (uint, error) {
	var key []byte
	var err error
	for i := uint(0); i < mapSize; i++ {
		key, offset, err = readMapKey(buffer, offset)
		if err != nil {
			return 0, err
		}
		switch bytesToKeyString(key) {
		case "autonomous_system_number":
			result.AutonomousSystemNumber, offset, err = readUInt32(buffer, offset)
			if err != nil {
				return 0, err
			}
		case "autonomous_system_organization":
			result.AutonomousSystemOrganization, offset, err = readString(buffer, offset)
			if err != nil {
				return 0, err
			}
		default:
			return 0, errors.New("unknown asn key: " + string(key))
		}
	}
	return offset, nil
}
