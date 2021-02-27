package service

import (
	"fmt"
	"github.com/hirochachacha/go-smb2"
	"net"
	"time"
)

// TYPE ServiceCheck
// PARAMS NONE
func SMBListSharesCheck(s *Service) (bool, error) {

	if s.Username == "" {
		s.CheckFailedWithReason("SMB username is null - this is not supported!")
		return false, nil
	}

	// BUILD A TCP CONNECTION
	address := fmt.Sprintf("%s:%d", s.Host, s.Port)
	conn, err := net.DialTimeout("tcp", address, time.Duration(5) * time.Second)
	if err != nil {
		s.CheckFailedWithError(err)
		return false, err
	}
	defer conn.Close()

	// CONFIGURE THE SMB CONNECTION, SET UP AUTH
	smbDialer := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User: s.Username,
			Password: s.Password,
		},
	}

	// MAKE THE SMB CONNECTION W/ CONFIGURED AUTH
	smbConn, err := smbDialer.Dial(conn)
	if err != nil || smbConn == nil {
		s.CheckFailedWithError(err)
		return false, err
	}
	defer smbConn.Logoff()

	// TRY TO LIST SHARES - MAKE SURE WE'RE AUTHED
	_, err = smbConn.ListSharenames()
	if err != nil {
		s.CheckFailedWithError(err)
		return false, err
	}

	s.CheckPassed()
	return true, nil
}