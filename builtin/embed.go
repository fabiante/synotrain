package builtin

import (
	"embed"
	"errors"
	"github.com/fabiante/synotrain/app"
	"gopkg.in/yaml.v3"
	"io"
)

//go:embed files/*
var fs embed.FS

func Get(path string) ([]byte, error) {
	return fs.ReadFile("files/" + path)
}

type synonymResource struct {
	Desc     string
	Synonyms app.SynonymGroup
}

func UnmarshalSynonymFile(reader io.Reader) ([]app.SynonymGroup, error) {
	decoder := yaml.NewDecoder(reader)

	var groups []app.SynonymGroup

	for {
		var res synonymResource

		if err := decoder.Decode(&res); err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				return nil, err
			}
		}

		groups = append(groups, res.Synonyms)
	}

	return groups, nil
}
