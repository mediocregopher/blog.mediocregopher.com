package chat

import (
	"strconv"
	"strings"

	"github.com/mediocregopher/radix/v4"
)

func parseStreamEntryID(str string) (radix.StreamEntryID, error) {

	split := strings.SplitN(str, "-", 2)
	if len(split) != 2 {
		return radix.StreamEntryID{}, errInvalidMessageID
	}

	time, err := strconv.ParseUint(split[0], 10, 64)
	if err != nil {
		return radix.StreamEntryID{}, errInvalidMessageID
	}

	seq, err := strconv.ParseUint(split[1], 10, 64)
	if err != nil {
		return radix.StreamEntryID{}, errInvalidMessageID
	}

	return radix.StreamEntryID{Time: time, Seq: seq}, nil
}
