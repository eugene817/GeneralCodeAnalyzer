import ast
import json
import sys

def analyze_code(code: str):
    issues = []
    reccomendations = []

    try:
        tree = ast.parse(code)
        for node in ast.walk(tree):
            if isinstance(node, ast.FunctionDef) and node.name == "factorial":
                reccomendations.append("Consider using iterations")
    except Exception as e:
        issues.append(f"Error AST analyzis: {e}")


    result = {
            "issues": issues,
            "reccomendations": reccomendations
            }
    print(json.dumps(result))

if __name__ == "__main__":
    code = sys.stdin.read()
    analyze_code(code)
