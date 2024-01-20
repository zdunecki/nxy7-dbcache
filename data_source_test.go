package livesessiondbcache_test

import (
	livesessiondbcache "cache"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-faker/faker/v4"
)

type FakeDataSource struct {
	m           map[string]User
	mutex       sync.RWMutex
	accessCount atomic.Int64
	// flag indicating whether Get calls to this source should return errors
	shouldFail bool
}

var _ livesessiondbcache.DataSource[User] = &FakeDataSource{}

func MakeFakeDataSource(userAmount uint32, shouldFail bool) *FakeDataSource {
	f := FakeDataSource{
		m:          make(map[string]User),
		shouldFail: shouldFail,
	}
	for i := uint32(0); i < userAmount; i++ {
		f.m[fmt.Sprintf("%v", i)] = GenerateRandomUser()
	}
	return &f
}

// Get that simulates network latency
func (f *FakeDataSource) Get(key string) (*User, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	f.accessCount.Add(1)

	waitTime := rand.Uint32()%900 + 200
	time.Sleep(time.Millisecond * time.Duration(waitTime))

	if f.shouldFail {
		return nil, fmt.Errorf("simulated source error")
	}

	u, ok := f.m[key]
	if ok {
		return &u, nil
	} else {
		return nil, nil
	}
}

func (f *FakeDataSource) AccessCount() int {
	return int(f.accessCount.Load())
}
func (f *FakeDataSource) GetAllKeys() []string {
	var keys []string
	for k := range f.m {
		keys = append(keys, k)
	}
	return keys
}

// Example data retrieved by DataSource[User] used in tests
type User struct {
	Name string
	Age  uint32
}

func GenerateRandomUser() User {
	age := rand.Uint32() % 100
	name := faker.Name()
	return User{
		Age:  age,
		Name: name,
	}
}
