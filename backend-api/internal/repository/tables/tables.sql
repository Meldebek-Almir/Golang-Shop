CREATE TABLE IF NOT EXISTS users(
    user_id SERIAL PRIMARY KEY,
    nickname VARCHAR(30) NOT NULL UNIQUE,
    age INTEGER NOT NULL,
    gender VARCHAR(7) NOT NULL,
    first_name VARCHAR (30) NOT NULL,
    last_name VARCHAR (30) NOT NULL,
    email VARCHAR (100) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    role VARCHAR(10) NOT NULL,
    created TIMESTAMP NOT NULL,
    updated TIMESTAMP NOT NULL
);


CREATE TABLE IF NOT EXISTS products (
  product_id SERIAL PRIMARY KEY,
  product_name VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  price INTEGER NOT NULL,
  rating FLOAT NOT NULL,
  image_url VARCHAR(255) NOT NULL,
  user_id INTEGER NOT NULL,
  available_quantity INTEGER NOT NULL DEFAULT 0,
  total_quantity_sold INTEGER NOT NULL DEFAULT 0,
  FOREIGN KEY (user_id) REFERENCES users(user_id)
);
CREATE TABLE IF NOT EXISTS ratings (
  rating_id SERIAL PRIMARY KEY,
  product_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  rating INTEGER NOT NULL,
  FOREIGN KEY (product_id) REFERENCES products(product_id),
  FOREIGN KEY (user_id) REFERENCES users(user_id),
  UNIQUE (user_id, product_id)
);


CREATE TABLE IF NOT EXISTS comments (
  comment_id SERIAL PRIMARY KEY,
  product_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  message TEXT NOT NULL,
  FOREIGN KEY (product_id) REFERENCES products(product_id),
  FOREIGN KEY (user_id) REFERENCES users(user_id)
);



CREATE TABLE IF NOT EXISTS orders (
  order_id SERIAL PRIMARY KEY,
  product_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  quantity INTEGER NOT NULL,
  total_price INTEGER DEFAULT 100,
  order_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);