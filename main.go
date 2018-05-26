package main

import (
	"context"
	"fmt"
	"flag"
	"log"
	"net/url"
	
	"github.com/vmware/govmomi"
	//"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"
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
		
	GetHosts(ctx, c)
	GetNetworks(ctx, c)
}