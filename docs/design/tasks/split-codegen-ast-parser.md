Your workspace is /home/hamid/Other/projects/foreit-wt/wave3. Always use absolute paths. Never read or edit outside it. Do NOT run any git command.

ROLE: You are the delegated coder. The rule in CLAUDE.md/AGENTS.md that forbids the primary agent from editing code applies to Claude, not to you: edit the files yourself. Do NOT call any other agent. Use US spelling in comments.

# Split forge/codegen/ast_parser.go (1501 lines) by responsibility; make option parsing table-driven

Module root: /home/hamid/Other/projects/foreit-wt/wave3/forge (run go commands from there).
Only edit files in forge/codegen (and add tests there). Other agents edit other packages concurrently: ignore build errors outside forge/codegen.

## Step 1 — pure move (commit-sized, no behaviour change)
Move declarations out of ast_parser.go into new files in the same package. Do not rename anything or change any body; only package clauses and imports change.
- ast_parser.go: `ASTParser`, `NewASTParser`, `ParseDirectory`, `ParseFile`, `embedsSchema`, `extractModelDefinition`, `findMethod`, `collectAssignedExprs`, `resolveAssignedExpr`, `formatSelectorExpr`.
- ast_fields.go: `extractFields`, `extractFieldFromCall`, `findFieldBuilderInChain`, `extractOptionsFromChain`, `extractOptionsFromVariadicArgs`, `isFieldBuilder`, `mapFieldTypeToGoType`, `extractOptionFromMethod`, `extractDefaultValue`, `buildValidationTag`.
- ast_relations.go: `extractRelations`, `extractRelationFromExpr`, `extractRelationOptionsFromVariadicArgs`, `findRelationBuilderInChain`, `extractRelationOptionsFromChain`, `extractRelationOptionFromMethod`.
- ast_meta.go: `extractMeta`, `extractStringSliceFromExpr`, `extractIndexesFromExpr`, `extractConstraintsFromExpr`, `extractUniqueTogetherFromExpr`.
- ast_hooks.go: `extractHooks`, `extractHooksFromExpr`, `extractHooksFromCompositeLiteral`, `extractHooksFromCallChain`, `setHookFromBuilderMethod`, `setHookValue`, `extractHookReference`.
- ast_literals.go: `extractStringArg`, `extractIntArg`, `extractFloatArg`, `extractBoolArg`, `extractBoolArgAt`, `extractStringFromExpr`, `extractIntFromExpr`, `extractBoolFromExpr`.
Run the acceptance commands; everything must pass before step 2.

## Step 2 — table-driven option parsing (D8)
`extractOptionFromMethod` is 154 lines of a `switch methodName` where most cases store one argument under one key. Replace the repeated cases with a table: `var fieldOptionParsers = map[string]func(p *ASTParser, call *ast.CallExpr, options map[string]interface{})` built from small helpers such as `stringOption(key)`, `intOption(key)`, `floatOption(key)`, `boolOption(key)` / `trueOption(key)` for zero-arg flags. Keep cases with special logic as named functions in the table. The table is a package-level variable that is never mutated after initialization (that is allowed). `extractOptionFromMethod` becomes a lookup plus call, under 20 lines.

Before changing it, write a table-driven test `TestExtractOptionFromMethod_AllOptions` that parses a small model source (use `go/parser` on a string, the way existing codegen tests build ASTs — read them first) exercising every method name the switch handles, and asserts the resulting options map. Run it against the old code (it must pass), then refactor and keep it passing.

## Acceptance (from the module root)
    gofmt -l codegen                                  # empty
    go build ./codegen/... && go vet ./codegen/...
    staticcheck ./codegen/...                         # no output
    go test -race ./codegen/... ./cli/... -count=1
    wc -l codegen/ast*.go                             # every file under 600 lines

Report (max 20 lines): files and line counts, the new test, each command result.
