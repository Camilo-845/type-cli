package words

import (
	"testing"
)

func TestEmbeddedLanguages(t *testing.T) {
	langs := ListLanguages()
	if len(langs) == 0 {
		t.Fatal("expected at least one embedded language")
	}

	found := map[string]bool{}
	for _, l := range langs {
		found[l] = true
	}

	expected := []string{"english", "english_1k", "spanish"}
	for _, e := range expected {
		if !found[e] {
			t.Errorf("expected language %q not found in %v", e, langs)
		}
	}
}

func TestLoadSpanish(t *testing.T) {
	obj, err := LoadLanguage("spanish")
	if err != nil {
		t.Fatalf("LoadLanguage(spanish): %v", err)
	}

	if obj.Name != "spanish" {
		t.Errorf("name = %q, want %q", obj.Name, "spanish")
	}

	if len(obj.Words) == 0 {
		t.Fatal("spanish has no words")
	}

	if !obj.OrderedByFrequency {
		t.Error("spanish should have orderedByFrequency=true")
	}

	if len(obj.AdditionalAccents) == 0 {
		t.Error("spanish should have additionalAccents")
	}
}

func TestLoadEnglish(t *testing.T) {
	obj, err := LoadLanguage("english")
	if err != nil {
		t.Fatalf("LoadLanguage(english): %v", err)
	}

	if len(obj.Words) < 100 {
		t.Errorf("english has %d words, expected >= 100", len(obj.Words))
	}
}

func TestGenerate(t *testing.T) {
	words := Generate(10, "spanish")
	if len(words) != 10 {
		t.Errorf("Generate(10) returned %d words", len(words))
	}

	for _, w := range words {
		if w == "" {
			t.Error("Generate returned empty word")
		}
	}
}

func TestGenerateUnknownFallsBack(t *testing.T) {
	words := Generate(5, "nonexistent")
	if len(words) != 5 {
		t.Errorf("expected 5 words, got %d", len(words))
	}
}

func TestLanguageMeta(t *testing.T) {
	meta, ok := LanguageMeta("spanish")
	if !ok {
		t.Fatal("LanguageMeta(spanish) not found")
	}

	if !meta.OrderedByFrequency {
		t.Error("spanish meta should have orderedByFrequency")
	}
}

func TestNames(t *testing.T) {
	names := Names()
	if len(names) == 0 {
		t.Fatal("Names() returned empty list")
	}
}

func TestLoadSpanish1K(t *testing.T) {
	obj, err := LoadLanguage("spanish_1k")
	if err != nil {
		t.Fatalf("LoadLanguage(spanish_1k): %v", err)
	}

	if obj.Name != "spanish_1k" {
		t.Errorf("name = %q, want %q", obj.Name, "spanish_1k")
	}

	if len(obj.Words) < 500 {
		t.Errorf("spanish_1k has %d words, expected >= 500", len(obj.Words))
	}

	if !obj.OrderedByFrequency {
		t.Error("spanish_1k should have orderedByFrequency=true")
	}
}
