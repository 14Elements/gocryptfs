//+build linux

package matrix

import (
	"encoding/binary"
	"io/ioutil"
	"os"
	"syscall"
	"testing"

	"github.com/hanwen/go-fuse/fuse"

	"github.com/rfjakob/gocryptfs/tests/test_helpers"
)

// Check that readdir(3) returns valid inode numbers in the directory entries
func TestReaddirInodes(t *testing.T) {
	// Create test file
	path := test_helpers.DefaultPlainDir + "/TestReaddirInodes"
	err := ioutil.WriteFile(path, []byte("foobar"), 0600)
	if err != nil {
		t.Fatal(err)
	}
	// open mountpoint dir
	d, err := os.Open(test_helpers.DefaultPlainDir)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer d.Close()
	buf := make([]byte, 1000)
	// readdir(3) use getdents64(2) internally which returns linux_dirent64
	// structures. We don't have readdir(3) so we call getdents64(2) directly.
	n, err := syscall.Getdents(int(d.Fd()), buf)
	if n == 0 {
		t.Skipf("cannot test on empty directory")
	}
	// The inode number of the first directory entry ("TestReaddirInodes" or
	// any other file left-over from earlier tests) is in the first 8
	// bytes of the buffer
	inode := binary.LittleEndian.Uint64(buf)
	if inode == 0 || inode == fuse.FUSE_UNKNOWN_INO {
		t.Errorf("got invalid inode number: %d = 0x%x", inode, inode)
	}
}
