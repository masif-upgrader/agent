package main

import (
	"github.com/kataras/iris"
	log "github.com/sirupsen/logrus"
	"net"
)

func startRestServer(sock string) (*iris.Application, error) {
	log.WithFields(log.Fields{"socket": sock}).Info("Listening on *nix socket")

	server, errLs := net.Listen("unix", sock)
	if errLs != nil {
		return nil, errLs
	}

	app := iris.New()

	// TODO

	go app.Run(iris.Listener(server), iris.WithoutStartupLog)

	return app, nil
}
