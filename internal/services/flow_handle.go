package services

import (
	goflow "github.com/s8sg/goflow/v1"
)

func (app *CmsApp) StartFlow(fs *goflow.FlowService) {
	go func() {
		if err := fs.Start(); err != nil {
			Logger.Error("go-flow service start error")
			panic(err)
		}
	}()

}
