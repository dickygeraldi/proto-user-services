package validation

import (
	"protoUserService/pkg/services/api/v1/logging"
	"regexp"
)

var re = regexp.MustCompile("select|insert|update|alter|delete")
var digitCheck = regexp.MustCompile(`^[0-9]+$`)
var checkUsername = regexp.MustCompile("^[a-z0-9]+(?:_[a-z0-9]+)*$")

// Function for validation user services request
func RegistrationRequest(api, username, password, datetime, dataIp string) (string, bool) {
	if api != "" && username != "" && password != "" {
		if re.MatchString(api) == true || re.MatchString(username) || re.MatchString(password) == true {
			go logging.SetLogging("Registration Account", dataIp, "SQL Injection", "Warning message", "warning", "SQL Injection in this request", datetime)
			return "Coba lagi nanti, status pendaftaran di pending", false
		} else {
			if len(username) >= 4 && len(username) <= 10 {
				if checkUsername.MatchString(username) == true {
					return "Validasi berhasil", true
				} else {
					return "Username htidak valid", false
				}
			} else {
				return "nilai username harus lebih dari 4", false
			}
		}
	} else {
		return "Semua data harus diisi", false
	}
}
