package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/crc-org/admin-helper/pkg/constants"
	"github.com/crc-org/admin-helper/pkg/hosts"
	"github.com/crc-org/admin-helper/pkg/logging"
	"github.com/crc-org/admin-helper/pkg/types"
)

func Mux(hosts *hosts.Hosts) http.Handler {
	mux := http.NewServeMux()
	logger := logging.GetLogger()
	mux.HandleFunc("/version", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, constants.Version)
	})
	mux.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "post only", http.StatusBadRequest)
			return
		}
		var req types.AddRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err := hosts.Add(req.IP, req.Hosts)
		logger.LogModification(logging.Modification{
			Operation: "add",
			IP:        req.IP,
			Hosts:     req.Hosts,
			Caller:    r.RemoteAddr,
			Error:     err,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/remove", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "post only", http.StatusBadRequest)
			return
		}
		var req types.RemoveRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err := hosts.Remove(req.Hosts)
		logger.LogModification(logging.Modification{
			Operation: "remove",
			Hosts:     req.Hosts,
			Caller:    r.RemoteAddr,
			Error:     err,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/clean", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "post only", http.StatusBadRequest)
			return
		}
		var req types.CleanRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err := hosts.Clean()
		logger.LogModification(logging.Modification{
			Operation: "clean",
			Caller:    r.RemoteAddr,
			Error:     err,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	return mux
}
