package main

import (
	"context"
	"net/http"
	"encoding/json"
	"github.com/swarmpit/agent/setup"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	"github.com/gorilla/mux"
	"io"
	"log"
	"time"
)

// A wedged volume plugin hangs the daemon's volume list; give up so one bad
// node can't stall swarmpit's cluster-wide listing.
const volumesTimeout = 10 * time.Second

// Volumes returns this node's volumes as the engine's /volumes does. Docker
// only lists the daemon it is asked, so swarmpit asks every node's agent.
func Volumes(cli *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), volumesTimeout)
		defer cancel()
		resp, err := cli.VolumeList(ctx, volume.ListOptions{})
		if err != nil {
			log.Printf("ERROR: Cannot list volumes: %s\n", err)
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func Info(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(setup.GetArgs())
}

func Logs(cli *client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var options = container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Timestamps: true,
			Details:    true,
		}

		params := mux.Vars(r)
		container := params["container"]

		query := r.URL.Query()
		since := query.Get("since")
		if since != "" {
			options.Since = since
		}

		resp, err := cli.ContainerLogs(context.Background(), container, options)

		if err != nil {
			// A container that has gone away is routine: swarmpit keeps tailing
			// a task's logs for a while after the task is replaced. Report it
			// without the ERROR level, which otherwise floods the log every
			// couple of seconds for as long as the tail is open.
			if errdefs.IsNotFound(err) {
				if routerArg.Debug.Http {
					log.Printf("DEBUG: Container gone, no logs to serve: %s\n", err)
				}
			} else {
				log.Printf("ERROR: Cannot obtain container logs: %s\n", err)
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer resp.Close()
		content, err := io.ReadAll(resp)
		if err != nil {
			log.Printf("ERROR: Cannot read container logs: %s\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(string(content))
	}
}
