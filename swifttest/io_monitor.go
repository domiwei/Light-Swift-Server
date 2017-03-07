package swifttest

import (
//	"bytes"
//	"encoding/json"
//	"io"
//	"os"
//	"sync"
)

type IOMonitor struct {
		dirtyBytesThreshold int
}

func (monitor IOMonitor) checkWrite(objr *objectResource) interface{} {
		if objr.container.DirtyDataBytes > monitor.dirtyBytesThreshold {
			// TODO: begin to write data
		}

		return nil
}
