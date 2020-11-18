package main

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var admin int
var tokens = make(map[string]*TokenDetails)

//CreateToken
func CreateToken(userId uint64, person string) (*TokenDetails, error) {
	td := &TokenDetails{
		AtExpires: time.Now().Add(time.Minute * 20),
		person: person,
		ID: userId,
		}
	_ = os.Setenv("ACCESS_SECRET", "jdnfksdmfksd")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["exp"] = time.Now().Add(time.Minute * 20).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	var err error
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func checkToken(token ,person string) error{
	for _, value := range tokens {
		if value.AccessToken==token && person== value.person{
			if time.Now().Before(value.AtExpires) {
				return nil
			}else {
				delete(tokens, string(value.ID))
				return errors.New("Expired token ")
			}
		}
	}
	return errors.New("no such token ")
}

func checkAdmin(pass string)  error{
	content, err := ioutil.ReadFile("AdminPass.txt")
	if err!=nil{
		log.Fatal("Error in reading user & password file :", err.Error())
	}
	if pass==string(content) {
		return nil
	}
	return errors.New("Wrong password ")
}