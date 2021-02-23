package service

import (

	logging "github.com/op/go-logging"
	"strings"
)

var log = logging.MustGetLogger("main")

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
	ServiceCheckData  map[string]interface{} // DEPENDS ON THE DATA REQUIRED BY THE CHECK, HENCE MAP
	Points int
}

// SERVICE METHODS

// DISPATCH THE SERVICE CHECK DEPENDING ON THE TYPE OF CHECK SPECIFIED
func (s Service) DispatchServiceCheck() (bool, error) {

	return s.ServiceCheck(&s)
}

// LOG A FAILED SERVICE CHECK
func (s Service) CheckFailed() {
	log.Infof("Score check for service \"%s\" failed\n", s.Name)
	s.checkFailed()
}

func (s Service) CheckFailedWithError(err error) {
	log.Infof("Score check for service \"%s\" failed with error: %s\n", s.Name, err)
	s.checkFailed()
}

func (s Service) CheckFailedWithReason(msg string) {
	log.Infof("Score check for service \"%s\" failed with reason: %s\n", s.Name, msg)
	s.checkFailed()
}

// LOG A SUCCESSFUL CHECK
func (s Service) CheckPassed() {
	log.Infof("Score check for service \"%s\" succeeded\n", s.Name)
	s.checkPassed()
}

// ACTUALLY UPDATE SCORING
func (s Service) checkFailed() {

}

// ACTUALLY UPDATE SCORING
func (s Service) checkPassed() {

}