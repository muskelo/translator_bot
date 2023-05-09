package libretranslate

import "testing"

func TestTranslation(t *testing.T) {
	translation, err := New("http://localhost:5000/")
	if err != nil {
		t.Fatalf("TestTranslation error=%v", err)
	}
	{
		err := translation.LoadSupportedLanguages()
		if err != nil {
			t.Fatalf("Test LoadSupportedLanguages, err=%v", err)
		}
		if len(translation.Langs) == 0 {
			t.Fatal("Test LoadSupportedLanguages, supported languages not loaded")
		}
	}
	{
		code := translation.Langs[0].Code
		v := translation.IsSupported(code)
		if v != true {
			t.Fatal("Test IsSupported, get false, want true")
		}
	}
	{
		code := "XXXXX"
		v := translation.IsSupported(code)
		if v != false {
			t.Fatal("Test IsSupported, get true, want false")
		}
	}
	{
		_, langCode, err := translation.Detect("Привет, что делаешь? Хорошо поел")
		if err != nil {
			t.Fatalf("Test Detect error=%v", err)
		}
		if langCode != "ru" {
			t.Fatalf("Test IsSupported, get %v, want ru", langCode)
		}
	}
	{
		res, err := translation.Translate("ru", "en", "Привет, что делаешь? Хорошо поел")
		if err != nil {
			t.Fatalf("Test Detect error=%v", err)
		}
		want := "Hey, what are you doing? Good ate"
		if res != want {
			t.Fatalf("Test IsSupported, get %v, want ", want)
		}
	}
}
