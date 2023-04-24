package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"shop/frontend-client/internal/models"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	productResp, err := http.Get("http://localhost:9090/api/products")
	if err != nil {
		log.Println(err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	var product []models.Product
	err = json.NewDecoder(productResp.Body).Decode(&product)
	if err != nil {
		log.Println(err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	var user models.User
	token, err := r.Cookie("jwt_token")
	if err == nil {
		req, err := http.NewRequest("POST", "http://localhost:9090/api/user/info", nil)
		if err != nil {
			log.Println(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		req.AddCookie(token)

		client := &http.Client{
			Timeout: 10 * time.Second,
		}
		userResp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		if userResp.StatusCode == http.StatusOK {

			err = json.NewDecoder(userResp.Body).Decode(&user)
			if err != nil {
				log.Println(err)
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
		} else if userResp.StatusCode == http.StatusUnauthorized {
			cookie := &http.Cookie{
				Name:   "jwt_token",
				Value:  "",
				Path:   "/",
				MaxAge: -1,
			}
			http.SetCookie(w, cookie)
		} else {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
	}

	RenderTemplate(w, "index.html", models.UserDataProduct{
		User:    user,
		Product: product,
	})
}
