package app

type Data struct {
	Synonyms []SynonymGroup
}

func NewData() *Data {
	return &Data{
		Synonyms: make([]SynonymGroup, 0),
	}
}
