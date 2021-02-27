package service

import (
	"fmt"
	"net"
	"time"
)

// CHECK IF A TCP CONNECTION CAN BE ESTABLISHED
// DOES NOT REQUIRE ANY SERVICE CHECK DATA
func SimpleTCPCheck(s *Service) (bool, error) {
	tcpConnection, err := net.DialTimeout(s.TransportProtocol, fmt.Sprintf("%s:%d", s.Host, s.Port), time.Duration(5) * time.Second)

	// IF THERE'S AN ERROR, RETURN FALSE AND THE ERROR
	if err != nil {
		s.CheckFailed()
		return false, err
	} else {
		tcpConnection.Close()
		return true, nil
	}

}
