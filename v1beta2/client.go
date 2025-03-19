package v1beta1

import (
	"context"
	"fmt"
	kesselv1beta2 "github.com/project-kessel/inventory-api/api/kessel/inventory/v1beta2"
	nethttp "net/http"

	"github.com/authzed/grpcutil"
	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Inventory interface{}

type InventoryClient struct {
	ResourceServiceClient kesselv1beta2.KesselResourceServiceClient
	gRPCConn              *grpc.ClientConn
	tokenClient           *TokenClient
}

type InventoryHttpClient struct {
	ResourceClient kesselv1beta2.KesselResourceServiceHTTPClient
	tokenClient    *TokenClient
}

var (
	_ Inventory = &InventoryHttpClient{}
	_ Inventory = &InventoryClient{}
)

func New(config *Config) (*InventoryClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.EmptyDialOption{})
	var tokencli *TokenClient
	if config.enableOIDCAuth {
		tokencli = NewTokenClient(config)
	}

	if config.insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		tlsConfig, err := grpcutil.WithSystemCerts(grpcutil.VerifyCA)
		if err != nil {
			return nil, err
		}
		opts = append(opts, tlsConfig)
	}

	conn, err := grpc.NewClient(
		config.url,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return &InventoryClient{
		ResourceServiceClient: kesselv1beta2.NewKesselResourceServiceClient(conn),
		gRPCConn:              conn,
		tokenClient:           tokencli,
	}, err
}

func NewHttpClient(ctx context.Context, config *Config) (*InventoryHttpClient, error) {
	var tokencli *TokenClient
	if config.enableOIDCAuth {
		tokencli = NewTokenClient(config)
	}

	var opts []http.ClientOption
	if config.httpUrl != "" {
		opts = append(opts, http.WithEndpoint(config.httpUrl))
	}

	if !config.insecure {
		opts = append(opts, http.WithTLSConfig(config.tlsConfig))
	}

	client, err := http.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &InventoryHttpClient{
		ResourceClient: kesselv1beta2.NewKesselResourceHTTPClient(client),
		tokenClient:    tokencli,
	}, nil
}

func (a InventoryClient) GetTokenCallOption() ([]grpc.CallOption, error) {
	var opts []grpc.CallOption
	opts = append(opts, grpc.EmptyCallOption{})
	token, err := a.tokenClient.GetToken()
	if err != nil {
		return nil, err
	}
	if a.tokenClient.Insecure {
		opts = append(opts, WithInsecureBearerToken(token.AccessToken))
	} else {
		opts = append(opts, WithBearerToken(token.AccessToken))
	}

	return opts, nil
}

func (a InventoryHttpClient) GetTokenHTTPOption() ([]http.CallOption, error) {
	var opts []http.CallOption
	token, err := a.tokenClient.GetToken()
	if err != nil {
		return nil, err
	}
	header := nethttp.Header{}
	header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	opts = append(opts, http.Header(&header))
	return opts, nil
}
