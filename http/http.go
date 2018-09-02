package http

import (
	"fmt"
	"net/http"
	"github.com/qnib/go-rfxbridge/rfx"
	"github.com/urfave/cli"
	"github.com/urfave/negroni"
	"log"
	"strings"
)

type Server struct {
	ctx *cli.Context
	devs rfx.Devices
}

func NewServer(ctx *cli.Context) Server {
	m := map[string]string{}
	for _, obj := range strings.Split(ctx.GlobalString("dev-map"), ",") {
		if obj == "" {
			continue
		}
		kv := strings.Split(obj, ":")
		m[kv[0]] = kv[1]
	}
	return Server{
		ctx: ctx,
		devs: rfx.NewDevices(ctx.GlobalString("usb"), ctx.Bool("debug"), m),
	}
}

func RunServer(ctx *cli.Context) {
	s := NewServer(ctx)
	s.Run()
}

func (s *Server) GetDev(w http.ResponseWriter, req *http.Request) {
	keys, ok := req.URL.Query()["device"]

	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(401)
		fmt.Fprintf(w, "Url Param 'device' is missing")
		return
	}
	key := keys[0]
	d, err := s.devs.GetKey(key)
	if err != nil {
		// Server error
		w.WriteHeader(501)
		fmt.Fprintf(w, "Device '%s' not found: %s", key, err.Error())
	}
	fmt.Fprintf(w, fmt.Sprintf("%s\n",d))
}

func (s *Server) GetAll(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, s.devs.String())
}

func (s *Server) Run() {
	go s.devs.WatchRFX()
	go s.devs.UpdateData()
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.GetAll)
	mux.HandleFunc("/get", s.GetDev)

	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(mux)
	addr := s.ctx.GlobalString("listen-addr")
	log.Printf("Server http: %s", addr)
	http.ListenAndServe(addr, n)
}