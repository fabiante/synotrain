package builtin

import (
	"embed"
	"errors"
	"io"

	"github.com/fabiante/synotrain/app"
	"gopkg.in/yaml.v3"
)

//go:embed files/*
var fs embed.FS

func Get(path string) ([]byte, error) {
	return fs.ReadFile("files/" + path)
}

func UnmarshalSynonymFile(reader io.Reader) ([]app.SynonymGroup, error) {
	decoder := yaml.NewDecoder(reader)

	var groups []app.SynonymGroup

	for {
		var group app.SynonymGroup

		if err := decoder.Decode(&group); err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				return nil, err
			}
		}

		groups = append(groups, group)
	}

	return groups, nil
}
