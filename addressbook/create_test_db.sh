#!/bin/sh

docker exec -i addressbook_db_1 mysql -uroot -ppassword <<EOF
CREATE DATABASE addressbook_test;
GRANT ALL PRIVILEGES ON addressbook_test.* to 'dev'@'%'; 
FLUSH PRIVILEGES;
EOF

docker exec -i addressbook_db_1 mysql -udev -ppassword addressbook_test <<EOF
CREATE TABLE people (
  id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(30) NOT NULL,
  email VARCHAR(50)
)
EOF
