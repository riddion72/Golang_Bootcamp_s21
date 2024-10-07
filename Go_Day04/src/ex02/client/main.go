package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type RequestBody struct {
	Money      int    `json:"money"`
	CandyType  string `json:"candyType"`
	CandyCount int    `json:"candyCount"`
}

type ResponseBody struct {
	Change int    `json:"change"`
	Thanks string `json:"thanks"`
	Error  string `json:"error"`
}

func err_code(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getClient() (client *http.Client) {
	cp := x509.NewCertPool()
	data, err := os.ReadFile("../cert/minica.pem")
	err_code(err)
	cp.AppendCertsFromPEM(data)
	c, err := tls.LoadX509KeyPair("../cert/localhost/cert.pem", "../cert/localhost/key.pem")
	err_code(err)
	certs := []tls.Certificate{c}

	if len(certs) == 0 {
		client = &http.Client{}
	} else {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					Certificates: certs,
					RootCAs:      cp,
					// InsecureSkipVerify: true,
				},
			},
		}
	}

	return
}

func getFlags() (request RequestBody) {
	candyType := flag.String("k", "", "Candy Type")
	candyCount := flag.Int("c", 0, "Count of Candy")
	moneyAmount := flag.Int("m", 0, "Amount of money")
	flag.Parse()
	request.CandyCount = *candyCount
	request.CandyType = *candyType
	request.Money = *moneyAmount
	return
}

func readServer(data []byte) (ans ResponseBody) {
	err := json.Unmarshal(data, &ans)
	err_code(err)

	return
}

func printAnswer(resp ResponseBody) {
	if resp.Error == "" {
		fmt.Printf("%s Your change is %d\n", resp.Thanks, resp.Change)
	} else {
		fmt.Printf("error: %s\n", resp.Error)
	}
}

func main() {
	client := getClient()
	reqBody := getFlags()
	jsonRequest, _ := json.Marshal(reqBody)

	resp, err := client.Post("https://localhost:3333/buy_candy", "application/json", bytes.NewBuffer(jsonRequest))
	err_code(err)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	err_code(err)

	ansBody := readServer(body)

	printAnswer(ansBody)

	// log.Printf("Status: %s  Body: %s\n", resp.Status, string(body))
}
