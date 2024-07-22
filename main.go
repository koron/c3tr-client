package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/koron-go/jsonhttpc"
)

type Request struct {
	Prompt string `json:"prompt"`

	NPredict      int     `json:"n_predict"`
	RepeatPenalty float64 `json:"repeat_penalty"`
	Temperature   float64 `json:"temperature"`
	TopP          float64 `json:"top_p"`
}

type Response struct {
	Content string `json:"content"`

	GenerationSettings struct {
		DynatempExponent       float64       `json:"dynatemp_exponent"`
		DynatempRange          float64       `json:"dynatemp_range"`
		FrequencyPenalty       float64       `json:"frequency_penalty"`
		Grammar                string        `json:"grammar"`
		IgnoreEOS              bool          `json:"ignore_eos"`
		LogitBias              []interface{} `json:"logit_bias"`
		MinKeep                float64       `json:"min_keep"`
		MinP                   float64       `json:"min_p"`
		Mirostat               float64       `json:"mirostat"`
		MirostatETA            float64       `json:"mirostat_eta"`
		MirostatTAU            float64       `json:"mirostat_tau"`
		Model                  string        `json:"model"`
		NCtx                   int           `json:"n_ctx"`
		NDiscard               int           `json:"n_discard"`
		NKeep                  int           `json:"n_keep"`
		NPredict               int           `json:"n_predict"`
		NProbs                 int           `json:"n_probs"`
		PenalizeNL             bool          `json:"penalize_nl"`
		PenaltyPromptTokens    []interface{} `json:"penalty_prompt_tokens"`
		PresensePenalty        float64       `json:"presence_penalty"`
		RepeatLastN            int64         `json:"repeat_last_n"`
		RepeatPenalty          float64       `json:"repeat_penalty"`
		Samplers               []string      `json:"samplers"`
		Seed                   int64         `json:"seed"`
		Stop                   []string      `json:"stop"`
		Stream                 bool          `json:"stream"`
		Temperature            float64       `json:"temperature"`
		TfsZ                   float64       `json:"tfs_z"`
		TopK                   float64       `json:"top_k"`
		TopP                   float64       `json:"top_p"`
		TypicalP               float64       `json:"typical_p"`
		UsePenaltyPromptTokens bool          `json:"use_penalty_prompt_tokens"`
	} `json:"generation_settings"`

	IDSlot int    `json:"id_slot"`
	Model  string `json:"model"`
	Prompt string `json:"prompt"`

	Stop         bool   `json:"stop"`
	StoppedEOS   bool   `json:"stopped_eos"`
	StoppedLimit bool   `json:"stopped_limit"`
	StoppedWord  bool   `json:"stopped_word"`
	StoppingWord string `json:"stopping_word"`

	Timings struct {
		PredictedMS         float64 `json:"predicted_ms"`
		PredictedN          int     `json:"predicted_n"`
		PredictedPerSecond  float64 `json:"predicted_per_second"`
		PredictedPerTokenMS float64 `json:"predicted_per_token_ms"`

		PromptMS         float64 `json:"prompt_ms"`
		PromptN          int     `json:"prompt_n"`
		PromptPerSecond  float64 `json:"prompt_per_second"`
		PromptPerTokenMS float64 `json:"prompt_per_token_ms"`
	} `json:"timiings"`

	TokensCached    int `json:"tokens_cached"`
	TokensEvaluated int `json:"tokens_evaluated"`
	TokensPredicted int `json:"tokens_predicted"`

	Truncated bool `json:"truncated"`
}

func main() {
	var (
		verbose bool
		mode    string
		param   PromptParam
		req     Request
		res     Response
	)

	flag.BoolVar(&verbose, "verbose", false, `verbose messages`)
	flag.StringVar(&mode, "mode", "", `translation mode: EtoJ, JtoE or auto (default)`)
	flag.IntVar(&req.NPredict, "n_predict", -1, `number of predict`)
	flag.Float64Var(&req.RepeatPenalty, "repeat_penalty", 1.0, `repeat penalty`)
	flag.Float64Var(&req.Temperature, "temperature", 0.0, `temperature`)
	flag.Float64Var(&req.TopP, "top_p", 0.0, `top P`)
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("no prompts to translate")
	}

	content := strings.TrimSpace(flag.Arg(0))
	mode, err := regulateMode(mode, content)
	if err != nil {
		log.Fatal(err)
	}

	param.Mode = mode
	param.WritingStyle = EducationCasual
	param.Content = content
	prompt, err := param.Generate()
	if err != nil {
		log.Fatalf("failed to generate prompt: %s", err)
	}

	if verbose {
		log.Printf("mode is %q", param.Mode)
		log.Printf("prompt is... %q", prompt)
	}

	req.Prompt = prompt
	err = jsonhttpc.Do(context.Background(), "POST", "http://127.0.0.1:8080/completions", &req, &res)
	if err != nil {
		log.Fatalf("failed to request: %s", err)
	}

	fmt.Println(res.Content)
}
