package fileWallet

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"heyuanlong/blockchain-step/crypto"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type StoreI interface {
	// Loads and decrypts the key from disk.
	GetKey(addr crypto.Address,dir string, filename string, auth string) (*Key, error)
	// Writes and encrypts the key.
	StoreKey(filename string, k *Key, auth string) error
	// Joins filename with the key directory unless it is already absolute.
	JoinPath(dir string, filename string) string
}

//-----------------------------------------------------------------------------

type StoreFile struct {
}
func NewStoreFile()*StoreFile{
	return &StoreFile{
	}
}

func (ts *StoreFile) GetKey(addr crypto.Address,dir string, filename string, auth string) (*Key, error){
	path := filepath.Join(dir, filename)
	fd, err := os.Open(path)
	if err != nil {
		return nil,err
	}
	defer fd.Close()

	fi,err :=fd.Stat()
	if err != nil {
		return nil,err
	}

	// Skip any non-key files from the folder
	if ts.nonKeyFile(fi) {
		log.Trace("Ignoring file on account scan", "path", path)
		return nil,errors.New("Ignoring file")
	}

	key := new(Key)
	if err := json.NewDecoder(fd).Decode(key); err != nil {
		return nil,err
	}
	return key,nil
}

func (ts *StoreFile) StoreKey(filename string, key *Key, auth string) error {
	content, err := json.Marshal(key)
	if err != nil {
		return err
	}
	return ts.writeKeyFile(filename, content)
}

func (ts *StoreFile) JoinPath(dir string, filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}

	filename = ts.keyFileName(filename)
	return filepath.Join(dir, filename)
}

//-----------------下面都是文件存储key的辅助函数-------------------------------------------------

func (ts *StoreFile) writeTemporaryKeyFile(file string, content []byte) (string, error) {
	// Create the keystore directory with appropriate permissions
	// in case it is not present yet.
	const dirPerm = 0700
	if err := os.MkdirAll(filepath.Dir(file), dirPerm); err != nil {
		return "", err
	}
	// Atomic write: create a temporary hidden file first
	// then move it into place. TempFile assigns mode 0600.
	f, err := ioutil.TempFile(filepath.Dir(file), "."+filepath.Base(file)+".tmp")
	if err != nil {
		return "", err
	}
	if _, err := f.Write(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		return "", err
	}
	f.Close()
	return f.Name(), nil
}

func (ts *StoreFile) writeKeyFile(file string, content []byte) error {
	name, err := ts.writeTemporaryKeyFile(file, content)
	if err != nil {
		return err
	}
	return os.Rename(name, file)
}

// keyFileName implements the naming convention for keyfiles:
// UTC--<created_at UTC ISO8601>-<address hex>
func (ts *StoreFile) keyFileName(keyAddr string) string {
	t := time.Now().UTC()
	return fmt.Sprintf("UTC--%s--%s", ts.toISO8601(t), keyAddr)
}

func (ts *StoreFile) toISO8601(t time.Time) string {
	var tz string
	name, offset := t.Zone()
	if name == "UTC" {
		tz = "Z"
	} else {
		tz = fmt.Sprintf("%03d00", offset/3600)
	}
	return fmt.Sprintf("%04d-%02d-%02dT%02d-%02d-%02d.%09d%s",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), tz)
}

// nonKeyFile ignores editor backups, hidden files and folders/symlinks.
func (ts *StoreFile) nonKeyFile(fi os.FileInfo) bool {
	// Skip editor backups and UNIX-style hidden files.
	if strings.HasSuffix(fi.Name(), "~") || strings.HasPrefix(fi.Name(), ".") {
		return true
	}
	// Skip misc special files, directories (yes, symlinks too).
	if fi.IsDir() || fi.Mode()&os.ModeType != 0 {
		return true
	}
	return false
}

