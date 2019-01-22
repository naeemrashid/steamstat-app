package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/naeemkhan12/golang-moving-average"
	"github.com/naeemkhan12/kettle"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"steamstats-app/model"
	"steamstats-app/utils"
	"strconv"
	"sync"
	"syscall"
	"time"
)

const (
	APIKEY     = "B05DAE446F74CDDADAFFEFFE6314B39F"
	DATEFORMAT = "2006-1-6"
	DBMS       = "mysql"
	// TODO: replace user with steamuser
	DBMS_ARGS = "steamuser:steamuser@/steamstats?charset=utf8&parseTime=True&loc=Local"
)

var (
	steamClient *kettle.Client
	db          *gorm.DB
	gameStat    = struct {
		sync.RWMutex
		m map[int64]*GameStat
	}{m: make(map[int64]*GameStat)}
)

func init() {
	var err error
	log.Println("Connecting to Database...")
	db, err = gorm.Open(DBMS, DBMS_ARGS)
	if err != nil {
		fmt.Println(err)
		panic("failed to connect to database")
	}
	//Migrate the schema
	db.AutoMigrate(&model.Game{}, &model.PeakPlayer{}, &model.GameDetail{}, &model.Price{})
	httpClient := http.DefaultClient
	steamClient = kettle.NewClient(httpClient, APIKEY)
	err = fetchGames()
	if err != nil {
		log.Println("Error occured: ", err)
		panic(err)
	}
	initGameStatMap()
	updateGamesStat()
	updatePeakPlayers()
}

func main() {
	go serve()
	updateStatsEvery(1 * time.Hour)
	updatePeakPlayerEvery(24 * time.Hour)
	waitForSignal()
	err := db.Close()
	if err != nil {
		panic("Error closing database connection")
	}
}
func initGameStatMap() {
	log.Println("Initializeing Gamestats table...")
	var games []model.Game
	db.Find(&games)
	gameStat.Lock()
	defer gameStat.Unlock()
	for _, game := range games {
		gameStat.m[game.ID] = &GameStat{
			GameID:         game.ID,
			CurrentPlayers: 0,
			PeakPlayers:    0,
			Avg24hr:        movingaverage.New(24),
			Avg48hr:        movingaverage.New(48),
		}
	}
}
func updateGameStat(id int64) error {
	stat, _, err := steamClient.ISteamUserStatsService.GetNumberOfCurrentPlayers(id)
	if err != nil {
		return err
	}
	gameStat.Lock()
	defer gameStat.Unlock()
	if stat.Result == 1 {
		pCount := stat.PlayerCount
		if game, ok := gameStat.m[id]; ok {
			pCountP := game.CurrentPlayers
			if pCountP == 0 {
				game.CurrentPlayers = pCount
				pCountP = pCount
			}
			if pCount > game.PeakPlayers {
				game.PeakPlayers = pCount
			}
			avg := utils.AvgChangePercnt(float64(pCountP), float64(pCount))
			game.Avg24hr.Add(movingaverage.Values{Time: time.Now(), Value: avg})
			game.Avg48hr.Add(movingaverage.Values{Time: time.Now(), Value: avg})
			game.CurrentPlayers = pCount
		}

	}
	return nil
}
func updateGamesStat() {
	var games []model.Game
	db.Find(&games)
	for _, game := range games {
		updateGameStat(game.ID)
	}
}
func moving24hrAvg(id int64) float64 {
	gameStat.RLock()
	defer gameStat.RUnlock()
	if game, ok := gameStat.m[id]; ok {
		return game.Avg24hr.Avg()
	}
	return 0
}
func get48hrChange(id int64) []movingaverage.Values {
	gameStat.RLock()
	defer gameStat.RUnlock()
	if game, ok := gameStat.m[id]; ok {
		return game.Avg48hr.Values
	}
	return nil
}
func updateStatsEvery(d time.Duration) {
	go func() {
		for range time.Tick(d) {
			updateGamesStat()
		}
	}()
}

func waitForSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}

type GameStat struct {
	GameID         int64
	CurrentPlayers int64
	PeakPlayers    int64
	Avg24hr        *movingaverage.MovingAverage
	Avg48hr        *movingaverage.MovingAverage
}

func updatePeakPlayers() {
	gameStat.Lock()
	for _, v := range gameStat.m {
		peakP := v.PeakPlayers
		v.PeakPlayers = 0
		db.Save(&model.PeakPlayer{GameID: v.GameID, PeakPlayers: peakP, Date: time.Now()})
	}
	defer gameStat.Unlock()
}
func updatePeakPlayerEvery(d time.Duration) {
	go func() {
		for range time.Tick(d) {
			updatePeakPlayers()
		}
	}()
}

