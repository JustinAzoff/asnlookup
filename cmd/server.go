package cmd

import (
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	"golang.org/x/net/context"

	"github.com/JustinAzoff/asnlookup/asndb"
	pb "github.com/JustinAzoff/asnlookup/pb"
	"github.com/spf13/cobra"
)

var Bind string

func init() {
	serverCmd.Flags().StringVarP(&Bind, "bind", "b", ":50051", "Address:port to bind to")
	RootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "gRPC Server",
	Run: func(cmd *cobra.Command, args []string) {
		b, err := asndb.NewAsnBackend("asn.db", "asnames.json")
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Listening on %s", Bind)
		lis, err := net.Listen("tcp", Bind)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		server := asnlookupServer{backend: b}

		// Creates a new gRPC server
		s := grpc.NewServer()
		pb.RegisterAsnlookupServer(s, &server)
		s.Serve(lis)
	},
}

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
func (s *asnlookupServer) LookupManyBatch(stream pb.Asnlookup_LookupManyBatchServer) error {
	for {
		batch, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		response := &pb.LookupReplyBatch{}
		for _, req := range batch.Requests {
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
			response.Replies = append(response.Replies, rep)
		}
		if err := stream.Send(response); err != nil {
			return err
		}
	}
}
