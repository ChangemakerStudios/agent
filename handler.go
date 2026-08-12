package main

import (
	"context"
	"net/http"
	"encoding/json"
	"github.com/swarmpit/agent/setup"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	"github.com/gorilla/mux"
	"io"
	"log"
)

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
