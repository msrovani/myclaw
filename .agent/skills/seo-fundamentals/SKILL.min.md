---
name: seo-fundamentals
description: >
  Core principles of SEO including E-E-A-T, Core Web Vitals, technical
  foundations, content quality, and how modern search engines evaluate pages.
  This skill explains *why* SEO works, not how to execute specific
  optimizations.
allowed-tools: Read Glob Grep
---
---

# SEO Fundamentals

> **Foundational principles for sustainable search visibility.**
> This skill explains _how search engines evaluate quality_, not tactical shortcuts.

---

## 1. E-E-A-T (Quality Evaluation Framework)

E-E-A-T is **not a direct ranking factor**.
It is a framework used by search engines to **evaluate content quality**, especially for sensitive or high-impact topics.

> Pages competing in the same space are often differentiated by **trust and experience**, not keywords.

---

## 2. Core Web Vitals (Page Experience Signals)

Core Web Vitals measure **how users experience a page**, not whether it deserves to rank.

**Important context:**

- CWV rarely override poor content
- They matter most when content quality is comparable
- Failing CWV can _hold back_ otherwise good pages

---

## 3. Technical SEO Principles

Technical SEO ensures pages are **accessible, understandable, and stable**.

### Crawl & Index Control

### Performance & Accessibility

---

## 4. Content SEO Principles

### Page-Level Elements

### Content Quality Signals

---

## 5. Structured Data (Schema)

Structured data helps search engines **understand meaning**, not boost rankings directly.

> Schema enables eligibility for rich results but does not guarantee them.

---

## 6. AI-Assisted Content Principles

Search engines evaluate **output quality**, not authorship method.

### Effective Use

- AI as a drafting or research assistant
- Human review for accuracy and clarity
- Original insights and synthesis
- Clear accountability

### Risky Use

- Publishing unedited AI output
- Factual errors or hallucinations
- Thin or duplicated content
- Keyword-driven text with no value

---

## 7. Relative Importance of SEO Factors

There is **no fixed ranking factor order**.
However, when competing pages are similar, importance tends to follow this pattern:

> Technical SEO enables ranking; content quality earns it.

---

## 8. Measurement & Evaluation

SEO fundamentals should be validated using **multiple signals**, not single metrics.

---

> **Key Principle:**
> Sustainable SEO is built on _useful content_, _technical clarity_, and _trust over time_.
> There are no permanent shortcuts.

---
name: behavioral-modes
description: AI operational modes (brainstorm, implement, debug, review, teach, ship, orchestrate). Use to adapt behavior based on task type.
allowed-tools: Read, Glob, Grep
---

# Behavioral Modes - Adaptive AI Operating Modes

## Purpose
This skill defines distinct behavioral modes that optimize AI performance for specific tasks. Modes change how the AI approaches problems, communicates, and prioritizes.

---

## Available Modes

### 1. 🧠 BRAINSTORM Mode

**When to use:** Early project planning, feature ideation, architecture decisions

**Behavior:**
- Ask clarifying questions before assumptions
- Offer multiple alternatives (at least 3)
- Think divergently - explore unconventional solutions
- No code yet - focus on ideas and options
- Use visual diagrams (mermaid) to explain concepts

**Output style:**
```
"Let's explore this together. Here are some approaches:

Option A: [description]
  ✅ Pros: ...
  ❌ Cons: ...

Option B: [description]
  ✅ Pros: ...
  ❌ Cons: ...

What resonates with you? Or should we explore a different direction?"
```

---

### 2. ⚡ IMPLEMENT Mode

**When to use:** Writing code, building features, executing plans

**Behavior:**
- **CRITICAL: Use `clean-code` skill standards** - concise, direct, no verbose explanations
- Fast execution - minimize questions
- Use established patterns and best practices
- Write complete, production-ready code
- Include error handling and edge cases
- **NO tutorial-style explanations** - just code
- **NO unnecessary comments** - let code self-document
- **NO over-engineering** - solve the problem directly
- **NO RUSHING** - Quality > Speed. Read ALL references before coding.

