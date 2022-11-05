package main

import (
	"net"
	"flag"
	"context"
	"log"
	"fmt"

	pf "github.com/ZhengjunHUO/go-playground/protobuf/protob"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8080, "The port server listen to")
)

type server struct {
	pf.UnimplementedGetterServer
}

func (s *server) ShowK8S(ctx context.Context, req *pf.Requete) (*pf.K8SInfo, error) {
	log.Println("Recieve request: ", req.GetId())
	return &pf.K8SInfo {
		Name: "huo",
		Size: 6,
		Ismanaged: true,
		Cni: &pf.Cni{
			Name: "Cilium",
			IsOverlayed: false,
			IsDirectRouting: true,
		},
	}, nil
}

func main() {
	flag.Parse()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Bind to port %d failed: %v", *port, err)
	}

	log.Println("Starting server ...")
	s := grpc.NewServer()
	pf.RegisterGetterServer(s, &server{})
	if err := s.Serve(l); err != nil {
		log.Fatalln("Server up failed: ", err)
	}
}
