package service

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// HELPER FUNCTION - IS STATUS CODE GOOD OR BAD
func smtpResponseIsGood(response *[]byte) bool {

	// PARSE OUT THE STATUS CODE
	responseStr := string(*response)
	statusCodeStr := strings.Split(responseStr, " ")[0]
	statusCode, err := strconv.Atoi(statusCodeStr)
	if err != nil {
		return false
	}

	// CHECK IF THE STATUS CODE IS IN THE RANGE OF GOOD STATUS CODES
	if statusCode >= 200 && statusCode < 400 {
		return true
	} else {
		return false
	}


}

// TYPE ServiceCheck
// PARAMS destinationUser
func SMTPSendCheck (s *Service) (bool, error) {

	var conn net.Conn
	var err error
	var address string
	var response []byte

	// BUILD A TCP CONNECTION TO THE SERVER
	address = fmt.Sprintf("%s:%d", s.Host, s.Port)
	conn, err  = net.Dial("tcp", address)
	if err != nil {
		s.CheckFailedWithError(err)
		return false, err
	}
	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.Debug(err)
		s.CheckFailedWithError(err)
		return false, err
	}
	defer conn.Close()


	// GET THE BANNER
	response = make([]byte, 1024)
	_, err = conn.Read(response[:])
	if err != nil {
		s.CheckFailedWithError(err)
		return false, err
	}

	// SEND THE HELO
	helo := "helo blacklight-scoreserver.local\r\n"
	fmt.Fprintf(conn, helo)

	// GET THE HELO RESPONSE
	response = make([]byte, 1024)
	_, err = conn.Read(response[:])
	if err != nil {
		s.CheckFailedWithError(err)
		return false, err
	}
	if !smtpResponseIsGood(&response) {
		reason := "Bad SMTP status code: " + string(response)
		s.CheckFailedWithReason(reason)
		return false, nil
	}

	// SET SENDER
	mailFrom := "MAIL FROM: admin@blacklight-scoreserver.local\r\n"
	fmt.Fprintf(conn, mailFrom)

	// GET THE SENDER RESPONSE
	response = make([]byte, 1024)
	_, err = conn.Read(response[:])
	if err != nil {
		s.CheckFailedWithError(err)
		return false, err
	}
	if !smtpResponseIsGood(&response) {
		reason := "Bad SMTP status code: " + string(response)
		s.CheckFailedWithReason(reason)
		return false, nil
	}

	// SET DESTINATION
	destinationUser, exists := s.ServiceCheckData["destinationUser"].(string)
	if !exists {
		destinationUser = "root"
	}
	mailTo := fmt.Sprintf("RCPT TO: %s\r\n", destinationUser)
	fmt.Fprintf(conn, mailTo)

	// GET THE DESTINATION RESPONSE
	response = make([]byte, 1024)
	_, err = conn.Read(response[:])
	if err != nil {
		s.CheckFailedWithError(err)
		return false, err
	}
	if !smtpResponseIsGood(&response) {
		reason := "Bad SMTP status code: " + string(response)
		s.CheckFailedWithReason(reason)
		return false, nil
	}

	// SET DATA
	data := "DATA\r\n"
	fmt.Fprintf(conn, data)

	// GET DATA RESPONSE
	response = make([]byte, 1024)
	_, err = conn.Read(response[:])
	if err != nil {
		s.CheckFailedWithError(err)
		return false, err
	}
	if !smtpResponseIsGood(&response) {
		reason := "Bad SMTP status code: " + string(response)
		s.CheckFailedWithReason(reason)
		return false, nil
	}

	// SEND MESSAGE
	sendData, exists := s.ServiceCheckData["msgData"].(string)
	if !exists {
		sendData = "Scoreserver check"
	}
	message := fmt.Sprintf("%s\r\n.\r\n", sendData)
	fmt.Fprintf(conn, message)

	// GET RESPONSE
	response = make([]byte, 1024)
	_, err = conn.Read(response[:])
	if err != nil {
		s.CheckFailedWithError(err)
		return false, err
	}
	if !smtpResponseIsGood(&response) {
		reason := "Bad SMTP status code: " + string(response)
		s.CheckFailedWithReason(reason)
		return false, nil
	}

	// PULL OUT THE MESSAGE ID SO WE CAN USE IT LATE

	// CLOSE THE CONNECTION
	fmt.Fprintf(conn, "quit\r\n")

	// DONE
	response = make([]byte, 1024)
	_, err = conn.Read(response[:])
	if err != nil {
		s.CheckFailedWithError(err)
		return false, err
	}

	s.CheckPassed()
	return true, nil
}
