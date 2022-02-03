package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type Session struct {
	VisitCount int `json:"visitCount"`
}

type Store interface {
	Get(string) (Session, error)
	Set(string, Session) error
}

type redisStore struct {
	client *redis.Client
}

func NewRedisStore() Store {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Falied to ping Redis: %v", err)
	}

	return &redisStore{
		client: client,
	}
}

func (r redisStore) Set(id string, session Session) error {
	bs, err := json.Marshal(session)
	if err != nil {
		return errors.Wrap(err, "failed to save session to redis")
	}

	if err := r.client.Set(id, bs, 0).Err(); err != nil {
		return errors.Wrap(err, "failed to save session to redis")
	}
	return nil
}

func (r redisStore) Get(id string) (Session, error) {
	var session Session
	bs, err := r.client.Get(id).Bytes()
	if err != nil {
		return session, errors.Wrap(err, "failed to get session from redis")
	}

	if err := json.Unmarshal(bs, &session); err != nil {
		return session, errors.Wrap(err, "failed to unmarshall session data")
	}
	return session, nil

}

func main() {
	fmt.Println("Hello sessions")
	sessionStore := NewRedisStore()
	sessionStore.Set("s1", Session{})
	fmt.Println(sessionStore.Get("s1"))
}
