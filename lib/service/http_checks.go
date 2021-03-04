package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// TYPE ServiceCheck
// ServiceCheckData: relativeURL, expectedCode
func HTTPGetStatusCodeCheck(s *Service) (bool, error) {
	var relativeURL, expectedCodeStr string
	var expectedCode int
	var response *http.Response
	var exists bool
	var err error


	// GET DATA FROM THE HASH MAP
	relativeURL, exists = s.ServiceCheckData["relativeURL"].(string)
	if !exists {
		reason := "No reason specified"
		s.CheckFailedWithReason(reason)
		return false, nil
	}
	expectedCodeStr, exists = s.ServiceCheckData["expectedCode"].(string)
	if !exists {
		reason := "No status code was specified with expectedCode"
		s.CheckFailedWithReason(reason)
		return false, nil
	}
	expectedCode, _ = strconv.Atoi(expectedCodeStr)

	// RUN THE GET REQUEST
	response, err = http.Get(fmt.Sprintf("http://%s:%d%s", s.Host, s.Port, relativeURL))
	if err != nil {
		reason := fmt.Sprintf("Couldn't GET the service - %s", err)
		s.CheckFailedWithReason(reason)
		return false, err
	}

	// CHECK THE STATUS CODE
	if response.StatusCode == expectedCode {
		s.CheckPassed()
		return true, nil
	} else {
		s.CheckFailedWithReason("Status code mismatch!")
		return false, nil
	}

}

// TYPE ServiceCheck
// ServiceCheckData: fullURL, expectedContents
func HTTPGetContentsCheck(s *Service) (bool, error) {

	var url, expectedContents string
	var response* http.Response
	var exists bool
	var err error


	// GET DATA FROM THE HASHMAP
	url, exists = s.ServiceCheckData["fullURL"].(string)
	if !exists {
		s.CheckFailedWithReason("No url specified")
		return false, nil
	}

	expectedContents, exists = s.ServiceCheckData["expectedContents"].(string)
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
		s.CheckFailedWithReason("response didn't contain expected string")
		return false, nil
	}
}
