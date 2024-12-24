package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"sort"
	"strings"

	"ex00/client/config"
	"ex00/client/internal/models"

	"github.com/google/uuid"
)

func MethodGet(inPut []string) error {
	if len(inPut) > 1 {
		return errors.New("Too mach arguments")
	}
	if err := uuid.Validate(inPut[0]); err != nil {
		return errors.New("Key is not a proper UUID4")
	}
	var req models.GetRequestModel
	req.UUID4 = inPut[0]
	jsonRequest, _ := json.Marshal(req)

	return sender(jsonRequest, "/req_get")
}

func MethodSet(inPut []string) error {
	if len(inPut) < 2 {
		return errors.New("Invalid arguments")
	}
	if err := uuid.Validate(inPut[0]); err != nil {
		return errors.New("Key is not a proper UUID4")
	}
	var req models.SetRequestModel
	req.UUID4 = inPut[0]
	req.Value = strings.Join(inPut[1:], " ")
	jsonRequest, _ := json.Marshal(req)

	// fmt.Println(string(jsonRequest))
	return sender(jsonRequest, "/req_set")
}

func MethodDelete(inPut []string) error {
	if len(inPut) >= 2 {
		return errors.New("Too mach arguments")
	}
	if err := uuid.Validate(inPut[0]); err != nil {
		return errors.New("Key is not a proper UUID4")
	}
	var req models.GetRequestModel
	req.UUID4 = inPut[0]
	jsonRequest, _ := json.Marshal(req)

	return sender(jsonRequest, "/req_delete")
}

func PrintErorr(r error) {
	if r != nil {
		log.Println("Error: ", r)
	}
}

func sender(jsonRequest []byte, end string) error {
	resp, err := http.Post("http://"+config.Flager.Addres+end, "application/json", bytes.NewBuffer(jsonRequest))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// var body []byte
	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	return err
}

func CompeareAddresses(inPut []string) bool {
	bufferAddresses := config.HB.Address
	sort.Strings(inPut)
	sort.Strings(bufferAddresses)
	return reflect.DeepEqual(inPut, bufferAddresses)
}

func TryReconnect() {
	noNodes := true
	for _, address := range config.HB.Address {
		// log.Println("reconn >", address, config.Flager.Addres)
		if address != config.Flager.Addres {
			config.Flager.Addres = address
			noNodes = false
			break
		}
	}
	// log.Println("Trying to connect to >", config.Flager.Addres)
	if noNodes {
		log.Fatalln("No more available database nodes")
	}

	// Выводим сообщение о переподключении
	fmt.Printf("Reconnected to a database of Warehouse 13 at %s\n", config.Flager.Addres)
	TryConnect(false)
}

func TryConnect(fist bool) {
	r, err := http.Get("http://" + config.Flager.Addres + "/HeartBeat")
	if err != nil {
		log.Fatal(err)
	}
	json.NewDecoder(r.Body).Decode(&config.HB)
	if config.HB.MainAddres != config.Flager.Addres {
		config.Flager.Addres = config.HB.MainAddres
		// log.Println("Main >", config.Flager.Addres)
		TryConnect(false)
		return
	}
	if fist {
		fmt.Println("Connected to a database of Warehouse 13 at ", config.Flager.Addres)
	}
	PrintAddress()
}

func PrintAddress() {
	fmt.Println("Known nodes:")
	for _, node := range config.HB.Address {
		fmt.Println(node)
	}
	if config.HB.Replication > len(config.HB.Address) {
		fmt.Printf("WARNING: cluster size (%d) is smaller than a replication factor (%d)!\n", len(config.HB.Address), config.HB.Replication)
	}
}
