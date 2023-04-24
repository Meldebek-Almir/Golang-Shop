package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"shop/frontend-client/internal/models"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		RenderTemplate(w, "create.html", nil)
	} else if r.Method == http.MethodPost {
		r.ParseForm()

		file, header, err := r.FormFile("image_url")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()
		fileExt := filepath.Ext(header.Filename)
		newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), fileExt)
		filePath := filepath.Join("./templates/static/", "images", newFilename)
		outFile, err := os.Create(filePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer outFile.Close()
		_, err = io.Copy(outFile, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		imageURL := fmt.Sprintf("/static/images/%s", newFilename)

		price, _ := strconv.Atoi(r.FormValue("price"))
		quantity, _ := strconv.Atoi(r.FormValue("available_quantity"))
		product := models.Product{
			ProductName:       r.FormValue("product_name"),
			Description:       r.FormValue("description"),
			Price:             price,
			ImageUrl:          imageURL,
			AvailableQuantity: quantity,
		}

		var productJson bytes.Buffer

		err = json.NewEncoder(&productJson).Encode(&product)
		if err != nil {
			log.Println(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		token, err := r.Cookie("jwt_token")
		if err != nil {
			ErrorHandler(w, http.StatusUnauthorized)
			return
		}

		req, err := http.NewRequest("POST", "http://localhost:9090/api/products", &productJson)
		if err != nil {
			log.Println(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(token)

		client := http.Client{
			Timeout: 10 * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else {
			ErrorHandler(w, resp.StatusCode)
			return
		}
	}
}

func More(w http.ResponseWriter, r *http.Request) {
	idstr := chi.URLParam(r, "product_id")
	resp, err := http.Get("http://localhost:9090/api/products/" + idstr)
	if err != nil {
		log.Println(err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var product models.Product
	if resp.StatusCode == http.StatusOK {

		err = json.NewDecoder(resp.Body).Decode(&product)
		if err != nil {
			log.Println(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

	} else {
		ErrorHandler(w, resp.StatusCode)
		return
	}

	RenderTemplate(w, "more.html", product)
}

func Search(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	query := r.FormValue("query")
	resp, err := http.Get("http://localhost:9090/api/search/" + query)
	if err != nil {
		log.Println(err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	var product []models.Product
	var user models.User
	if resp.StatusCode == 200 {
		err = json.NewDecoder(resp.Body).Decode(&product)
		if err != nil {
			log.Println(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
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
		} else {
			ErrorHandler(w, http.StatusUnauthorized)
			return
		}

	} else {
		fmt.Println(resp.StatusCode)
	}

	RenderTemplate(w, "index.html", models.UserDataProduct{
		User:    user,
		Product: product,
	})
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idstr := chi.URLParam(r, "product_id")
	token, err := r.Cookie("jwt_token")

	if err != nil {
		log.Println(err)
		ErrorHandler(w, http.StatusUnauthorized)
		return
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:9090/api/products/%s/delete", idstr), nil)
	if err != nil {
		log.Println(err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(token)

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else if resp.StatusCode != 0 {
		ErrorHandler(w, resp.StatusCode)
		return
	} else {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
}
