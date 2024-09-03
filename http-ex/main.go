package main

import (
	"context"
	"fmt"
	kessel "github.com/project-kessel/inventory-api/api/kessel/inventory/v1beta1/resources"
	"inventory-client-go/v1beta1"
)

func main() {

	client, err := v1beta1.NewHttpClient(context.Background(),
		v1beta1.NewConfig(v1beta1.WithHTTPUrl("localhost:8081")))
	v1beta1.WithTLSInsecure(true)
	//v1beta1.WithAuthEnabled("svc-test", "", "http://localhost:8084/realms/redhat-external/protocol/openid-connect/token"),
	//v1beta1.WithHTTPTLSConfig(tls.Config{})
	if err != nil {
		fmt.Println(err)
	}
	request := kessel.CreateRhelHostRequest{RhelHost: &kessel.RhelHost{
		Metadata: &kessel.Metadata{
			ResourceType: "rhel-host",
			Workspace:    "",
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
	//optts, err := client.GetTokenHTTPOption()
	resp, err := client.RhelHostServiceClient.CreateRhelHost(context.Background(), &request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}