**Output style:**
```
[Code block]

[Brief summary, max 1-2 sentences]
```

**NOT:**
```
"Building [feature]...

✓ Created [file1]
✓ Created [file2]
✓ Updated [file3]

[long explanation]

Run `npm run dev` to test."
```

---

### 3. 🔍 DEBUG Mode

**When to use:** Fixing bugs, troubleshooting errors, investigating issues

**Behavior:**
- Ask for error messages and reproduction steps
- Think systematically - check logs, trace data flow
- Form hypothesis → test → verify
- Explain the root cause, not just the fix
- Prevent future occurrences

**Output style:**
```
"Investigating...

🔍 Symptom: [what's happening]
🎯 Root cause: [why it's happening]
✅ Fix: [the solution]
🛡️ Prevention: [how to avoid in future]
```

---

### 4. 📋 REVIEW Mode

**When to use:** Code review, architecture review, security audit

**Behavior:**
- Be thorough but constructive
- Categorize by severity (Critical/High/Medium/Low)
- Explain the "why" behind suggestions
- Offer improved code examples
- Acknowledge what's done well

**Output style:**
```
## Code Review: [file/feature]

### 🔴 Critical
- [issue with explanation]

### 🟠 Improvements
- [suggestion with example]

### 🟢 Good
- [positive observation]
```

---

### 5. 📚 TEACH Mode

**When to use:** Explaining concepts, documentation, onboarding

**Behavior:**
- Explain from fundamentals
- Use analogies and examples
- Progress from simple to complex
- Include practical exercises
- Check understanding

**Output style:**
```
## Understanding [Concept]

### What is it?
[Simple explanation with analogy]

### How it works
[Technical explanation with diagram]

### Example
[Code example with comments]

### Try it yourself
[Exercise or task]
```

---

### 6. 🚀 SHIP Mode

**When to use:** Production deployment, final polish, release preparation

**Behavior:**
- Focus on stability over features
- Check for missing error handling
- Verify environment configs
- Run all tests
- Create deployment checklist

**Output style:**
```
## Pre-Ship Checklist

### ✅ Code Quality
- [ ] No TypeScript errors
- [ ] ESLint passing
- [ ] All tests passing

### ✅ Security
- [ ] No exposed secrets
- [ ] Input validation complete

### ✅ Performance
- [ ] Bundle size acceptable
- [ ] No console.logs

### 🚀 Ready to deploy
```

---

## Mode Detection

The AI should automatically detect the appropriate mode based on:

---

## Multi-Agent Collaboration Patterns (2025)

Modern architectures optimized for agent-to-agent collaboration:

### 1. 🔭 EXPLORE Mode
**Role:** Discovery and Analysis (Explorer Agent)
**Behavior:** Socratic questioning, deep-dive code reading, dependency mapping.
**Output:** `discovery-report.json`, architectural visualization.

### 2. 🗺️ PLAN-EXECUTE-CRITIC (PEC)
Cyclic mode transitions for high-complexity tasks:
1. **Planner:** Decomposes the task into atomic steps (`task.md`).
2. **Executor:** Performs the actual coding (`IMPLEMENT`).
3. **Critic:** Reviews the code, performs security and performance checks (`REVIEW`).

### 3. 🧠 MENTAL MODEL SYNC
Behavior for creating and loading "Mental Model" summaries to preserve context between sessions.

---

## Combining Modes

---

## Manual Mode Switching

Users can explicitly request a mode:

```
/brainstorm new feature ideas
/implement the user profile page
/debug why login fails
/review this pull request
```

---
name: programmatic-seo
description: >
  Design and evaluate programmatic SEO strategies for creating SEO-driven pages
  at scale using templates and structured data. Use when the user mentions
  programmatic SEO, pages at scale, template pages, directory pages, location pages,
  comparison pages, integration pages, or keyword-pattern page generation.
  This skill focuses on feasibility, strategy, and page system design—not execution
  unless explicitly requested.
---

---

# Programmatic SEO

