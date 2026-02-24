---
name: slm-finetuning-unsloth
description: Guide for fine-tuning Small Language Models (SLMs) using the Unsloth framework. Use this skill for highly efficient, low-memory training of models like Phi-4, Llama 3.2 (1B/3B), and Mistral. It covers LoRA/QLoRA techniques, dataset preparation, and exporting to GGUF format for local use.
---

# SLM Fine-Tuning with Unsloth

Unsloth makes fine-tuning 2x faster and uses 70% less memory compared to standard Hugging Face implementations.

## Key Features

- **Speed**: Up to 30x faster than traditional training.
- **Memory**: Train 7B models on 16GB VRAM.
- **Native Support**: Llama, Mistral, Gemma, Phi models.

## Implementation Steps

### 1. Setup

```python
from unsloth import FastLanguageModel
import torch

model, tokenizer = FastLanguageModel.from_pretrained(
    model_name = "unsloth/Llama-3.2-3B-Instruct",
    max_seq_length = 2048,
    load_in_4bit = True,
)
```

### 2. Parameter Tuning

- **LoRA Rank (R)**: Default to 16 or 32.
- **Alpha**: Usually 2x Rank.
- **Target Modules**: `q_proj`, `k_proj`, `v_proj`, `o_proj`, `gate_proj`, `up_proj`, `down_proj`.

### 3. Dataset Prep

- Formats: Alpaca, ShareGPT, or Raw Text.
- Use `SFTTrainer` from TRL for simplified training loops.

## Export & Quantization

Export the final model to `GGUF` for use in Ollama or LM Studio:

```python
model.save_pretrained_gguf("model", tokenizer, quantization_method = "q4_k_m")
```

## Best Practices

- **Small Batches**: 2-4 with gradient accumulation.
- **Learning Rate**: 2e-4 is a common starting point.
- **Validation**: Monitor loss curves, not just accuracy.