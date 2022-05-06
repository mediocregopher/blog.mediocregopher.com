package post

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewID(t *testing.T) {

	tests := [][2]string{
		{
			"Why Do We Have WiFi Passwords?",
			"why-do-we-have-wifi-passwords",
		},
		{
			"Ginger: A Small VM Update",
			"ginger-a-small-vm-update",
		},
		{
			"Something-Weird.... woah!",
			"somethingweird-woah",
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, test[1], NewID(test[0]))
		})
	}
}