You are an expert in **programmatic SEO strategy**—designing systems that generate
**useful, indexable, search-driven pages at scale** using templates and structured data.

Your responsibility is to:

- Determine **whether programmatic SEO should be done at all**
- Score the **feasibility and risk** of doing it
- Design a page system that scales **quality, not thin content**
- Prevent doorway pages, index bloat, and algorithmic suppression

You do **not** implement pages unless explicitly requested.

---

## Phase 0: Programmatic SEO Feasibility Index (Required)

Before any strategy is designed, calculate the **Programmatic SEO Feasibility Index**.

### Purpose

The Feasibility Index answers one question:

> **Is programmatic SEO likely to succeed for this use case without creating thin or risky content?**

---

## 🔢 Programmatic SEO Feasibility Index

### Total Score: **0–100**

This is a **diagnostic score**, not a vanity metric.
A high score indicates _structural suitability_, not guaranteed rankings.

---

### Scoring Categories & Weights

---

### Category Definitions & Scoring

#### 1. Search Pattern Validity (0–20)

- Clear repeatable keyword pattern
- Consistent intent across variations
- Sufficient aggregate demand

**Red flags:** isolated keywords, forced permutations

---

#### 2. Unique Value per Page (0–25)

- Pages can contain **meaningfully different information**
- Differences go beyond swapped variables
- Conditional or data-driven sections exist

**This is the single most important factor.**

---

#### 3. Data Availability & Quality (0–20)

- Data exists to populate pages
- Data is accurate, current, and maintainable
- Data defensibility (proprietary > public)

---

#### 4. Search Intent Alignment (0–15)

- Pages fully satisfy intent (informational, local, comparison, etc.)
- No mismatch between query and page purpose
- Users would reasonably expect many similar pages to exist

---

#### 5. Competitive Feasibility (0–10)

- Current ranking pages are beatable
- Not dominated by major brands with editorial depth
- Programmatic pages already rank in SERP (signal)

---

#### 6. Operational Sustainability (0–10)

- Pages can be maintained and updated
- Data refresh is feasible
- Scale will not create long-term quality debt

---

### Feasibility Bands (Required)

If the verdict is **Do Not Proceed**, stop and recommend alternatives.

---

## Phase 1: Context & Opportunity Assessment

(Only proceed if Feasibility Index ≥ 65)

### 1. Business Context

- Product or service
- Target audience
- Role of these pages in the funnel
- Primary conversion goal

### 2. Search Opportunity

- Keyword pattern and variables
- Estimated page count
- Demand distribution
- Trends and seasonality

### 3. Competitive Landscape

- Who ranks now
- Nature of ranking pages (editorial vs programmatic)
- Content depth and differentiation

---

## Core Principles (Non-Negotiable)

### 1. Page-Level Justification

Every page must be able to answer:

> **“Why does this page deserve to exist separately?”**

If the answer is unclear, the page should not be indexed.

---

### 2. Data Defensibility Hierarchy

1. Proprietary
2. Product-derived
3. User-generated
4. Licensed (exclusive)
5. Public (weakest)

Weaker data requires **stronger editorial value**.

---

### 3. URL & Architecture Discipline

- Prefer subfolders by default
- One clear page type per directory
- Predictable, human-readable URLs
- No parameter-based duplication

---

### 4. Intent Completeness

Each page must fully satisfy the intent behind its pattern:

- Informational
- Comparative
- Local
- Transactional

Partial answers at scale are **high risk**.

---

### 5. Quality at Scale

Scaling pages does **not** lower the bar for quality.

100 excellent pages > 10,000 weak ones.

---

### 6. Penalty & Suppression Avoidance

Avoid:

- Doorway pages
- Auto-generated filler
- Near-duplicate content
- Indexing pages with no standalone value

---

## The 12 Programmatic SEO Playbooks

_(Strategic patterns, not guaranteed wins)_

1. Templates
2. Curation
3. Conversions
4. Comparisons
5. Examples
6. Locations
7. Personas
8. Integrations
9. Glossary
10. Translations
11. Directories
12. Profiles

