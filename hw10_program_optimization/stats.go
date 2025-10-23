package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	re, err := regexp.Compile("\\." + domain + "$")
	if err != nil {
		return nil, err
	}
	u, err := getUsers(r, re)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, re)
}

type users []*User

func getUsers(r io.Reader, re *regexp.Regexp) (result users, err error) {
	if err != nil {
		return
	}
	result = make(users, 0, 100_000)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		bytes := scanner.Bytes()
		user := User{}
		if err = unmarshal(&bytes, &user); err != nil {
			return
		}
		matched := re.Match([]byte(user.Email))

		if matched {
			result = append(result, &user)
		}
	}
	return
}

func unmarshal(data *[]byte, user *User) error {
	return user.UnmarshalJSON(*data)
}

func countDomains(u users, re *regexp.Regexp) (DomainStat, error) {
	result := make(DomainStat)
	for _, pUser := range u {
		user := *pUser
		matched := re.Match([]byte(user.Email))

		if matched {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}
	return result, nil
}
