#!/usr/bin/env bash

#pulling mysql from docker and executing the image
sudo su
docker pull mysql/mysql-server
docker run --name=mysql1 -d mysql/mysql-server:tag

#creating sql user and logging in
mysql -u root -p
echo
echo Please provide the credentials for the database.
echo Username :
read -r username
echo Password :
read -r password
GRANT ALL PRIVILEGES ON *.* TO "$username"@'localhost' IDENTIFIED BY "$password";
/q
mysql -u "$username" -p
echo "$password"

#creating the db
CREATE DATABASE kdb;
USE kdb;
CREATE TABLE user (id int(10) not null primary key auto_increment, username varchar(30) not null unique, token varchar(30) unique, access_level varchar(10), password varchar(100), first_name varchar(30), last_name varchar(30), email varchar(50));
CREATE TABLE cluster (id int(10) not null primary key auto_increment, cluster_name varchar(30) unique, kafka_version varchar(10), active_controllers int(10));
CREATE TABLE broker (id int(10) not null primary key auto_increment, host varchar(100) unique, port int(10), created_at datetime, cluster_id int(10) not null);
echo Tables created :
DESCRIBE user;
DESCRIBE cluster;
DESCRIBE broker;
\q