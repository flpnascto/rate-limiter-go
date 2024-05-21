package audit

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func LoadTest() {
	loadScenarios := []int{10, 25, 50, 100} // Diferentes níveis de carga
	for _, numRequests := range loadScenarios {
		fmt.Printf(" === Iniciando teste com %d requisições ===\n", numRequests)
		runLoadTest(numRequests)
		time.Sleep(10 * time.Second) // Intervalo entre cenários para evitar sobrecarga contínua
	}
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
				log.Println("Error making request:", err)
				errorLock.Lock()
				errorCount++
				errorLock.Unlock()
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println("Error reading response body:", err)
				errorLock.Lock()
				errorCount++
				errorLock.Unlock()
				return
			}

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

			fmt.Println("Http Code", resp.StatusCode, "message:", string(body), "at ", time.Now().Local())
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond) // Simula tráfego real com atraso aleatório

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
