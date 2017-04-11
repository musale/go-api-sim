package phone

import (
	"errors"
	"regexp"
	"strings"
)

// CheckValid returns whether phone is valid
func CheckValid(num string) (string, error) {
	bad := []string{"\t", "\n", " ", ",", "-", "(", ")", ".", "'", "\""}
	if num == "" || len(num) < 5 {
		return "", errors.New("number too short")
	}
	for i := range bad {
		num = strings.Replace(num, bad[i], "", -1)
	}
	var number string
	if num[0:1] == "+" {
		number, err = isInternational(num)
		if err != nil {
			return "", err
		}
	} else {
		number, err = isKenyan(num)
		if err != nil {
			return "", err
		}
	}
	return number, nil
}

func isInternational(num string) (bool, string) {
	if num[1:4] == "254" {
		return isKenyan(num)
	} else {
		match, err := regexp.MatchString("^\\+{1}[0-9]{7,15}$", num)
		if err != nil {
			return false, ""
		}
		if match == false {
			return false, ""
		}
	}
	return true, num
}

func isKenyan(n string) (bool, string) {
	pattern := "^[7]{1}[0-9]{8}$"

	if n[0:1] == "+" || n[0:1] == "0" {
		n = n[1:]
	}
	if n[0:3] == "254" {
		n = n[3:]
	}
	if n[0:1] == "0" {
		n = n[1:]
	}
	match, _ := regexp.MatchString(pattern, n)
	if match == false {
		return false, ""
	}
	num := "+254" + n
	return true, num
}
