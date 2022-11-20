package main

import (
	"context"
	"fmt"

	sw "github.com/ionos-cloud/uptrends-go"
)

func main() {
	auth := context.WithValue(context.Background(), sw.ContextBasicAuth, sw.BasicAuth{
		UserName: "1d14e32965a74c34a2f6641510f067f0",
		Password: "X6nZkH+TXrV9Org952BbyosEHZZpfdw7",
	})

	client := sw.NewAPIClient(sw.NewConfiguration())
	m, _, err := client.MonitorApi.MonitorGetMonitor(auth, "c8c90ddd-91a0-4556-bfb4-41a5fcd30e37", &sw.MonitorApiMonitorGetMonitorOpts{})
	if err != nil {
		panic(err)
	}

	fmt.Println(m.Name)
}
