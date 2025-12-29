package controller

import "math/rand/v2"

func getBackoffSeconds(baseBackoff, attempt int) int {
	jitter := rand.IntN(baseBackoff)
	backoff := baseBackoff*(1<<attempt) + jitter
	if backoff > maxBackoffSeconds {
		return maxBackoffSeconds
	}
	return backoff
}
