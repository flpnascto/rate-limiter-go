package audit_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/flpnascto/rate-limiter-go/audit"
	"github.com/spf13/viper"
)

func TestAudit(t *testing.T) {
	t.Log("TestAudit")

	numRequests := 100
	fmt.Printf(" === Iniciando teste com %d requisições ===\n", numRequests)
	audit.LoadTest(numRequests)
	time.Sleep(setInterval() * time.Second)
}

func setInterval() time.Duration {
	var interval time.Duration
	if viper.GetInt("IpBlockTime") > viper.GetInt("TokenBlockTime") {
		interval = time.Duration(viper.GetInt("IpBlockTime"))
	} else {
		interval = time.Duration(viper.GetInt("TokenBlockTime"))
	}
	return interval
}
