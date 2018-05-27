package main

import (
	"context"
	"fmt"
	"log"

	
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

type Publisher interface{
	AddReceiver(r Receiver)
	RemoveReceiver(r Receiver)
	BroadCast()
}

type Receiver interface{
	Update(hss []mo.HostSystem)
}

type Printer interface{
	Display()
}

type InfoVMware struct {
	Name 		string
	receivers 	[]Receiver
	ctx			context.Context
	client		*govmomi.Client	
}

func (info *InfoVMware)AddReceiver(r Receiver){
	info.receivers = append(info.receivers, r)
}

func (info *InfoVMware)RemoveReceiver(r Receiver){
	for i, receiver := range info.receivers {
		if receiver == r {
			info.receivers = append(info.receivers[:i], info.receivers[i + 1:]...)
			return
		}
	}
}

func (info *InfoVMware)BroadCast(){
	hss, err := info.GetHosts(info.ctx, info.client)
	if err != nil {
		log.Fatal(err)
	}
	for _, receiver := range info.receivers {
		receiver.Update(hss)
	}
}

func (info *InfoVMware)GetHosts(ctx context.Context, c *govmomi.Client)([]mo.HostSystem, error){
	m := view.NewManager(c.Client)
	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"HostSystem"}, true)
	if err != nil {
		return nil, err
	}
	
	defer v.Destroy(ctx)

	var hss []mo.HostSystem
	fmt.Printf("Got hosts...\n")
	err = v.Retrieve(ctx, []string{"HostSystem"}, []string{"summary"}, &hss)
	if err != nil {
		return nil, err
	}
	return hss, nil
}

func NewInfoVMware(name string, ctx context.Context, client *govmomi.Client)*InfoVMware{
	return &InfoVMware{
		Name: name,
		ctx:  ctx,
		client: client,
	}
}