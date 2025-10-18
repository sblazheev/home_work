package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
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
	u, err := getUsers(r, domain)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users []*User

func getUsers(r io.Reader, domain string) (result users, err error) {
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
		if strings.Contains(user.Email, domain) {
			result = append(result, &user)
		}
	}
	return
}

func unmarshal(data *[]byte, user *User) error {
	return user.UnmarshalJSON(*data)
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	for _, pUser := range u {
		user := *pUser
		if strings.Contains(user.Email, domain) {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}
	return result, nil
}
