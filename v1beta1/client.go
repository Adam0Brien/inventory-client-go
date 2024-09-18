package v1beta1

import (
	"context"
	"fmt"
	"github.com/authzed/grpcutil"
	"github.com/go-kratos/kratos/v2/transport/http"
	kesselrel "github.com/project-kessel/inventory-api/api/kessel/inventory/v1beta1/relationships"
	kessel "github.com/project-kessel/inventory-api/api/kessel/inventory/v1beta1/resources"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	nethttp "net/http"
)

type Inventory interface {
}

type inventoryClient struct {
	K8sClusterService                              kessel.KesselK8SClusterServiceClient
	K8SPolicyIsPropagatedToK8SClusterServiceClient kesselrel.KesselK8SPolicyIsPropagatedToK8SClusterServiceClient
	PolicyServiceClient                            kessel.KesselK8SPolicyServiceClient
	RhelHostServiceClient                          kessel.KesselRhelHostServiceClient
	gRPCConn                                       *grpc.ClientConn
	tokenClient                                    *TokenClient
}

type inventoryHttpClient struct {
	K8sClusterService                                  kessel.KesselK8SClusterServiceHTTPClient
	K8SPolicyIsPropagatedToK8SClusterServiceHTTPClient kesselrel.KesselK8SPolicyIsPropagatedToK8SClusterServiceHTTPClient
	PolicyServiceClient                                kessel.KesselK8SPolicyServiceHTTPClient
	RhelHostServiceClient                              kessel.KesselRhelHostServiceHTTPClient
	tokenClient                                        *TokenClient
}

var _ Inventory = &inventoryHttpClient{}
var _ Inventory = &inventoryClient{}

func New(config *Config) (*inventoryClient, error) {
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

	return &inventoryClient{
		K8sClusterService: kessel.NewKesselK8SClusterServiceClient(conn),
		K8SPolicyIsPropagatedToK8SClusterServiceClient: kesselrel.NewKesselK8SPolicyIsPropagatedToK8SClusterServiceClient(conn),
		PolicyServiceClient:                            kessel.NewKesselK8SPolicyServiceClient(conn),
		RhelHostServiceClient:                          kessel.NewKesselRhelHostServiceClient(conn),
		gRPCConn:                                       conn,
		tokenClient:                                    tokencli,
	}, err
}

func NewHttpClient(ctx context.Context, config *Config) (*inventoryHttpClient, error) {
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

	return &inventoryHttpClient{
		K8sClusterService: kessel.NewKesselK8SClusterServiceHTTPClient(client),
		K8SPolicyIsPropagatedToK8SClusterServiceHTTPClient: kesselrel.NewKesselK8SPolicyIsPropagatedToK8SClusterServiceHTTPClient(client),
		PolicyServiceClient:   kessel.NewKesselK8SPolicyServiceHTTPClient(client),
		RhelHostServiceClient: kessel.NewKesselRhelHostServiceHTTPClient(client),
		tokenClient:           tokencli,
	}, nil
}

func (a inventoryClient) GetTokenCallOption() ([]grpc.CallOption, error) {
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

func (a inventoryHttpClient) GetTokenHTTPOption() ([]http.CallOption, error) {
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
