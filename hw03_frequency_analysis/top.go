package hw03frequencyanalysis

import (
	"regexp"
	"sort"
)

var re = regexp.MustCompile(`\S+`)

func Top10(text string) []string {
	wordCounts := make(map[string]int)
	// Используем регулярное выражение для разделения слов
	words := re.FindAllString(text, -1)
	for _, word := range words {
		wordCounts[word]++
	}
	type wordFreq struct {
		word  string
		count int
	}
	wordList := make([]wordFreq, 0, len(wordCounts))
	for word, count := range wordCounts {
		wordList = append(wordList, wordFreq{word, count})
	}
	// Сортируем по убыванию частоты
	sort.Slice(wordList, func(i, j int) bool {
		if wordList[i].count == wordList[j].count {
			return wordList[i].word < wordList[j].word
		}
		return wordList[i].count > wordList[j].count
	})
	// Формируем результат из топ-10 слов
	result := make([]string, 0, 10)
	for i := 0; i < len(wordList) && i < 10; i++ {
		result = append(result, wordList[i].word)
	}
	return result
}
