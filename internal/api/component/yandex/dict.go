package yandex

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	dictAPIURL = "https://dictionary.yandex.net/api/v1/dicservice.json"
	dictLang   = "en-ru"
)

type DictAPIClient struct {
	ApiKey string `yaml:"yandex_dict_api_key"`
}

type dictTranslateWordResult struct {
	Def []struct {
		Tr []struct {
			Text string `json:"text"`
			Syn  []struct {
				Text string `json:"text"`
			} `json:"syn"`
		} `json:"tr"`
	} `json:"def"`
}

func (dc *DictAPIClient) TranslateWord(word string) (translations []string, err error) {
	req, err := http.NewRequest(http.MethodGet, dictAPIURL+"/lookup", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("key", dc.ApiKey)
	q.Add("lang", dictLang)
	q.Add("text", word)
	req.URL.RawQuery = q.Encode()

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	} else if res.StatusCode != http.StatusOK {
		return nil, errors.New("invalid request")
	}

	var result dictTranslateWordResult
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	for _, def := range result.Def {
		for _, tr := range def.Tr {
			translations = append(translations, tr.Text)
			for _, syn := range tr.Syn {
				translations = append(translations, syn.Text)
			}
		}
	}

	return translations, nil
}
