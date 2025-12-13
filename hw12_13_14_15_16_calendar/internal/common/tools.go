package common

import "strings"

func ExtractAddr(addr string) (res string) {
	res = addr
	addrSplit := strings.Split(addr, "@")
	if len(addrSplit) > 1 {
		res = addrSplit[len(addrSplit)-1]
	}
	return res
}
