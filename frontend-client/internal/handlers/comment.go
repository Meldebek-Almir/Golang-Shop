package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"shop/frontend-client/internal/models"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	com := models.CommentProduct{
		Message: r.FormValue("message"),
	}
	id := (r.URL.Query().Get("id"))
	var message bytes.Buffer
	err := json.NewEncoder(&message).Encode(&com)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	token, err := r.Cookie("jwt_token")
	if err == nil {
		req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:9090/api/products/%s/comments", id), &message)
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
		if userResp.StatusCode != http.StatusOK && userResp.StatusCode != 0 {
			ErrorHandler(w, userResp.StatusCode)
			return
		} else if userResp.StatusCode == 0 {
			ErrorHandler(w, http.StatusUnauthorized)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/product/%s", id), http.StatusSeeOther)
	} else {
		ErrorHandler(w, http.StatusUnauthorized)
		return
	}
}
