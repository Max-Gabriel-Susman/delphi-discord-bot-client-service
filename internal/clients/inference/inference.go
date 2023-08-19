package inference

type InferenceGenerationParameters struct {
	BestOf               int      `json:"best_of"`
	DecoderInputeDetails bool     `json:"decoder_input_details"`
	Details              bool     `json:"details"`
	DoSample             bool     `json:"do_sample"`
	MaxNewTokens         int      `json:"max_new_tokens"`
	RepetitionPenalty    float32  `json:"repetition_penalty"`
	ReturnFullText       bool     `json:"return_full_text"`
	Seed                 *string  `json:"seed"` // verify type of this example value was null
	Stop                 []string `json:"stop"`
	Temperature          float32  `json:"temperature"`
	TopK                 float32  `json:"top_k"` // verify type, example vaulue was int but peers were floats
	TopP                 float32  `json:"top_p"`
	Truncate             *string  `json:"truncate"` // verify type of this example value was null
	TypicalP             float32  `json:"typical_p"`
	Watermark            bool     `json:"watermark"`
	// make sure we're not missing any fields: TODO
}

type GenerateInferenceRequest struct {
	Inputs     string                        `json:"inputs"`
	Parameters InferenceGenerationParameters `json:"parameters"`
}

type BestOfSequence struct {
	FinishReason       string           `json:"finish_reason"`
	GeneratedInference string           `json:"generated_text"`
	GeneratedTokens    int              `json:"generated_tokens"`
	Prefill            []PrefillElement `json:"prefill"`
	Seed               int              `json:"seed"`
	Tokens             []Token          `json:"tokens"`
}

type PrefillElement struct { // TODO: perhaps a better name can be found from the falcon docs?
	ID      int     `json:"id"`
	LogProb float32 `json:"logprob"`
	Text    string  `json:"text"`
}

type Token struct {
	ID      int     `json:"id"`
	LogProb float32 `json:"logprob"`
	Special bool    `json:"special"`
	Text    string  `json:"text"`
}

type GeneratedInferenceResponseDetails struct {
	BestOfSequences []BestOfSequence `json:"best_of_sequences"`
	FinishReason    string           `json:"finish_reason"`
	GeneratedTokens int              `json:"generated_tokens"`
	Prefill         []PrefillElement `json:"prefill"`
	Seed            int              `json:"seed"`
	Tokens          []Token          `json:"tokens"`
}

type GeneratedInferenceResponse struct {
	Details            GeneratedInferenceResponseDetails `json:"details"`
	GeneratedInference string                            `json:"generated_text"`
}

func (resp *GeneratedInferenceResponse) ConcatenateTokens() string {
	var concatenated string
	for _, token := range resp.Details.Tokens {
		concatenated += token.Text + " "
	}
	return concatenated
}

type ErrorResponse struct {
	Error     string `json:"error"`
	ErrorType string `json:"error_type"`
}
