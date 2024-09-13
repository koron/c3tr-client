package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/koron-go/jsonhttpc"
	"github.com/mattn/go-isatty"
	"github.com/peterh/liner"
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

var (
	verbose    bool
	entrypoint string

	reqTmpl Request
)

func translate(text string, mode string, writingStyle string, subStyles map[string]string) (string, error) {
	param := PromptParam{
		Mode:         mode,
		WritingStyle: writingStyle,
		Content:      text,
	}
	prompt, err := param.Generate()
	if err != nil {
		return "", fmt.Errorf("failed to generate prompt: %w", err)
	}

	if verbose {
		log.Printf("mode is %q", param.Mode)
		log.Printf("prompt is... %q", prompt)
	}

	req := reqTmpl
	req.Prompt = prompt
	var res Response
	err = jsonhttpc.Do(context.Background(), "POST", entrypoint, &req, &res)
	if err != nil {
		return "", fmt.Errorf("failed to request: %w", err)
	}
	return res.Content, nil
}

func reverseMode(mode string) string {
	switch mode {
	case EnglishToJapanese:
		return JapaneseToEnglish
	case JapaneseToEnglish:
		return EnglishToJapanese
	default:
		return EnglishToJapanese
	}
}

func isTerminal(fd uintptr) bool {
	return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}

func canBeInteractive() bool {
	return isTerminal(os.Stdin.Fd()) && isTerminal(os.Stdout.Fd())
}

func main() {
	var (
		modeOpt   string
		wstyleOpt string
		iteration int
	)

	// Core options.
	flag.BoolVar(&verbose, "verbose", false, `verbose messages`)
	flag.StringVar(&entrypoint, "entrypoint", "http://127.0.0.1:8080/completions", `entrypoint`)
	flag.StringVar(&modeOpt, "mode", "", `translation mode: EtoJ, JtoE or auto (default)`)
	flag.StringVar(&wstyleOpt, "writingstyle", Technical, `writing style`)

	// Request options.
	flag.IntVar(&reqTmpl.NPredict, "n_predict", -1, `number of predict`)
	flag.Float64Var(&reqTmpl.RepeatPenalty, "repeat_penalty", 1.0, `repeat penalty`)
	flag.Float64Var(&reqTmpl.Temperature, "temperature", 0.0, `temperature`)
	flag.Float64Var(&reqTmpl.TopP, "top_p", 0.0, `top P`)

	// Operation mode options.
	flag.IntVar(&iteration, "iteration", 0, "number of times to repeat the reverse translation. -1 means to repeat until the translation matches the translation history.")

	flag.Parse()

	// Validate option values.
	if !IsValidWritingStyles(wstyleOpt) {
		log.Fatalf("unknown %q writingstyle. please choose one from following: %s",
			wstyleOpt, strings.Join(ValidWritingStyles, ", "))
	}

	// Start interative translation mode, if without any arguments and the
	// os.Stdin and os.Stdout are connected a TTY.
	if flag.NArg() == 0 {
		if !canBeInteractive() {
			log.Fatal("no text to translate. for enabling interactive mode, do not use pipes nor redirects")
		}
		err := interactiveMode(modeOpt, wstyleOpt, nil)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err := singleSentenceMode(flag.Arg(0), modeOpt, wstyleOpt, nil, iteration)
	if err != nil {
		log.Fatal(err)
	}
}

// Operation mode implementations

func singleSentenceMode(textOrig, modeOrig, wstyle string, subStyles map[string]string, iteration int) error {
	text := strings.TrimSpace(textOrig)
	mode, err := regulateMode(modeOrig, text)
	if err != nil {
		return err
	}

	// One-shot translation.
	if iteration == 0 {
		translation, err := translate(text, mode, wstyle, subStyles)
		if err != nil {
			return err
		}
		fmt.Println(translation)
		return nil
	}

	// Repeat the translation with specifying the number of times.
	if iteration > 0 {
		for i := range iteration + 1 {
			translation, err := translate(text, mode, wstyle, subStyles)
			if err != nil {
				return fmt.Errorf("failed at #%d: %w", i, err)
			}
			fmt.Printf("#%d\t%s\n", i, translation)
			text = translation
			mode = reverseMode(mode)
		}
		return nil
	}

	// Repeat the translation until the same result is obtained.
	seen := map[string]struct{}{}
	for i := 0; ; i++ {
		if _, ok := seen[text]; ok {
			break
		}
		seen[text] = struct{}{}
		translation, err := translate(text, mode, wstyle, subStyles)
		if err != nil {
			log.Fatalf("failed at #%d: %s", i, err)
		}
		fmt.Printf("#%d\t%s\n", i, translation)
		text = translation
		mode = reverseMode(mode)
	}
	return nil
}

func interactiveMode(modeOrig, writingStyle string, subStyles map[string]string) error {
	state := liner.NewLiner()
	defer state.Close()
	for {
		line, err := state.Prompt("c3tr> ")
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		text := strings.TrimSpace(line)
		mode, err := regulateMode(modeOrig, text)
		if err != nil {
			return err
		}
		translation, err := translate(text, mode, writingStyle, nil)
		if err != nil {
			return err
		}
		fmt.Println(translation)
		state.AppendHistory(text)
	}
}
