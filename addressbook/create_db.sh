#!/bin/sh

docker exec -i addressbook_db_1 mysql -uroot -ppassword <<EOF
CREATE DATABASE addressbook;
CREATE USER 'dev'@'%' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON addressbook.* to 'dev'@'%'; 
FLUSH PRIVILEGES;
EOF

docker exec -i addressbook_db_1 mysql -udev -ppassword addressbook <<EOF
CREATE TABLE people (
  id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(30) NOT NULL,
  email VARCHAR(50)
)
EOF
