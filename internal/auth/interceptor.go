package auth

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UnaryInterceptor for validating the bearer token
func (s *serverAPI) validateBearerTokenInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// Extract metadata from incoming context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	// Get the authorization tokens from metadata
	tokens := md.Get("authorization")
	if len(tokens) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing authorization token")
	}

	// Validate the token
	connectionToken := tokens[0]
	if connectionToken != fmt.Sprintf("Bearer %s", s.connectionToken) {
		return nil, status.Error(codes.Unauthenticated, "invalid authorization token")
	}

	// Call the handler to proceed with the actual RPC
	return handler(ctx, req)
}
