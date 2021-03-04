package service

// A MAP THAT LISTS ALL POSSIBLE ServiceChecks SO THAT THE FRONT END CAN USE IT TO
// CONSTRUCT OBJECTS. ALL ServiceChecks MUST BE REGISTERED HERE
var ServiceChecks = map[string]ServiceCheck {
	"LDAPUserQuery": LDAPUserQuery,
	"MinioLoginCheck": MinioLoginCheck,
	"MinioBucketCheck": MinioBucketCheck,
	"FTPCRUDCheck": FTPCRUDCheck,
	"HTTPGetContentsCheck": HTTPGetContentsCheck,
	"HTTPGetStatusCodeCheck": HTTPGetStatusCodeCheck,
	"POP3BasicAuthCheck": POP3BasicAuthCheck,
	"SMBListSharesCheck": SMBListSharesCheck,
	"SMTPSendCheck": SMTPSendCheck,
	"SSHLoginCheck": SSHLoginCheck,
	"SimpleTCPCheck": SimpleTCPCheck,
	"VNCConnectCheck": VNCConnectCheck,
}

var ServiceCheckParams = map[string][]string {
	"LDAPUserQuery": {"domainName"},
	"MinioLoginCheck": {},
	"MinioBucketCheck": {},
	"FTPCRUDCheck": {"filename"},
	"HTTPGetContentsCheck": {"fullURL", "expectedContents"},
	"HTTPGetStatusCodeCheck": {"relativeURL", "expectedCode"},
	"POP3BasicAuthCheck": {},
	"SMBListSharesCheck": {},
	"SMTPSendCheck": {"destinationUser"},
	"SSHLoginCheck": {},
	"SimpleTCPCheck": {},
	"VNCConnectCheck": {},
}