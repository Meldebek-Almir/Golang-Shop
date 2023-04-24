package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"shop/backend-api/internal/jwt"
	"shop/backend-api/internal/models"
)

func (c *Client) Login(w http.ResponseWriter, r *http.Request) {
	var credintails models.Credintails
	err := json.NewDecoder(r.Body).Decode(&credintails)
	if err != nil {
		log.Println(err)
		return
	}
	user, err := c.authService.IsValidUser(credintails.Identifier, credintails.Password)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expiry := time.Now().Add(time.Minute * 120)
	expiryStr := expiry.Format(time.RFC3339)
	claims := map[string]string{
		"exp":  expiryStr,
		"sub":  strconv.Itoa(user.ID),
		"role": user.Role,
	}

	jwtToken, err := jwt.CreateJWT(claims)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Value:   jwtToken,
		Name:    "jwt_token",
		Expires: expiry,
	})

	w.WriteHeader(200)
	fmt.Fprintf(w, "hello token: %s", jwtToken)
}

func (c *Client) Register(w http.ResponseWriter, r *http.Request) {
	var user models.UserRegister
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		return
	}

	if user.Role != "seller" && user.Role != "client" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = c.authService.CreateUser(&models.User{
		Nickname:  user.Nickname,
		Age:       user.Age,
		Gender:    user.Gender,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		Role:      user.Role,
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(200)
}

func (c *Client) GetInfoUser(w http.ResponseWriter, r *http.Request) {
	strId := r.Context().Value(keyUserType(keyUser))

	str, ok := strId.(string)
	if !ok {
		log.Println("empty strId")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userId, _ := strconv.Atoi(str)
	user, err := c.authService.GetUserById(userId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
