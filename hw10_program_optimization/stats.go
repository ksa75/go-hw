package hw10programoptimization

import (
	"encoding/json"
	"io"
	"strings"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	domain = "." + strings.ToLower(domain) // Precompute and lowercase once

	result := make(DomainStat, 1000)
	decoder := json.NewDecoder(r)

	var user User
	for {
		if err := decoder.Decode(&user); err == io.EOF {
			break
		}

		email := strings.ToLower(user.Email)
		domainPart := email[strings.LastIndexByte(email, '@')+1:]

		if strings.HasSuffix(domainPart, domain) {
			result[domainPart]++
		}
	}

	return result, nil
}
