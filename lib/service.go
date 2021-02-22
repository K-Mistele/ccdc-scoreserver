package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// DATA TYPES FOR THE SERVICE DATA TYPE

// THE SERVICE TYPE
type Service struct {
	Host              string // IP ADDRESS OR HOSTNAME
	Port              int
	Name              string // HUMAN-READABLE SERVICE NAME
	TransportProtocol string // STRING, tcp OR udp. POSSIBLY NOT GOING TO BE USED.
	Username          string
	Password          string
	ServiceCheckType  string            // STRING THAT'S A SUPPORTED SERVICE CHECK TYPE
	ServiceCheckData  map[string]string // DEPENDS ON THE DATA REQUIRED BY THE CHECK, HENCE INTERFACE. MOST LIKELY TO BE
	//	NIL OR A STRUCT THOUGH
	Points int
}

// DEFINE A TYPE OF FUNCTION THAT TAKES A SERVICE AS A PARAM
type ServiceCheck func(s *Service) (bool, error)

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

// MAKE A GET REQUEST AND CHECK FOR A STATUS CODE
// REQUIRES A url: string AND expectedCode: string AS SERVICE CHECK DATA
func HTTPGETStatusCode(s *Service) (bool, error) {

	log.Printf("Making HTTP GET status code service check for %s\n", s.Name)

	var url string
	var expectedCode int
	var exists bool
	var err error

	// MAKE SURE THERE'S A URL SPECIFIED
	url, exists = s.ServiceCheckData["url"]
	if !exists {
		log.Fatalf("No URL provided for HTTP status check for %s:%d (%s)\n", s.Host, s.Port, s.Name)
	}

	// CHECK FOR THE EXPECTED CODE, ASSUME 200 IF NOT SPECIFIED
	code, exists := s.ServiceCheckData["expectedCode"]
	if !exists {
		expectedCode = 200
	} else {
		expectedCode, err = strconv.Atoi(code)
		if err != nil {
			log.Fatalf("Bad 'expected HTTP status code' %s for service %s\n", code, s.Name)
		}
	}

	// MAKE A GET REQUEST
	response, err := http.Get(url)
	defer response.Body.Close()

	// IF WE CAN'T GET THE URL, THEN RETURN FALSE
	if err != nil {
		return false, err
	}

	if response.StatusCode == expectedCode {
		log.Printf("HTTP GET status code check for %s came back %t\n", s.Name, true)
		return true, nil
	} else {
		log.Printf("HTTP GET status code check for %s came back %t\n", s.Name, false)
		return false, nil
	}

}

// MAKE A GET REQUEST AND CHECK FOR CONTENT.
// REQUIRES url: string AND expectedContent: string AS SERVICE CHECK DATA
func HTTPGETContent(s *Service) (bool, error) {
	log.Printf("Making HTTP GET content service check for %s\n", s.Name)

	var url, content, expectedContent string
	var err error
	var exists bool
	var response *http.Response

	// MAKE SURE THERE'S A URL AND AN EXPECTED CONTENT SPECIFIED
	url, exists = s.ServiceCheckData["url"]
	if !exists {
		log.Fatalf("No URL provided for HTTP content check for service %s\n", s.Name)
	}

	// MAKE SURE THERE'S AN EXPECTED CONTENT VALUE
	expectedContent, exists = s.ServiceCheckData["expectedContent"]
	if !exists {
		log.Fatalf("No expected content provided for service %s\n", s.Name)
	}

	// MAKE THE GET REQUEST
	response, err = http.Get(url)
	if err != nil {
		log.Printf("HTTP GET content check for %s failed: %e\n", s.Name, err)
		return false, err
	}
	defer response.Body.Close()

	// READ IN THE BODY OF THE REQUEST AND THROW INTO A STRING
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("HTTP GET content check for %s failed: %e\n", s.Name, err)
		return false, err
	}
	content = string(body)

	// MAKE THE SERVICE CHECK
	if strings.Contains(content, expectedContent) {
		log.Printf("HTTP GET content check for %s succeeded!\n", s.Name)
		return true, nil
	} else {
		log.Printf("HTTP GET content check for %s failed!\n", s.Name)
		return false, nil
	}

}

// CREATE A MAP THAT MAPS STRINGS TO SERVICECHECK FUNCTIONS
var checks = map[string]ServiceCheck{
	"tcp":           SimpleTCPCheck,
	"http/get/code": HTTPGETStatusCode,
	"http/get/content": HTTPGETContent,
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
