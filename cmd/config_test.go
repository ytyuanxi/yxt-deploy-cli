package cmd

import (
	"fmt"
	"testing"
)

func TestReadConfigTest(t *testing.T) {
	_, err := readConfigTest()
	if err != nil {
		errStr := fmt.Sprintf("%s", err)
		t.Errorf(errStr)
	}

}

func TestReadConfig(t *testing.T) {
	_, err := readConfig()
	if err != nil {
		errStr := fmt.Sprintf("%s", err)
		t.Errorf(errStr)
	}

}
