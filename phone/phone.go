package phone

import (
	"regexp"
	"strings"

	"github.com/etowett/go-api-sim/utils"
)

func IsValid(num string) (bool, string) {
	var valid bool
	bad := []string{"\t", "\n", " ", ",", "-", "(", ")", ".", "'", "\""}
	if num == "" || len(num) < 2 {
		return false, ""
	}
	for i := range bad {
		num = strings.Replace(num, bad[i], "", -1)
	}
	var number string
	if num[0:1] == "+" {
		valid, number = isInternational(num)
	} else {
		valid, number = isKenyan(num)
	}
	return valid, number
}

func isInternational(num string) (bool, string) {
	if num[1:4] == "254" {
		return isKenyan(num)
	} else {
		match, err := regexp.MatchString("^\\+{1}[0-9]{7,15}$", num)
		if err != nil {
			utils.Logger.Println("regexp err ", err)
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
