#!/bin/bash

#Need to run this script as a root
if [ "$EUID" -ne 0 ]
  then echo "Please run as root (with sudo command)"
  exit
fi

echo initialization of the database...
#check whether the db is running
pid=`lsof -i:3306 -t`
if [ $pid == '']
then
  echo Please start mysql sever
else
  echo Mysql running with pid: $pid
fi

#creating sql user and logging in
echo Please provide the credentials for the database.
echo Username :
read -r username
echo Password :
read -r password
#echo Please provide the root password - default procedure is to skip this

mysql -u "$username" -p "$password" <<PATCH_SCRIPT
#creating the cloud tables
USE kdb;
CREATE TABLE kdb.secret (
	Name varchar(100) NOT NULL,
	OwnerId INT NOT NULL,
	Provider varchar(100) NOT NULL,
	`Type` INT NULL,
	CreatedOn DATETIME DEFAULT NOW() NOT NULL,
	CreatedBy INT NOT NULL,
	ModifiedOn DATETIME DEFAULT NOW() NULL,
	ModifiedBy INT NULL,
	Activated BOOLEAN DEFAULT FALSE NOT NULL,
	Deleted BOOLEAN DEFAULT FALSE NOT NULL,
	Encrpted BOOLEAN DEFAULT FALSE NOT NULL,
	Tags varchar(100) NULL,
	ID BIGINT NOT NULL AUTO_INCREMENT,
	CONSTRAINT secret_PK PRIMARY KEY (ID),
	CONSTRAINT secret_FK FOREIGN KEY (OwnerId) REFERENCES kdb.`user`(id) ON DELETE CASCADE ON UPDATE CASCADE
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_general_ci;
CREATE UNIQUE INDEX secret_OwnerId_IDX USING HASH ON kdb.secret (OwnerId);
\q;
PATCH_SCRIPT

echo kdb database initialized successfully with secret tables