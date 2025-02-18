package builtin

import (
	"bytes"
	"github.com/fabiante/synotrain/app"
)

// Data builds a new Data object based on the builtin data defined in the "builtin/files" directory.
func Data() *app.Data {
	data := app.NewData()

	file, err := Get("de.yml")
	if err != nil {
		panic(err)
	}

	groups, err := UnmarshalSynonymFile(bytes.NewReader(file))
	if err != nil {
		panic(err)
	}

	data.Synonyms = groups

	return data
}
