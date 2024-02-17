package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/huskar-t/opcda"
)

func main() {
	opcda.Initialize()
	defer opcda.Uninitialize()
	host := "localhost"
	progID := "Matrikon.OPC.Simulation.1"
	tags := []string{
		"Bucket Brigade.ArrayOfReal8",
		"Bucket Brigade.ArrayOfString",
		"Bucket Brigade.Boolean",
		"Bucket Brigade.Int1",
		"Bucket Brigade.Int2",
		"Bucket Brigade.Int4",
		"Bucket Brigade.Real4",
		"Bucket Brigade.Real8",
		"Bucket Brigade.String",
		"Bucket Brigade.Time",
		"Bucket Brigade.UInt1",
		"Bucket Brigade.UInt2",
		"Bucket Brigade.UInt4",
	}
	value := []interface{}{
		[]float64{1.2, 2.3, 3.4},
		[]string{"hello", "world"},
		true,
		int8(1),
		int16(2),
		int32(3),
		float32(5.5777),
		float64(6.6777777777),
		"hello",
		time.Now(),
		uint8(7),
		uint16(8),
		uint32(9),
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
	time.Sleep(time.Second)
	// group async write
	serverHandles := make([]uint32, len(itemList))
	for i, item := range itemList {
		serverHandles[i] = item.GetServerHandle()
	}
	ch := make(chan *opcda.WriteCompleteCallBackData, 100)
	err = group.RegisterWriteComplete(ch)
	if err != nil {
		log.Fatalf("register write complete callback failed: %s\n", err)
	}
	finish := make(chan struct{})
	go func() {
		for {
			select {
			case data := <-ch:
				tagList := make([]string, len(data.ItemClientHandles))
				for i, handle := range data.ItemClientHandles {
					for _, item := range itemList {
						if item.GetClientHandle() == handle {
							tagList[i] = item.GetItemID()
						}
					}
				}
				fmt.Printf("write complete received\ntransaction id: %d\ngroup handle: %d\nmasterError: %v\nitems: [%s]\n", data.TransID, data.GroupHandle, data.MasterErr, strings.Join(tagList, ","))
				for i, err := range data.Errors {
					if err != nil {
						log.Printf("async write item %s failed: %s\n", tagList[i], err)
					}
				}
				close(finish)
				return
			}
		}
	}()
	transID := uint32(1)
	_, errs, err = group.AsyncWrite(serverHandles, value, transID)
	if err != nil {
		log.Fatalf("async write failed: %s\n", err)
	}
	for i, err := range errs {
		if err != nil {
			log.Fatalf("async write item %s failed: %s\n", tags[i], err)
		}
	}
	select {
	case <-finish:
	}
}
