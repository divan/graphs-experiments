package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"github.com/divan/graph-experiments/cmd/data_generator/p2p"
	"github.com/divan/graph-experiments/generator/basic"
	"github.com/divan/graph-experiments/generator/net"
	"github.com/divan/graph-experiments/graph"
)

type Generator interface {
	Generate() *graph.Data
}

func main() {
	var (
		dataKind         = flag.String("type", "net", "Example random IPs network (net, line)")
		netHosts         = flag.Int("net.hosts", 20, "Number of hosts for net generator")
		netConns         = flag.Int("net.connections", 4, "Number of connections between hosts for net generator")
		p2pSendN         = flag.Int("p2psend.N", 3, "Number of peers to propagate (0..N of peers)")
		p2pSendDelay     = flag.Duration("p2psend.delay", 10*time.Millisecond, "Delay for each step")
		p2pSendTTL       = flag.Int("p2psend.ttl", 10, "Message TTL")
		p2pSendStartNode = flag.String("p2psend.startNode", "192.168.1.2", "IP address of node initiating the sending")
		output           = flag.String("o", "data.json", "Output filename for network data")
		p2pOutput        = flag.String("p", "propagation.json", "Output filename for p2p sending data")
	)
	flag.Parse()

	// Prepare output files for writing
	netFd, err := os.Create(*output)
	if err != nil {
		log.Fatal("Open file for writing failed:", err)
	}
	defer netFd.Close()

	p2pFd, err := os.Create(*p2pOutput)
	if err != nil {
		log.Fatal("Open file for writing failed:", err)
	}
	defer p2pFd.Close()

	var generator Generator
	switch *dataKind {
	case "net":
		generator = net.NewDummyGenerator(*netHosts, *netConns, "192.168.1.1", net.Uniform)
	case "line":
		generator = basic.NewLineGenerator(10)
	case "circle":
		generator = basic.NewCircleGenerator(10)
	}

	data := generator.Generate()
	err = json.NewEncoder(netFd).Encode(data)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Written graph into", *output)

	if *dataKind == "p2psend" {
		sendData := p2p.SimulatePropagation(data, *p2pSendN, *p2pSendTTL, *p2pSendDelay, *p2pSendStartNode)
		err = json.NewEncoder(p2pFd).Encode(sendData)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Written p2p sending data into", *p2pOutput)
	}
}
