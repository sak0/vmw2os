package main

import (
	"context"
	"fmt"
	"log"
	
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

func GetNetworks(ctx context.Context, c *govmomi.Client){
	m := view.NewManager(c.Client)
	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"Network"}, true)
	if err != nil {
		log.Fatal(err)
	}
	defer v.Destroy(ctx)
	
	var networks []mo.Network
	err = v.Retrieve(ctx, []string{"Network"}, nil, &networks)
	if err != nil {
		log.Fatal(err)
	}

	for _, net := range networks {
		fmt.Printf("%s: %s\n", net.Name, net.Reference())
	}
}