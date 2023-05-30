package uuid

import (
	"sync"
	"time"
)

const (
	snowflakeMaxSerialNumber uint16 = uint16(0b00001111)<<8 + 0xff
	snowflake63Bit                  = ^(uint64(1) << 63)
)

var (
	snowflakeMutex        sync.Mutex
	snowflakeTimestamp    int64 = time.Now().UTC().Unix()
	snowflakeSerialNumber uint16
	snowflakeGroupID      byte
	snowflakeMechineID    byte
	hexBytes              []byte
)

func init() {
	hexBytes = make([]byte, 10+('z'-'a'+1)+('Z'-'A'+1))
	n := 0
	for i := byte('0'); i <= '9'; i++ {
		hexBytes[n] = i
		n++
	}
	for i := byte('a'); i <= 'z'; i++ {
		hexBytes[n] = i
		n++
	}
	for i := byte('A'); i <= 'Z'; i++ {
		hexBytes[n] = i
		n++
	}
}

func SetSnowflakeGroupID(id byte) {
	snowflakeGroupID = id & 0b0011111
}

func SetSnowflakeMechineID(id byte) {
	snowflakeMechineID = id & 0b0011111
}

func SnowflakeID() (n uint64) {
	timestamp := time.Now().UTC().Unix()
	// Number of generated in this millisecond.
	var serialNumber uint16
	snowflakeMutex.Lock()
	if timestamp == snowflakeTimestamp {
		snowflakeSerialNumber++
		// Too many
		if snowflakeSerialNumber > snowflakeMaxSerialNumber {
			snowflakeSerialNumber = 0
			timestamp++
			snowflakeTimestamp = timestamp
		}
	} else {
		snowflakeTimestamp = timestamp
	}
	serialNumber = snowflakeSerialNumber
	snowflakeMutex.Unlock()
	// 12bit id serial number
	n |= uint64(serialNumber)
	// 5bit mechine id
	n |= uint64(snowflakeMechineID) << 12
	// 5bit group id
	n |= uint64(snowflakeGroupID) << 17
	// 41bit timestamp
	n |= uint64(timestamp) << 22
	// 1bit 0
	n &= snowflake63Bit
	return n
}

func SnowflakeIDString() string {
	id := SnowflakeID()
	b := make([]byte, 20)
	i := 0
	m := uint64(0)
	for {
		m = id % uint64(len(hexBytes))
		b[i] = hexBytes[m]
		i++
		id = id / uint64(len(hexBytes))
		if id == 0 {
			break
		}
	}
	return string(b[:i])
}
