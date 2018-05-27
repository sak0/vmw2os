package httpapi

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/tabwriter"
	
	"github.com/vmware/govmomi/units"
	"github.com/vmware/govmomi/vim25/mo"
)

type Server struct {
	port	int
	Vchosts []mo.HostSystem
	Vcnets  []mo.Network
}

func NewServer(port int)*Server{
	return &Server{
		port: port,
	}
}

func (s *Server)TestFunc(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Visit test link from %s\n", r.RemoteAddr)
}

func (s *Server)Update(hss []mo.HostSystem, nets []mo.Network){
	s.Vchosts = hss
	s.Vcnets  = nets
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

func (s *Server)Run(){
	mux := http.NewServeMux()
	mux.HandleFunc("/test", http.HandlerFunc(s.TestFunc))
	mux.HandleFunc("/hosts", http.HandlerFunc(s.HostsFunc))
	//mux.HandleFunc("/nets", http.HandlerFunc(s.NetsFunc))
	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(s.port), mux))
}