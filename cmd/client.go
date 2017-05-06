package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	pb "github.com/JustinAzoff/asnlookup/pb"
)

const (
	address = "localhost:50051"
)

//type AsnlookupClient interface {
//        Hello(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*HelloReply, error)
//        Lookup(ctx context.Context, in *LookupRequest, opts ...grpc.CallOption) (*LookupReply, error)
//        LookupMany(ctx context.Context, opts ...grpc.CallOption) (Asnlookup_LookupManyClient, error)
//        LookupBatch(ctx context.Context, in *LookupRequestBatch, opts ...grpc.CallOption) (*LookupReplyBatch, error)
//}

func client() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewAsnlookupClient(conn)
	client.Hello(context.Background(), &pb.Empty{})

	stream, err := client.LookupMany(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	waitc := make(chan struct{})
	go func() {
		for {
			rec, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a record: %v", err)
			}
			fmt.Printf("%s\t%s\t%d\t%s\t%s\n", rec.Prefix, rec.Address, rec.As, rec.Owner, rec.Cc)
		}
	}()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		ip := scanner.Text()
		req := &pb.LookupRequest{Address: ip}
		err := stream.Send(req)
		if err != nil {
			log.Fatalf("Failed to send: %v", req)
		}
	}
	stream.CloseSend()
	<-waitc
}
