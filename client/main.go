package main

// tag::import[]
import (
	pb "github.com/cmceniry/login-grpcls/directorycontents"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	// end::import[]

	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	// tag::tls[]
	certificate, err := tls.LoadX509KeyPair(
		"../login-glss/certs/client.crt",
		"../login-glss/certs/client.key",
	)
	if err != nil {
		log.Fatalf("could not load client key pair: %s", err)
	}
	caCert, err := ioutil.ReadFile("../login-glss/certs/CA.crt")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	// end::tls[]

	// tag::credentials[]
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certificate},
		RootCAs:      caCertPool,
	})
	// end::credentials[]

	// tag::dial[]
	conn, err := grpc.Dial("localhost:4270", grpc.WithTransportCredentials(creds))
	// end::dial[]
	if err != nil {
		log.Fatalf("connect failure: %s", err)
	}
	defer conn.Close()
	// tag::listerclient[]
	c := pb.NewListerClient(conn)
	// end::listerclient[]

	// tag::ls[]
	files, err := c.LS(context.Background(), &pb.Path{Name: os.Args[1]})
	// end::ls[]
	if err != nil {
		log.Fatalf("LS failure: %s", err)
	}
	// tag::recvloop[]
	for {
		f, err := files.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("LS file failure: %s", err)
		}
		fmt.Printf("%s %10d %s %s\n", f.Mode, f.Size, f.Modtime, f.Name)
	}
	// end::recvloop[]
}
