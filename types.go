package hilbert

import "time"

type Cache[T any] interface {
	Set(key string, value T, ttl ...time.Duration)
	Get(key string) (T, bool)
}

const DefaultTTL time.Duration = 0
