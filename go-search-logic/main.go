package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"

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

type TFIDFIndex map[string]map[int]TermStats

type InvertedIndex map[string][]int

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
		stemmed := english.Stem(word, false)
		words[i] = stemmed
	}

	return words
}

func buildInvertedIndex(docs []Document) InvertedIndex {
	index := make(InvertedIndex)

	for docID, doc := range docs {
		tokens := tokenizeAndStem(doc.Title)

		for _, token := range tokens {
			index[token] = append(index[token], docID)
		}
	}

	return index
}

func calculateTFIDF(docs []Document, invertedIndex InvertedIndex) TFIDFIndex {
	tfidfIndex := make(TFIDFIndex)

	N := float64(len(docs))

	df := make(map[string]int)
	for _, doc := range docs {
		termsInDoc := make(map[string]bool)

		tokens := tokenizeAndStem(doc.Title)

		for _, token := range tokens {
			termsInDoc[token] = true
		}

		for term := range termsInDoc {
			df[term]++
		}
	}

	for docID, doc := range docs {
		tokens := tokenizeAndStem(doc.Title)

		tf := make(map[string]float64)
		for _, token := range tokens {
			tf[token]++
		}

		for term, freq := range tf {
			idf := math.Log(N / float64(1+df[term]))

			if _, exists := tfidfIndex[term]; !exists {
				tfidfIndex[term] = make(map[int]TermStats)
			}
			tfidfIndex[term][docID] = TermStats{TF: freq, IDF: idf}
		}
	}

	return tfidfIndex
}

func processQuery(query string) []string {
	tokens := tokenizeAndStem(query)
	return tokens
}

// func calculateCosineSimilarity(queryVector map[string]float64, docVector map[int]TermStats) float64 {
//     dotProduct := 0.0
//     queryMagnitude := 0.0
//     docMagnitude := 0.0

//     // Calculate dot product and magnitudes
//     for term, queryTFIDF := range queryVector {
//         docTFIDF, exists := docVector[term]
//         if exists {
//             dotProduct += queryTFIDF * docTFIDF.TF * docTFIDF.IDF
//         }
//         queryMagnitude += queryTFIDF * queryTFIDF
//     }

//     for _, docTFIDF := range docVector {
//         docMagnitude += docTFIDF.TF * docTFIDF.TF * docTFIDF.IDF * docTFIDF.IDF
//     }

//     // Calculate cosine similarity
//     if queryMagnitude == 0 || docMagnitude == 0 {
//         return 0.0
//     }
//     return dotProduct / (math.Sqrt(queryMagnitude) * math.Sqrt(docMagnitude))
// }

func rankDocuments(queryTokens []string, tfidfIndex TFIDFIndex) []int {
	// Calculate the query vector
	queryVector := make(map[string]float64)
	for _, token := range queryTokens {
		queryVector[token]++
	}

	// Calculate cosine similarity for each document
	similarityScores := make(map[int]float64)
	for term, queryTFIDF := range queryVector {
		if docScores, exists := tfidfIndex[term]; exists {
			for docID, docTFIDF := range docScores {
				similarityScores[docID] += queryTFIDF * docTFIDF.TF * docTFIDF.IDF
			}
		}
	}

	// Rank documents based on similarity scores
	var rankedDocs []int
	for docID := range similarityScores {
		rankedDocs = append(rankedDocs, docID)
	}

	return rankedDocs
}

func main() {
	userQuery := "array"
	queryTokens := processQuery(userQuery)
	_ = queryTokens

	dir := "../code_snippets"

	fileList, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory: ", err)
		return
	}

	for _, file := range fileList {
		filePath := fmt.Sprintf("../code_snippets/%s", file.Name())
		doc, err := loadJSON(filePath)
		if err != nil {
			fmt.Println("Error loading JSON:", err)
			return
		}

		docs = append(docs, doc)
	}

	invertedIndex := buildInvertedIndex(docs)

	tfidfIndex := calculateTFIDF(docs, invertedIndex)

	rankedDocs := rankDocuments(queryTokens, tfidfIndex)

	fmt.Println("Ranked Document IDs:", rankedDocs)

	fmt.Println(docs[rankedDocs[0]].Title)
}
