package cmd

import (
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	"golang.org/x/net/context"

	"github.com/JustinAzoff/asnlookup/asndb"
	pb "github.com/JustinAzoff/asnlookup/pb"
)

const (
	port = ":50051"
)

//type AsnlookupServer interface {
//        Hello(context.Context, *Empty) (*HelloReply, error)
//        Lookup(context.Context, *LookupRequest) (*LookupReply, error)
//        LookupMany(Asnlookup_LookupManyServer) error
//        LookupBatch(context.Context, *LookupRequestBatch) (*LookupReplyBatch, error)
//}

type asnlookupServer struct {
	backend *asndb.AsnBackend
}

func (s *asnlookupServer) Hello(ctx context.Context, empty *pb.Empty) (*pb.HelloReply, error) {
	log.Printf("Hello called!")
	return &pb.HelloReply{
		Message: "Hello!",
	}, nil
}
func (s *asnlookupServer) Lookup(ctx context.Context, req *pb.LookupRequest) (*pb.LookupReply, error) {
	rec, err := s.backend.Lookup(req.Address)
	if err != nil {
		return nil, err
	}
	return &pb.LookupReply{
		Address: rec.IP,
		Prefix:  rec.Prefix,
		As:      int32(rec.AS),
		Owner:   rec.Owner,
		Cc:      rec.CC,
	}, nil
}
func (s *asnlookupServer) LookupMany(stream pb.Asnlookup_LookupManyServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		rec, err := s.backend.Lookup(req.Address)
		if err != nil {
			return err
		}
		rep := &pb.LookupReply{
			Address: rec.IP,
			Prefix:  rec.Prefix,
			As:      int32(rec.AS),
			Owner:   rec.Owner,
			Cc:      rec.CC,
		}
		if err := stream.Send(rep); err != nil {
			return err
		}
	}
}
func (s *asnlookupServer) LookupBatch(ctx context.Context, req *pb.LookupRequestBatch) (*pb.LookupReplyBatch, error) {
	return nil, fmt.Errorf("not yet")
}

func server() {
	b, err := asndb.NewAsnBackend("asn.db", "asnames.json")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on %s", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := asnlookupServer{backend: b}

	// Creates a new gRPC server
	s := grpc.NewServer()
	pb.RegisterAsnlookupServer(s, &server)
	s.Serve(lis)
}
