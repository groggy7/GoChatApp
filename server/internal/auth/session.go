package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Session struct {
	Username  string
	Expiry    time.Time
	SessionID string
}

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       1,
})

var sessions map[string]*Session

func GetSessions() map[string]*Session {
	return sessions
}

func InitSessionServer() {
	sessions = make(map[string]*Session)
	keys, err := client.Keys(context.Background(), "*").Result()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, key := range keys {
		value, err := client.Get(context.Background(), key).Result()
		if err != nil {
			fmt.Println(err)
			continue
		}

		var session Session
		if err := json.Unmarshal([]byte(value), &session); err != nil {
			log.Println(err)
		}

		sessions[key] = &session
	}
}

func CreateSession(username string) *Session {
	logger := log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags)

	session := &Session{
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

func GetSession(username string) (*Session, error) {
	sessionJSON, err := client.Get(context.Background(), username).Result()
	if err != nil {
		return nil, err
	}

	var session Session
	if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
		return nil, err
	}

	return &session, nil
}
