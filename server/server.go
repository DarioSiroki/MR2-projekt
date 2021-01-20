package main

import (
	"context"
	"encoding/binary"
	"encoding/json"

	"fmt"
	"log"
	"net"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

var rdb *redis.Client

func registerUser(c net.Conn) string {
	user := User{}
	tmp := make([]byte, 256)
	n, err := c.Read(tmp)
	if err != nil {
		fmt.Println("read error:", err)
	}
	err = json.Unmarshal(tmp[:n], &user)

	if err != nil {
		log.Print(err)
		return "error"
	}

	fmt.Println("Registering ", user.Email, user.Pubkey)
	err = rdb.Set(ctx, user.Email, user.Pubkey, 0).Err()
	if err != nil {
		panic(err)
	}

	return "success"
}

func getUsers() string {
	// return user list
	keys := rdb.Keys(ctx, "*")
	json, _ := json.Marshal(keys.Val())
	return string(json)
}

func getKey(c net.Conn) string {
	keyreq := KeyRequest{}
	tmp := make([]byte, 256)
	n, err := c.Read(tmp)
	if err != nil {
		fmt.Println("read error:", err)
	}
	err = json.Unmarshal(tmp[:n], &keyreq)
	redisval := rdb.Get(ctx, keyreq.Email)
	if redisval.Err() == redis.Nil {
		return "No such e-mail"
	}
	return redisval.Val()

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
			result = getUsers()
			break
		case 3:
			result = getKey(c)
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
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	fmt.Println("TCP server starting")
	for {
		connection, err := ln.Accept()
		if err != nil {
			log.Fatal("Error accepting connection", err)
		}
		go handler(connection)
	}
}
