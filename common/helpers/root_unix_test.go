//go:build !windows

package helpers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunningAsRoot(t *testing.T) {
	assert.Equal(t, os.Geteuid() == 0, RunningAsRoot())
}
