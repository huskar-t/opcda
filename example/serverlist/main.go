package main

import (
	"fmt"

	"github.com/huskar-t/opcda"
)

func main() {
	opcda.Initialize()
	defer opcda.Uninitialize()
	host := "localhost"
	serverInfos, err := opcda.GetOPCServers(host)
	if err != nil {
		panic(err)
	}
	for _, info := range serverInfos {
		fmt.Printf("ProgID: %s, ClsStr: %s, VerIndProgID: %s\n", info.ProgID, info.ClsStr, info.VerIndProgID)
	}
}
