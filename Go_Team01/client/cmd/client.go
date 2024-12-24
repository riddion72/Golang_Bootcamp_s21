package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"ex00/client/config"
	"ex00/client/internal/models"
	us "ex00/client/internal/usecase"
)

func init() {
	config.Flager.H = flag.String("H", "127.0.0.1", "Print host")
	config.Flager.P = flag.String("P", "8765", "Print port")
	flag.Parse()
	config.Flager.Addres = *config.Flager.H + ":" + *config.Flager.P

	us.TryConnect(true)
}

func main() {
	// go http.ListenAndServe(flager.addres, nil)
	defer func() {
		log.Println("server stopped")
		if err := recover(); err != nil {
			log.Println("Unknown panic happend: ", err)
		}
	}()
	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		for {
			ans, err := http.Get("http://" + config.Flager.Addres + "/HeartBeat")
			if err != nil {
				// log.Print("Error conecn, tryReconnect", err)
				us.TryReconnect()
				continue
			}
			var ansBeat models.HeartBeat
			body, err := io.ReadAll(ans.Body)
			if err != nil {
				log.Print(err)
				continue
			}
			json.Unmarshal(body, &ansBeat)
			if !us.CompeareAddresses(ansBeat.Address) {
				log.Println("Address changed...")
				config.HB = ansBeat
				us.PrintAddress()
			}
			// log.Println(config.HB.Address)
			time.Sleep(time.Second * 5)
		}
	}()
	for scanner.Scan() {
		input := scanner.Text()
		comand := strings.Split(input, " ")
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			continue
		}
		switch comand[0] {
		case "SET":
			err := us.MethodSet(comand[1:])
			us.PrintErorr(err)
			// fmt.Println("SET")
		case "GET":
			err := us.MethodGet(comand[1:])
			us.PrintErorr(err)
			// fmt.Println(comand[1:])
		case "DELETE":
			err := us.MethodDelete(comand[1:])
			us.PrintErorr(err)
			// fmt.Println("DELETE")
		default:
			if comand[0] != "" {
				fmt.Printf("Unknown command: %s\n", comand[0])
			}
		}
	}

}
