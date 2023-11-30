package toolkit

import (
	"errors"
	"testing"
)

func TestTools_CreateSlug(t *testing.T) {
	var testTools Tools

	// Test case 1: Non-empty string
	input1 := "Hello World"
	expectedOutput1 := "hello-world"
	output1, err := testTools.CreateSlug(input1)
	if err != nil {
		t.Errorf("CreateSlug(%q) returned an unexpected error: %v", input1, err)
	}
	if output1 != expectedOutput1 {
		t.Errorf("CreateSlug(%q) = %q, want %q", input1, output1, expectedOutput1)
	}

	// Test case 2: Empty string
	input2 := ""
	expectedError2 := errors.New("Empty string not Allowed")
	_, err2 := testTools.CreateSlug(input2)
	if err2 == nil {
		t.Errorf("CreateSlug(%q) did not return an error, expected: %v", input2, expectedError2)
	} else if err2.Error() != expectedError2.Error() {
		t.Errorf("CreateSlug(%q) returned an unexpected error: %v, expected: %v", input2, err2, expectedError2)
	}

	// Test case 3: String with special characters
	input3 := "Hello@World!"
	expectedOutput3 := "hello-world"
	output3, err3 := testTools.CreateSlug(input3)
	if err3 != nil {
		t.Errorf("CreateSlug(%q) returned an unexpected error: %v", input3, err3)
	}
	if output3 != expectedOutput3 {
		t.Errorf("CreateSlug(%q) = %q, want %q", input3, output3, expectedOutput3)
	}
}
