import os
import re

def minify_markdown(content):
    """
    Minifica o conteúdo markdown para consumo de IA.
    - Mantém YAML frontmatter.
    - Remove comentários HTML.
    - Remove tabelas de navegação e mapas de conteúdo (human-only).
    - Remove seções de "Selective Reading" e "Learning Path".
    - Condensa espaços em branco.
    """
    # Preservar YAML frontmatter
    frontmatter_match = re.match(r'^---\s*\n(.*?)\n---\s*\n', content, re.DOTALL)
    frontmatter = ""
    body = content
    if frontmatter_match:
        frontmatter = frontmatter_match.group(0)
        body = content[len(frontmatter):]

    # Remover Comentários
    body = re.sub(r'<!--.*?-->', '', body, flags=re.DOTALL)

    # Remover seções específicas que incham o contexto sem adicionar lógica
    sections_to_remove = [
        r'## 🎯 Selective Reading Rule.*?\n(?=##|#|$)',
        r'## 📑 Content Map.*?\n(?=##|#|$)',
        r'## 🔗 Related Skills.*?\n(?=##|#|$)',
        r'## ✅ Decision Checklist.*?\n(?=##|#|$)',
        r'## 📚 Learning Path.*?\n(?=##|#|$)',
        r'## 📊 Impact Priority Guide.*?\n(?=##|#|$)',
        r'## 🚀 Quick Decision Tree.*?\n(?=##|#|$)',
        r'## 📖 Section Details.*?\n(?=##|#|$)'
    ]

    for section in sections_to_remove:
        body = re.sub(section, '', body, flags=re.DOTALL | re.IGNORECASE)

    # Remover tabelas Markdown (geralmente usadas para índices)
    body = re.sub(r'\|.*\|.*\n\|[- :|]*\|.*\n(\|.*\|.*\n)*', '', body)

    # Limpeza de linhas vazias excessivas
    body = re.sub(r'\n\s*\n', '\n\n', body)

    return (frontmatter + body).strip()

def bulk_convert(src_base, dest_base):
    print(f"Iniciando conversão de {src_base} para {dest_base}...")
    converted_count = 0

    for root, dirs, files in os.walk(src_base):
        if "SKILL.md" in files:
            # Caminho relativo para manter a estrutura
            rel_path = os.path.relpath(root, src_base)
            dest_dir = os.path.join(dest_base, rel_path)

            os.makedirs(dest_dir, exist_ok=True)

            src_file = os.path.join(root, "SKILL.md")
            dest_file = os.path.join(dest_dir, "SKILL.min.md")

            try:
                with open(src_file, 'r', encoding='utf-8') as f:
                    content = f.read()

                optimized = minify_markdown(content)

                with open(dest_file, 'w', encoding='utf-8') as f:
                    f.write(optimized)

                converted_count += 1
                if converted_count % 50 == 0:
                    print(f"Processadas {converted_count} skills...")
            except Exception as e:
                print(f"Erro ao processar {src_file}: {e}")

    print(f"Concluído! {converted_count} skills convertidas para .min.md")

if __name__ == "__main__":
    # Caminhos absolutos baseados no seu projeto
    BASE_PATH = "C:/Users/msrov/OneDrive/Área de Trabalho/openclauw"
    SRC = os.path.join(BASE_PATH, ".agent/skills")
    DEST = os.path.join(BASE_PATH, ".agent2/skills")

    bulk_convert(SRC, DEST)
