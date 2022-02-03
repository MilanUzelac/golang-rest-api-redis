package main

import (
	"encoding/json"
	"fmt"
	"log"

	"math/rand"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type Person struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
}

type PersonStore interface {
	Get(string) (Person, error)
	Set(string, Person) error
}

type redisStore struct {
	client *redis.Client
}

func NewRedisStore() PersonStore {
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

func (r redisStore) Set(id string, person Person) error {
	bs, err := json.Marshal(person)
	if err != nil {
		return errors.Wrap(err, "failed to save session to redis")
	}

	if err := r.client.Set(id, bs, 0).Err(); err != nil {
		return errors.Wrap(err, "failed to save session to redis")
	}
	return nil
}

func (r redisStore) Get(id string) (Person, error) {
	var person Person
	bs, err := r.client.Get(id).Bytes()
	if err != nil {
		return person, errors.Wrap(err, "failed to get session from redis")
	}

	if err := json.Unmarshal(bs, &person); err != nil {
		return person, errors.Wrap(err, "failed to unmarshall session data")
	}
	return person, nil

}

var dict = []int32("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateHashValue(n int) string {
	b := make([]int32, n)
	for i := range b {
		b[i] = dict[rand.Intn(len(dict))]
	}
	return string(b)
}

func main() {

	per1 := Person{
		FirstName: "John",
		LastName:  "Doe",
		Age:       31,
	}

	fmt.Println("Hello sessions")
	sessionStore := NewRedisStore()
	sessionStore.Set("person1", per1)
	fmt.Println(sessionStore.Get("person1"))
	rand := GenerateHashValue(20)
	fmt.Println("Random -> " + rand)
}
