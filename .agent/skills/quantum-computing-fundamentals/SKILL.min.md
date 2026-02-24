---
name: quantum-computing-fundamentals
description: Technical introduction to quantum computing principles and basic algorithm development. Use this skill when exploring Qiskit (IBM), Q# (Microsoft), or Cirq (Google) for quantum circuit design, qubit manipulation, and proof-of-concept quantum applications.
---

# Quantum Computing Fundamentals

Understanding the transition from classical bits to quantum qubits.

## Core Concepts

- **Superposition**: Qubits existing in multiple states simultaneously.
- **Entanglement**: Correlation between qubits where the state of one determines the other.
- **Interference**: Amplifying correct answers while canceling out incorrect ones.

## Key Algorithms

### 1. Shor's Algorithm

Revolutionary for factoring large numbers (security impact).

### 2. Grover's Algorithm

Unstructured search with quadratic speedup.

### 3. VQE (Variational Quantum Eigensolver)

Used in chemistry and material science simulation.

## Workflow (Qiskit Example)

1. **Initialize**: Define quantum and classical registers.
2. **Gates**: Apply Hadamard (H), CNOT (CX), and Phase (Z) gates.
3. **Measurement**: Collapse quantum states into classical bits.
4. **Execution**: Run on simulator or real IBM Quantum hardware.

```python
from qiskit import QuantumCircuit
qc = QuantumCircuit(2)
qc.h(0)     # Superposition
qc.cx(0, 1) # Entanglement
qc.measure_all()
```

## Security & Post-Quantum (PQC)

Prepare for the "Quantum Apocalypse" by transitioning to Lattice-based or Hash-based cryptography (NIST standards).