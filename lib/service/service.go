package service

import (
	"log"
)

// DEFINE A TYPE OF FUNCTION THAT TAKES A SERVICE AS A PARAM
type ServiceCheck func(s *Service) (bool, error)

// THE SERVICE TYPE
type Service struct {
	Host              string // IP ADDRESS OR HOSTNAME
	Port              int
	Name              string // HUMAN-READABLE SERVICE NAME
	TransportProtocol string // STRING, tcp OR udp. POSSIBLY NOT GOING TO BE USED.
	Username          string
	Password          string
	ServiceCheck      ServiceCheck      // STRING THAT'S A SUPPORTED SERVICE CHECK TYPE
	ServiceCheckData  map[string]string // DEPENDS ON THE DATA REQUIRED BY THE CHECK, HENCE MAP
	Points int
}

// SERVICE METHODS

// DISPATCH THE SERVICE CHECK DEPENDING ON THE TYPE OF CHECK SPECIFIED
func (s Service) DispatchServiceCheck() (bool, error) {

	return s.ServiceCheck(&s)
}

// LOG A FAILED SERVICE CHECK
func (s Service) CheckFailed() {
	log.Printf("Service check for %s failed\n", s.Name)
}

// LOG A SUCCESSFUL CHECK
func (s Service) CheckPassed() {
	log.Printf("Service check for %s succeeded\n", s.Name)
}