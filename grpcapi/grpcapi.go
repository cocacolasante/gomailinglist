package grpcapi

import (
	"context"
	"database/sql"
	"gomailinglist/mdb"
	"log"
	pb "mailinglist/proto"
	"net"
	"time"


	"google.golang.org/grpc"
)

type MailServer struct {
	pb.UnimplementedMailingListServiceServer
	db *sql.DB
}

func pbEntryToMdbEntry(pbEntry *pb.EmailEntry) mdb.EmailEntry {
	t := time.Unix(pbEntry.ConfirmedAt, 0)

	return mdb.EmailEntry{
		Id: pbEntry.Id,
		Email: pbEntry.Email,
		ConfirmedAt: &t,
		OptOut: pbEntry.OptOut,
	}
}

func mdbEntryToPbEntry(mdbEntry *mdb.EmailEntry) pb.EmailEntry {
	return pb.EmailEntry{
		Id: mdbEntry.Id,
		Email: mdbEntry.Email,
		ConfirmedAt: mdbEntry.ConfirmedAt.Unix(),
		OptOut: mdbEntry.OptOut,
	}
}

func emailResponse(db *sql.DB, email string ) (*pb.emailResponse, error) {
	entry, err := mdb.GetEmail(db, email)
	if err != nil {
		return &pb.emailResponse{}, err
	}

	if entry != nil {
		return &pb.emailResponse{}, nil
	}
	res := mdbEntryToPbEntry(entry)
	return &pb.emailResponse{EmailEntry: &res}, nil
}

func (s *MailServer) GetEmail(ctx context.Context, request *pb.GetEmailRequest) (*pb.GetEmailResponse, error ){
	log.Printf("GRPC GetEmail: %v \n", request)
	return emailResponse(s.db, request.EmailAddr)
	
}
func (s *MailServer) GetBatchEmail(ctx context.Context, request *pb.GetBatchEmailRequest) (*pb.GetEmailBatchResponse, error ){
	log.Printf("GRPC GetEmail: %v \n", request)

	params := mdb.GetEmailBatchQueryParams{
		Page: int(request.Page),
		Count: int(request.Count),
	}

	mdbEntries, err := mdb.GetEmailBatch(s.db, params)
	if err != nil {
		return &pb.GetBatchEmailResponse{}, err
	}

	pbEntries := make([]*pb.EmailEntry, 0, len(mdbEntries))

	for i :=0; i < len(mdbEntries); i++{
		entry := mdbEntryToPbEntry(&mdbEntries[i])
		pbEntries = append(pbEntries, entry)


	}

	return &pb.GetEmailBatchResponse{EmailEntries: pbEntries}

}

func (s *MailServer) CreateEmail(ctx context.Context, request *pb.CreateEmailRequest) (*pb.GetEmailResponse, error ){
	log.Printf("GRPC CreateEmail: %v \n", request)
	err := mdb.CreateEmail(s.db, request.EmailAddr)

	if err != nil {
		return &pb.EmailResponse{}, err
	}

	return emailResponse(s.db, request.EmailAddr)

	
}
func (s *MailServer) UpdateEmail(ctx context.Context, request *pb.UpdateEmailRequest) (*pb.GetEmailResponse, error ){
	log.Printf("GRPC UpdateEmail: %v \n", request)

	entry := pbEntryToMdbEntry(request.EmailEntry)

	err := mdb.UpdateEmail(s.db, entry)

	if err != nil {
		return &pb.EmailResponse{}, err
	}

	return emailResponse(s.db, entry.Email)

	
}

func (s *MailServer) DeleteEmail(ctx context.Context, request *pb.DeleteEmailRequest) (*pb.GetEmailResponse, error ){
	log.Printf("GRPC DeleteEmail: %v \n", request)


	err := mdb.DeleteEmail(s.db, request.Email)

	if err != nil {
		return &pb.EmailResponse{}, err
	}

	return emailResponse(s.db, request.EmailAddr)


	
}

func Serve(db *sql.DB, bind string) {
	listener, err := net.Listen("tcp", bind)

	if err != nil {
		log.Fatalf("Cannot connect to %v\n", bind)

	}

	gRPCServer := grpc.NewServer()
	mailServer := MailServer{db: db}

	pb.RegisterMailingListServiceServer(gRPCServer, &mailServer)

	log.Printf("gRPC Api server listening to %v\n", bind)
	if err := gRPCServer.Serve(listener); err != nil {
		log.Fatalf("Cannot bind server due to error %v\n", err )
	}


}
