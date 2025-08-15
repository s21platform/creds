package service

import (
	"context"

	pb "github.com/s21platform/creds/pkg/creds"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server реализует gRPC сервер для работы с учетными данными
type Server struct {
	pb.UnimplementedCredentialsServiceServer
	credentials map[string]map[string]string
}

// NewServer создает новый экземпляр сервера с предустановленными учетными данными
func NewServer() *Server {
	return &Server{
		credentials: map[string]map[string]string{
			"NOTIFICATION": {
				"NOTIFICATION_SERVICE_HOST": "https://api.example.com",
				"NOTIFICATION_SERVICE_PORT": "443",
				"NOTIFICATION_SERVICE_API_KEY": "secret-key-123",
			},
			"DATABASE": {
				"DATABASE_URL": "postgres://user:pass@localhost:5432/db",
			},
			"REDIS": {
				"REDIS_HOST": "localhost",
				"REDIS_PORT": "6379",
			},
		},
	}
}

// GetCreds возвращает значения для запрошенных переменных
func (s *Server) GetCreds(ctx context.Context, req *pb.GetCredsRequest) (*pb.GetCredsResponse, error) {
	if len(req.Names) == 0 {
		return nil, status.Error(codes.InvalidArgument, "names list is empty")
	}

	response := &pb.GetCredsResponse{
		Credentials: make([]*pb.Credential, 0, len(req.Names)),
	}

	for _, name := range req.Names {
		if value, ok := s.credentials[name]; ok {
			for key, value := range value {
				response.Credentials = append(response.Credentials, &pb.Credential{
					Name:  key,
					Value: value,
				})
			}
		}
	}

	return response, nil
}
