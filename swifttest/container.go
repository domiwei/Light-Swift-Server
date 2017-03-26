package swifttest

import (
	"fmt"
	"light-swift-server/io"
	"path"
	"sort"
	"strings"
	"time"
	"sync"
	"net/http"
)

// The Key type represents an item stored in an container.
type Key struct {
	Key          string `json:"name"`
	LastModified string `json:"last_modified"`
	Size         int64  `json:"bytes"`
	// ETag gives the hex-encoded MD5 sum of the contents,
	// surrounded with double-quotes.
	ETag        string `json:"hash"`
	ContentType string `json:"content_type"`
	// Owner        Owner
}

type Container struct {
	Metadata
	Name    string
	Ctime   time.Time

	objRWLock *sync.RWMutex // Lock for object map
	Objects map[string]*Object
	ioMonitor *IOMonitor // It may be nil when data is not dirty enough
	ioMonitorRunning bool

	Bytes   int
	DirtyDataBytes int `dirty + clean = Bytes`
	CleanDataBytes int
	dirtyBytesThreshold int // DirtyThreshold
}

func (c Container) list(delimiter string, marker string, prefix string, parent string) (resp []interface{}) {
	var tmp orderedObjects

	// first get all matching objects and arrange them in alphabetical order.
	for _, obj := range c.Objects {
		if strings.HasPrefix(obj.Name, prefix) {
			tmp = append(tmp, obj)
		}
	}
	sort.Sort(tmp)

	var prefixes []string
	for _, obj := range tmp {
		if !strings.HasPrefix(obj.Name, prefix) {
			continue
		}

		isPrefix := false
		name := obj.Name
		if parent != "" {
			if path.Dir(obj.Name) != path.Clean(parent) {
				continue
			}
		} else if delimiter != "" {
			if i := strings.Index(obj.Name[len(prefix):], delimiter); i >= 0 {
				name = obj.Name[:len(prefix)+i+len(delimiter)]
				if prefixes != nil && prefixes[len(prefixes)-1] == name {
					continue
				}
				isPrefix = true
			}
		}

		if name <= marker {
			continue
		}

		if isPrefix {
			prefixes = append(prefixes, name)

			resp = append(resp, Subdir{
				Subdir: name,
			})
		} else {
			resp = append(resp, obj)
		}
	}

	return
}

// validContainerName returns whether name is a valid bucket name.
// Here are the rules, from:
// http://docs.openstack.org/api/openstack-object-storage/1.0/content/ch_object-storage-dev-api-storage.html
//
// Container names cannot exceed 256 bytes and cannot contain the / character.
//
func validContainerName(name string) bool {
	if len(name) == 0 || len(name) > 256 {
		return false
	}
	for _, r := range name {
		switch {
		case r == '/':
			return false
		default:
		}
	}
	return true
}

// Save container to disk
func saveContainerToDisk(userName string, containerName string, c *Container) {
	path := fmt.Sprintf("./testData/%s/%s", userName, containerName)
	io.Save(path, c)
}

// Container constructor
func NewContainer(name string) *Container {
	return &Container{
			Name:    name,
			Objects: make(map[string]*Object),
			Metadata: Metadata{
				Meta: make(http.Header),
			},
			objRWLock: new(sync.RWMutex),
		}
}


// Check if need to flush dirty data
func (c *Container) checkIfNeedIOMonitor() interface{} {
	if c.DirtyDataBytes > c.dirtyBytesThreshold {
		// Check if monitor is running
		if c.ioMonitorRunning == true {
			return nil
		}

		c.objRWLock.Lock()
		defer c.objRWLock.Unlock()

		// Check again before new a monitor
		if c.ioMonitorRunning == true {
			return nil
		}
		c.ioMonitor = &IOMonitor{c.dirtyBytesThreshold, c}
		go c.ioMonitor.flushDirtyData()
	}

	return nil
}
