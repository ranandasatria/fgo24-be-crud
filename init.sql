-- Active: 1750753580331@@127.0.0.1@5433@backend
CREATE DATABASE backend;

CREATE TABLE users(
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL
)

select * from users;

update users set name = 'rananda' where id = 2;