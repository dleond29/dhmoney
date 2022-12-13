CREATE DATABASE digitalmoneyhouse;
USE digitalmoneyhouse;
CREATE TABLE users(id INT NOT NULL PRIMARY KEY AUTO_INCREMENT, dni INT, phone INT);
CREATE TABLE accounts(id INT NOT NULL PRIMARY KEY AUTO_INCREMENT, user_id int not null, auth_id VARCHAR(255), cvu VARCHAR(22), alias VARCHAR(255), balance DECIMAL(15, 2) DEFAULT "0.00");
CREATE TABLE transactions(id INT NOT NULL PRIMARY KEY AUTO_INCREMENT, account_id int not null, origin_cvu VARCHAR(22),  destination_cvu VARCHAR(22),  description VARCHAR(50), amount DECIMAL(15, 2), date_time datetime, type VARCHAR(20));
CREATE TABLE cards(id INT NOT NULL PRIMARY KEY AUTO_INCREMENT, account_id int not null, pan VARCHAR(20), holder_name VARCHAR(255), expiration_date datetime, cid VARCHAR(4), type VARCHAR(20));