#!/bin/bash

#Need to run this script as a root
if [ "$EUID" -ne 0 ]
  then echo "Please run as root (with sudo command)"
  exit
fi

echo initialization of the database... 
#downloading locally
yum update
yum install wget -y
wget http://repo.mysql.com/mysql-community-release-el7-5.noarch.rpm
echo mysql package is downloaded...
rpm -ivh mysql-community-release-el7-5.noarch.rpm
yum update
echo updating yum...
yum install mysql-server -y
echo installing mysql-server...
systemctl start mysqld
echo starting mysql-server...

#creating sql user and logging in
echo Please provide the credentials for the database.
echo Username :
read -r username
echo Password :
read -r password
echo Please provide the root password - default procedure is to skip this
mysql -u root <<MYSQL_SCRIPT
GRANT ALL PRIVILEGES ON *.* TO "$username"@'localhost'IDENTIFIED BY "$password";
\q
MYSQL_SCRIPT

#logging into the user
mysql -u "$username" -p"$password" <<INIT_SCRIPT
#creating the db
CREATE DATABASE kdb;
USE kdb;
CREATE TABLE user (id int(10) not null primary key auto_increment, username varchar(30) not null unique, token varchar(30) unique, access_level varchar(10), password varchar(100), first_name varchar(30), last_name varchar(30), email varchar(50));
CREATE TABLE cluster (id int(10) not null primary key auto_increment, cluster_name varchar(30) unique, kafka_version varchar(10), active_controllers int(10));
CREATE TABLE broker (id int(10) not null primary key auto_increment, host varchar(100) unique, port int(10), created_at datetime, cluster_id int(10) not null);
CREATE TABLE broker_metrics (id int(20) not null primary key auto_increment, timestamp int(20) not null, host varchar(100) not null, partitions int(10), leaders int(10), act_controller int(5), offline_part int(10), under_replicated int(10), bytes_in int(15) default 0, bytes_out int(15) default 0, mesg_rate decimal(20,2) default 0.0, isr_exp_rate decimal(20,2), isr_shrink_rate decimal(20,2), send_time int(15), queue_time int(15), remote_time int(15), local_time int(15), total_time int(15), net_proc_avg_idle_perc decimal(10,5), max_lag decimal(15,2), unclean_lead_elec int(10), failed_fetch_rate int(10), failed_prod_rate int(10), total_messages int(20));
\q
INIT_SCRIPT

echo kdb database initialized successfully with user, cluster and broker tables

