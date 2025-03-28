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
        opts = append(opts, openai.WithModel("deepseek-chat"))
		// deepseek-reasoner does not support successive user or assistant messages (messages[1] and messages[2] in your input). You should interleave the user/assistant messages in the message sequence.)
    }
    return openai.Translate(q, target, oa.APIKey, opts...)
}

func init() {
	translate.Register(&OpenAI{})
}
