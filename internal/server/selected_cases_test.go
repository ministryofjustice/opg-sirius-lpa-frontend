package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectedCases(t *testing.T) {
	SetSelectedCaseUids([]string{"case1", "case2"})
	assert.Equal(t, GetSelectedCaseUids(), []string{"case1", "case2"})
	ClearSelectedCaseUids()
	assert.Equal(t, GetSelectedCaseUids(), []string{})
}
