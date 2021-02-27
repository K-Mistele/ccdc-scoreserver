package service

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// HELPER FUNCTION
func pop3ResponseIsGood(response *[]byte) bool {

	if strings.Contains(string(*response), "ERR") {
		return false
	} else {
		return true
	}

}

// TYPE ServiceCheck
// PARAMS NONE
// NOTE ONLY SUPPORTS AUTH WITH "USER" AND "PASS", DOES NOT SUPPORT TLS.
func POP3BasicAuthCheck(s *Service) (bool, error) {

	var conn net.Conn
	var err error
	var address string
	var response []byte

	// BUILD A TCP CONNECTION TO THE IMAP SERVER
	address = fmt.Sprintf("%s:%d", s.Host, s.Port)
	conn, err = net.DialTimeout("tcp", address, time.Duration(5) * time.Second)
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

	if !pop3ResponseIsGood(&response) {
		reason := fmt.Sprintf("Bad IMAP response - %s", string(response))
		s.CheckFailedWithReason(reason)
		return false, nil
	}

	// LOGIN
	user := fmt.Sprintf("USER %s\r\n", s.Username)
	fmt.Fprintf(conn, user)

	// GET LOGIN RESPONSE
	response = make([]byte, 1024)
	_, err = conn.Read(response[:])
	if err != nil {
		s.CheckFailedWithError(err)
		return false, err
	}

	if !pop3ResponseIsGood(&response) {
		reason := fmt.Sprintf("Bad IMAP response - %s", string(response))
		s.CheckFailedWithReason(reason)
		return false, nil
	}

	// FINISH LOGGING IN
	password := fmt.Sprintf("PASS %s\r\n", s.Password)
	fmt.Fprintf(conn, password)

	// GET LOGIN RESPONSE
	response = make([]byte, 1024)
	_, err = conn.Read(response[:])
	if err != nil {
		s.CheckFailedWithError(err)
		return false, err
	}

	if !pop3ResponseIsGood(&response) {
		reason := fmt.Sprintf("Bad IMAP response - %s", string(response))
		s.CheckFailedWithReason(reason)
		return false, nil
	}

	// LIST MAILS
	fmt.Fprintf(conn, "LIST\r\n")

	// GET LOGIN RESPONSE
	response = make([]byte, 1024)
	_, err = conn.Read(response[:])
	if err != nil {
		s.CheckFailedWithError(err)
		return false, err
	}

	if !pop3ResponseIsGood(&response) {
		reason := fmt.Sprintf("Bad IMAP response - %s", string(response))
		s.CheckFailedWithReason(reason)
		return false, nil
	}

	// SIGN OUT
	fmt.Fprintf(conn, "quit\r\n")

	// GET LOGIN RESPONSE
	response = make([]byte, 1024)
	_, err = conn.Read(response[:])
	if err != nil {
		s.CheckFailedWithError(err)
		return false, err
	}

	if !pop3ResponseIsGood(&response) {
		reason := fmt.Sprintf("Bad IMAP response - %s", string(response))
		s.CheckFailedWithReason(reason)
		return false, nil
	}

	s.CheckPassed()
	return true, nil

}
