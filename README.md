# steamstat-app


## Build Instructions

+ setiing up Golang
``` 
wget https://golang.org/doc/install?download=go1.11.4.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.11.4.linux-amd64.tar.gz

# Add /usr/local/go/bin to the PATH environment variable. You can do this by adding this line to your /etc/profile (for a system-wide installation) or $HOME/.profile: or $HOME/.bashrc

export PATH=$PATH:/usr/local/go/bin

```
+ Setup Dep: A dependecy management tool for Go.
```
 curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
```
+ setup MYSQL 
```
sudo apt-get install mysql-common mysql-server mysql-client -y
# setup mysql user 
CREATE USER 'naeem'@'localhost' IDENTIFIED BY 'naeem';
CREATE DATABSE steamstats;
GRANT ALL PRIVILEGES ON steamstats. * TO 'naeem'@'localhost';
```

+ Install project dependencies
```
dep ensure
```
+ Build/Run Project
```
go run main.go
```
