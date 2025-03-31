package openai

import (	
	"log"
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
		log.Printf("Cache HIT for key: %s", cacheKey)
		return cached, nil
	}

    // https://github.com/metatube-community/metatube-sdk-go/pull/143
//     systemPrompt := `You are a professional translator for adult video content.
// Rules:
// 1. Use official translations for actor/actress names if available, otherwise keep them unchanged
// 2. Do not invent translations for names without official versions
// 3. Maintain any numbers, dates, and measurements in their original format
// 4. Translate naturally and fluently, avoiding word-for-word translation
// 5. Do not output any explanations or notes
// 6. Only output the translation`
    systemPrompt := `你是一名专业的成人影片内容翻译专家，请严格遵守以下规则：
1. 演员名称、术语如有官方译名必须使用，否则保持原名
2. 确保译文自然流畅，符合目标语言表达习惯，禁止逐字直译或机器翻译腔
3. 仅输出最终翻译结果，不得包含任何形式的解释/注释/说明
4. 输出必须为纯净的翻译文本，不得包含任何非翻译内容`

	opts := []openai.Option{
		openai.WithFrom(source),
		openai.WithUrl("https://api.deepseek.com"),
		openai.WithModel("deepseek-chat"),  // deepseek-reasoner
		openai.WithSystemPrompt(systemPrompt),
	}

	result, err = openai.Translate(q, target, oa.APIKey, opts...)
	if err != nil {
		log.Printf("failed: %s", err)
		return "", err
	}
	sanitizeResult := sanitizeText(result)
	log.Printf("new translation(sanitize): %s", sanitizeResult)

	translationCache.Add(cacheKey, result)
	log.Printf("Cache STORED for key: %s", cacheKey)
	return result, nil
}

func init() {
	translate.Register(&OpenAI{})
}

func sanitizeText(text string) string {
    text = strings.ReplaceAll(text, "\n\n", "\n")
    return text
}