package dao

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// User which store User information
type User struct {
	UUID      string    `bson:"uuid,omitempty"`
	Name      string    `bson:"name,omitempty"`
	Email     string    `bson:"email,omitempty"`
	Password  string    `bson:"password,omitempty"`
	CreatedAt time.Time `bson:"created_at_time"`
}

// Session which store Session information
type Session struct {
	UUID      string    `bson:"uuid,omitempty"`
	Email     string    `bson:"email,omitempty"`
	UserID    string    `bson:"userid"`
	CreatedAt time.Time `bson:"created_at_time"`
}

// CreateSession create a new session for an existing user
func (user *User) CreateSession() (session Session, err error) {
	collection := client.Database(dbname).Collection("sessions")
	session = Session{UUID: createUUID(), Email: user.Email, UserID: user.UUID, CreatedAt: time.Now()}
	_, err = collection.InsertOne(context.Background(), session)
	return
}

// Session get the session for an existing user
func (user *User) Session() (session Session, err error) {
	collection := client.Database(dbname).Collection("sessions")
	filter := bson.M{"UserID": user.UUID}
	session = Session{}
	err = collection.FindOne(context.Background(), filter).Decode(&session)
	if err != nil {
		log.Fatal(err)
		return
	}
	return
}

// Check if session is valid in the database
func (session *Session) Check() (valid bool, err error) {
	collection := client.Database(dbname).Collection("sessions")
	filter := bson.M{"uuid": session.UUID}
	err = collection.FindOne(context.Background(), filter).Decode(session)
	if err != nil {
		valid = false
		return
	}
	if session.UUID != "" {
		valid = true
	}
	return
}

// DeleteByUUID delete session from database
func (session *Session) DeleteByUUID() (err error) {
	collection := client.Database(dbname).Collection("sessions")
	filter := bson.M{"UUID": session.UUID}
	_, err = collection.DeleteOne(context.Background(), filter)
	return
}

// SessionDeleteAll delete all sessions from database
func SessionDeleteAll() (err error) {
	collection := client.Database(dbname).Collection("sessions")
	_, err = collection.DeleteMany(context.Background(), bson.D{})
	return
}

// User get the user from the session
func (session *Session) User() (user User, err error) {
	collection := client.Database(dbname).Collection("users")

	user = User{}
	filter := bson.M{"uuid": session.UserID}
	err = collection.FindOne(context.Background(), filter).Decode(&user)
	return
}

// Create a new user, save user info into the database
func (user *User) Create() (err error) {
	collection := client.Database(dbname).Collection("users")
	user.UUID = createUUID()
	user.Password = Encrypt(user.Password)
	_, err = collection.InsertOne(context.Background(), user)
	return
}

// Delete user from database
func (user *User) Delete() (err error) {
	collection := client.Database(dbname).Collection("users")
	filter := bson.M{"uuid": user.UUID}
	_, err = collection.DeleteOne(context.Background(), filter)
	return
}

// Update user information in the database
func (user *User) Update() (err error) {
	collection := client.Database(dbname).Collection("users")
	filter := bson.M{"uuid": user.UUID}
	collection.FindOneAndUpdate(context.Background(), filter, user)
	return
}

// UserDeleteAll delete all users from database
func UserDeleteAll() (err error) {
	collection := client.Database(dbname).Collection("users")
	_, err = collection.DeleteMany(context.Background(), bson.D{})
	return err
}

// Users get all users in the database and returns it
func Users() (users []User, err error) {
	collection := client.Database(dbname).Collection("users")

	cur, _ := collection.Find(context.Background(), bson.D{})
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		user := User{}
		err = cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
			return
		}
		users = append(users, user)
	}
	if err = cur.Err(); err != nil {
		log.Fatal(err)
	}
	return
}

// UserByEmail get a single user given the email
func UserByEmail(email string) (user User, err error) {
	collection := client.Database(dbname).Collection("users")

	user = User{}
	filter := bson.M{"email": email}
	err = collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// UserByUUID get a single user given the UUID
func UserByUUID(uuid string) (user User, err error) {
	collection := client.Database(dbname).Collection("users")

	user = User{}
	filter := bson.M{"uuid": uuid}
	err = collection.FindOne(context.Background(), filter).Decode(&user)
	return
}
