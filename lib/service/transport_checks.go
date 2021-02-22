package service

import (
	"fmt"
	"log"
	"net"
)

// CHECK IF A TCP CONNECTION CAN BE ESTABLISHED
// DOES NOT REQUIRE ANY SERVICE CHECK DATA
func SimpleTCPCheck(s *Service) (bool, error) {
	tcpConnection, err := net.Dial(s.TransportProtocol, fmt.Sprintf("%s:%d", s.Host, s.Port))

	// IF THERE'S AN ERROR, RETURN FALSE AND THE ERROR
	if err != nil {
		log.Printf("Error on service check %s -  %s:%s - %s\n", s.Name, s.Host, s.Port, err)
		return false, err
	} else {
		tcpConnection.Close()
		return true, nil
	}

}
