package service

import (

	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// TYPE ServiceCheck
// ServiceCheckData: url, expectedCode
func HTTPGetStatusCodeCheck(s *Service) (bool, error) {
	var url, expectedCodeStr string
	var expectedCode int
	var response *http.Response
	var exists bool
	var err error


	// GET DATA FROM THE HASH MAP
	url, exists = s.ServiceCheckData["url"]
	if !exists {
		log.Errorf(" No url specified for service check for %s\n", s.Name)
		s.CheckFailed()
		return false, nil
	}
	expectedCodeStr, exists = s.ServiceCheckData["expectedCode"]
	if !exists {
		log.Errorf("No expectedCode specified for serviceCheck for %s\n", s.Name)
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

	var url, expectedContents string
	var response* http.Response
	var exists bool
	var err error


	// GET DATA FROM THE HASHMAP
	url, exists = s.ServiceCheckData["url"]
	if !exists {
		log.Errorf("No url specified for service check for %s\n", s.Name)
		s.CheckFailed()
		return false, nil
	}

	expectedContents, exists = s.ServiceCheckData["expectedContents"]
	if !exists {
		log.Errorf("no expectedContents specified for service check for %s\n", s.Name)
		s.CheckFailed()
		return false, nil
	}

	// MAKE THE REQUEST
	response, err = http.Get(url)
	if err != nil {
		s.CheckFailed()
		return false, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {

		log.Errorf("failed to read HTTP GET response in HTTPGetContentsCheck")
		s.CheckFailed()
		return false, err
	}
	data := string(body)

	if strings.Contains(data, expectedContents) {
		s.CheckPassed()
		return true, nil
	} else {
		s.CheckFailed()
		return false, nil
	}
}

func HTTPGetForm(s *Service) (bool, error) {

	log.Warningf("HTTPGetForm not implemented!")
	return true, nil
}

func HTTPPutStatusCodeCheck(s *Service) (bool, error) {
	log.Warningf("HTTPPutStatusCodeCheck not implemented")
	return true, nil
}

func HTTPPostStatusCodeCheck(s *Service) (bool, error) {

	log.Warningf("HTTPPostStatusCodeCheck is not implemented")
	return true, nil
}

func HTTPPostRedirectCheck(s *Service) (bool, error) {
	log.Warningf("HTTPPostRedirectCheck is not implemented")
	return true, nil
}
