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

func CreateSession(username string) *Session {
	logger := log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags)
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		logger.Println("Failed to connect to Redis:", err)
		return nil
	}

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

	val, err := client.Get(ctx, username).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(val)

	return session
}

func GetSession(username string) (*Session, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})

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

/*func (s *Session) isExpired() bool {
	return s.expiry.Before(time.Now())
}*/
