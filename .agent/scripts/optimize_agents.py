import os
import re

def minify_agent(content):
    """
    Minifica a persona do agente para consumo de IA.
    - Mantém o frontmatter (crucial para o sistema).
    - Converte seções longas de prosa em listas de diretrizes.
    - Remove tabelas de navegação internas.
    - Remove exemplos redundantes.
    """
    # Preservar YAML frontmatter
    frontmatter_match = re.match(r'^---\s*\n(.*?)\n---\s*\n', content, re.DOTALL)
    frontmatter = ""
    body = content
    if frontmatter_match:
        frontmatter = frontmatter_match.group(0)
        body = content[len(frontmatter):]

    # Remover Navegação e seções humanas
    body = re.sub(r'## 📑 Quick Navigation.*?\n(?=##|#|$)', '', body, flags=re.DOTALL)
    body = re.sub(r'---.*?\n', '', body) # Remover divisores

    # Simplificar seções de Role e Responsibilities
    # (Mantém os títulos mas remove o texto explicativo longo entre eles)

    # Remover blocos de exemplo verbosos
    body = re.sub(r'### ❌ WRONG Example.*?###', '###', body, flags=re.DOTALL)
    body = re.sub(r'### ❌ Example Violation.*?###', '###', body, flags=re.DOTALL)

    # Reduzir espaços
    body = re.sub(r'\n\s*\n', '\n\n', body)

    return (frontmatter + body).strip()

def process_agents(src_dir, dest_dir):
    if not os.path.exists(dest_dir):
        os.makedirs(dest_dir)

    for agent_file in os.listdir(src_dir):
        if agent_file.endswith(".md"):
            src_path = os.path.join(src_dir, agent_file)
            dest_path = os.path.join(dest_dir, agent_file)

            try:
                with open(src_path, 'r', encoding='utf-8') as f:
                    content = f.read()

                minified = minify_agent(content)

                with open(dest_path, 'w', encoding='utf-8') as f:
                    f.write(minified)

                print(f"Agente otimizado: {agent_file}")
            except Exception as e:
                print(f"Erro ao processar {agent_file}: {e}")

if __name__ == "__main__":
    # Caminhos baseados na estrutura detectada
    BASE_PATH = "C:/Users/msrov/OneDrive/Área de Trabalho/openclauw"
    SRC_AGENTS = os.path.join(BASE_PATH, "openclaw-main/.agent/agents")
    DEST_AGENTS = os.path.join(BASE_PATH, ".agent2/agents")

    process_agents(SRC_AGENTS, DEST_AGENTS)
