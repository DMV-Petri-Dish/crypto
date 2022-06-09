// Package disk implements the ability to read and write blocks to disk
// writing each block to a separate block numbered file.
package disk

import "os"

// Disk represents the serialization implementation for reading and storing
// blocks in their own separate files on disk. This implements the database.Storage interface
type Disk struct {
	dbPath string
}

// New constructs a Disk value for use.
func New(dbPath string) (*Disk, error) {
	if err := os.MkdirAll(dbPath, 0755); err != nil {
		return nil, err
	}

	return &Disk{dbPath: dbPath}, nil
}
