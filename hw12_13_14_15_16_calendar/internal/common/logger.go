//revive:disable
package common

import "time"

type LogEntry struct {
	IP        string
	Date      time.Time
	Path      string
	Proto     string
	Method    string
	UserAgent string
	Status    int
	Latency   int
}
