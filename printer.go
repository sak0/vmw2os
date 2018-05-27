package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	
	"github.com/vmware/govmomi/units"
	"github.com/vmware/govmomi/vim25/mo"
)

type CmdInterface struct{
	Hosts []mo.HostSystem
	Nets  []mo.Network
}

func (cmd *CmdInterface)Update(hss []mo.HostSystem, nets []mo.Network){
	cmd.Hosts = hss
	cmd.Nets  = nets
}

func (cmd *CmdInterface)Display(){
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, "Name:\tUsed CPU:\tTotal CPU:\tFree CPU:\tUsed Memory:\tTotal Memory:\tFree Memory:\t\n")
	for _, hs := range cmd.Hosts {
		totalCPU := int64(hs.Summary.Hardware.CpuMhz) * int64(hs.Summary.Hardware.NumCpuCores)
		freeCPU := int64(totalCPU) - int64(hs.Summary.QuickStats.OverallCpuUsage)
		freeMemory := int64(hs.Summary.Hardware.MemorySize) - (int64(hs.Summary.QuickStats.OverallMemoryUsage) * 1024 * 1024)
		fmt.Fprintf(tw, "%s\t", hs.Summary.Config.Name)
		fmt.Fprintf(tw, "%d\t", hs.Summary.QuickStats.OverallCpuUsage)
		fmt.Fprintf(tw, "%d\t", totalCPU)
		fmt.Fprintf(tw, "%d\t", freeCPU)
		fmt.Fprintf(tw, "%d\t", units.ByteSize(hs.Summary.QuickStats.OverallMemoryUsage))
		fmt.Fprintf(tw, "%d\t", units.ByteSize(hs.Summary.Hardware.MemorySize) / 1024 / 1024)
		fmt.Fprintf(tw, "%d\t", units.ByteSize(freeMemory) / 1024 / 1024)
		fmt.Fprintf(tw, "\n")
		if hs.Config != nil && hs.Config.Network != nil {
			for _, pg := range hs.Config.Network.Portgroup {
				fmt.Fprintf(tw, "%s - %s\n", pg.Key, pg.Vswitch)
			}
		}
	}
	_ = tw.Flush()
}