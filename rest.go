package main

import (
	"net"
	"os"

	"github.com/kataras/iris/v12"
	v1 "github.com/masif-upgrader/agent/v1"
	"github.com/masif-upgrader/common"
	log "github.com/sirupsen/logrus"
)

func startRestServer(sock string) (*iris.Application, error) {
	log.WithFields(log.Fields{"socket": sock}).Info("Listening on *nix socket")

	server, errLs := net.Listen("unix", sock)
	if errLs != nil {
		return nil, errLs
	}

	os.Chmod(sock, 0770)

	app := iris.New()

	app.Get("/v1/load", getV1Load)

	go app.Run(iris.Listener(server), iris.WithoutStartupLog)

	return app, nil
}

func getV1Load(ctx iris.Context) {
	ctx.JSON(&v1.Load{
		queryStats.queryLoad(),
		actionsStats[common.PkgMgrInstall].queryLoad(),
		actionsStats[common.PkgMgrUpdate].queryLoad(),
		actionsStats[common.PkgMgrConfigure].queryLoad(),
		actionsStats[common.PkgMgrRemove].queryLoad(),
		actionsStats[common.PkgMgrPurge].queryLoad(),
		errorStats.queryLoad(),
	})
}
