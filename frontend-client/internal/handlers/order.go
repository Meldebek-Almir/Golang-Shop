package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"shop/frontend-client/internal/models"
)

func Order(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	quantity, err := strconv.Atoi(r.FormValue("order"))
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest)
		return
	}

	com := models.Order{
		Quantity: quantity,
	}
	strId := (r.URL.Query().Get("id"))
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Println(err)
		ErrorHandler(w, http.StatusBadRequest)
		return
	}
	com.ProductId = id
	var message bytes.Buffer
	err = json.NewEncoder(&message).Encode(&com)
	if err != nil {
		log.Println(err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	fmt.Println("orderId ", com.OrderId)
	fmt.Println("productId", com.ProductId)
	fmt.Println("UserId", com.UserId)
	fmt.Println("Quantity", com.Quantity)
	token, err := r.Cookie("jwt_token")
	if err == nil {
		fmt.Println(strId)
		req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:9090/api/products/%s/order", strId), &message)
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
			ErrorHandler(w, userResp.StatusCode)
			return
		}
		fmt.Println("i work")
		if userResp.StatusCode != http.StatusOK && userResp.StatusCode != 0 {
			ErrorHandler(w, userResp.StatusCode)
			return
		} else if userResp.StatusCode == 0 {
			ErrorHandler(w, http.StatusUnauthorized)
			return
		}
		fmt.Println("i work2")

		http.Redirect(w, r, fmt.Sprintf("/product/%s", strId), http.StatusSeeOther)
	} else {
		ErrorHandler(w, http.StatusUnauthorized)
		return
	}
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	var order []models.Order
	token, err := r.Cookie("jwt_token")
	if err == nil {
		req, err := http.NewRequest("GET", "http://localhost:9090/api/orders", nil)
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
		body, err := io.ReadAll(userResp.Body)
		if err != nil {
			log.Println(err)
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
		fmt.Println(string(body))

		err = json.Unmarshal([]byte(body), &order)
		if err != nil {
			log.Println(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		fmt.Println(order)
		RenderTemplate(w, "order.html", order)
	} else {
		ErrorHandler(w, http.StatusUnauthorized)
		return
	}
}
