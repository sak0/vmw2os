package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	//"github.com/vmware/govmomi/vim25/types"
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

func GetHostNetwork(ctx context.Context, c *govmomi.Client){
	m := view.NewManager(c.Client)
	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"HostSystem"}, true)
	if err != nil {
		log.Fatal(err)
	}
	defer v.Destroy(ctx)
	
	var hss []mo.HostSystem
	err = v.Retrieve(ctx, []string{"HostSystem"}, []string{"config.network.portgroup", "summary"}, &hss)
	if err != nil {
		log.Fatal(err)
	}

	for _, hs := range hss {
		fmt.Printf("***%s***\n", hs.Summary.Config.Name)
		tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
		fmt.Fprintf(tw, "Name:\tVswtich:\tVlanId:\t\n")
		for _, pg := range hs.Config.Network.Portgroup {
			fmt.Fprintf(tw, "%s\t", pg.Spec.Name)
			fmt.Fprintf(tw, "%s\t", pg.Spec.VswitchName)
			fmt.Fprintf(tw,"%d\t", pg.Spec.VlanId)
			fmt.Fprintf(tw, "\n")
		}
		tw.Flush()
		fmt.Println()
	}
}