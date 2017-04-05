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
		ownerContainer *Container // Point to container
}

func (monitor *IOMonitor) flushDirtyData() interface{} {
	monitor.ownContainer.objRWLock.RLock()
	defer monitor.ownContainer.objRWLock.RUnlock()

	for objname, data := range monitor.ownContainer.Objects {
			
	}
	return nil
}

func (monitor *IOMonitor) writeDataToDisk(c *Container, objname *string, data *object) {
	path := fmt.Sprintf("./testData/%s/%s", c.userName, c.containerName)
	io.Save(path, data)
}
