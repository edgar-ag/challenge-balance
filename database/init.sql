CREATE DATABASE IF NOT EXISTS challange;
 
 USE challange;

 CREATE TABLE customer (
    id INT AUTO_INCREMENT PRIMARY KEY,
    account_number VARCHAR(255) NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE transaction (
    id INT AUTO_INCREMENT PRIMARY KEY,
    customer_id INT,
    txn VARCHAR(255) NOT NULL,
    txn_date VARCHAR(5) NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES customer(id)
);