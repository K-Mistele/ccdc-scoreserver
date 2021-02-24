package service

import (
	"fmt"
	ssh "github.com/helloyi/go-sshclient"
)

// TYPE ServiceCheck
// PARAMS NONE
func SSHLoginCheck(s *Service) (bool, error) {

	var address string
	var err error

	// DIAL UP A CONNECTION
	address = fmt.Sprintf("%s:%d", s.Host, s.Port)
	sshClient ,err := ssh.DialWithPasswd(address, s.Username, s.Password)
	if err != nil {
		s.CheckFailedWithError(err)
		return false, err
	}
	defer sshClient.Close()

	// TRY TO EXECUTE AN IDEMPOTENT COMMAND
	err = sshClient.Cmd("id").Cmd("whoami").Cmd("hostname").Run()
	if err != nil {
		s.CheckFailedWithError(err)
		return false, err
	} else {
		s.CheckPassed()
		return true, nil
	}

}