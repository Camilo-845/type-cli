package words

import (
	"testing"
)

func TestEmbeddedLanguages(t *testing.T) {
	langs := ListLanguages()
	if len(langs) < 10 {
		t.Fatalf("expected at least 10 embedded languages, got %d", len(langs))
	}

	found := map[string]bool{}
	for _, l := range langs {
		found[l] = true
	}

	expected := []string{"english", "english_1k", "spanish", "french", "german", "code_python"}
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
}

func TestLoadEnglish(t *testing.T) {
	obj, err := LoadLanguage("english")
	if err != nil {
		t.Fatalf("LoadLanguage(english): %v", err)
	}

	if len(obj.Words) < 100 {
		t.Errorf("english has %d words, expected >= 100", len(obj.Words))
	}

	if !obj.OrderedByFrequency {
		t.Error("english should have orderedByFrequency=true")
	}
}

func TestLoadGerman(t *testing.T) {
	obj, err := LoadLanguage("german")
	if err != nil {
		t.Fatalf("LoadLanguage(german): %v", err)
	}

	if len(obj.AdditionalAccents) == 0 {
		t.Error("german should have additionalAccents")
	}
}

func TestGenerate(t *testing.T) {
	cfg := GeneratorConfig{}
	words := Generate(10, "spanish", cfg)
	if len(words) != 10 {
		t.Errorf("Generate(10) returned %d words", len(words))
	}

	for _, w := range words {
		if w == "" {
			t.Error("Generate returned empty word")
		}
	}
}

func TestGenerateWithPunctuation(t *testing.T) {
	cfg := GeneratorConfig{Punctuation: true}
	words := Generate(50, "english", cfg)
	if len(words) != 50 {
		t.Errorf("Generate(50) returned %d words", len(words))
	}
	for _, w := range words {
		if w == "" {
			t.Error("Generate returned empty word")
		}
	}
}

func TestGenerateWithLazyMode(t *testing.T) {
	cfg := GeneratorConfig{LazyMode: true}
	words := Generate(10, "spanish", cfg)
	if len(words) != 10 {
		t.Errorf("Generate(10,lazy) returned %d words", len(words))
	}
}

func TestGenerateUnknownFallsBack(t *testing.T) {
	cfg := GeneratorConfig{}
	words := Generate(5, "nonexistent", cfg)
	if len(words) != 5 {
		t.Errorf("expected 5 words, got %d", len(words))
	}
}

func TestLanguageMeta(t *testing.T) {
	meta, ok := LanguageMeta("english")
	if !ok {
		t.Fatal("LanguageMeta(english) not found")
	}

	if !meta.OrderedByFrequency {
		t.Error("english meta should have orderedByFrequency")
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
}

func TestZipfDistribution(t *testing.T) {
	// Smoke test: should never panic or return out of bounds
	for i := 0; i < 1000; i++ {
		idx := zipfRandomIndex(200)
		if idx < 0 || idx >= 200 {
			t.Errorf("zipfRandomIndex(200) = %d, out of bounds", idx)
		}
	}
}

func TestReplaceAccents(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"café", "cafe"},
		{"tú", "tu"},
		{"México", "Mexico"},
		{"über", "uber"},
		{"straße", "strasse"},
		{"façade", "facade"},
	}

	for _, tt := range tests {
		result := ReplaceAccents(tt.input, nil)
		if result != tt.expected {
			t.Errorf("ReplaceAccents(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestWordsetRandomWord(t *testing.T) {
	ws := NewWordset([]string{"the", "be", "of", "and", "a", "to", "in"})
	for i := 0; i < 100; i++ {
		word := ws.RandomWord(false)
		found := false
		for _, w := range ws.Words {
			if w == word {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("RandomWord returned unknown word: %q", word)
		}
	}
}

func TestWordsetNextWord(t *testing.T) {
	ws := NewWordset([]string{"a", "b", "c"})

	w1 := ws.NextWord()
	w2 := ws.NextWord()
	w3 := ws.NextWord()

	if w1 != "a" || w2 != "b" || w3 != "c" {
		t.Errorf("NextWord sequence = %q,%q,%q, want a,b,c", w1, w2, w3)
	}

	// Should wrap around
	w4 := ws.NextWord()
	if w4 != "a" {
		t.Errorf("Wrap NextWord = %q, want a", w4)
	}
}

func TestWordsetShuffledWord(t *testing.T) {
	ws := NewWordset([]string{"a", "b", "c", "d", "e"})
	seen := make(map[string]bool)
	for i := 0; i < 5; i++ {
		seen[ws.ShuffledWord()] = true
	}
	if len(seen) != 5 {
		t.Errorf("ShuffledWord should return all words, got %d unique", len(seen))
	}
}

func TestGeneratorNoRepeat(t *testing.T) {
	lang := &LanguageObject{
		Name:              "test",
		Words:             []string{"the", "be", "of", "and", "a"},
		OrderedByFrequency: false,
	}
	cfg := GeneratorConfig{}
	gen := NewGenerator(lang, cfg)
	words := gen.GenerateWords(100)

	// Check no sequential repeats
	for i := 1; i < len(words); i++ {
		if words[i] == words[i-1] {
			t.Errorf("sequential repeat at %d: %q", i, words[i])
		}
	}
}

func TestGeneratorPunctuation(t *testing.T) {
	lang := &LanguageObject{
		Name:              "english",
		Words:             []string{"the", "be", "of", "and", "a", "to", "in", "he", "have", "it"},
		OrderedByFrequency: true,
	}
	cfg := GeneratorConfig{Punctuation: true}
	gen := NewGenerator(lang, cfg)
	words := gen.GenerateWords(100)

	if len(words) != 100 {
		t.Errorf("GenerateWords(100) returned %d", len(words))
	}

	for _, w := range words {
		if w == "" {
			t.Error("GenerateWords returned empty word")
		}
	}
}
