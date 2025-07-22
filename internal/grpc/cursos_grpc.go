package grpc

import pb "github.com/dsniels/storage-service/proto"

type CursosGrpc struct {
	pb.UnimplementedCursosProtoServiceServer
}

func NewCursosGrpc() {
}
