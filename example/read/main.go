package main

import (
	"log"
	"time"

	"github.com/huskar-t/opcda"
	"github.com/huskar-t/opcda/com"
)

func main() {
	err := com.Initialize()
	if err != nil {
		panic(err)
	}
	defer com.Uninitialize()
	host := "localhost"
	progID := "Matrikon.OPC.Simulation.1"
	tags := []string{
		"Random.ArrayOfReal8",
		"Random.ArrayOfString",
		"Random.Boolean",
		"Random.Int1",
		"Random.Int2",
		"Random.Int4",
		"Random.Int8",
		"Random.Qualities",
		"Random.Real4",
		"Random.Real8",
		"Random.String",
		"Random.Time",
		"Random.UInt1",
		"Random.UInt2",
		"Random.UInt4",
		"Random.UInt8",
		"Read Error.UInt4",
		// Errors on purpose to show error handling
		"Write Only.UInt4",
		"Bucket Brigade.Real4",
	}
	server, err := opcda.Connect(progID, host)
	if err != nil {
		log.Fatalf("connect to opc server failed: %s\n", err)
	}
	defer server.Disconnect()
	groups := server.GetOPCGroups()
	group, err := groups.Add("group1")
	if err != nil {
		log.Fatalf("add group failed: %s\n", err)
	}
	items := group.OPCItems()
	itemList, errs, err := items.AddItems(tags)
	if err != nil {
		log.Fatalf("add items failed: %s\n", err)
	}
	for i, err := range errs {
		if err != nil {
			log.Fatalf("add item %s failed: %s\n", tags[i], err)
		}
	}
	// Wait for the OPC server to be ready
	time.Sleep(time.Second * 2)
	// group sync read
	log.Println("group sync read")
	serverHandles := make([]uint32, len(itemList))

	for i, item := range itemList {
		serverHandles[i] = item.GetServerHandle()
	}
	status, resultErrs, err := group.SyncRead(opcda.OPC_DS_CACHE, serverHandles)
	if err != nil {
		log.Fatalf("sync read failed: %s\n", err)
	}
	for i, item := range status {
		if resultErrs[i] == nil {
			log.Printf("%s:\t%s\t%d\t%v\n", tags[i], item.Timestamp, item.Quality, item.Value)
		} else {
			log.Printf("%s: error %v", tags[i], resultErrs[i])
		}
	}
	// item read
	log.Println("item read")
	for i, item := range itemList {
		value, quality, timestamp, err := item.Read(opcda.OPC_DS_CACHE)
		if err != nil {
			log.Printf("%s: error %v\n", tags[i], err)
		} else {
			log.Printf("%s:\t%s\t%d\t%v\n", tags[i], timestamp, quality, value)
		}
	}
}
