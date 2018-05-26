package httpapi

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	port	int
}

func NewServer(port int)*Server{
	return &Server{
		port: port,
	}
}

func (s *Server)TestFunc(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Visit test link from %s\n", r.RemoteAddr)
}

func (s *Server)Run(){
	mux := http.NewServeMux()
	mux.HandleFunc("/test", http.HandlerFunc(s.TestFunc))
	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(s.port), mux))
}