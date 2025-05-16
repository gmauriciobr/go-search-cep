package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	GET_CEP = "http://viacep.com.br/ws/%s/json/"
)

func main() {

	file, err := os.Create("cep.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := io.Writer(file)

	semaphore := make(chan struct{}, 100)
	line := make(chan string)
	var wg sync.WaitGroup

	go func() {
		for m := range line {
			writer.Write([]byte(fmt.Sprintf(m + "\n")))
		}
	}()

	for a := 5; a <= 9; a++ {
		fmt.Printf("Lopp A: %d\n", a)
		for b := 1; b <= 9; b = b + 1 {
			fmt.Printf("Lopp B: %d\n", b)
			for c := 1; c <= 9; c = c + 1 {
				fmt.Printf("Lopp C: %d\n", c)
				for d := 1; d <= 9; d = d + 1 {
					fmt.Printf("Lopp D: %d\n", d)
					for z := 100; z < 999; z = z + 1 {
						wg.Add(1)
						semaphore <- struct{}{}

						go func(aA int, bB int, cC int, dD int, zZ int) {
							defer wg.Done()
							defer func() { <-semaphore }()

							cepString := fmt.Sprintf("0%s%s%s%s%s", strconv.Itoa(aA), strconv.Itoa(bB), strconv.Itoa(cC), strconv.Itoa(dD), strconv.Itoa(zZ))

							str, err := searchCep(cepString)
							if err == nil && str != nil {
								line <- *str
							}
						}(a, b, c, d, z)
					}
				}
			}
		}
	}

	wg.Wait()
	close(line)

}

func searchCep(cep string) (*string, error) {
	var jsonResponse map[string]interface{}
	client := http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf(GET_CEP, cep))
	if err != nil {
		fmt.Println("Erro ao fazer requisição:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Erro na requisição:", resp.StatusCode)
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&jsonResponse)
	if jsonResponse["erro"] != nil {
		return nil, err
	}

	respBytes, _ := json.Marshal(jsonResponse)
	respSring := string(respBytes)

	return &respSring, nil
}
