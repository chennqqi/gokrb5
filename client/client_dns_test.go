// +build dns

package client

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"gopkg.in/jcmturner/gokrb5.v4/config"
	"gopkg.in/jcmturner/gokrb5.v4/keytab"
	"gopkg.in/jcmturner/gokrb5.v4/testdata"
	"testing"
)

func TestResolveKDC(t *testing.T) {
	//ns := os.Getenv("DNSUTILS_OVERRIDE_NS")
	//if ns == "" {
	//	os.Setenv("DNSUTILS_OVERRIDE_NS", testdata.TEST_NS)
	//}
	c, _ := config.NewConfigFromString(testdata.TEST_KRB5CONF)
	c.LibDefaults.DNSLookupKDC = true
	var cl Client
	cl.WithConfig(c)
	count, res, err := cl.resolveKDC(c.LibDefaults.DefaultRealm, true)
	if err != nil {
		t.Errorf("error resolving KDC via DNS TCP: %v", err)
	}
	assert.Equal(t, 5, count, "Number of SRV records not as expected: %v", res)
	assert.Equal(t, count, len(res), "Map size does not match: %v", res)
	t.Logf("res: %v", res)
	expected := []string{
		"kdc.test.gokrb5:88",
		"kdc1a.test.gokrb5:88",
		"kdc2a.test.gokrb5:88",
		"kdc1b.test.gokrb5:88",
		"kdc2b.test.gokrb5:88",
	}
	for _, s := range expected {
		var found bool
		for _, v := range res {
			if s == v {
				found = true
				break
			}
		}
		assert.True(t, found, "Record %s not found in results", s)
	}
	c.LibDefaults.DNSLookupKDC = false
	_, res, err = cl.resolveKDC(c.LibDefaults.DefaultRealm, true)
	if err != nil {
		t.Errorf("error resolving KDCs from config: %v", err)
	}
	assert.Equal(t, "10.80.88.88:88", res[1], "KDC not read from config as expected")
}

func TestClient_Login_DNSKDCs(t *testing.T) {
	//ns := os.Getenv("DNSUTILS_OVERRIDE_NS")
	//if ns == "" {
	//	os.Setenv("DNSUTILS_OVERRIDE_NS", testdata.TEST_NS)
	//}
	c, _ := config.NewConfigFromString(testdata.TEST_KRB5CONF)
	// Set to lookup KDCs in DNS
	c.LibDefaults.DNSLookupKDC = true
	//Blank out the KDCs to ensure they are not being used
	c.Realms = []config.Realm{}

	b, err := hex.DecodeString(testdata.TESTUSER1_KEYTAB)
	kt, _ := keytab.Parse(b)
	cl := NewClientWithKeytab("testuser1", "TEST.GOKRB5", kt)
	cl.WithConfig(c)

	err = cl.Login()
	if err != nil {
		t.Errorf("Error on logging in using DNS lookup of KDCs: %v\n", err)
	}
}