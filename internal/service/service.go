package service

import (
	"context"
	"encoding/json"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/s21platform/creds/internal/model"
	pb "github.com/s21platform/creds/pkg/creds"
)

const (
	tokenKey = "token"
)

// Service реализует gRPC сервер для работы с учетными данными
type Service struct {
	pb.UnimplementedCredentialsServiceServer
	dbR DbRepo
}

// NewServer создает новый экземпляр сервера с предустановленными учетными данными
func New(dbR DbRepo) *Service {
	return &Service{
		dbR: dbR,
	}
}

// validateToken проверяет токен из метаданных запроса
func (s *Service) validateToken(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("ff to get token")

		return status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	tokens := md.Get(tokenKey)
	if len(tokens) == 0 {
		log.Println("nos such")

		return status.Error(codes.Unauthenticated, "token is not provided")
	}

	token, err := s.dbR.GetToken(ctx, tokens[0])
	if err != nil {
		log.Println("Failed to get token")
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	if token != tokens[0] {
		log.Println("not eq")

		return status.Error(codes.Unauthenticated, "invalid token")
	}

	return nil
}

// GetCreds возвращает значения для запрошенных переменных
func (s *Service) GetCreds(ctx context.Context, req *pb.GetCredsRequest) (*pb.GetCredsResponse, error) {
	// Проверяем токен перед выполнением основной логики
	if err := s.validateToken(ctx); err != nil {
		log.Println("Fa")

		return nil, err
	}

	if len(req.Names) == 0 {
		log.Println("Failn")

		return nil, status.Error(codes.InvalidArgument, "names list is empty")
	}

	response := &pb.GetCredsResponse{
		Credentials: make([]*pb.Credential, 0, len(req.Names)),
	}

	for _, name := range req.Names {
		credsData, err := s.dbR.GetCreds(ctx, name)
		if err != nil {
			log.Println("Faile", name, err)

			return nil, err
		}

		var creds model.Creds
		if err := json.Unmarshal(credsData.Data, &creds); err != nil {
			log.Println("Failed to ")

			return nil, err
		}
		for _, cred := range creds.Creds {
			response.Credentials = append(response.Credentials, &pb.Credential{
				Name:  cred.Name,
				Value: cred.Value,
			})
		}
	}

	return response, nil
}
