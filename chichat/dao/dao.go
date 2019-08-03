package dao

import (
	"context"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	host   = "mongodb://localhost:27017"
	dbname = "chitchat_db"
)

var client *mongo.Client

func init() {
	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(host))
	if err != nil {
		log.Fatal(err)
	}
	return
}

// create a random UUID with from RFC 4122
func createUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatalln("Cannot generate UUID", err)
	}

	// 0x40 is reserved variant from RFC 4122
	u[8] = (u[8] | 0x40) & 0x7F
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

// hash plaintext with SHA-1
func Encrypt(plainText string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plainText)))
	return
}
