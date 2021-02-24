package service

import (
	"fmt"
	adAuth "github.com/korylprince/go-ad-auth/v3"
	"strings"
)

// TYPE ServiceCheck
// PARAMS domainName - the AD domain
func LDAPUserQuery(s *Service) (bool, error) {

	domainName, exists := s.ServiceCheckData["domainName"].(string)
	if !exists {
		s.CheckFailedWithReason("No domain name was specified!")
		return false, nil
	}

	// BUILD THE BASE DN
	baseDN := fmt.Sprintf("OU=Users,DC=%s,DC=%s", strings.Split(domainName, ".")[0], strings.Split(domainName, ".")[1])

	// BUILD AN LDAP CONFIG
	config := &adAuth.Config{
		Server:   s.Host,
		Port:     s.Port,
		BaseDN:   baseDN,
		Security: adAuth.SecurityNone,
	}

	// AUTHENTICATE TO LDAP
	status, err := adAuth.Authenticate(config, s.Username, s.Password)
	if err != nil {
		s.CheckFailedWithError(err)
		return false, err
	}
	if !status {
		s.CheckFailedWithReason("Failed to connect!")
		return false, err
	}

	// IF WE MADE IT HERE THEN THE CHECK PASSED
	s.CheckPassed()
	return true, nil
}
