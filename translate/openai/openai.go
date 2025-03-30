package openai

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	openai "github.com/chu-shen/openai-translator"
	"github.com/metatube-community/metatube-sdk-go/translate"
)

var _ translate.Translator = (*OpenAI)(nil)

var translationCache = expirable.NewLRU[string, string](
	1000,
	nil,
	7*24*time.Hour,
)

type OpenAI struct {
	APIKey string `json:"openai-api-key"`
}

func (oa *OpenAI) Translate(q, source, target string) (result string, err error) {
	cacheKey := strings.Join([]string{
		q,
		source,
		target,
	}, "|")

	if cached, ok := translationCache.Get(cacheKey); ok {
		return cached, nil
	}

    // https://github.com/metatube-community/metatube-sdk-go/pull/143
    systemPrompt := `You are a professional translator for adult video content.
Rules:
1. Use official translations for actor/actress names if available, otherwise keep them unchanged
2. Do not invent translations for names without official versions
3. Maintain any numbers, dates, and measurements in their original format
4. Translate naturally and fluently, avoiding word-for-word translation
5. Do not add any explanations or notes
6. Only output the translation`

	opts := []openai.Option{
		openai.WithFrom(source),
		openai.WithUrl("https://api.deepseek.com"),
		openai.WithModel("deepseek-chat"),  // deepseek-reasoner
		openai.WithSystemPrompt(systemPrompt),
	}

	result, err = openai.Translate(q, target, oa.APIKey, opts...)
	if err != nil {
		return "", err
	}

	translationCache.Add(cacheKey, result)
	return result, nil
}

func init() {
	translate.Register(&OpenAI{})
}