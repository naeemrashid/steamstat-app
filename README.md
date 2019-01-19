# steamstat-app


## Build Instructions

+ setiing up Golang
``` 
# Download Golang version v1.11.4
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
CREATE USER 'steamuser'@'localhost' IDENTIFIED BY 'steamuser';
CREATE DATABSE steamstats;
GRANT ALL PRIVILEGES ON steamstats. * TO 'steamuser'@'localhost';
```
+ Install project dependencies
```
dep ensure
```
+ Build/Run Project
```
go build -o steamstatserver main.go
./steamstatserver
```
## How it works
The idea is to fetch all games and related details available on steam using steampowered api. And fetch each game statistics(current players in the game using [this api call](https://api.steampowered.com/ISteamUserStats/GetNumberOfCurrentPlayers/v1/?appid=730) ) after every hour and maintain a local cache for average change in player count for each game and perform different calculation e.g calculate moving average, calculate peak players, change in moving average in 24hr and change in 48hr etc. Record peak current players for each game and save it into database with the date by querying current players each hour and saving the peak value every 24 hours. In other words record peak players of the game each day.
api provides four interfaces 
```
 GET    /api/v1/gamestats/trending 
 GET    /api/v1/gamestats/topgames 
 GET    /api/v1/gamestats/toprecords 
 GET    /api/v1/gamestats/games/:id 
```
+ Trending Games: Trending games are selected based upon change in game players on hourly bases. Games with highest positive change are at the top is decending order with respect to average change. Average change of each game is calculated using this formula.
```
func AvgChangePercnt(initialValue, finalValue float64)float64{
	if initialValue == 0{
		return 0
	}
	return  (finalValue-initialValue)/initialValue * 100
}
```
+ Top Games By current Players: Top games are sorted in desecing order with respect to current players count.
+ Top Records: Games with highest peak players are placed at the top, the idea is to save peak value of current players each day and select the games with the peak value is descending order.
+ Game Details: Provides different details available on steam powered api about each games e.g, title, description, about the game, reviews, Header Image etc.
## Future Plans/Roadmap
+ In depth analysis of games can be performed based on stored data.
+ Application fetch the list of game on start and do not update frequently, to update game list you need to restart the app, which is a limitation but not difficult to fix, if thats required and necessary do let us know otherwise you can restart the app and it will fetch updated list, as apps list on steam is not updated so frequently.
+ Currently most of the games statistics are stored in-memory(hash-map) which can be moved to database, its was a design decision that most of the stuff we wont need for fututre use so we should keep those into in-memory e.g current players, 24 hr average, 48hr change etc.
+ Application use two dependency projects e.g https://github.com/naeemkhan12/golang-moving-average,https://github.com/naeemkhan12/kettle. I have made changes according to application needs and both projects are maintaned by myself, so we are good here.
## Limitations 
+ steamcharts.com has **hours played** count which we skipped as steampowered api does not provide any interface to query hours played unless you have players account ids(which we do not have ones).
+ steamcharts.com has a top graph with entry of online Players and in-game players, we had the in-game player count but we didn't had the online players count so we skipped that table, again it was not in the api straight-forward you need player account ids to query player status etc.
