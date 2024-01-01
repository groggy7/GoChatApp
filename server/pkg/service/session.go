package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"server/pkg/model"
	"time"

	"github.com/redis/go-redis/v9"
)

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       1,
})

func GetClient() *redis.Client {
	return client
}

func CreateSession(username string) *model.Session {
	logger := log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags)

	session := &model.Session{
		Username:  username,
		Expiry:    time.Now().Add(3 * 24 * time.Hour),
		SessionID: fmt.Sprintf("%d", rand.Int63n(1000000000)),
	}

	sessionJSON, err := json.Marshal(session)
	if err != nil {
		logger.Println(err)
		panic(err)
	}

	ctx := context.Background()
	if err := client.Set(ctx, username, sessionJSON, 0).Err(); err != nil {
		log.Println(err)
		panic(err)
	}

	return session
}
