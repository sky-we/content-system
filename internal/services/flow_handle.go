package services

import (
	"fmt"
	goflow "github.com/s8sg/goflow/v1"
)

func (app *CmsApp) StartFlow(fs *goflow.FlowService) {
	go func() {
		if err := fs.Start(); err != nil {
			fmt.Println("go-flow service starting...")
			panic(err)
		}
	}()

}
