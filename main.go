package main

import (
	"context"
	"fmt"
	"flag"
	"log"
	"net/url"
	"os"
	"text/tabwriter"
	
	"github.com/vmware/govmomi"
	//"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/units"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

var urlFlag = flag.String("url", "administrator@vsphere.local:ZXCVbnm,@172.16.70.19", "url info")

const (
	ipaddr   = "172.16.70.19"
	username = "administrator@vsphere.local"
	password = "ZXCVbnm,"
	insecure = true
)

func NewClient(ctx context.Context, u *url.URL) (*govmomi.Client, error) {
	u.User = url.UserPassword(username, password)

	// Connect and log in to ESX or vCenter
	return govmomi.NewClient(ctx, u, insecure)
}


func main(){
	ctx := context.Background()
	flag.Parse()
	urlstr := username + ":" + password + "@" + ipaddr
	
	u, err := soap.ParseURL(urlstr)
	if err != nil {
		fmt.Printf("ParseURL failed %v", err)
		return
	} else {
		fmt.Printf("Connecting to %v\n", u)
	}
	
	c, err := NewClient(ctx, u)
	defer c.Logout(ctx)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Got client: %v\n", c)
	}
	
	m := view.NewManager(c.Client)
	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"HostSystem"}, true)
	if err != nil {
		log.Fatal(err)
	}
	
	defer v.Destroy(ctx)
	
	var hss []mo.HostSystem
	err = v.Retrieve(ctx, []string{"HostSystem"}, []string{"summary"}, &hss)
	if err != nil {
		log.Fatal(err)
	}
	
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, "Name:\tUsed CPU:\tTotal CPU:\tFree CPU:\tUsed Memory:\tTotal Memory:\tFree Memory\t:\n")
	for _, hs := range hss {
		totalCPU := int64(hs.Summary.Hardware.CpuMhz) * int64(hs.Summary.Hardware.NumCpuCores)
		freeCPU := int64(totalCPU) - int64(hs.Summary.QuickStats.OverallCpuUsage)
		freeMemory := int64(hs.Summary.Hardware.MemorySize) - (int64(hs.Summary.QuickStats.OverallMemoryUsage) * 1024 * 1024)
		fmt.Fprintf(tw, "%s\t", hs.Summary.Config.Name)
		fmt.Fprintf(tw, "%d\t", hs.Summary.QuickStats.OverallCpuUsage)
		fmt.Fprintf(tw, "%d\t", totalCPU)
		fmt.Fprintf(tw, "%d\t", freeCPU)
		fmt.Fprintf(tw, "%d\t", units.ByteSize(hs.Summary.QuickStats.OverallMemoryUsage))
		fmt.Fprintf(tw, "%d\t", units.ByteSize(hs.Summary.Hardware.MemorySize))
		fmt.Fprintf(tw, "%d\t\n", freeMemory)
	}
	_ = tw.Flush()
}