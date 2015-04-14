package contnet

import (
	"bufio"
	"encoding/gob"
	"github.com/kardianos/osext"
	"log"
	"os"
	"path/filepath"
)

const (
	__errSnapshotFile = "Failed to take a snapshot because snapshot file could not be created."
	__errSnapshotJson = "Failed to take a snapshot because object failed to serialize."
	__snapshotSaved   = "Snapshot successfully created."

	__errRestoreFile   = "Failed to restore object from snapshot file because file could not be opened."
	__errRestoreJson   = "Failed to restore object from snapshot file because it failed to deserialize."
	__snapshotRestored = "Object successfully restored from snapshot."
)

func __fullPath(basePath, filename string) string {
	if basePath == "" {
		basePath, _ = osext.ExecutableFolder()
	}
	return filepath.Join(basePath, filename+".bkp")
}

func __snapshot(path, filename string, object interface{}) error {
	// create new snapshot file
	fullpath := __fullPath(path, filename)
	file, err := os.Create(fullpath)
	if err != nil {
		log.Print(__errSnapshotFile, err.Error())
		return err
	}
	defer file.Close()

	bufferedWriter := bufio.NewWriter(file)

	err = gob.NewEncoder(bufferedWriter).Encode(object)
	if err != nil {
		log.Print(__errSnapshotJson, err.Error())
		return err
	}

	bufferedWriter.Flush()

	log.Print(__snapshotSaved)
	return nil
}

func __restoreFromSnapshot(path, filename string, object interface{}) (interface{}, error) {
	fullpath := __fullPath(path, filename)
	file, err := os.Open(fullpath)
	if err != nil {
		log.Print(__errRestoreFile, err.Error())
		return nil, err
	}
	defer file.Close()

	bufferedReader := bufio.NewReader(file)

	err = gob.NewDecoder(bufferedReader).Decode(object)
	if err != nil {
		log.Print(__errRestoreJson, err.Error())
		return nil, err
	}

	log.Print(__snapshotRestored)
	return object, nil
}
