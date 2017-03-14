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
		monitor_owner *Container // Point to container
}

func (monitor *IOMonitor) flushDirtyData() interface{} {

		return nil
}
