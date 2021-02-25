package service

import (
	"context"
	"fmt"
	vnc "github.com/kward/go-vnc"
	"github.com/kward/go-vnc/rfbflags"
	"net"
)

// TYPE ServiceCheck
// PARAMS NONE
func VNCConnectCheck(s *Service) (bool, error) {

	var address string
	var err error
	var conn net.Conn

	address = fmt.Sprintf("%s:%d", s.Host, s.Port)

	// BUILD A TCP CONNECTION
	conn, err = net.Dial("tcp", address)
	if err != nil {
		s.CheckFailedWithError(err)
		return false, nil
	}
	defer conn.Close()

	// NEGOTIATE A VNC CONNECTION. THIS WILL FAIL IF AUTH FAILS
	vncConfig := vnc.NewClientConfig(s.Password)	// IF SERVER HAS NO PASSWORD, THIS IS STILL FINE. SERVER WILL IGNORE IT
	vncClient, err := vnc.Connect(context.Background(), conn, vncConfig)
	if err != nil {
		s.CheckFailedWithError(err)
		return false, nil
	}

	// REQUEST A FRAME BUFFER UPDATE. THIS MAKES SURE THE PROTOCOL IS WORKING RIGHT
	w, h := vncClient.FramebufferWidth(), vncClient.FramebufferHeight()
	if err := vncClient.FramebufferUpdateRequest(rfbflags.RFBTrue, 0, 0, w, h); err != nil {
		s.CheckFailedWithError(err)
		return false, nil
	}

	s.CheckPassed()
	return true, nil
}
