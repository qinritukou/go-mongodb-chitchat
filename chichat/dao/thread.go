package dao

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// Thread which store thread information
type Thread struct {
	UUID      string    `bson:"uuid,omitempty"`
	Topic     string    `bson:"topic,omitempty"`
	UserID    string    `bson:"userid"`
	CreatedAt time.Time `bson:"created_at_time"`
}

// Post which store post information
type Post struct {
	UUID      string    `bson:"uuid,omitempty"`
	Body      string    `bson:"body,omitempty"`
	ThreadID  string    `bson:"threadid"`
	UserID    string    `bson:"userid"`
	CreatedAt time.Time `bson:"created_at_time"`
}

// CreatedAtDate format the CreatedAt date to display nicely on the screen
func (thread *Thread) CreatedAtDate() string {
	return thread.CreatedAt.Format("2006-01-02 15:04:05")
}

// CreatedAtDate use the CreatedAtDate function
func (post *Post) CreatedAtDate() string {
	return post.CreatedAt.Format("2006-01-02 15:04:05")
}

// Posts get posts to a thread
func (thread *Thread) Posts() (posts []Post, err error) {
	collection := client.Database(dbname).Collection("posts")
	filter := bson.M{"threadid": thread.UUID}
	cur, _ := collection.Find(context.Background(), filter)
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		post := Post{}
		err = cur.Decode(&post)
		if err != nil {
			log.Fatal(err)
			return
		}
		posts = append(posts, post)
	}
	if err = cur.Err(); err != nil {
		log.Fatal(err)
	}
	return
}

// CreateThread create a new thread
func (user *User) CreateThread(topic string) (conv Thread, err error) {
	collection := client.Database(dbname).Collection("threads")

	threadToInsert := Thread{UUID: createUUID(), Topic: topic, UserID: user.UUID, CreatedAt: time.Now()}
	_, err = collection.InsertOne(context.Background(), threadToInsert)
	return
}

// CreatePost create a new post to a thread
func (user *User) CreatePost(conv Thread, body string) (post Post, err error) {
	collection := client.Database(dbname).Collection("posts")

	postToInsert := Post{UUID: createUUID(), Body: body, UserID: user.UUID, ThreadID: conv.UUID, CreatedAt: time.Now()}
	_, err = collection.InsertOne(context.Background(), postToInsert)
	return
}

// Threads Get all threads in the database and returns it
func Threads() (threads []Thread, err error) {
	collection := client.Database(dbname).Collection("threads")

	cur, _ := collection.Find(context.Background(), bson.D{})
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		thread := Thread{}
		err = cur.Decode(&thread)
		if err != nil {
			log.Fatal(err)
			return
		}
		threads = append(threads, thread)
	}
	if err = cur.Err(); err != nil {
		log.Fatal(err)
	}
	return
}

// ThraedByUUID get a thread by the uuid
func ThraedByUUID(uuid string) (conv Thread, err error) {
	collection := client.Database(dbname).Collection("threads")

	filter := bson.M{"uuid": uuid}
	conv = Thread{}
	err = collection.FindOne(context.Background(), filter).Decode(&conv)
	if err != nil {
		log.Fatal(err)
		return
	}
	return
}

// User get the user who started this thread
func (thread *Thread) User() (user User, err error) {
	collection := client.Database(dbname).Collection("users")

	filter := bson.M{"uuid": thread.UserID}
	user = User{}
	err = collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		log.Fatal(err)
		return
	}
	return
}

// User get the user who wrote the post
func (post *Post) User() (user User, err error) {
	collection := client.Database(dbname).Collection("users")

	filter := bson.M{"uuid": post.UserID}
	user = User{}
	err = collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		log.Fatal(err)
		return
	}
	return
}
