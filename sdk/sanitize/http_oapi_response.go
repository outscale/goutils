//nolint
package sanitize

func init() {
	serviceFields[Response]["oapi"] = map[string]string{
		"AdminPassword": Sensitive,
		"BgpKey": Sensitive,
		"City": PII,
		"Country": PII,
		"Email": PII,
		"FirstName": PII,
		"JobTitle": PII,
		"LastName": PII,
		"MobileNumber": PII,
		"PhoneNumber": PII,
		"PrivateKey": Sensitive,
		"SecretKey": Sensitive,
		"StateProvince": PII,
		"UserEmail": PII,
		"UserId": PII,
		"UserName": PII,
		"ZipCode": PII,
	}
}