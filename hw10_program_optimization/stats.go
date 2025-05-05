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

	result := make(DomainStat)
	decoder := json.NewDecoder(r)

	var user User
	for {
		if err := decoder.Decode(&user); err == io.EOF {
			break
		}

		email := strings.ToLower(user.Email)
		at := strings.LastIndexByte(email, '@')
		if at == -1 {
			continue // Skip invalid email addresses
		}
		if at == len(email)-1 {
			continue // Skip emails without domain part
		}
		domainPart := email[at+1:]

		if strings.HasSuffix(domainPart, domain) {
			result[domainPart]++
		}
	}

	return result, nil
}
