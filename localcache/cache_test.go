package localcache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

const TTL = 3

type localCacheTestSuite struct {
	suite.Suite
	localcache *localCache
}

func (suite *localCacheTestSuite) SetupTest() {
	suite.localcache = New(TTL).(*localCache)
}

func (suite *localCacheTestSuite) TeardownTest() {
	suite.localcache = nil
}

func TestCacheTestSuite(t *testing.T) {
	suite.Run(t, new(localCacheTestSuite))
}

func (suite *localCacheTestSuite) TestCacheGetNilIfKeyDoesNotExist() {
	value, ok := suite.localcache.Get("test")

	suite.Require().Nil(value)
	suite.Require().False(ok)
}

func (suite *localCacheTestSuite) TestCacheGetValueIfKeyExists() {
	testCases := []struct {
		name	string
		key   string
		value interface{}
	} {
		{ name: "string", key: "stringKey", value: "value" },
		{ name: "int", key: "intKey", value: 1 },
		{ name: "float", key: "floatKey", value: 1.1 },
		{ name: "bool", key: "boolKey", value: true },
		{ name: "struct", key: "structKey", value: struct{ name string }{ name: "test" } },
		{ name: "slice", key: "sliceKey", value: []string{"test"} },
		{ name: "map", key: "mapKey", value: map[string]string{"test": "test"} },
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			suite.localcache.Set(tc.key, tc.value)

			value, ok := suite.localcache.Get(tc.key)

			suite.Require().Equal(tc.value, value)
			suite.Require().True(ok)
		})
	}
}

func (suite *localCacheTestSuite) TestCacheSetValueSuccessfully() {
	suite.localcache.data = map[string]interface{}{ "key": "value1" }

	suite.localcache.Set("key", "value2")
	suite.Require().Equal("value2", suite.localcache.data["key"])
}

func (suite *localCacheTestSuite) TestCacheDataExpiredAfterTTL() {
	suite.localcache.Set("key", "value")
	time.Sleep((TTL + 1) * time.Second)

	_, ok := suite.localcache.Get("key")
	suite.Require().False(ok)
}