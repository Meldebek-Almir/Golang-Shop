package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"shop/frontend-client/internal/models"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	var user models.User
	token, err := r.Cookie("jwt_token")
	if err == nil {
		req, err := http.NewRequest("POST", "http://localhost:9090/api/user/info", nil)
		if err != nil {
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
		err = json.NewDecoder(userResp.Body).Decode(&user)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		if userResp.StatusCode != 200 && userResp.StatusCode != 0 {
			ErrorHandler(w, userResp.StatusCode)
			return
		} else if userResp.StatusCode == 0 {
			ErrorHandler(w, http.StatusUnauthorized)
			return
		}

		RenderTemplate(w, "profile.html", user)
	} else {
		ErrorHandler(w, http.StatusUnauthorized)
		return
	}
}
