package livesessiondbcache_test

import (
	"cache"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanCacheData(t *testing.T) {
	dataSource := MakeFakeDataSource(10, false)
	testCache := livesessiondbcache.MakeCache[User](dataSource)

	keys := dataSource.GetAllKeys()
	first := keys[0]
	user := dataSource.m[first]

	t.Run("can retrieve data", func(t *testing.T) {
		u, err := testCache.Get(first)
		assert.Nil(t, err)
		assert.Equal(t, user, *u)
	})

	t.Run("makes one request per user", func(t *testing.T) {
		for i := 10; i < 10; i++ {
			u, err := testCache.Get(first)
			assert.Nil(t, err)
			assert.Equal(t, user, *u)
		}
		assert.Equal(t, 1, dataSource.AccessCount())
	})

	t.Run("returns nil for invalid key", func(t *testing.T) {
		u, err := testCache.Get("random invalid key")
		assert.Nil(t, err)
		assert.Nil(t, u)
	})

	t.Run("retrieves correct data for every user", func(t *testing.T) {
		for k, expectedUser := range dataSource.m {
			u, err := testCache.Get(k)
			assert.Nil(t, err)
			assert.Equal(t, *u, expectedUser)
		}
	})

}

func TestCachePassesErrorsFromDataSource(t *testing.T) {
	dataSource := MakeFakeDataSource(10, true)
	testCache := livesessiondbcache.MakeCache[User](dataSource)

	keys := dataSource.GetAllKeys()
	for _, key := range keys {
		u, err := testCache.Get(key)
		assert.Nil(t, u)
		assert.NotNil(t, err)
	}

}

func TestConcurrentReads(t *testing.T) {
	testCases := []struct {
		name          string
		userAmount    int
		requestAmount int
	}{
		{"requested test case", 100, 10},
		{"big amount of users", 10000, 10},
		{"big amount of requests", 10, 10000},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, makeConcurrencyTest(testCase.userAmount, testCase.requestAmount))
	}
}

func makeConcurrencyTest(userAmount, requestAmount int) func(t *testing.T) {
	return func(t *testing.T) {
		dataSource := MakeFakeDataSource(uint32(userAmount), false)
		testCache := livesessiondbcache.MakeCache[User](dataSource)
		keys := dataSource.GetAllKeys()

		var wg sync.WaitGroup
		for i := range keys {
			for k := 0; k < requestAmount; k++ {
				wg.Add(1)
				go func(key string) {
					defer wg.Done()

					u, err := testCache.Get(key)
					assert.Nil(t, err)
					assert.NotNil(t, u)
				}(keys[i])
			}
		}
		wg.Wait()
		assert.Equal(t, userAmount, dataSource.AccessCount())
	}
}
