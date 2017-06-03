package handlers

import "net/http/pprof"

func SetPprofMux() {
	Mux.GetFunc("/debug/pprof", pprof.Index)
	Mux.GetFunc("/debug/pprof/cmdline", pprof.Cmdline)
	Mux.GetFunc("/debug/pprof/symbol", pprof.Symbol)
	Mux.PostFunc("/debug/pprof/symbol", pprof.Symbol)
	Mux.GetFunc("/debug/pprof/profile", pprof.Profile)
	Mux.Get("/debug/heap", pprof.Handler("heap"))
	Mux.Get("/debug/goroutine", pprof.Handler("goroutine"))
	Mux.Get("/debug/block", pprof.Handler("block"))
	Mux.Get("/debug/threadcreate", pprof.Handler("threadcreate"))
}
