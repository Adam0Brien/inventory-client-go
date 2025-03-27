package main

import (
	"context"
	"fmt"

	kessel "github.com/project-kessel/inventory-api/api/kessel/inventory/v1beta1/resources"
	"github.com/project-kessel/inventory-client-go/v1beta1"
)

func main() {
	client, err := v1beta1.New(v1beta1.NewConfig(
		v1beta1.WithgRPCUrl("localhost:9081"),
		v1beta1.WithTLSInsecure(true),
		v1beta1.WithAuthEnabled("", "", ""),
	))
	if err != nil {
		fmt.Println(err)
	}

	request := kessel.CreateRhelHostRequest{RhelHost: &kessel.RhelHost{
		Metadata: &kessel.Metadata{
			ResourceType: "rhel-host",
			WorkspaceId:  "",
		},
		ReporterData: &kessel.ReporterData{
			ReporterType:       kessel.ReporterData_ACM,
			ReporterInstanceId: "service-account-svc-test",
			ConsoleHref:        "www.example.com",
			ApiHref:            "www.example.com",
			LocalResourceId:    "1",
			ReporterVersion:    "0.1",
		},
	}}

	opts, err := client.GetTokenCallOption()
	if err != nil {
		fmt.Println(err)
	}
	_, err = client.RhelHostServiceClient.CreateRhelHost(context.Background(), &request, opts...)
	if err != nil {
		fmt.Println(err)
	}
}
