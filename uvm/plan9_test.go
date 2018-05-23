// +build functional,p9

// To run: go test -v -tags "functional p9"

package uvm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Microsoft/hcsshim/internal/schema2"
)

// TestPlan9 tests adding/removing Plan9 shares to/from a v2 Linux utility VM
// TODO: This is a very basic test currently. Need multiple shares and so-on.
//       Can be iterated on later.
func TestPlan9(t *testing.T) {
	uvmID := "TestPlan9"
	uvm := createLCOWUVM(t, uvmID)
	defer uvm.Terminate()

	dir := createTempDir(t)
	dirU := strings.ToUpper(dir)
	defer os.RemoveAll(dir)
	var iterations uint32 = 64
	for i := 0; i < int(iterations); i++ {
		if err := uvm.AddPlan9(dirU, fmt.Sprintf("/tmp/%s", filepath.Base(dir)), schema2.VPlan9FlagNone); err != nil {
			t.Fatalf("AddPlan9 failed: %s", err)
		}
	}
	if len(uvm.plan9Shares) != 1 {
		t.Fatalf("Should only be one Plan9 entry")
	}
	if _, ok := uvm.plan9Shares[dir]; ok {
		t.Fatalf("should not found as upper case")
	}
	if _, ok := uvm.plan9Shares[strings.ToLower(dir)]; !ok {
		t.Fatalf("not found!")
	}
	if uvm.plan9Shares[strings.ToLower(dir)].refCount != iterations {
		t.Fatalf("iteration mismatch: %d %d", iterations, uvm.plan9Shares[strings.ToLower(dir)].refCount)
	}

	// Remove them all
	for i := 0; i < int(iterations); i++ {
		if err := uvm.RemovePlan9(dir); err != nil {
			t.Fatalf("RemovePlan9 failed: %s", err)
		}
	}
	if len(uvm.plan9Shares) != 0 {
		t.Fatalf("Should not be any plan9 entries remaining")
	}

}
