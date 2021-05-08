package gtools

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strconv"
)

func IP2Long(ip string) int64 {
	var ipLong uint32
	err := binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &ipLong)
	if err != nil {
		log.Println(fmt.Sprintf("[IP2Long] exception, err: %v", err))
		return 0
	}

	return int64(ipLong)
}

func Long2IP4(ipLong int64) string {
	// need to do two bit shifting and “0xff” masking
	b0 := strconv.FormatInt((ipLong>>24)&0xff, 10)
	b1 := strconv.FormatInt((ipLong>>16)&0xff, 10)
	b2 := strconv.FormatInt((ipLong>>8)&0xff, 10)
	b3 := strconv.FormatInt(ipLong&0xff, 10)

	return b0 + "." + b1 + "." + b2 + "." + b3
}
