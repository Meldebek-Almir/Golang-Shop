package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"shop/frontend-client/internal/models"
)

func Signin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		RenderTemplate(w, "sign-in.html", nil)
	case http.MethodPost:
		r.ParseForm()
		credintails := models.Credintails{
			Identifier: r.FormValue("identity"),
			Password:   r.FormValue("password"),
		}
		body, err := json.Marshal(&credintails)
		if err != nil {
			fmt.Println(err.Error())
		}
		res, err := http.Post("http://localhost:9090/login", "application/json", bytes.NewReader(body))
		if res == nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		if err != nil || res.StatusCode != 200 {
			ErrorHandler(w, res.StatusCode)
			w.WriteHeader(res.StatusCode)
			return
		}
		token := res.Cookies()

		http.SetCookie(w, &http.Cookie{
			Expires: token[0].Expires,
			Value:   token[0].Value,
			Name:    "jwt_token",
		})
		http.Redirect(w, r, "/", http.StatusFound)

	}
}

func Signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		RenderTemplate(w, "sign-up.html", nil)
	case http.MethodPost:
		r.ParseForm()

		age, err := strconv.Atoi(r.FormValue("age"))
		if err != nil {
			fmt.Println(err)
		}
		credintails := models.UserRegister{
			Nickname:        r.FormValue("nickname"),
			Age:             age,
			Gender:          r.FormValue("gender"),
			FirstName:       r.FormValue("first_name"),
			LastName:        r.FormValue("last_name"),
			Email:           r.FormValue("email"),
			Password:        r.FormValue("password"),
			Role:            r.FormValue("role"),
			ConfirmPassword: r.FormValue("confirm"),
		}
		body, err := json.Marshal(&credintails)
		if err != nil {
			fmt.Println(err.Error())
		}
		req, err := http.NewRequest(http.MethodPost, "http://localhost:9090/register", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println(err)
		}

		client := http.Client{
			Timeout: 10 * time.Second,
		}
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		if res.StatusCode != 200 {
			w.WriteHeader(res.StatusCode)
			return
		}
		http.Redirect(w, r, "/sign-in", http.StatusFound)

	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "jwt_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
}
