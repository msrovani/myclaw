import os
import json
import re
from pathlib import Path

def extract_lite_info(skill_md_path):
    """Extrai apenas ID e uma descrição curta para o índice JSON."""
    if not os.path.exists(skill_md_path):
        return None

    with open(skill_md_path, "r", encoding="utf-8") as f:
        content = f.read()

    # Extrair Nome do Frontmatter
    name_match = re.search(r'name:\s*(.*)', content)
    name = name_match.group(1).strip() if name_match else os.path.basename(os.path.dirname(skill_md_path))

    # Extrair primeira frase da descrição (limite 100 chars)
    desc_match = re.search(r'description:\s*(.*)', content)
    description = desc_match.group(1).strip() if desc_match else "No description"
    if len(description) > 100:
        description = description[:97] + "..."

    return {
        "id": name.lower().replace(" ", "-"),
        "name": name,
        "description": description,
        "path": os.path.relpath(skill_md_path, start=os.path.join(os.getcwd(), ".agent2")).replace("\\", "/")
    }

def update_index():
    base_dir = Path("C:/Users/msrov/OneDrive/Área de Trabalho/openclauw/.agent2/skills")
    index_path = base_dir.parent / "skills_index_lite.json"

    skills_index = []

    for root, dirs, files in os.walk(base_dir):
        if "SKILL.min.md" in files:
            info = extract_lite_info(os.path.join(root, "SKILL.min.md"))
            if info:
                skills_index.append(info)

    with open(index_path, "w", encoding="utf-8") as f:
        json.dump(skills_index, f, indent=2)

    print(f"Índice Lite atualizado com {len(skills_index)} skills em {index_path}")

if __name__ == "__main__":
    update_index()
