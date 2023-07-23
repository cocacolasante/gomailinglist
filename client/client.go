package main

import (
	"context"
	pb "gomailinglist/proto"
	"log"
	"time"

	"github.com/alexflint/go-arg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func logResponse(res *pb.EmailResponse, err error ){
	if err != nil {
		log.Fatalf("Error with %v\n", err)

	}

	if res.EmailEntry == nil {
		log.Printf("Not Found")
	}else {
		log.Printf("response: %v \n", res.EmailEntry)

	}
}

func CreateEmail(client pb.MailingListServiceClient, addr string) {
	log.Println("create EMail")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.CreateEmail(ctx, &pb.CreateEmailRequest{EmailAddr: addr})
	logResponse(res, err)
	

}
func GetEmail(client pb.MailingListServiceClient, addr string) {
	log.Println("get EMail")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GetEmail(ctx, &pb.GetEmailRequest{EmailAddr: addr})
	logResponse(res, err)


}
func GetBatchEmail(client pb.MailingListServiceClient, count int, page int) {
	log.Println("get batch email")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GetEmail(ctx, &pb.GetEmailRequest{EmailAddr: addr})
	if err != nil {
		log.Fatalf("Error %v \n", err)
	}

	for i := 0; i < len(res.EmailEntries); i++{
		log.Printf("Item [%v of %v], %s", i +1, len(res.EmailEntries), res.EmailEntries[i])

	}



}

func UpdateEmail(client pb.MailingListServiceClient, entry pb.EmailEntry) {
	log.Println("update EMail")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.UpdateEmail(ctx, &pb.UpdateEmailRequest{EmailAddr: &entry})
	logResponse(res, err)
	

}
func DeleteEmail(client pb.MailingListServiceClient, addr string) {
	log.Println("delete EMail")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.DeleteEmail(ctx, &pb.DeleteEmailRequest{EmailAddr: addr})
	logResponse(res, err)
	

}

var args struct {
	GrpcAddr string `arg:"env:MAILINGLIST_GRPC_ADDR`
}

func main() {
	arg.MustParse(&args)

	if args.GrpcAddr == "" {
		args.GrpcAddr = ":8081"
	}

	conn, err := grpc.Dial(args.GrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect %v", err)

	}

	defer conn.Close()

	client := pb.NewMailingListServiceClient(conn)
	
}