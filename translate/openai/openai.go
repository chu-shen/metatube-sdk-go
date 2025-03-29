package openai

import (
	"fmt"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	openai "github.com/chu-shen/openai-translator"

	"github.com/metatube-community/metatube-sdk-go/translate"
)

var _ translate.Translator = (*OpenAI)(nil)

type OpenAI struct {
	APIKey string `json:"openai-api-key"`
	cache  *expirable.Cache[string, string] // 使用带TTL的缓存
}

func (oa *OpenAI) Translate(q, source, target string) (result string, err error) {
    key := oa.cacheKey(q, source, target)
	if val, ok := oa.cache.Get(key); ok {
		return val, nil
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
    }
    if true {
        opts = append(opts, openai.WithUrl("https://api.deepseek.com"))
        opts = append(opts, openai.WithModel("deepseek-chat")) // deepseek-reasoner
        opts = append(opts, openai.WithSystemPrompt(systemPrompt))
    }
    result, err = openai.Translate(q, target, oa.APIKey, opts...)
	if err != nil {
		return "", err
	}

	oa.cache.Add(key, result)
	
	return result, nil
}

func init() {
	translate.Register(&OpenAI{
		APIKey: os.Getenv("OPENAI_API_KEY"),
		cache: expirable.NewLRU[string, string](
			1000,
			nil,
			7*24*time.Hour,
		),
	})
}

func (oa *OpenAI) cacheKey(q, source, target string) string {
	return fmt.Sprintf("%s|%s|%s", source, target, q)
}