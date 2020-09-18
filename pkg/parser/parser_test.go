package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParserParseSingleIP(t *testing.T) {
	sampleData := `8.8.8.8:53 1600459170 NOERROR docs.bugbounty.com.  A
	docs.bugbounty.com. A 185.199.111.153`

	var domain string
	var ip []string
	var nxalias string
	var rcode string

	err := Parse(strings.NewReader(sampleData), func(Domain string, IP []string, NxAlias string, RCode string, Resolver string) {
		domain = Domain
		ip = IP
		nxalias = NxAlias
		rcode = RCode
	})
	require.Nil(t, err, "Could not parse sample data")
	require.Equal(t, "docs.bugbounty.com", domain, "Could not get domain")
	require.Equal(t, []string{"185.199.111.153"}, ip, "Could not get ip")
	require.Equal(t, "", nxalias, "Could not get nxalias")
	require.Equal(t, "NOERROR", rcode, "Could not get response code")
}

func TestParserParseMultipleDomains(t *testing.T) {
	sampleData := `
8.8.8.8:53 1600459170 NOERROR docs.bugbounty.com.  A
docs.bugbounty.com. A 185.199.111.153

8.8.8.8:53 1600459170 NOERROR docs.bugbounty.com.  A
docs.hackerone.com. A 185.199.111.152`

	var domain []string
	var ip []string
	err := Parse(strings.NewReader(sampleData), func(Domain string, IP []string, NxAlias string, RCode string, Resolver string) {
		domain = append(domain, Domain)
		ip = append(ip, IP[0])
	})
	require.Nil(t, err, "Could not parse sample data")
	require.Equal(t, []string{"docs.bugbounty.com", "docs.hackerone.com"}, domain, "Could not get domain")
	require.Equal(t, []string{"185.199.111.153", "185.199.111.152"}, ip, "Could not get ip")
}

func TestParserParseMultipleIPCNAME(t *testing.T) {
	sampleData := `
8.8.8.8:53 1600459170 NOERROR docs.bugbounty.com.  A
docs.hackerone.com. CNAME hacker0x01.github.io.
hacker0x01.github.io. A 185.199.111.153
hacker0x01.github.io. A 185.199.108.153
hacker0x01.github.io. A 185.199.109.153
hacker0x01.github.io. A 185.199.110.153`

	var domain string
	var ip []string
	err := Parse(strings.NewReader(sampleData), func(Domain string, IP []string, NxAlias string, RCode string, Resolver string) {
		domain = Domain
		ip = IP
	})
	require.Nil(t, err, "Could not parse sample data")
	require.Equal(t, "docs.hackerone.com", domain, "Could not get domain")
	require.Equal(t, []string{"185.199.111.153", "185.199.108.153", "185.199.109.153", "185.199.110.153"}, ip, "Could not get ip")
}

func TestParserParseMultipleCNAMEIP(t *testing.T) {
	sampleData := `
	8.8.8.8:53 1600459170 NOERROR docs.bugbounty.com.  A
	docs.bugbounty.com. CNAME bugbounty.github.io.
bugbounty.github.io. CNAME bugbounty-local.herokudns.io.
bugbounty-local.herokudns.io. A 185.199.111.153`

	var domain string
	var ip []string
	err := Parse(strings.NewReader(sampleData), func(Domain string, IP []string, NxAlias string, RCode string, Resolver string) {
		domain = Domain
		ip = IP
	})
	require.Nil(t, err, "Could not parse sample data")
	require.Equal(t, "docs.bugbounty.com", domain, "Could not get domain")
	require.Equal(t, []string{"185.199.111.153"}, ip, "Could not get ip")
}
