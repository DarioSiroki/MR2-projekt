package main

import (
	"context"
	"encoding/binary"
	"encoding/json"

	"fmt"
	"io"
	"log"
	"net"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

var rdb *redis.Client

func registerUser(c net.Conn) string {
	user := User{}

	buf := make([]byte, 0, 4096) // big buffer
	tmp := make([]byte, 256)
	for {
		n, err := c.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		buf = append(buf, tmp[:n]...)
	}

	err := json.Unmarshal(buf, &user)

	if err != nil {
		log.Print(err)
		return "error"
	}

	err = rdb.Set(ctx, user.Email, user.Pubkey, 0).Err()
	if err != nil {
		panic(err)
	}

	return "success"
}

func getUsers(c net.Conn) string {
	// return user list
	return "success"
}

func getUserByEmail() {
	// return user by email
}

func handler(c net.Conn) {
	for {
		actionBytes := make([]byte, 2)
		c.Read(actionBytes)
		action := binary.BigEndian.Uint16(actionBytes)

		var result string
		switch action {
		case 1:
			result = registerUser(c)
			break
		case 2:
			result = getUsers(c)
			break
		}

		fmt.Println(result)
		c.Write([]byte(result))
		break
	}
	c.Close()
}

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	for {
		connection, err := ln.Accept()
		if err != nil {
			log.Fatal("Error accepting connection", err)
		}
		go handler(connection)
	}
}
