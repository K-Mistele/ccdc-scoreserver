package service

import (
	"log"
	"net/http"
	"strconv"
)

// TYPE ServiceCheck
// ServiceCheckData: url, expectedCode
func HTTPGetStatusCodeCheck(s *Service) (bool, error) {
	var url, expectedCodeStr string
	var expectedCode int
	var response *http.Response
	var exists bool
	var err error

	log.Printf("Running HTTP GET Status Code check for %s\n", s.Name)

	// GET DATA FROM THE HASH MAP
	url, exists = s.ServiceCheckData["url"]
	if !exists {
		log.Printf("[ERROR] No url specified for service check for %s\n", s.Name)
		s.CheckFailed()
		return false, nil
	}
	expectedCodeStr, exists = s.ServiceCheckData["expectedCode"]
	if !exists {
		log.Printf("[ERROR] No expectedCode specified for serviceCheck for %s\n", s.Name)
		return false, nil
	}
	expectedCode, _ = strconv.Atoi(expectedCodeStr)

	// RUN THE GET REQUEST
	response, err = http.Get(url)
	if err != nil {
		s.CheckFailed()
		return false, err
	}

	// CHECK THE STATUS CODE
	if response.StatusCode == expectedCode {
		s.CheckPassed()
		return true, nil
	} else {
		s.CheckFailed()
		return false, nil
	}

}

// TYPE ServiceCheck
// ServiceCheckData: url, expectedContents
func HTTPGetContentsCheck(s *Service) (bool, error) {

	return true, nil
}

func HTTPPutStatusCodeCheck(s *Service) (bool, error) {
	return true, nil
}

func HTTPPostStatusCodeCheck(s *Service) (bool, error) {
	return true, nil
}

func HTTPPostRedirectCheck(s *Service) (bool, error) {
	return true, nil
}
