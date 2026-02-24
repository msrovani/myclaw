import os
import re

def minify_workflow(content):
    """
    Minifica os workflows (slash commands) para consumo de IA.
    - Mantém o frontmatter (descrição do comando).
    - Remove tabelas explicativas extensas.
    - Condensa protocolos em passos atômicos.
    - Remove exemplos de saída redundantes.
    """
    # Preservar YAML frontmatter
    frontmatter_match = re.match(r'^---\s*\n(.*?)\n---\s*\n', content, re.DOTALL)
    frontmatter = ""
    body = content
    if frontmatter_match:
        frontmatter = frontmatter_match.group(0)
        body = content[len(frontmatter):]

    # Remover tabelas de agentes e modos (muito tokens)
    body = re.sub(r'\|.*\|.*\n\|[- :|]*\|.*\n(\|.*\|.*\n)*', '', body)

    # Simplificar Protocolos
    body = re.sub(r'## 🔴 STRICT 2-PHASE ORCHESTRATION.*?\n(?=##|#|$)', '## PHASE PROTOCOL\n1. Plan (project-planner)\n2. Approve (User)\n3. Execute (Parallel Specialists)\n', body, flags=re.DOTALL)

    # Remover seções de "Available Agents" (já estão no ARCHITECTURE_LITE)
    body = re.sub(r'## Available Agents.*?\n(?=##|#|$)', '', body, flags=re.DOTALL)

    # Condensar regras críticas
    body = re.sub(r'> ⚠️.*?(\n\n|$)', '', body, flags=re.DOTALL)

    # Reduzir espaços
    body = re.sub(r'\n\s*\n', '\n\n', body)

    return (frontmatter + body).strip()

def process_workflows(src_dir, dest_dir):
    if not os.path.exists(dest_dir):
        os.makedirs(dest_dir)

    for workflow_file in os.listdir(src_dir):
        if workflow_file.endswith(".md"):
            src_path = os.path.join(src_dir, workflow_file)
            dest_path = os.path.join(dest_dir, workflow_file)

            try:
                with open(src_path, 'r', encoding='utf-8') as f:
                    content = f.read()

                minified = minify_workflow(content)

                with open(dest_path, 'w', encoding='utf-8') as f:
                    f.write(minified)

                print(f"Workflow otimizado: {workflow_file}")
            except Exception as e:
                print(f"Erro ao processar {workflow_file}: {e}")

if __name__ == "__main__":
    BASE_PATH = "C:/Users/msrov/OneDrive/Área de Trabalho/openclauw"
    SRC_WORKFLOWS = os.path.join(BASE_PATH, ".agent/workflows")
    DEST_WORKFLOWS = os.path.join(BASE_PATH, ".agent2/workflows")

    process_workflows(SRC_WORKFLOWS, DEST_WORKFLOWS)
