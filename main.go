package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var (
	client = &http.Client{}
)

var (
	Token     string
	BotPrefix string

	config *configStruct
)

type configStruct struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
}

func ReadConfig() error {
	fmt.Println("Reading config file...")
	file, err := ioutil.ReadFile("./config.json")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(string(file))

	err = json.Unmarshal(file, &config)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	Token = config.Token
	BotPrefix = config.BotPrefix

	return nil

}

var BotId string
var goBot *discordgo.Session

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	BotId = u.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Bot is running !")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotId {
		return
	}

	if strings.Split(m.Content, " ")[0] == "!search" {
		playerSearch := strings.Split(m.Content, "!search ")[1]
		fmt.Println("searching stats for", playerSearch)
		_, _ = s.ChannelMessageSend(m.ChannelID, search(playerSearch))
	}

	if strings.Split(m.Content, " ")[0] == "!liststats" {
		playerSearch := strings.Split(m.Content, "!liststats ")[1]
		fmt.Println("listing stats for", playerSearch)
		_, _ = s.ChannelMessageSend(m.ChannelID, listStats(playerSearch))
	}
}

func search(playerName string) string {
	stats, success := searchPlayer(playerName)
	if !success {
		return playerName + " could not be found."
	}
	ppg := fmt.Sprint(stats.PPG)
	apg := fmt.Sprint(stats.APG)
	rpg := fmt.Sprint(stats.RPG)
	fg := fmt.Sprint(stats.FG)
	return stats.NAME + " averages " + ppg + " points per game on " + fg + "% shooting, " + apg + " assists per game, and " + rpg + " rebounds per game."
}

func listStats(playerName string) string {
	stats, success := searchPlayer(playerName)
	if !success {
		return playerName + " could not be found."
	}
	ppg := fmt.Sprint(stats.PPG) + " ppg\n"
	apg := fmt.Sprint(stats.APG) + " apg\n"
	rpg := fmt.Sprint(stats.RPG) + " rpg\n"
	fg := fmt.Sprint(stats.FG) + " fg%\n"
	return stats.NAME + "\n" + ppg + fg + apg + rpg
}

func searchPlayer(playerName string) (PlayerStats, bool) {
	playerName = strings.ReplaceAll(strings.Trim(playerName, "\n"), " ", "+")

	searchReq, err := http.NewRequest("GET", "https://site.web.api.espn.com/apis/search/v2?region=us&lang=en&section=nba&limit=10&page=1&query="+playerName+"&dtciVideoSearch=true&iapPackages=ESPN_PLUS_UFC_PPV_284,ESPN_PLUS,ESPN_PLUS_MLB&type=promoted,team,player,league,article,live,upcoming,replay,clips", nil)

	if err != nil {
		fmt.Println(err)
		return PlayerStats{"", 0, 0, 0, 0}, false
	}

	searchRes, err := client.Do(searchReq)
	if err != nil {
		fmt.Println(err)
		return PlayerStats{"", 0, 0, 0, 0}, false
	}
	defer searchRes.Body.Close()

	searchBody, err := ioutil.ReadAll(searchRes.Body)
	if err != nil {
		fmt.Println(err)
		return PlayerStats{"", 0, 0, 0, 0}, false
	}

	var searchResult Search
	json.Unmarshal(searchBody, &searchResult)

	var url string
	for i := 0; i < len(searchResult.Results); i++ {
		if searchResult.Results[i].Type == "player" {
			url = searchResult.Results[i].Contents[0].Link.Web
		}
	}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
		return PlayerStats{"", 0, 0, 0, 0}, false
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return PlayerStats{"", 0, 0, 0, 0}, false
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return PlayerStats{"", 0, 0, 0, 0}, false
	}

	name := strings.Split(string(body), ",\"dspNm\":\"")[1]
	name = strings.Split(name, "\"")[0]

	if !strings.Contains(string(body), "\"lbl\":\"PTS\",\"val\":\"") {
		return PlayerStats{NAME: "", PPG: 0, APG: 0, RPG: 0, FG: 0}, false
	}
	ppg := strings.Split(string(body), "\"lbl\":\"PTS\",\"val\":\"")[1]
	ppg = strings.Split(ppg, "\"")[0]
	ppgFloat, _ := strconv.ParseFloat(ppg, 32)

	apg := strings.Split(string(body), "\"lbl\":\"AST\",\"val\":\"")[1]
	apg = strings.Split(apg, "\"")[0]
	apgFloat, _ := strconv.ParseFloat(apg, 32)

	rpg := strings.Split(string(body), "\"lbl\":\"REB\",\"val\":\"")[1]
	rpg = strings.Split(rpg, "\"")[0]
	rpgFloat, _ := strconv.ParseFloat(rpg, 32)

	fg := strings.Split(string(body), "\"lbl\":\"FG%\",\"val\":\"")[1]
	fg = strings.Split(fg, "\"")[0]
	fgFloat, _ := strconv.ParseFloat(fg, 32)

	return PlayerStats{name, float32(ppgFloat), float32(apgFloat), float32(rpgFloat), float32(fgFloat)}, true
}

func main() {
	//fmt.Print("Enter a player name: ")
	//var playerSearch string
	//
	//input := bufio.NewReader(os.Stdin)
	//
	//playerSearch, err := input.ReadString('\n')
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//searchPlayer(playerSearch)

	err := ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	Start()

	<-make(chan struct{})
	return
}
