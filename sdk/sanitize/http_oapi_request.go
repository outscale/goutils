//nolint
package sanitize

func init() {
	serviceFields[Request]["oapi"] = map[string]string{
		"BgpKey": Sensitive,
		"City": PII,
		"Country": PII,
		"Email": PII,
		"FirstName": PII,
		"JobTitle": PII,
		"LastName": PII,
		"Login": PII,
		"MobileNumber": PII,
		"NewUserEmail": PII,
		"NewUserName": PII,
		"Password": Sensitive,
		"PhoneNumber": PII,
		"PrivateKey": Sensitive,
		"SecretKey": Sensitive,
		"StateProvince": PII,
		"UserEmail": PII,
		"UserName": PII,
		"ZipCode": PII,
	}
}