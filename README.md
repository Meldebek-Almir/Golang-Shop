# Golang Shop

Golang Shop is an e-commerce web application that allows users to browse and purchase items online. It has a registration and authorization system, a search and filter feature, a rating and commenting system, and different roles for admins, clients, and sellers.
# API Endpoints

    - POST /register: register a new user and send data to the database.
    - POST /login: check user credentials against the database and authorize the user.
    - POST /api/products: add a new product to the database (seller only).
    - POST /api/products/{product_id}/comments: add a comment for a specific product (client only).
    - POST /api/products/{product_id}/rate: rate a specific product (client only).
    - POST /api/products/{product_id}/order: purchase a specific product (client only).
    - GET /api/products: get all available products.
    - GET /api/products/{product_id}: get a specific product by ID.
    - GET /api/search/{query}: search for products by name.
    - GET /api/filter: filter products by price and rating.
    - POST /api/user/info: get user information (authenticated users only).

# Frontend Endpoints

    - GET /: the welcome page.
    - GET /sign-in: sign-in page.
    - POST /sign-in: submit sign-in form.
    - GET /sign-up: sign-up page.
    - POST /sign-up: submit sign-up form.
    - GET /logout: log out the user.
    - GET /profile: user profile page.
    - GET /product/{product_id}: product details page.
    - GET /create-product: create product page (seller only).
    - POST /create-product: submit new product form.
    - GET /search: search page.
    - GET /: the welcome page.

# Installation and Run 

    Clone the repository: git clone https://github.com/expose443/CandyShop

    1) cd CandyShop

    2) make postgres

    3) cd backend-api

    4) ./secret-key.sh

    4) go run ./cmd 

    5) cd ../frontend-client

    6) go run ./cmd 

    frontend: http://localhost:8080
