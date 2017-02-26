package phone

import (
	"regexp"
	"strings"
)

func IsValid(num string) (bool, string) {
	bad := []string{"\t", "\n", " ", ",", "-", "(", ")", ".", "'", "\""}
	if num == "" || len(num) < 2 {
		return false, ""
	}

	for i := range bad {
		num = strings.Replace(num, bad[i], "", -1)
	}

	var valid bool
	var number string
	if num[0:1] == "+" {
		valid, number = isInternational(num)
	} else {
		valid, number = isKenyan(num)
	}
	return valid, number
}

func isInternational(num string) (bool, string) {
	if num[0:3] == "254" {
		valid, n := isKenyan(num)
		return valid, n
	} else {
		pattern := "^+{1}[0-9]{7,13}$"
		match, _ := regexp.MatchString(pattern, num)
		if match == false {
			return false, ""
		}
	}
	return true, num
}

func isKenyan(n string) (bool, string) {
	pattern := "^[0]{1}[7]{1}[0-9]{8}$|^[7]{1}[0-9]{8}$|^+254[7]{1}[0-9]{8}$|^254[7]{1}[0-9]{8}$"
	match, _ := regexp.MatchString(pattern, n)
	if match == false {
		return false, ""
	}
	if n[0:1] == "+" {
		n = n[1:]
	}
	if n[0:3] != "254" {
		if n[0:1] == "0" {
			n = "254" + n[1:]
		} else {
			n = "254" + n
		}
	}
	if len(n) != 12 {
		return false, ""
	}
	num := "+" + n
	return true, num
}
