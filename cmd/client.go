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
	"github.com/spf13/cobra"
)

//type AsnlookupClient interface {
//        Hello(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*HelloReply, error)
//        Lookup(ctx context.Context, in *LookupRequest, opts ...grpc.CallOption) (*LookupReply, error)
//        LookupMany(ctx context.Context, opts ...grpc.CallOption) (Asnlookup_LookupManyClient, error)
//        LookupBatch(ctx context.Context, in *LookupRequestBatch, opts ...grpc.CallOption) (*LookupReplyBatch, error)
//}

const (
	batchSize = 200
)

var Connect string

func init() {
	clientCmd.Flags().StringVarP(&Connect, "connect", "c", "localhost:50051", "Address:port to connect to")
	RootCmd.AddCommand(clientCmd)
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "gRPC client",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial(Connect, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		client := pb.NewAsnlookupClient(conn)
		client.Hello(context.Background(), &pb.Empty{})

		stream, err := client.LookupManyBatch(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		waitc := make(chan struct{})
		go func() {
			for {
				resp, err := stream.Recv()
				if err == io.EOF {
					// read done.
					close(waitc)
					return
				}
				if err != nil {
					log.Fatalf("Failed to receive a record: %v", err)
				}
				for _, rec := range resp.Replies {
					fmt.Printf("%s\t%s\t%d\t%s\t%s\n", rec.Prefix, rec.Address, rec.As, rec.Owner, rec.Cc)
				}
			}
		}()
		scanner := bufio.NewScanner(os.Stdin)
		batch := &pb.LookupRequestBatch{}
		for scanner.Scan() {
			ip := scanner.Text()
			req := &pb.LookupRequest{Address: ip}
			batch.Requests = append(batch.Requests, req)
			if len(batch.Requests) >= batchSize {
				err := stream.Send(batch)
				if err != nil {
					log.Fatalf("Failed to send: %v", batch)
				}
				batch.Requests = batch.Requests[:0]
			}
		}
		if len(batch.Requests) > 0 {
			err := stream.Send(batch)
			if err != nil {
				log.Fatalf("Failed to send: %v", batch)
			}
		}
		stream.CloseSend()
		<-waitc
	},
}
