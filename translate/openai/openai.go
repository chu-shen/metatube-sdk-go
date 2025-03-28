package openai

import (
	openai "github.com/zijiren233/openai-translator"

	"github.com/metatube-community/metatube-sdk-go/translate"
)

var _ translate.Translator = (*OpenAI)(nil)

type OpenAI struct {
	APIKey string `json:"openai-api-key"`
}

func (oa *OpenAI) Translate(q, source, target string) (result string, err error) {
    opts := []openai.Option{
        openai.WithFrom(source),
    }
    if true {
        opts = append(opts, openai.WithUrl("https://api.deepseek.com"))
        opts = append(opts, openai.WithModel("deepseek-reasoner")) // deepseek-chat
    }
    return openai.Translate(q, target, oa.APIKey, opts...)
}

func init() {
	translate.Register(&OpenAI{})
}
