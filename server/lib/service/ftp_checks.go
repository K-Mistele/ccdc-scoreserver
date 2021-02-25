package service

import (
	"bytes"
	"fmt"
	"github.com/dutchcoders/goftp"
)

// TYPE ServiceCheck
// PARAMS filename (optional)
func FTPCRUDCheck (s *Service) (bool, error) {

	var address, filename  string
	var exists bool
	var ftp *goftp.FTP
	var err error

	address = fmt.Sprintf("%s:%d", s.Host, s.Port)
	filename, exists = s.ServiceCheckData["filename"].(string)
	if !exists {
		filename = "scoreserver-check.txt"
	}

	// CONNECT TO THE SERVER
	ftp, err = goftp.Connect(address)
	if err != nil {
		s.CheckFailedWithError(err)
		return false, err
	}
	defer ftp.Close()

	//LOGIN
	err = ftp.Login(s.Username, s.Password)
	if err != nil {
		s.CheckFailedWithError(err)
		return false, err
	}

	// PUT IN PASSIVE
	_, err = ftp.Pasv()

	// PUT A FILE
	data := bytes.NewBufferString("Scoreserver check")
	err = ftp.Stor(filename, data)
	if err != nil {
		s.CheckFailedWithError(err)
		return false, err
	}

	// DELETE THE FILE
	err = ftp.Dele(filename)
	if err != nil {
		s.CheckFailedWithError(err)
		return false, err
	}

	s.CheckPassed()
	return true, nil

}