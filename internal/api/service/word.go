package service

import (
	"bufio"
	"errors"
	"io"
	"log"
	"math/rand"
	"russian-learns-english/internal/api/component/yandex"
	"strings"
	"sync"
)

type WordService interface {
	GetWordList() (wl []Word)
	GetRandomWord() Word
	LoadWordList(reader io.Reader) (int, error)
	CheckWordTranslation(ID WordID, translations []string) (bool, error)
}

type WordID = int64

type Word struct {
	ID           WordID   `json:"id"`
	Original     string   `json:"original"`
	Translations []string `json:"translations,omitempty"`
}

func (w Word) TranslationExist(translation string) bool {
	for _, tr := range w.Translations {
		if tr == translation {
			return true
		}
	}

	return false
}

type wordService struct {
	wordList   map[WordID]Word
	wordListMu *sync.Mutex
	yd         *yandex.DictAPIClient
}

func (w *wordService) GetWordList() (wl []Word) {
	w.wordListMu.Lock()
	defer w.wordListMu.Unlock()

	for _, word := range w.wordList {
		wl = append(wl, word)
	}

	return wl
}

func (w *wordService) CheckWordTranslation(ID WordID, translations []string) (bool, error) {
	w.wordListMu.Lock()
	word, ok := w.wordList[ID]
	w.wordListMu.Unlock()
	if !ok {
		return false, errors.New("word not found")
	}

	for _, v := range translations {
		if word.TranslationExist(v) {
			return true, nil
		}
	}

	return false, nil
}

func (w *wordService) GetRandomWord() Word {
	IDs := make([]WordID, 0)
	for ID, _ := range w.wordList {
		IDs = append(IDs, ID)
	}

	w.wordListMu.Lock()
	defer w.wordListMu.Unlock()
	word := w.wordList[IDs[WordID(rand.Int()%len(IDs))]]
	word.Translations = nil
	return word
}

func (w *wordService) LoadWordList(reader io.Reader) (int, error) {
	scanner := bufio.NewScanner(reader)

	wg := &sync.WaitGroup{}
	for scanner.Scan() {
		orig := strings.Trim(scanner.Text(), " ")
		wg.Add(1)
		go func() {
			defer wg.Done()

			translations, err := w.yd.TranslateWord(orig)
			if err != nil {
				log.Printf("Failed to et translations for %s: %v", orig, err)
			}

			wID := w.newWordID()
			w.wordListMu.Lock()
			defer w.wordListMu.Unlock()
			w.wordList[wID] = Word{
				ID:           wID,
				Original:     orig,
				Translations: translations,
			}
		}()
	}
	wg.Wait()

	totalLoaded := len(w.wordList)
	if totalLoaded == 0 {
		return 0, errors.New("no words")
	}

	return totalLoaded, nil
}

func (w *wordService) newWordID() WordID {
	w.wordListMu.Lock()
	defer w.wordListMu.Unlock()
	for {
		ID := WordID(rand.Int())
		if _, ok := w.wordList[ID]; ok {
			continue
		}
		return ID
	}
}

func NewWordService(yandexDict *yandex.DictAPIClient) WordService {
	return &wordService{
		yd:         yandexDict,
		wordList:   make(map[WordID]Word),
		wordListMu: &sync.Mutex{},
	}
}
