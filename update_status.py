import json
from datetime import datetime

with open("ops/codex-24x7/state/todo-current.json") as f:
    todos = json.load(f)

# The goal is to close the TODO in `forge/log/encoder.go`
# We removed it, so it's resolved.
