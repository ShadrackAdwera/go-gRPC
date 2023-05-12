package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	GRPC_GATEWAY_USER_AGENT_HEADER = "grpcgateway-user-agent"
	X_FORWARDED_FOR_HEADER         = "x-forwarded-for"
	USER_AGENT_HEADER              = "user-agent"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func extractMetadata(ctx context.Context) *Metadata {
	meta := &Metadata{}
	if mtd, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := mtd.Get(GRPC_GATEWAY_USER_AGENT_HEADER); len(userAgents) > 0 {
			meta.UserAgent = userAgents[0]
		}
		if clientIps := mtd.Get(X_FORWARDED_FOR_HEADER); len(clientIps) > 0 {
			meta.ClientIP = clientIps[0]
		}

		if userAgents := mtd.Get(USER_AGENT_HEADER); len(userAgents) > 0 {
			meta.UserAgent = userAgents[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		meta.ClientIP = p.Addr.String()
	}

	return meta
}
