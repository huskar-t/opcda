package main

import (
	"log"
	"time"

	"github.com/huskar-t/opcda"
)

func main() {
	opcda.Initialize()
	defer opcda.Uninitialize()
	host := "localhost"
	progID := "Matrikon.OPC.Simulation.1"
	tags := []string{
		//"Random.ArrayOfReal8",
		//"Random.ArrayOfString",
		//"Random.Boolean",
		//"Random.Int1",
		//"Random.Int2",
		//"Random.Int4",
		//"Random.Int8",
		//"Random.Qualities",
		//"Random.Real4",
		//"Random.Real8",
		//"Random.String",
		//"Random.Time",
		//"Random.UInt1",
		//"Random.UInt2",
		//"Random.UInt4",
		//"Random.UInt8",
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
	status, err := group.SyncRead(opcda.OPC_DS_CACHE, serverHandles)
	if err != nil {
		log.Fatalf("sync read failed: %s\n", err)
	}
	for i, item := range status {
		log.Printf("%s:\t%s\t%d\t%v\n", tags[i], item.Timestamp, item.Quality, item.Value)
	}
	// item read
	log.Println("item read")
	for i, item := range itemList {
		value, quality, timestamp, err := item.Read(opcda.OPC_DS_CACHE)
		if err != nil {
			log.Fatalf("read item failed: %s\n", err)
		}
		log.Printf("%s:\t%s\t%d\t%v\n", tags[i], timestamp, quality, value)
	}
}
