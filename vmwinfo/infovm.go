package vmwinfo

import (
	"context"
	"fmt"
	"log"
	"time"

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
	Update(Info)
}

type Printer interface{
	Display()
}

type InfoVMware struct {
	Name 			string
	receivers 		[]Receiver
	ctx				context.Context
	client			*govmomi.Client
	period			time.Duration
	updateC			chan Info
	stopC			chan string
}

type Info struct {
	Hosts []mo.HostSystem
	Nets  []mo.Network
	Dss   []mo.Datastore
	Vms   []mo.VirtualMachine
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

func (info *InfoVMware)BroadCast(){}

func (info *InfoVMware)GetDss(ctx context.Context, c *govmomi.Client)([]mo.Datastore, error){
	m := view.NewManager(c.Client)
	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"Datastore"}, true)
	if err != nil {
		return nil, err
	}
	defer v.Destroy(ctx)

	var dss []mo.Datastore
	fmt.Printf("<GetDss %v> Got datastores...\n", time.Now())
	start := time.Now()
	err = v.Retrieve(ctx, []string{"Datastore"}, []string{"summary"}, &dss)
	//err = v.Retrieve(ctx, []string{"HostSystem"}, nil, &hss)
	//err = v.Retrieve(ctx, []string{"HostSystem"}, []string{"summary"}, &hss)
	if err != nil {
		return nil, err
	}
	fmt.Printf("<GetDss %v>List Datastore spend %v.\n",time.Now(),  time.Since(start))
	return dss, nil
}


func (info *InfoVMware)GetHosts(ctx context.Context, c *govmomi.Client)([]mo.HostSystem, error){
	m := view.NewManager(c.Client)
	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"HostSystem"}, true)
	if err != nil {
		return nil, err
	}
	defer v.Destroy(ctx)

	var hss []mo.HostSystem
	fmt.Printf("<GetHosts %v> Got hosts...\n", time.Now())
	start := time.Now()
	err = v.Retrieve(ctx, []string{"HostSystem"}, []string{"config.network.portgroup", "summary"}, &hss)
	//err = v.Retrieve(ctx, []string{"HostSystem"}, nil, &hss)
	//err = v.Retrieve(ctx, []string{"HostSystem"}, []string{"summary"}, &hss)
	if err != nil {
		return nil, err
	}
	fmt.Printf("<GetHosts %v>List HostSystem spend %v.\n",time.Now(),  time.Since(start))
	return hss, nil
}

func (info *InfoVMware)GetNetworks(ctx context.Context, c *govmomi.Client)([]mo.Network, error){
	
	m := view.NewManager(c.Client)
	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"Network"}, true)
	if err != nil {
		log.Fatal(err)
	}
	defer v.Destroy(ctx)
	
	var networks []mo.Network
	start := time.Now()
	err = v.Retrieve(ctx, []string{"Network"}, nil, &networks)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("<GetNetworks %v>List Network spend %v.\n",time.Now(),  time.Since(start))
	
	return networks, nil
}

func (info *InfoVMware)GetVms(ctx context.Context, c *govmomi.Client)([]mo.VirtualMachine, error){
	m := view.NewManager(c.Client)
	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		log.Fatal(err)
	}
	defer v.Destroy(ctx)
	
	var vms []mo.VirtualMachine
	start := time.Now()
	err = v.Retrieve(ctx, []string{"VirtualMachine"}, []string{"summary"}, &vms)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("<GetVms %v>List VirtualMachine spend %v.\n",time.Now(),  time.Since(start))
	
	return vms, nil
}

func (info *InfoVMware)Collect()Info{
	hss, err := info.GetHosts(info.ctx, info.client)
	if err != nil{
		log.Fatal(err)
	}
	nets, err := info.GetNetworks(info.ctx, info.client)
	if err != nil{
		log.Fatal(err)
	}
	dss, err := info.GetDss(info.ctx, info.client)
	if err != nil{
		log.Fatal(err)
	}
	vms, err := info.GetVms(info.ctx, info.client)
	if err != nil{
		log.Fatal(err)
	}	
	return Info{
		Hosts : hss,
		Nets  : nets,
		Dss   : dss,
		Vms   : vms,
	}
}

func (info *InfoVMware)Run(){
	go func(){
		for {
			select {
				case packet := <-info.updateC:
					fmt.Printf("<RunLoop %v>Receive hosts and networks update info.\n", time.Now())
					for _, receiver := range info.receivers {
						receiver.Update(packet)
					}
				case <-info.stopC:
					return	
			}
		}
	}()
	
	info.updateC <- info.Collect()
	
	ticker := time.NewTicker(info.period)
	go func(){
		for {
			select {
				case <-ticker.C:
					fmt.Printf("<RunLoop %v>Collect hosts and networks info.\n", time.Now())
					info.updateC <- info.Collect()
				case <-info.stopC:
					return	
			}	
		}
	}()
	

}

func NewInfoVMware(name string, ctx context.Context, client *govmomi.Client, period time.Duration, stopC chan string)*InfoVMware{
	return &InfoVMware{
		Name: 			name,
		ctx:  			ctx,
		client: 		client,
		period: 		period,
		updateC:		make(chan Info),
		stopC:			stopC,
	}
}