package main

import (
	scraper "ak/scraper"
	types "ak/types"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
)

func greet(ec *nats.EncodedConn, m *nats.Msg) {
	b := types.Greet{Message: "Hello from ak microservice!"}
	ec.Publish(m.Reply, &b)
}

func recentAnime(ec *nats.EncodedConn, m *nats.Msg) {
	var data types.Payload
	json.Unmarshal(m.Data, &data)
	fmt.Println(string(m.Data))
	pageNumber, err := strconv.Atoi(data.Data.PageNumber)
	if err != nil {
		ec.Publish(m.Reply, scraper.GetRecentAnime(0))
	} else {
		ec.Publish(m.Reply, scraper.GetRecentAnime(pageNumber))
	}
}

func getAnime(ec *nats.EncodedConn, m *nats.Msg) {
	var data types.Payload
	json.Unmarshal(m.Data, &data)
	fmt.Println(string(m.Data))
	pageNumber, err := strconv.Atoi(data.Data.PageNumber)
	if err != nil {
		ec.Publish(m.Reply, scraper.GetAnime(data.Data.Url, 0))
	} else {
		ec.Publish(m.Reply, scraper.GetAnime(data.Data.Url, pageNumber))
	}
}

func getPlayerData(ec *nats.EncodedConn, m *nats.Msg) {
	var data types.Payload
	json.Unmarshal(m.Data, &data)
	fmt.Println(string(m.Data))
	ec.Publish(m.Reply, scraper.GetPlayerData(data.Data.Url))
}

func main() {
	// Connect to a server
	nc, err := nats.Connect("nats.default.svc:4222",
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ReconnectWait(time.Second*3),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected")
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	// Responding to a request message
	ec.QueueSubscribe("ak.greet", "ak", func(m *nats.Msg) {
		go greet(ec, m)
	})
	fmt.Println("Listening on subject 'ak.greet'...")
	ec.QueueSubscribe("ak.recent-anime", "ak", func(m *nats.Msg) {
		go recentAnime(ec, m)
	})
	fmt.Println("Listening on subject 'ak.recent-anime'...")
	ec.QueueSubscribe("ak.anime", "ak", func(m *nats.Msg) {
		go getAnime(ec, m)
	})
	fmt.Println("Listening on subject 'ak.anime'...")
	ec.QueueSubscribe("ak.player-data", "ak", func(m *nats.Msg) {
		go getPlayerData(ec, m)
	})
	fmt.Println("Listening on subject 'ak.player-data'...")
	runtime.Goexit()
}
