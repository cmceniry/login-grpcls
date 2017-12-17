package main

// tag::import[]
import (
	pb "github.com/cmceniry/login-grpcls/directorycontents"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	// end::import[]

	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
)

// tag::implement[]
type server struct{}

func (s *server) LS(p *pb.Path, fileInfoStream pb.Lister_LSServer) error {
	// end::implement[]
	// tag::walk[]
	err := filepath.Walk(p.Name, func(path string, info os.FileInfo, err error) error {
		// end::walk[]
		if err != nil {
			return err
		}
		// tag::respond[]
		f := &pb.File{
			Name:    info.Name(),
			Size:    info.Size(),
			Mode:    info.Mode().String(),
			Modtime: info.ModTime().Format("Jan _2 15:04"),
		}
		err = fileInfoStream.Send(f)
		// end::respond[]
		if err != nil {
			return err
		}
		if info.IsDir() && path != p.Name {
			return filepath.SkipDir
		}
		return nil
	})
	return err
}

// tag::listen[]
func main() {
	l, err := net.Listen("tcp", ":4270")
	// end::listen[]
	if err != nil {
		log.Fatalf("fail listen: %s", err)
	}

	// tag::tls[]
	certificate, err := tls.LoadX509KeyPair(
		"../login-glss/certs/server.crt",
		"../login-glss/certs/server.key",
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
		ServerName:   "localhost",
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})
	// end::credentials[]

	// tag::wiring[]
	s := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterListerServer(s, &server{})
	// end::wiring[]
	// tag::reflect[]
	reflection.Register(s)
	// end::reflect[]
	log.Printf("Starting server")
	// tag::start[]
	err = s.Serve(l)
	// end::start[]
	if err != nil {
		log.Fatalf("fail grpc: %s", err)
	}
}

// end::wiring[]
