---
name: local-llm-orchestration
description: Guidance for configuring and optimizing local Large Language Model (LLM) engines such as Ollama, vLLM, and LM Studio. Use this skill when the user wants to set up, fine-tune, or serve LLMs locally for privacy, cost-reduction, or low-latency applications. It covers hardware requirements, quantization levels, and API serving best practices.
---

# Local LLM Orchestration

This skill provides comprehensive guidance for running and optimizing Large Language Models on local hardware.

## Core Engines

### 1. Ollama

The easiest way to get up and running on macOS, Linux, and Windows.

- **Serving**: `ollama serve`
- **Running**: `ollama run <model-name>`
- **API**: Standard OpenAI-compatible endpoint at `http://localhost:11434`.

### 2. vLLM (High Throughput)

Preferred for production-grade local serving with PagedAttention.

- **Deployment**: `python -m vllm.entrypoints.openai.api_server --model <model>`
- **Optimization**: Use `--quantization awq` or `fp8` for memory efficiency.

### 3. LM Studio / AnythingLLM

GUI-based solutions for testing and structured document querying.

## Hardware & Quantization

- **VRAM is King**: Most LLMs require GPU memory.
- **4-bit (GGUF/AWQ)**: The "sweet spot" for performance/memory.
- **Rule of Thumb**:
  - 8GB VRAM → 7B-8B models (quantized).
  - 24GB VRAM → 30B-35B models (quantized) or 2x 8B models.
  - 80GB VRAM → 70B models.

## Workflow Patterns

1. **Environmental Audit**: Check CUDA/ROCm availability.
2. **Model Selection**: Choose based on task (Llama-3 for general, DeepSeek-Coder for dev, Mistral for chat).
3. **Quantization Choice**: Balance quality vs. token-per-second (TPS).
4. **API Integration**: Link the local endpoint to the application code.

## Optimization Strategies

- **Flash Attention 2**: Enable for faster processing.
- **Context Length**: Limit to avoid out-of-memory (OOM) errors.
- **Temperature/Top-P**: Tuning for deterministic vs. creative outputs.