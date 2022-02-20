package service

import (
	"testing"

	"github.com/ivan-bokov/go-pdns/internal/stacktrace"
	"github.com/ivan-bokov/go-pdns/internal/storage/sqlite"
	"github.com/stretchr/testify/assert"
)

var service *Service

func init() {
	storage := sqlite.New(":memory:")
	err := storage.CreateTable()
	if err != nil {
		panic(stacktrace.Wrap(err))
	}
	service = New(storage, true)
}

func TestService_AddDomainKey(t *testing.T) {
	k1 := &KeyData{
		Flags:     257,
		Active:    true,
		Published: true,
		Content:   "Private-key-format: v1.2\\nAlgorithm: 5 (RSASHA1)\\nModulus: qpe9fxlN4dBT38cLPWtqljZhcJjbqRprj9XsYmf2/uFu4kA5sHYrlQY7H9lpzGJPRfOAfxShBpKs1AVaVInfJQ==\\nPublicExponent: AQAB\\nPrivateExponent: Ad3YogzXvVDLsWuAfioY571QlolbdTbzVlhLEMLD6dSRx+xcZgw6c27ak2HAH00iSKTvqK3AyeaK8Eqy/oJ5QQ==\\nPrime1: wo8LZrdU2y0xLGCeLhwziQDDtTMi18NEIwlx8tUPnhs=\\nPrime2: 4HcuFqgo7NOiXFvN+V2PT+QaIt2+oi6D2m8/qtTDS78=\\nExponent1: GUdCoPbi9JM7l1t6Ud1iKMPLqchaF5SMTs0UXAuous8=\\nExponent2: nzgKqimX9f1corTAEw0pddrwKyEtcu8ZuhzFhZCsAxM=\\nCoefficient: YGNxbulf5GTNiIu0oNKmAF0khNtx9layjOPEI0R4/RY=",
	}
	k2 := &KeyData{
		Flags:     256,
		Active:    true,
		Published: true,
		Content:   "Private-key-format: v1.2\\nAlgorithm: 5 (RSASHA1)\\nModulus: tY2TAMgL/whZdSbn2aci4wcMqohO24KQAaq5RlTRwQ33M8FYdW5fZ3DMdMsSLQUkjGnKJPKEdN3Qd4Z5b18f+w==\\nPublicExponent: AQAB\\nPrivateExponent: BB6xibPNPrBV0PUp3CQq0OdFpk9v9EZ2NiBFrA7osG5mGIZICqgOx/zlHiHKmX4OLmL28oU7jPKgogeuONXJQQ==\\nPrime1: yjxe/iHQ4IBWpvCmuGqhxApWF+DY9LADIP7bM3Ejf3M=\\nPrime2: 5dGWTyYEQRBVK74q1a64iXgaNuYm1pbClvvZ6ccCq1k=\\nExponent1: TwM5RebmWeAqerzJFoIqw5IaQugJO8hM4KZR9A4/BTs=\\nExponent2: bpV2HSmu3Fvuj7jWxbFoDIXlH0uJnrI2eg4/4hSnvSk=\\nCoefficient: e2uDDWN2zXwYa2P6VQBWQ4mR1ZZjFEtO/+YqOJZun1Y=",
	}
	assert.Equal(t, service.AddDomainKey("unit.test.", k1), nil)
	assert.Equal(t, service.AddDomainKey("unit.test.", k2), nil)
}

func TestService_CreateSlaveDomain(t *testing.T) {
	assert.Equal(t, service.CreateSlaveDomain("10.0.0.1", "example.com."), nil)
}

func TestFeedRecord(t *testing.T) {
	//service.StartTransaction()
	rr := &DNSResourceRecord{
		Qname:   "example.com.",
		Content: "ns1.example.com. hostmaster.example.com. 2013013441 7200 3600 1209600 300",
		TTL:     300,
		Qtype:   "SOA",
		Qclass:  "IN",
	}
	assert.Equal(t, service.FeedRecord(rr, ""), nil, "Не удалось записать")
	rr = &DNSResourceRecord{
		Qname:   "replace.example.com.",
		Content: "127.0.0.1",
		TTL:     300,
		Qtype:   "A",
		Qclass:  "IN",
	}
	assert.Equal(t, service.FeedRecord(rr, ""), nil, "Не удалось записать")
}
