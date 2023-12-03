package main

// import (
//     "sync"
//     "sort"
// )
import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"

	"time"

	"github.com/kljensen/snowball/english"
)

type Document struct {
	Title       string   `json:"title"`
	Description []string `json:"description"`
	Answer      []string `json:"answer"`
}

type TermStats struct {
	TF  float64
	IDF float64
}

var docs []Document

type InvertedIndex map[string][]int
type TFIDFIndex map[string]map[int]TermStats
type IDF map[string][2]float64
type TermFrequency map[int]map[string]int

func loadJSON(filePath string) (Document, error) {
	var doc Document
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return doc, err
	}
	err = json.Unmarshal(fileContent, &doc)
	if err != nil {
		return doc, err
	}
	return doc, nil
}

func tokenizeAndStem(text string) []string {
	words := strings.Fields(text)

	for i, word := range words {
		if strings.ContainsRune(word, '?') {
			word = word[:len(word)-1]
		}
		stemmed := english.Stem(word, false)
		words[i] = stemmed
	}

	return words
}

func processQuery(query string) []string {
	tokens := tokenizeAndStem(query)
	return tokens
}

func buildIDF(docs []Document, queryTokens []string) (IDF, InvertedIndex, TermFrequency) {
	idf := make(IDF)
	index := make(InvertedIndex)
	termFrequency := make(TermFrequency)
	for docID, doc := range docs {
		tokens := tokenizeAndStem(doc.Title)
		for _, token := range tokens {
			termFreq, exists := termFrequency[docID]
			if !exists {
				termFrequency[docID] = termFreq
				termFrequency[docID] = make(map[string]int)
			}
			termFrequency[docID][token]++
			val, exists := idf[token]
			if !exists {
				idf[token] = val
				idf[token] = [2]float64{1, 0}
			} else {
				count := idf[token][0]
				idf[token] = [2]float64{count + 1, 0}
			}
			docArr := index[token]
			docArr = append(docArr, docID)
			index[token] = docArr
		}
	}
	for _, token := range queryTokens {
		val, exists := idf[token]
		if !exists {
			idf[token] = val
			idf[token] = [2]float64{1, 0}
		} else {
			count := idf[token][0]
			idf[token] = [2]float64{count + 1, 0}
		}
	}
	for key, _ := range idf {
		df := idf[key][0]
		idfValue := math.Log2(float64(len(docs)) / df)
		idf[key] = [2]float64{df, idfValue}
	}
	return idf, index, termFrequency
}

func getDocIdsForQueryTokens(queryTokens []string, invertedIndex InvertedIndex) InvertedIndex {
	queryDocsIndex := make(InvertedIndex)
	for _, token := range queryTokens {
		if _, exists := invertedIndex[token]; exists {
			queryDocsIndex[token] = invertedIndex[token]
		}
	}
	return queryDocsIndex
}

func main() {
	startTime := time.Now()
	userQuery := "pass-by-reference or pass-by-value java How do I avoid checking for nulls in Java?"
	queryTokens := processQuery(userQuery)
	_ = queryTokens

	dir := "../test_snippets"

	fileList, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory: ", err)
		return
	}

	for _, file := range fileList {
		filePath := fmt.Sprintf("../test_snippets/%s", file.Name())
		doc, err := loadJSON(filePath)
		if err != nil {
			fmt.Println("Error loading JSON:", err)
			return
		}
		docs = append(docs, doc)
	}

	idf, invertedIndex, termFrequency := buildIDF(docs, queryTokens)
	fmt.Println(len(idf))
	fmt.Println(len(invertedIndex))

	// for key, value := range idf {
	// 	if _, exists := invertedIndex[key]; !exists {
	// 		fmt.Println(key, value)
	// 	}
	// }
	// for key, value := range idf {
	// 	fmt.Println(key, value)
	// }
	// for key, value := range termFrequency {
	// 	fmt.Println(key, value)
	// }
	_ = termFrequency

	docIds := getDocIdsForQueryTokens(queryTokens, invertedIndex)
	fmt.Println(docIds)
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("Elapsed time: %v\n", elapsedTime)
}
