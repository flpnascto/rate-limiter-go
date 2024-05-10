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
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
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

			resp, err := client.Do(req)
			if err != nil {
				log.Println("Error making request:", err)
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println("Error reading response body:", err)
				return
			}

			fmt.Println("Http Code", resp.StatusCode, "message:", string(body), "at ", time.Now().Local())

		}()
	}
	wg.Wait()
}

func getIp() string {
	ips := []string{"192.168.0.1:1010", "192.168.0.2:1010", "192.168.0.3:1010"}
	randomIndex := rand.Intn(len(ips))
	return ips[randomIndex]
}
