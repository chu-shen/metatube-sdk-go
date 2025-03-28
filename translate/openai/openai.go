package openai

import (
	openai "github.com/chu-shen/openai-translator"

	"github.com/metatube-community/metatube-sdk-go/translate"
)

var _ translate.Translator = (*OpenAI)(nil)

type OpenAI struct {
	APIKey string `json:"openai-api-key"`
}

func (oa *OpenAI) Translate(q, source, target string) (result string, err error) {
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
        opts = append(opts, openai.WithModel("deepseek-reasoner")) // deepseek-chat
        opts = append(opts, openai.WithSystemPrompt(systemPrompt))
    }
    return openai.Translate(q, target, oa.APIKey, opts...)
}

func init() {
	translate.Register(&OpenAI{})
}
