package toolkit

import (
	"os"
	"testing"
)

func TestTools_CreateDir(t *testing.T) {
	var testTools Tools

	// Test case 1: Directory does not exist
	dir := "./testdata/newdir"
	err := testTools.CreateDir(dir)
	if err != nil {
		t.Errorf("CreateDir() returned an error: %v", err)
	}

	// Verify that the directory was created
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("CreateDir() failed to create the directory: %v", err)
	}

	// Test case 2: Directory already exists
	err = testTools.CreateDir(dir)
	if err != nil {
		t.Errorf("CreateDir() returned an error: %v", err)
	}

	// Verify that the directory still exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("CreateDir() deleted the existing directory: %v", err)
	}

	// Clean up: Remove the directory
	err = os.Remove(dir)
	if err != nil {
		t.Errorf("Failed to remove the directory: %v", err)
	}
}
