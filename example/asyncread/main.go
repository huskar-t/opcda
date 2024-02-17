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
	// group async read
	log.Println("group sync read")
	serverHandles := make([]uint32, len(itemList))

	for i, item := range itemList {
		serverHandles[i] = item.GetServerHandle()
	}
	ch := make(chan *opcda.ReadCompleteCallBackData, 100)
	err = group.RegisterReadComplete(ch)
	if err != nil {
		log.Fatalf("register read complete callback failed: %s\n", err)
	}
	finish := make(chan struct{})
	go func() {
		for {
			select {
			case data := <-ch:
				log.Printf("read complete received, transaction id: %d, group handle: %d, masterQuality: %d, masterError: %v\n", data.TransID, data.GroupHandle, data.MasterQuality, data.MasterErr)
				tag := ""
				for i := 0; i < len(data.ItemClientHandles); i++ {
					for _, item := range itemList {
						if item.GetClientHandle() == data.ItemClientHandles[i] {
							tag = item.GetItemID()
						}
					}
					log.Printf("%s:\t%s\t%d\t%v\n", tag, data.TimeStamps[i], data.Qualities[i], data.Values[i])
				}
				close(finish)
				return
			}
		}
	}()
	transID := uint32(1)
	_, errs, err = group.AsyncRead(serverHandles, transID)
	if err != nil {
		log.Fatalf("sync read failed: %s\n", err)
	}
	for _, err := range errs {
		if err != nil {
			log.Fatalf("sync read failed: %s\n", err)
		}
	}
	select {
	case <-finish:
	}
}
