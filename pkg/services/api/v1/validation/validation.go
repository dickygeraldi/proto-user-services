package validation

import (
	"protoUserService/pkg/services/api/v1/logging"
	"regexp"
)

var re = regexp.MustCompile("select|insert|update|alter|delete")
var digitCheck = regexp.MustCompile(`^[0-9]+$`)

// function validation sql injection validation
func LoginRequest(dataIp, datetime, api, numberPhone, password string) (string, bool) {
	if api != "" && numberPhone != "" && password != "" {
		if re.MatchString(api) == true || re.MatchString(numberPhone) == true || re.MatchString(password) == true {
			go logging.SetLogging("Registration Account", dataIp, "SQL Injection", "Warning message", "warning", "SQL Injection in this request", datetime)
			return "Coba lagi nanti, status pendaftaran di pending", false
		} else {
			if digitCheck.MatchString(numberPhone) {
				if len(numberPhone) > 11 && len(numberPhone) < 15 {
					if digitCheck.MatchString(numberPhone) == true {
						if numberPhone[0:2] == "62" {
							return "Validasi berhasil", true
						} else {
							return "Number Phone start from 62 digit", false
						}
					} else {
						return "Number phone harus berisi nilai digit", false
					}
				} else {
					return "Number Phone length must be 11 to 15 digit", false
				}
			} else {
				return "nilai number phone harus berisi angka", false
			}
		}
	} else {
		return "Data harus diisi", false
	}
}

// Function for validation user services request
func RegistrationRequest(api, numberPhone, fullName, password, datetime, dataIp string) (string, bool) {
	if api != "" && numberPhone != "" && fullName != "" && password != "" {
		if re.MatchString(api) == true || re.MatchString(numberPhone) == true || re.MatchString(fullName) == true || re.MatchString(password) == true {
			go logging.SetLogging("Registration Account", dataIp, "SQL Injection", "Warning message", "warning", "SQL Injection in this request", datetime)
			return "Coba lagi nanti, status pendaftaran di pending", false
		} else {
			if digitCheck.MatchString(numberPhone) {
				if len(numberPhone) > 11 && len(numberPhone) < 15 {
					if digitCheck.MatchString(numberPhone) == true {
						if numberPhone[0:2] == "62" {
							return "Valid Data", true
						} else {
							return "Number Phone start from 62 digit", false
						}
					} else {
						return "Number phone harus berisi nilai digit", false
					}
				} else {
					return "Number Phone length must be 11 to 15 digit", false
				}
			} else {
				return "nilai number phone harus berisi angka", false
			}
		}
	} else {
		return "Semua data harus diisi", false
	}
}
