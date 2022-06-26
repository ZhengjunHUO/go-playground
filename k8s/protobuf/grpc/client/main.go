package main

import (
	"fmt"
	"log"
	"flag"
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pf "github.com/ZhengjunHUO/playground/k8s/protobuf/protob"
)

var (
	ip   = flag.String("ip", "127.0.0.1", "target ip")
	port = flag.Int("port", 8080, "target port")
	id   = flag.String("id", "huok8s", "id of cluster to query")
)

func main() {
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *ip, *port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to %s: %v\n", addr, err)
	}
	defer conn.Close()
	client := pf.NewGetterClient(conn)	

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
	resp, err := client.ShowK8S(ctx, &pf.Requete{Id: *id})
	if err != nil {
		log.Fatalf("Failed to get response from server: %v\n", err)
	}

	log.Printf("Get cluster id [%s]'s info:\n%v\n", *id, resp)
}
