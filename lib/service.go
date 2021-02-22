package lib

import (
	"log"
	"net"
	"fmt"
)
// THE SERVICE TYPE
type Service struct {
	Host         		string 		// IP ADDRESS OR HOSTNAME
	Port         		int
	Name         		string 		// HUMAN-READABLE SERVICE NAME
	TransportProtocol   string 		// STRING, tcp OR udp. POSSIBLY NOT GOING TO BE USED.
	Username     		string
	Password     		string
	ServiceCheckType	string		// STRING THAT'S A SUPPORTED SERVICE CHECK TYPE
	ServiceCheckData	interface{}	// DEPENDS ON THE DATA REQUIRED BY THE CHECK, HENCE INTERFACE. MOST LIKELY TO BE
									//	NIL OR A STRUCT THOUGH
}

// DEFINE A TYPE OF FUNCTION THAT TAKES A SERVICE AS A PARAM
type ServiceCheck func(s* Service) (bool, error)


// CHECK IF A TCP CONNECTION CAN BE ESTABLISHED
func SimpleTCPCheck (s* Service) (bool, error) {
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

// CREATE A MAP THAT MAPS STRINGS TO SERVICECHECK FUNCTIONS
var checks = map[string] ServiceCheck {
	"tcp": SimpleTCPCheck,
}

// DISPATCH THE SERVICE CHECK DEPENDING ON THE TYPE OF CHECK SPECIFIED
func (s Service) DispatchServiceCheck() (bool, error) {

	// CHECK TO SEE IF THERE'S A SERVICE CHECK CALLBACK DEFINED FOR THE SERVICE TYPE
	callback, exists := checks[s.ServiceCheckType]
	if !exists {

		// IF NOT, FATAL ERROR
		log.Fatalf("Unable to perform check of type %s for service %s", s.ServiceCheckType, s.Name)
		return false, nil
	} else {

		// OTHERWISE, INVOKE THE CALLBACK AND RETURN
		return callback(&s)
	}
}




