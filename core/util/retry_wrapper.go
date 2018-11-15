package util

import (
	"github.com/wondywang/rpclookup/core/lg"
	"math/rand"
	"time"
)

//support

var logger *lg.Logger

func init() {
	rand.Seed(time.Now().UnixNano())
}

//most retry time: 3^attempts * sleep
func Retry(attempts int, sleep time.Duration, f func() error) error {
	if err := f(); err != nil {
		if s, ok := err.(stop); ok {
			// Return the original error for later checking
			return s.error
		}
		if attempts--; attempts > 0 {
			// Add some randomness to prevent creating a Thundering Herd
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter/2
			time.Sleep(sleep)
			return Retry(attempts, 2*sleep, f)
		}
		return err
	}
	return nil
}

type stop struct {
	error
}
