// internal/config/default.go
package config

var DefaultConfig = `# Revly Configuration

[llm]
# Set the LLM provider and base URL for API requests.
# The default provider is OpenRouter, but you can change it to any other provider.
# The base URL for that provider must also be set.

provider = "openrouter"
base_url = "https://openrouter.ai/api/v1"
api_base_url = "https://openrouter.ai/api/v1/chat/completions"
# SET YOUR API KEY using EXPORT LLM_API_KEY=<your-api-key> or set it in your environment variables.

[models]
# List all the models you want to use in the order of preference.
# The first model that is available will be used.
# You can use the model names from OpenRouter or any other provider you are using.
# Example:
# models = [
#   "qwen/qwen3-coder:free",
#   "qwen/qwen3-235b-a22b-2507:free",
# ]

models = [
"qwen/qwen3-coder:free",
"qwen/qwen3-235b-a22b-2507:free",
"moonshotai/kimi-k2:free",
"cognitivecomputations/dolphin-mistral-24b-venice-edition:free",
"tngtech/deepseek-r1t2-chimera:free",
"moonshotai/kimi-dev-72b:free",
"deepseek/deepseek-r1-0528-qwen3-8b:free",
"tencent/hunyuan-a13b-instruct:free",
"mistralai/mistral-small-3.2-24b-instruct:free",
"deepseek/deepseek-r1-0528:free",
"tngtech/deepseek-r1t-chimera:free",
"microsoft/mai-ds-r1:free",
"moonshotai/kimi-vl-a3b-thinking:free",
"nvidia/llama-3.1-nemotron-ultra-253b-v1:free",
]

[git]
show_diff = true
push_on_commit = false
`