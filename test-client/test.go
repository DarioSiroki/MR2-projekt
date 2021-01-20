package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
)

func main() {

	// CONNECT TO THE SERVER
	tcpAddress, err := net.ResolveTCPAddr("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Error resolving TCP address:", err.Error())
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddress)
	if err != nil {
		log.Fatal("Error connecting to the server:", err.Error())
	}

	// REGISTER A NEW USER WITH RANDOM EMAIL AND PUB KEY
	email := randomString(20)
	pubkey := randomString(20)
	user, _ := json.Marshal(User{Email: email, Pubkey: pubkey})
	key := make([]byte, 2)
	binary.BigEndian.PutUint16(key, 1)
	user = append(key, user...)

	_, err = conn.Write(user)
	if err != nil {
		log.Fatal("Error sending data to the server:", err.Error())
	}
	conn.Close()

	conn, err = net.DialTCP("tcp", nil, tcpAddress)
	if err != nil {
		log.Fatal("Error connecting to the server:", err.Error())
	}

}

func randomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}
