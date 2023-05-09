package libretranslate

import (
	"fmt"
	"net/http"

	"github.com/davidebianchi/go-jsonclient"
)

//
// -- Errors --
//

type NotSupportedError struct {
	LangCode string
}

func (err NotSupportedError) Error() string {
	return fmt.Sprintf("%v not supported", err.LangCode)
}

type NoTargetError struct {
	SourceLangCode string
	TragerLangCode string
}

func (err NoTargetError) Error() string {
	return fmt.Sprintf("%v is not available as a target language from %v ", err.TragerLangCode, err.SourceLangCode)
}

//
// -- Language --
//

type Language struct {
	Code    string
	Name    string
	Targets []string
}

func (lang Language) HasTarget(target string) bool {
	for i := range lang.Targets {
		if lang.Targets[i] == target {
			return true
		}
	}
	return false
}

//
// -- Translation --
//

func New(url string) (*Translation, error) {
	client, err := jsonclient.New(jsonclient.Options{
		BaseURL: url,
	})
	if err != nil {
		return nil, err
	}
	return &Translation{Client: client}, nil
}

type Translation struct {
	Langs  []Language
	Client *jsonclient.Client
}

func (t *Translation) LoadSupportedLanguages() error {
	req, err := t.Client.NewRequest(http.MethodGet, "languages", nil)
	if err != nil {
		return err
	}
	_, err = t.Client.Do(req, &t.Langs)
	return err
}

func (t *Translation) IsSupported(langCode string) bool {
	for i := range t.Langs {
		if t.Langs[i].Code == langCode {
			return true
		}
	}
	return false
}

func (t *Translation) GetLang(code string) Language {
	for i := range t.Langs {
		if t.Langs[i].Code == code {
			return t.Langs[i]
		}
	}
	return Language{}
}

func (t *Translation) Detect(text string) (float32, string, error) {
	type RequestData struct {
		Q string `json:"q"`
	}
	type ResponseData []struct {
		Confidence   float32
		LanguageCode string `json:"language"`
	}

	reqData := RequestData{text}
	req, err := t.Client.NewRequest(http.MethodPost, "detect", reqData)
	if err != nil {
		return 0, "", err
	}

	resData := ResponseData{}
	_, err = t.Client.Do(req, &resData)
	if err != nil {
		return 0, "", err
	}

	if len(resData) != 1 {
		return 0, "", fmt.Errorf("Expected 1 detected language")
	}
	return resData[0].Confidence, resData[0].LanguageCode, nil
}

func (t *Translation) checkTransalteLangs(sourceLangCode, targetLangCode string) error {
	if !t.IsSupported(targetLangCode) {
		return NotSupportedError{targetLangCode}
	}
	sourceLang := t.GetLang(sourceLangCode)
	if sourceLang.Code == "" {
		return NotSupportedError{sourceLangCode}
	}
	if !sourceLang.HasTarget(targetLangCode) {
		return NoTargetError{sourceLangCode, targetLangCode}
	}
	return nil
}

func (t *Translation) Translate(sourceLangCode, targetLangCode, text string) (string, error) {
	type RequestData struct {
		Source string `json:"source"`
		Target string `json:"target"`
		Q      string `json:"q"`
	}
	type ResponseData struct {
		TranslatedText string
	}

	err := t.checkTransalteLangs(sourceLangCode, targetLangCode)
	if err != nil {
		return "", err
	}

	reqData := RequestData{sourceLangCode, targetLangCode, text}
	req, err := t.Client.NewRequest(http.MethodPost, "/translate", reqData)
	if err != nil {
		return "", err
	}

	responseBody := ResponseData{}
	_, err = t.Client.Do(req, &responseBody)
	if err != nil {
		return "", err
	}

	return responseBody.TranslatedText, nil
}
