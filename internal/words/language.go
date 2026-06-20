package words

import (
	"embed"
	"encoding/json"
	"log"
)

//go:embed languages/*.json
var languageFS embed.FS

type LanguageObject struct {
	Name               string      `json:"name"`
	Words              []string    `json:"words"`
	RightToLeft        bool        `json:"rightToLeft,omitempty"`
	NoLazyMode         bool        `json:"noLazyMode,omitempty"`
	JoiningScript      bool        `json:"joiningScript,omitempty"`
	OrderedByFrequency bool        `json:"orderedByFrequency,omitempty"`
	AdditionalAccents  [][2]string `json:"additionalAccents,omitempty"`
	OriginalPunctuation bool       `json:"originalPunctuation,omitempty"`
	BCP47              string      `json:"bcp47,omitempty"`
}

var languageMeta = map[string]LanguageObject{}

func LoadLanguage(lang string) (*LanguageObject, error) {
	data, err := languageFS.ReadFile("languages/" + lang + ".json")
	if err != nil {
		return nil, err
	}

	var obj LanguageObject
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	return &obj, nil
}

func ListLanguages() []string {
	entries, err := languageFS.ReadDir("languages")
	if err != nil {
		return nil
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if len(name) > 5 && name[len(name)-5:] == ".json" {
			names = append(names, name[:len(name)-5])
		}
	}
	return names
}

func LanguageMeta(name string) (LanguageObject, bool) {
	meta, ok := languageMeta[name]
	return meta, ok
}

func init() {
	entries, err := languageFS.ReadDir("languages")
	if err != nil {
		log.Printf("failed to read embedded languages: %v", err)
		return
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if len(name) <= 5 || name[len(name)-5:] != ".json" {
			continue
		}

		langName := name[:len(name)-5]
		obj, err := LoadLanguage(langName)
		if err != nil {
			log.Printf("failed to load language %s: %v", langName, err)
			continue
		}

		if len(obj.Words) == 0 {
			log.Printf("language %s has no words, skipping", langName)
			continue
		}

		Register(obj.Name, obj.Words)
		languageMeta[obj.Name] = *obj
	}
}