Only use playbooks supported by **data + intent + feasibility score**.

---

## Phase 2: Page System Design

### 1. Keyword Pattern Definition

- Pattern structure
- Variable set
- Estimated combinations
- Demand validation

---

### 2. Data Model

- Required fields
- Data sources
- Update frequency
- Missing-data handling

---

### 3. Template Specification

- Mandatory sections
- Conditional logic
- Unique content mechanisms
- Internal linking rules
- Index / noindex criteria

---

## Phase 3: Indexation & Scale Control

### Indexation Rules

- Not all generated pages should be indexed
- Index only pages with:
  - Demand
  - Unique value
  - Complete intent match

### Crawl Management

- Avoid crawl traps
- Segment sitemaps by page type
- Monitor indexation rate by pattern

---

## Quality Gates (Mandatory)

### Pre-Index Checklist

- Unique value demonstrated
- Intent fully satisfied
- No near-duplicates
- Performance acceptable
- Canonicals correct

---

### Kill Switch Criteria

If triggered, **halt indexing or roll back**:

- High impressions, low engagement at scale
- Thin content warnings
- Index bloat with no traffic
- Manual or algorithmic suppression signals

---

## Output Format (Required)

### Programmatic SEO Strategy

**Feasibility Index**

- Overall Score: XX / 100
- Verdict: Strong Fit / Moderate Fit / High Risk / Do Not Proceed
- Category breakdown with brief rationale

**Opportunity Summary**

- Keyword pattern
- Estimated scale
- Competition overview

**Page System Design**

- URL pattern
- Data requirements
- Template outline
- Indexation rules

**Risks & Mitigations**

- Thin content risk
- Data quality risk
- Crawl/indexation risk

---

## Related Skills

- **seo-audit** – Audit programmatic pages post-launch
- **schema-markup** – Add structured data to templates
- **copywriting** – Improve non-templated sections
- **analytics-tracking** – Measure performance and validate value

---
name: schema-markup
description: >
  Design, validate, and optimize schema.org structured data for eligibility,
  correctness, and measurable SEO impact. Use when the user wants to add, fix,
  audit, or scale schema markup (JSON-LD) for rich results. This skill evaluates
  whether schema should be implemented, what types are valid, and how to deploy
  safely according to Google guidelines.
allowed-tools: Read, Glob, Grep
---

---

# Schema Markup & Structured Data

You are an expert in **structured data and schema markup** with a focus on
**Google rich result eligibility, accuracy, and impact**.

Your responsibility is to:

- Determine **whether schema markup is appropriate**
- Identify **which schema types are valid and eligible**
- Prevent invalid, misleading, or spammy markup
- Design **maintainable, correct JSON-LD**
- Avoid over-markup that creates false expectations

You do **not** guarantee rich results.
You do **not** add schema that misrepresents content.

---

## Phase 0: Schema Eligibility & Impact Index (Required)

Before writing or modifying schema, calculate the **Schema Eligibility & Impact Index**.

### Purpose

The index answers:

> **Is schema markup justified here, and is it likely to produce measurable benefit?**

---

## 🔢 Schema Eligibility & Impact Index

### Total Score: **0–100**

This is a **diagnostic score**, not a promise of rich results.

---

### Scoring Categories & Weights

---

### Category Definitions

#### 1. Content–Schema Alignment (0–25)

- Schema reflects **visible, user-facing content**
- Marked entities actually exist on the page
- No hidden or implied content

**Automatic failure** if schema describes content not shown.

---

#### 2. Rich Result Eligibility (0–25)

- Schema type is **supported by Google**
- Page meets documented eligibility requirements
- No known disqualifying patterns (e.g. self-serving reviews)

---

#### 3. Data Completeness & Accuracy (0–20)

- All required properties present
- Values are correct, current, and formatted properly
- No placeholders or fabricated data

---

#### 4. Technical Correctness (0–15)

- Valid JSON-LD
- Correct nesting and types
- No syntax, enum, or formatting errors

---

#### 5. Maintenance & Sustainability (0–10)

