package service

import (
	logging "github.com/op/go-logging"
	"sync"
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
	ServiceCheck      ServiceCheck           // STRING THAT'S A SUPPORTED SERVICE CHECK TYPE
	ServiceCheckData  map[string]interface{} // DEPENDS ON THE DATA REQUIRED BY THE CHECK, HENCE MAP
	Points            int
}

// SERVICE METHODS

// DISPATCH THE SERVICE CHECK DEPENDING ON THE TYPE OF CHECK SPECIFIED
func (s Service) DispatchServiceCheck(results *map[string]bool, wg *sync.WaitGroup){

	// RUN THE SERVICE CHECK
	serviceIsUp, _ := s.ServiceCheck(&s)

	// ADD THE RESULT TO THE MAP
	(*results)[s.Name] = serviceIsUp

	// NOTIFY THE WAITGROUP THAT WE'RE DONE
	wg.Done()
}

// LOG A FAILED SERVICE CHECK
func (s Service) CheckFailed() {
	log.Infof("Score check failed for service \"%s\" failed", s.Name)
}

func (s Service) CheckFailedWithError(err error) {
	log.Infof("Score check failed for service \"%s\" with error: %s", s.Name, err)
}

func (s Service) CheckFailedWithReason(msg string) {
	log.Infof("Score check failed for service \"%s\" with reason: %s", s.Name, msg)
}

// LOG A SUCCESSFUL CHECK
func (s Service) CheckPassed() {
	log.Infof("Score check succeeded for service \"%s\"", s.Name)
}
