package audit

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/spf13/viper"
)

func LoadTest() {
	loadScenarios := []int{10, 25, 50, 100, 500}
	for _, numRequests := range loadScenarios {
		fmt.Printf(" === Iniciando teste com %d requisições ===\n", numRequests)
		runLoadTest(numRequests)
		time.Sleep(setInterval() * time.Second)
	}
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

func runLoadTest(numRequests int) {
	var wg sync.WaitGroup

	var successCount, rateLimitedCount, errorCount int
	var successLock, rateLimitLock, errorLock sync.Mutex

	startTime := time.Now()

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := &http.Client{}
			req, err := http.NewRequest("GET", "http://localhost:8080/api/ping", nil)
			if err != nil {
				log.Println("Error creating request:", err)
				return
			}
			req.Header.Add("X-Forwarded-For", getIp())
			if withToken() {
				req.Header.Add("API_KEY", getToken())
			}

			resp, err := client.Do(req)
			if err != nil {
				errorLock.Lock()
				errorCount++
				errorLock.Unlock()
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusTooManyRequests {
				rateLimitLock.Lock()
				rateLimitedCount++
				rateLimitLock.Unlock()
			} else if resp.StatusCode == http.StatusOK {
				successLock.Lock()
				successCount++
				successLock.Unlock()
			} else {
				errorLock.Lock()
				errorCount++
				errorLock.Unlock()
			}

			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

		}()
	}
	wg.Wait()

	duration := time.Since(startTime)
	fmt.Printf("Teste concluído: %d requisições em %v\n", numRequests, duration)
	fmt.Printf("Sucesso: %d, Limitado: %d, Erros: %d\n", successCount, rateLimitedCount, errorCount)
}

func getIp() string {
	ips := []string{"192.168.0.1:1010", "192.168.0.2:1010", "192.168.0.3:1010"}
	randomIndex := rand.Intn(len(ips))
	return ips[randomIndex]
}

func getToken() string {
	tokens := []string{"abc123", "def456", "ghi789"}
	randomIndex := rand.Intn(len(tokens))
	return tokens[randomIndex]
}

func withToken() bool {
	randomIndex := rand.Intn(9) + 1
	return randomIndex%2 == 0
}
