package scrapper

import "strings"

var categories = map[string][]string{
    "backend":    {"backend", "back-end", "back end", "java developer", "python developer", "go developer", "node", "elixir"},
    "frontend":   {"frontend", "react", "vue", "angular", "next", "front-end", "front end", "web developer"},
    "fullstack":  {"fullstack", "full-stack", "full stack", "mern", "mean", "software engineer", "laravel"},
    "devops":     {"devops", "cloud", "aws", "azure", "gcp", "kubernetes", "docker", "infrastructure", "sre"},
    "data":       {"data engineer", "etl", "big data", "spark"},
    "ml_ai":      {"machine learning", "ai", "ml", "deep learning", "llm"},
    "mobile":     {"mobile", "android", "ios", "flutter", "react native"},
}

func getCategoriesFromTitle(title string) []string {
	var matchedCategories []string
	for category, keywords := range categories {
		for _, keyword := range keywords {
			if strings.Contains(strings.ToLower(title), strings.ToLower(keyword)) {
				matchedCategories = append(matchedCategories, category)
				break
			}
		}
	}
	return matchedCategories
}