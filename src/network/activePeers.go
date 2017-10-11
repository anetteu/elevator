package network

import (
	. "../config"
	"./localip"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var ActivePeers []string
var N_ELEVATORS int
var lostPeers = make(map[string]time.Time)

func UpdateActivePeers() {
	for {
		for ID, timestamp := range lostPeers {
			if time.Since(timestamp).Seconds() > UNCONNECTED_PEER_TIME_LIMIT {
				removeFromActivePeers(ID)
				delete(lostPeers, ID)
			}
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func StripPeers(peers PeerUpdate) PeerUpdate {
	for index, peer := range peers.Lost {
		peers.Lost[index] = strings.Split(peer, "-")[1]
	}
	for index, peer := range peers.Peers {
		peers.Peers[index] = strings.Split(peer, "-")[1]
	}
	if peers.New != "" {
		peers.New = strings.Split(peers.New, "-")[1]
	}
	return peers

}

func MapLostPeers(peers PeerUpdate) {
	if len(peers.Lost) > 0 {
		for _, ID := range peers.Lost {
			lostPeers[ID] = time.Now()
		}
	}
	if peers.New != "" {
		for ID := range lostPeers {
			if peers.New == ID {
				delete(lostPeers, ID)
			}
		}
		addToActivePeers(peers.New)
	}
}

func aliveElevators() string {
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}
	return id
}

func addToActivePeers(peer string) {
	inActivePeers := 0
	for _, elevator := range ActivePeers {
		if elevator == peer {
			inActivePeers = 1
		}
	}
	if inActivePeers == 0 {
		ActivePeers = append(ActivePeers, peer)
	}
}

func removeFromActivePeers(ID string) {
	for index, elevator := range ActivePeers {
		if elevator == ID {
			ActivePeers = removeIndexFromSlice(ActivePeers, index)
		}
	}
}

func removeIndexFromSlice(slice []string, index int) []string {
	slice[index] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}
