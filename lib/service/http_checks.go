package service

import (
	"bytes"
	"fmt"
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
	url, exists = s.ServiceCheckData["url"].(string)
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
	response, err = http.Get(fmt.Sprintf("http://%s:%d%s", s.Host, s.Port, url))
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
// ServiceCheckData: url, expectedContents
func HTTPGetContentsCheck(s *Service) (bool, error) {

	var url, expectedContents string
	var response* http.Response
	var exists bool
	var err error


	// GET DATA FROM THE HASHMAP
	url, exists = s.ServiceCheckData["url"].(string)
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

// TYPE ServiceCheck
// ServiceCheckData: url: string, data: callback(s* Service),
// 		headers map[string]string (should include content type and accept), expectedString
func HTTPPostContentsCheck(s *Service) (bool, error) {
	var url, expectedString string
	var data []byte
	var headers map[string]string
	var exists bool
	var err error
	var request* http.Request
	var client* http.Client
	var response* http.Response

	// GET DATA FROM THE HASH MAP
	url, exists = s.ServiceCheckData["url"].(string)
	if !exists {
		s.CheckFailedWithReason("No url specified")
		return false, nil
	}

	data, exists = s.ServiceCheckData["data"].([]byte)
	if !exists {
		data = make([]byte, 0)
	}

	headers = s.ServiceCheckData["headers"].(map[string]string)
	if !exists {
		// IF NO HEADERS EXIST, THEN MAKE IT AN EMPTY SLICE
		headers = make(map[string]string)
	}

	expectedString, exists = s.ServiceCheckData["expectedString"].(string)
	if !exists {
		s.CheckFailedWithReason("No expectedString to match was provided")
		return false, nil
	}

	// BUILD THE REQUEST
	request, err = http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		s.CheckFailedWithError(err)
		return false, err
	}

	// SET HEADERS
	for header, value := range headers {
		request.Header.Set(header, value)
	}

	// MAKE THE REQUEST
	client = &http.Client{}
	response, err = client.Do(request)
	if err != nil {
		s.CheckFailed()
		return false, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		s.CheckFailedWithError(err)
		return false, err
	}

	bodyStr := string(body)
	if strings.Contains(bodyStr, expectedString) {
		s.CheckPassed()
		return true, nil
	} else {
		s.CheckFailedWithReason("Response didn't contain expected string")
		return false, nil
	}


}

// TYPE ServiceCheck
// NOT IMPLEMENTED
func HTTPGetForm(s *Service) (bool, error) {

	log.Warningf("HTTPGetForm not implemented!")
	return true, nil
}

// TYPE ServiceCheck
// NOT IMPLEMENTED
func HTTPPutStatusCodeCheck(s *Service) (bool, error) {
	log.Warningf("HTTPPutStatusCodeCheck not implemented")
	return true, nil
}

// TYPE ServiceCheck
// NOT IMPLEMENTED
func HTTPPostForm(s *Service) (bool, error) {
	return true, nil
}

// TYPE ServiceCheck
// NOT IMPLEMENTED
func HTTPPostStatusCodeCheck(s *Service) (bool, error) {

	log.Warningf("HTTPPostStatusCodeCheck is not implemented")
	return true, nil
}

// TYPE ServiceCheck
// NOT IMPLEMENTED
func HTTPPostRedirectCheck(s *Service) (bool, error) {
	log.Warningf("HTTPPostRedirectCheck is not implemented")
	return true, nil
}
