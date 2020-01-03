package cmd

import (
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	"golang.org/x/net/context"

	"github.com/JustinAzoff/hostlookup/hostdb"
	pb "github.com/JustinAzoff/hostlookup/pb"
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
		b, err := hostdb.NewHostBackend("shrunken.csv.gz")
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Listening on %s", Bind)
		lis, err := net.Listen("tcp", Bind)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		server := hostlookupServer{backend: b}

		// Creates a new gRPC server
		s := grpc.NewServer()
		pb.RegisterHostlookupServer(s, &server)
		s.Serve(lis)
	},
}

//type hostlookupServer interface {
//        Hello(context.Context, *Empty) (*HelloReply, error)
//        Lookup(context.Context, *LookupRequest) (*LookupReply, error)
//        LookupMany(hostlookup_LookupManyServer) error
//        LookupBatch(context.Context, *LookupRequestBatch) (*LookupReplyBatch, error)
//}

type hostlookupServer struct {
	backend *hostdb.HostBackend
}

func (s *hostlookupServer) Hello(ctx context.Context, empty *pb.Empty) (*pb.HelloReply, error) {
	log.Printf("Hello called!")
	return &pb.HelloReply{
		Message: "Hello!",
	}, nil
}
func (s *hostlookupServer) Lookup(ctx context.Context, req *pb.LookupRequest) (*pb.LookupReply, error) {
	rec, err := s.backend.Lookup(req.Address)
	if err != nil {
		return nil, err
	}
	return &pb.LookupReply{
		Host: rec.Host,
	}, nil
}
func (s *hostlookupServer) LookupMany(stream pb.Hostlookup_LookupManyServer) error {
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
			Address: req.Address,
			Host:    rec.Host,
		}
		if err := stream.Send(rep); err != nil {
			return err
		}
	}
}
func (s *hostlookupServer) LookupManyBatch(stream pb.Hostlookup_LookupManyBatchServer) error {
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
				Address: req.Address,
				Host:    rec.Host,
			}
			response.Replies = append(response.Replies, rep)
		}
		if err := stream.Send(response); err != nil {
			return err
		}
	}
}
