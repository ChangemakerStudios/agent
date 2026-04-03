package main

import (
	"context"
	"net/http"
	"encoding/json"
	"github.com/swarmpit/agent/setup"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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
			log.Printf("ERROR: Cannot obtain container logs: %s\n", err)
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
