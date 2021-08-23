package main

import (
	_ "github.com/vultr/vultr-cloud-controller-manager/vultr"
	_ "k8s.io/component-base/metrics/prometheus/clientgo"
	_ "k8s.io/component-base/metrics/prometheus/version" // for version metrics registration
)

func main() {
	//rand.Seed(time.Now().UnixNano())
	//
	//command := app.NewCloudControllerManagerCommand()
	//
	//logs.InitLogs()
	//defer logs.FlushLogs()
	//
	//if err := command.Execute(); err != nil {
	//	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	//	os.Exit(1)
	//}
}
