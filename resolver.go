package dynamotest

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// TableNameResolver defines an interface for resolving a table name
type TableNameResolver interface {
	Resolve(tableName string) string
}

// DefaultTableNameResolver always returns the same name that is provided as aan input
type DefaultTableNameResolver struct{}

func (DefaultTableNameResolver) Resolve(tableName string) string {
	return tableName
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const defaultSuffixLen = 5

// RandomTableNameResolver attaches randomly generated suffix to the table name
type RandomTableNameResolver struct {
	SuffixLen int
	Seed      int64
}

func NewRandomTableNameResolver() *RandomTableNameResolver {
	return &RandomTableNameResolver{}
}

func (r *RandomTableNameResolver) Resolve(tableName string) string {
	var suffixLen = r.SuffixLen
	if suffixLen == 0 {
		suffixLen = defaultSuffixLen
	}

	return fmt.Sprintf("%s_%s", tableName, randStringBytes(suffixLen, r.Seed))
}

func randStringBytes(n int, seed int64) string {
	nseed := seed
	if nseed == 0 {
		nseed = time.Now().UnixNano()
	}

	rand.Seed(nseed)
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

type TimestampTableNameResolver struct {
	clock Clock
}

func NewTimestampTableNameResolver(clock Clock) *TimestampTableNameResolver {
	return &TimestampTableNameResolver{clock: clock}
}

func (r *TimestampTableNameResolver) Resolve(tableName string) string {
	return fmt.Sprintf("%s_%d", tableName, r.clock.Time().UnixNano())
}

type MemoizedTableNameResolver struct {
	resolver   TableNameResolver
	localCache map[string]string
	mutex      sync.Mutex
}

func NewMemoizedTableNameResolver(resolver TableNameResolver) *MemoizedTableNameResolver {
	return &MemoizedTableNameResolver{resolver: resolver, localCache: make(map[string]string)}
}

func (r *MemoizedTableNameResolver) Resolve(tableName string) string {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if resolver, ok := r.localCache[tableName]; ok {
		return resolver
	}

	r.localCache[tableName] = r.resolver.Resolve(tableName)
	return r.localCache[tableName]
}

type Clock interface {
	Time() time.Time
}

type FakeClock struct {
	FrozenTime time.Time
}

func (c FakeClock) Time() time.Time {
	return c.FrozenTime
}

type RealClock struct{}

func (RealClock) Time() time.Time {
	return time.Now()
}
