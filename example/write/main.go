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
	// group sync write
	log.Println("group sync write")
	serverHandles := make([]uint32, len(itemList))
	for i, item := range itemList {
		serverHandles[i] = item.GetServerHandle()
	}
	errs, err = group.SyncWrite(serverHandles, value)
	if err != nil {
		log.Fatalf("sync write failed: %s\n", err)
	}
	for i, err := range errs {
		if err != nil {
			log.Fatalf("sync write item %s failed: %s\n", tags[i], err)
		}
	}
	// item write
	log.Println("item write")
	for i, item := range itemList {
		err := item.Write(value[i])
		if err != nil {
			log.Fatalf("write item %s failed: %s\n", tags[i], err)
		}
	}
	log.Println("finish")
}
