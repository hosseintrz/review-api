package db

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrGeneratingToken = errors.New("error generating token")
)

type Store struct {
	Users       []*User
	Suggestions []*Suggestion
}

func NewStore() *Store {
	return &Store{
		Users:       []*User{},
		Suggestions: []*Suggestion{},
	}
}

type User struct {
	Username string
	Password string
	Token    string
}

type Suggestion struct {
	Username string `json:"user"`
	Text     string `json:"text"`
}

func (db *Store) GetUser(username string) (*User, error) {
	for _, user := range db.Users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

func (db *Store) AddUser(username, password string) (string, error) {
	token, err := generateToken(username)
	if err != nil {
		return "", ErrGeneratingToken
	}
	user := &User{
		Username: username,
		Password: password,
		Token:    token,
	}
	db.Users = append(db.Users, user)
	return token, nil
}

func (db *Store) GetUsernameFromToken(token string) (string, error) {
	for _, user := range db.Users {
		if user.Token == token {
			return user.Username, nil
		}
	}
	return "", ErrUserNotFound
}

func (db *Store) SubmitSuggest(username string, text string) {
	suggest := &Suggestion{username, text}
	db.Suggestions = append(db.Suggestions, suggest)
}

func (db *Store) GetSuggestions() []*Suggestion {
	return db.Suggestions
}

func generateToken(username string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(username), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	h := md5.New()
	h.Write(hash)
	return hex.EncodeToString(h.Sum(nil)), nil
}
