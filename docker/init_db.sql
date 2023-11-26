# rootユーザーのリモート接続を許可
RENAME USER root@'localhost' to root@'%';
CREATE DATABASE waroka CHARACTER SET = utf8;
