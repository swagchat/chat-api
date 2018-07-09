package handler

import "net/http/pprof"

func setPprofMux() {
	mux.GetFunc("/debug/pprof", pprof.Index)
	mux.GetFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.GetFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.PostFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.GetFunc("/debug/pprof/profile", pprof.Profile)
	mux.Get("/debug/heap", pprof.Handler("heap"))
	mux.Get("/debug/goroutine", pprof.Handler("goroutine"))
	mux.Get("/debug/block", pprof.Handler("block"))
	mux.Get("/debug/threadcreate", pprof.Handler("threadcreate"))
}
