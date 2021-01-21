package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
)

type User struct {
	Email  string `json:"email"`
	Pubkey string `json:"pubkey"`
}

type KeyRequest struct {
	Email string `json:"email"`
}

func main() {
	action, _ := strconv.Atoi(os.Args[1])

	// CONNECT TO THE SERVER
	tcpAddress, err := net.ResolveTCPAddr("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Error resolving TCP address:", err.Error())
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddress)
	if err != nil {
		log.Fatal("Error connecting to the server:", err.Error())
	}

	key := make([]byte, 2)
	binary.BigEndian.PutUint16(key, uint16(action))

	switch action {
	case 1:
		{
			// REGISTER A NEW USER WITH RANDOM EMAIL AND PUB KEY
			email := randomString(20) + "@" + randomString(8) + ".com"
			pubkey := randomString(20)
			fmt.Println(email, pubkey)
			user, _ := json.Marshal(User{Email: email, Pubkey: pubkey})
			user = append(key, user...)

			_, err = conn.Write(user)
			if err != nil {
				log.Fatal("Error sending data to the server:", err.Error())
			}

			break
		}
	case 2:
		{
			_, err = conn.Write(key)
			if err != nil {
				log.Fatal("Error sending data to the server:", err.Error())
			}
			break
		}
	case 3:
		{

			email := os.Args[2]
			send, _ := json.Marshal(KeyRequest{Email: email})
			send = append(key, send...)
			_, err = conn.Write(send)
			if err != nil {
				log.Fatal("Error sending data to the server:", err.Error())
			}
			break
		}
	}

	reply := make([]byte, 1024)
	_, err = conn.Read(reply)
	if err != nil {
		println("Write to server failed:", err.Error())
	}

	println("reply from server=", string(reply))

	conn.Close()
}

func randomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}
