package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

//go:generate easyjson stats.go
//easyjson:json
type UserEmail struct {
	Email string `json:"Email"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	suffix := strings.ToLower("." + domain)
	for scanner.Scan() {
		line := scanner.Bytes()
		user, err := getUser(line)
		if err != nil {
			return nil, fmt.Errorf("get user error: %w", err)
		}

		lowerEmail := strings.ToLower(user.Email)
		if strings.HasSuffix(lowerEmail, suffix) {
			atIndex := strings.LastIndexByte(lowerEmail, '@')
			if atIndex == -1 || atIndex == len(lowerEmail)-1 {
				continue
			}
			result[lowerEmail[atIndex+1:]]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func getUser(lineBytes []byte) (UserEmail, error) {
	var email UserEmail
	err := email.UnmarshalJSON(lineBytes)
	return email, err
}