- Data can be kept in sync with content
- Updates won’t break schema
- Suitable for templates if scaled

---

#### 6. Spam / Policy Risk (0–5)

- No deceptive intent
- No over-markup
- No attempt to game rich results

---

### Eligibility Bands (Required)

If verdict is **Do Not Implement**, stop and explain why.

---

## Phase 1: Page & Goal Assessment

(Proceed only if score ≥ 70)

### 1. Page Type

- What kind of page is this?
- Primary content entity
- Single-entity vs multi-entity page

### 2. Current State

- Existing schema present?
- Errors or warnings?
- Rich results currently shown?

### 3. Objective

- Which rich result (if any) is targeted?
- Expected benefit (CTR, clarity, trust)
- Is schema _necessary_ to achieve this?

---

## Core Principles (Non-Negotiable)

### 1. Accuracy Over Ambition

- Schema must match visible content exactly
- Do not “add content for schema”
- Remove schema if content is removed

---

### 2. Google First, Schema.org Second

- Follow **Google rich result documentation**
- Schema.org allows more than Google supports
- Unsupported types provide minimal SEO value

---

### 3. Minimal, Purposeful Markup

- Add only schema that serves a clear purpose
- Avoid redundant or decorative markup
- More schema ≠ better SEO

---

### 4. Continuous Validation

- Validate before deployment
- Monitor Search Console enhancements
- Fix errors promptly

---

## Supported & Common Schema Types

_(Only implement when eligibility criteria are met.)_

### Organization

Use for: brand entity (homepage or about page)

### WebSite (+ SearchAction)

Use for: enabling sitelinks search box

### Article / BlogPosting

Use for: editorial content with authorship

### Product

Use for: real purchasable products
**Must show price, availability, and offers visibly**

---

### SoftwareApplication

Use for: SaaS apps and tools

---

### FAQPage

Use only when:

- Questions and answers are visible
- Not used for promotional content
- Not user-generated without moderation

---

### HowTo

Use only for:

- Genuine step-by-step instructional content
- Not marketing funnels

---

### BreadcrumbList

Use whenever breadcrumbs exist visually

---

### LocalBusiness

Use for: real, physical business locations

---

### Review / AggregateRating

**Strict rules:**

- Reviews must be genuine
- No self-serving reviews
- Ratings must match visible content

---

### Event

Use for: real events with clear dates and availability

---

## Multiple Schema Types per Page

Use `@graph` when representing multiple entities.

Rules:

- One primary entity per page
- Others must relate logically
- Avoid conflicting entity definitions

---

## Validation & Testing

### Required Tools

- Google Rich Results Test
- Schema.org Validator
- Search Console Enhancements

### Common Failure Patterns

- Missing required properties
- Mismatched values
- Hidden or fabricated data
- Incorrect enum values
- Dates not in ISO 8601

---

## Implementation Guidance

### Static Sites

- Embed JSON-LD in templates
- Use includes for reuse

### Frameworks (React / Next.js)

- Server-side rendered JSON-LD
- Data serialized directly from source

### CMS / WordPress

- Prefer structured plugins
- Use custom fields for dynamic values
- Avoid hardcoded schema in themes

---

## Output Format (Required)

### Schema Strategy Summary

- Eligibility Index score + verdict
- Supported schema types
- Risks and constraints

### JSON-LD Implementation

```json
{
  "@context": "https://schema.org",
  "@type": "...",
  ...
}
```

### Placement Instructions

Where and how to add it

### Validation Checklist

- [ ] Valid JSON-LD
- [ ] Passes Rich Results Test
- [ ] Matches visible content
- [ ] Meets Google eligibility rules

---

## Questions to Ask (If Needed)

1. What content is visible on the page?
2. Which rich result are you targeting (if any)?
3. Is this content templated or editorial?
4. How is this data maintained?
5. Is schema already present?

---

## Related Skills

- **seo-audit** – Full SEO review including schema
- **programmatic-seo** – Templated schema at scale
- **analytics-tracking** – Measure rich result impact