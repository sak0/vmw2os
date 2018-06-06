package httpapi

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/tabwriter"
	
	"github.com/vmware/govmomi/units"
	"github.com/vmware/govmomi/vim25/mo"
	
	vmwinfo "github.com/sak0/vmw2os/vmwinfo"
)

var SingleSrv = new(Server)

type Server struct {
	port	int
	Vchosts []mo.HostSystem
	Vcnets  []mo.Network
	Vcdss   []mo.Datastore
	Vcvms   []mo.VirtualMachine
}

func NewServer(port int)*Server{
	if SingleSrv.port == 0 {
		SingleSrv = &Server{
			port : port,
		}
	}
	return SingleSrv
}

func (s *Server)TestFunc(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Visit test link from %s\n", r.RemoteAddr)
}

func (s *Server)Update(info vmwinfo.Info){
	s.Vchosts = info.Hosts
	s.Vcnets  = info.Nets
	s.Vcdss   = info.Dss
	s.Vcvms   = info.Vms
}

func (s *Server)HostsFunc(w http.ResponseWriter, r *http.Request){
	if s.Vchosts == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "There is no hosts collected yet.")
		return
	}
	
	tw := new(tabwriter.Writer).Init(w, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, "Name:\tUsed CPU:\tTotal CPU:\tFree CPU:\tUsed Memory:\tTotal Memory:\tFree Memory:\t\n")
	for _, hs := range s.Vchosts {
		totalCPU := int64(hs.Summary.Hardware.CpuMhz) * int64(hs.Summary.Hardware.NumCpuCores)
		freeCPU := int64(totalCPU) - int64(hs.Summary.QuickStats.OverallCpuUsage)
		freeMemory := int64(hs.Summary.Hardware.MemorySize) - (int64(hs.Summary.QuickStats.OverallMemoryUsage) * 1024 * 1024)
		fmt.Fprintf(tw, "%s\t", hs.Summary.Config.Name)
		fmt.Fprintf(tw, "%d\t", hs.Summary.QuickStats.OverallCpuUsage)
		fmt.Fprintf(tw, "%d\t", totalCPU)
		fmt.Fprintf(tw, "%d\t", freeCPU)
		fmt.Fprintf(tw, "%d\t", units.ByteSize(hs.Summary.QuickStats.OverallMemoryUsage))
		fmt.Fprintf(tw, "%d\t", units.ByteSize(hs.Summary.Hardware.MemorySize) / 1024 / 1024)
		fmt.Fprintf(tw, "%d\t\n", units.ByteSize(freeMemory) / 1024 / 1024)
	}
	tw.Flush()
}

func (s *Server)PortGroupsFunc(w http.ResponseWriter, r *http.Request){
	if s.Vchosts == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "There is no hosts collected yet.")
		return
	}
	tw := new(tabwriter.Writer).Init(w, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, "Host:\tName:\tVswtich:\tVlanId:\t\n")
	for _, host := range s.Vchosts{
		pgs := host.Config.Network.Portgroup
		for _, pg := range pgs {
			fmt.Fprintf(tw, "%s\t", host.Summary.Config.Name)
			fmt.Fprintf(tw, "%s\t", pg.Spec.Name)
			fmt.Fprintf(tw, "%s\t", pg.Spec.VswitchName)
			fmt.Fprintf(tw,"%d\t", pg.Spec.VlanId)
			fmt.Fprintf(tw, "\n")
		}
	}
	tw.Flush()
}

func (s *Server)DataStoreFunc(w http.ResponseWriter, r *http.Request){
	if s.Vcdss == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "There is no datastores collected yet.\n")
		return
	}
	tw := new(tabwriter.Writer).Init(w, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, "Name\tType\tCapacity\tFree\t\n")
	for _, ds := range s.Vcdss {
		fmt.Fprintf(tw, "%s\t", ds.Summary.Name)
		fmt.Fprintf(tw, "%s\t", ds.Summary.Type)
		fmt.Fprintf(tw, "%d\t", units.ByteSize(ds.Summary.Capacity))
		fmt.Fprintf(tw, "%d\t", units.ByteSize(ds.Summary.FreeSpace))
		fmt.Fprintf(tw, "\n")
	}
	tw.Flush()
}

func (s *Server)VmFunc(w http.ResponseWriter, r *http.Request){
	if s.Vcvms == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "There is no vm collected yet.\n")
		return
	}
	tw := new(tabwriter.Writer).Init(w, 0, 32, 2, ' ', 0)
	fmt.Fprintf(tw, "Name\tGuestFullName\t\n")
	for _, vm := range s.Vcvms {
		fmt.Fprintf(tw, "%s\t", vm.Summary.Config.Name)
		fmt.Fprintf(tw, "%s\t", vm.Summary.Config.GuestFullName)
		fmt.Fprintf(tw, "\n")
	}
	tw.Flush()
}

func (s *Server)Run(){
	mux := http.NewServeMux()
	mux.HandleFunc("/test", http.HandlerFunc(s.TestFunc))
	mux.HandleFunc("/hosts", http.HandlerFunc(s.HostsFunc))
	mux.HandleFunc("/pgs", http.HandlerFunc(s.PortGroupsFunc))
	mux.HandleFunc("/dss", http.HandlerFunc(s.DataStoreFunc))
	mux.HandleFunc("/vms",http.HandlerFunc(s.VmFunc))
	//mux.HandleFunc("/nets", http.HandlerFunc(s.NetsFunc))
	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(s.port), mux))
}