func fetchGames() error {
	log.Println("Downloading Games list from Steam API...")
	apps, _, err := steamClient.ISteamAppsService.GetAppList()
	if err != nil {
		log.Println("Unable to fetch apps list from steam")
		return err
	}
	// TODO: For testing purpose, only three games are tested
	//apps := []model.Game{model.Game{ID: 730, Name: "Counter Strike Global Offensive"},
	//	model.Game{ID: 570, Name: "Dota 2"}, model.Game{ID: 440, Name: "Team Fortress"}}
	log.Println("Adding 100 games details to database...")
	for _, app := range apps {
		details, _, err := steamClient.Store.AppDetails(app.AppID)
		if err == nil {
			db.Save(&model.Game{ID: app.AppID, Name: app.Name})
			db.Save(&model.GameDetail{
				GameID:              app.AppID,
				Title:               app.Name,
				Type:                details.Type,
				IsFree:              details.IsFree,
				DetailedDescription: details.DetailedDescription,
				AboutTheGame:        details.AboutTheGame,
				ShortDescription:    details.ShortDescription,
				SupportedLanguages:  details.SupportedLanguages,
				Reviews:             details.Reviews,
				HeaderImage:         details.HeaderImage,
				Website:             details.Website,
				Background:          details.Background,
			})
			db.Save(&model.Price{
				GameID:   app.AppID,
				Currency: details.PriceOverview.Currency,
				Initial:  details.PriceOverview.Initial,
				Final:    details.PriceOverview.Final,
			})
		}
	}
	return nil
}

func serve() {
	router := gin.Default()
	v1 := router.Group("/api/v1/gamestats")
	{
		v1.GET("/games/:id", gameDetails)
		v1.GET("/trending", trending)
		v1.GET("/topgames", topGamesByCP)
		v1.GET("/toprecords", topRecords)
	}
	router.Run()
}

func trending(c *gin.Context) {
	var games []model.Game
	var trendingGames []model.Trending
	db.Find(&games)
	gameStat.RLock()
	defer gameStat.RUnlock()
	for _, game := range games {
		trendingGames = append(trendingGames, model.Trending{GameID: game.ID,
			GameTitle:      game.Name,
			Change24hr:     gameStat.m[game.ID].Avg24hr.Avg(),
			Change48hr:     gameStat.m[game.ID].Avg48hr.Values,
			CurrentPlayers: gameStat.m[game.ID].CurrentPlayers})

	}
	sort.Slice(trendingGames[:], func(i, j int) bool {
		return trendingGames[i].Change24hr > trendingGames[j].Change24hr
	})
	// return top 500 entries
	if len(trendingGames) > 500 {
		trendingGames = trendingGames[:500]
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "trending": trendingGames})

}
func topGamesByCP(c *gin.Context) {
	var games []model.Game
	var topGames []model.TopGameByCP
	db.Find(&games)
	for _, game := range games {
		cp := int64(0)
		var last30Days []model.TimeSeries
		var peakPlayers []model.PeakPlayer
		gameStat.RLock()
		if value, ok := gameStat.m[game.ID]; ok {
			cp = value.CurrentPlayers
		}
		gameStat.RUnlock()
		db.Where("date BETWEEN ? AND ?", time.Now().AddDate(0, -1, 0).String(),
			time.Now().String()).Find(&peakPlayers)
		for _, val := range peakPlayers {
			last30Days = append(last30Days, model.TimeSeries{Time: val.Date.Format(DATEFORMAT), PeakPlayer: val.PeakPlayers})
		}
		topGames = append(topGames, model.TopGameByCP{GameID: game.ID,
			GameTitle:      game.Name,
			CurrentPlayers: cp,
			Last30Days:     last30Days,
		})
	}
	sort.Slice(topGames[:], func(i, j int) bool {
		return topGames[i].CurrentPlayers > topGames[j].CurrentPlayers
	})
	// return only first 500 entries
	if len(topGames) > 500 {
		topGames = topGames[:500]
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "topgames": topGames})
}

func topRecords(c *gin.Context) {
	var games []model.Game
	var topRecords []model.TopRecords
	db.Find(&games)
	gameStat.RLock()
	defer gameStat.RUnlock()
	for _, game := range games {
		pp := int64(0)
		var peakPlayer []model.PeakPlayer
		db.Model(&game).Related(&peakPlayer).Order("peak_players desc")
		if len(peakPlayer) > 0 {
			pp = peakPlayer[0].PeakPlayers
			topRecords = append(topRecords, model.TopRecords{GameID: game.ID,
				GameTitle:   game.Name,
				PeakPlayers: pp,
				Date:        peakPlayer[0].Date,
				Change48hr:  get48hrChange(game.ID),
			})
		}
	}
	sort.Slice(topRecords[:], func(i, j int) bool {
		return topRecords[i].PeakPlayers > topRecords[j].PeakPlayers
	})
	// return only first 500 entries
	if len(topRecords) > 500 {
		topRecords = topRecords[:500]
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "top_records": topRecords})
}
func gameDetails(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	var gameDetail model.GameDetail
	var price model.Price
	var transformedGameDetail model.TransformedGameDetail
	game := &model.Game{ID: int64(id)}
	db.Model(&game).Related(&gameDetail)
	db.Model(&game).Related(&price)
	gamePrice := model.GamePrice{
		Currency:        price.Currency,
		Initial:         price.Initial,
		Final:           price.Final,
		DiscountPercent: price.DiscountPercent,
	}
	transformedGameDetail = model.TransformedGameDetail{
		Title:               gameDetail.Title,
		Type:                gameDetail.Type,
		IsFree:              gameDetail.IsFree,
		DetailedDescription: gameDetail.DetailedDescription,
		AboutTheGame:        gameDetail.AboutTheGame,
		ShortDescription:    gameDetail.ShortDescription,
		SupportedLanguages:  gameDetail.SupportedLanguages,
		Reviews:             gameDetail.Reviews,
		HeaderImage:         gameDetail.HeaderImage,
		Website:             gameDetail.Website,
		Background:          gameDetail.Background,
		Price:               gamePrice,
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "details": transformedGameDetail})

}
