package builtin

import (
	"bytes"
	"github.com/matryer/is"
	"testing"
)

func TestUnmarshalSynonymFile(t *testing.T) {
	is := is.New(t)

	file, err := Get("test.yml")
	is.NoErr(err)

	groups, err := UnmarshalSynonymFile(bytes.NewReader(file))
	is.NoErr(err)
	is.Equal(2, len(groups))
}
