package main

import (
	"context"
	"fmt"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	//"sync"
	"time"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/soap"
	
	db "github.com/sak0/vmw2os/db"
	httpapi "github.com/sak0/vmw2os/api"
	vmwinfo "github.com/sak0/vmw2os/vmwinfo"
)

var (
	urlFlag = flag.String("url", "administrator@vsphere.local:ZXCVbnm,@172.16.70.19", "url info.")
	period  = flag.Duration("period", 30 * time.Second, "time period for collect vmware info.")
)

const (
	ipaddr   = "172.16.70.19"
	username = "administrator@vsphere.local"
	password = "ZXCVbnm,"
	insecure = true
	
	dbip     = "172.16.0.22"
	dbpass   = "huacloud"
	dbport	 = "3306"
	
	srvport  = 8888
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
		
	//GetHosts(ctx, c)
	//GetHostNetwork(ctx, c)
	
	/*Use mysql data*/	
	mc := db.MysqlConfig{
		Host 	 : dbip,
		Password : dbpass,
		Port     : dbport,
		User     : "root",
		Database : "vmw2os",
	}
	database, err := db.OpenDatabase(mc)
	if err != nil {
		log.Fatal(err)
	}
	database.Begin()
	fmt.Printf("God db: %v\n", database)
	type Cluster struct {
		Id		int		`db:"id"`
		Name	string	`db:"name"`
		VcId	int		`db:"vcenter_id"`
		VMs     int     `db:"vm_nums"`
	}
	var clusters []Cluster
	database.Select("*").From("cluster").Load(&clusters)
	fmt.Printf("From database %v\n", clusters)
	
	
	/* Publisher:  vminfo
	   Subscriber: cmd, srv */ 
	stopC := make(chan string)    
	vminfo := vmwinfo.NewInfoVMware("test", ctx, c, *period, stopC)
	
	var cmd = new(CmdInterface)
	srv := httpapi.NewServer(srvport)
	srv = httpapi.NewServer(9999)
	
	vminfo.AddReceiver(cmd)
	vminfo.AddReceiver(srv)
	vminfo.Run()
	
	//time.Sleep(1 * time.Second)
	//cmd.DisplayRaw()
	go srv.Run()
	
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	for sig := range sc{
		switch sig {
			case syscall.SIGTERM:
				close(stopC)
				time.Sleep(1 * time.Second)
				return
		}
	}
}