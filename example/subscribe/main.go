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
	ch := make(chan *opcda.DataChangeCallBackData, 100)
	go func() {
		for {
			select {
			case data := <-ch:
				log.Printf("data change received, transaction id: %d, group handle: %d, masterQuality: %d, masterError: %v\n", data.TransID, data.GroupServerHandle, data.MasterQuality, data.MasterErr)
				for i := 0; i < len(data.ItemClientHandles); i++ {
					tag := ""
					for _, item := range itemList {
						if item.GetClientHandle() == data.ItemClientHandles[i] {
							tag = item.GetItemID()
						}
					}
					log.Printf("item %s\ttimestamp: %s\tquality: %d\tvalue: %v\n", tag, data.TimeStamps[i], data.Qualities[i], data.Values[i])
				}
			}
		}
	}()
	err = group.RegisterDataChange(ch)
	if err != nil {
		log.Fatalf("register data change failed: %s\n", err)
	}
	select {}
}
