#!/bin/bash

#Need to run this script as a root
if [ "$EUID" -ne 0 ]
  then echo "Please run as root (with sudo command)"
  exit
fi

echo initialization of the database...
#check whether the db is running
pid=`lsof -i:3306 -t`
IF [$pid == '']
THEN
  echo Please start mysql sever
ELSE
  echo Mysql running with pid: $pid
FI

#creating sql user and logging in
echo Please provide the credentials for the database.
echo Username :
read -r username
echo Password :
read -r password
echo Please provide the root password - default procedure is to skip this

mysql -u "$username" -p "$password" <<PATCH_SCRIPT
#creating the cloud tables
USE kdb;
CREATE TABLE `secret` (
	`ID` INT NOT NULL AUTO_INCREMENT COMMENT 'unique id',
	`Name` VARCHAR COMMENT 'identifiable name',
	`OwnerId` INT NOT NULL COMMENT 'owner id correlation',
	`Provider` VARCHAR NOT NULL COMMENT 'cloud service provider',
	`Type` INT COMMENT 'secret type',
	`Key` VARCHAR NOT NULL COMMENT 'cloud secret',
	`CreatedOn` DATETIME NOT NULL DEFAULT 'SYSDATE' COMMENT 'create timestamp',
	`CreatedBy` INT NOT NULL COMMENT 'user correlation id',
	`ModifiedOn` DATETIME NOT NULL DEFAULT 'SYSDATE',
	`ModifiedBy` INT NOT NULL COMMENT 'user correlation id',
	`Activated` BOOLEAN NOT NULL DEFAULT 'FALSE' COMMENT 'activated flag',
	`Deleted` BOOLEAN NOT NULL DEFAULT 'FALSE' COMMENT 'deleted flag',
	`Encrypted` BOOLEAN NOT NULL DEFAULT 'FALSE' COMMENT 'encrypted flag',
	`Tags` VARCHAR COMMENT 'tags',
	UNIQUE KEY `IDX_Owner` (`OwnerId`) USING HASH,
	PRIMARY KEY (`ID`)
) ENGINE=InnoDB;
\q;
PATCH_SCRIPT

echo kdb database initialized successfully with secret tables