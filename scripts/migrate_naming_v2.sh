#!/bin/bash
# Automated migration script for v1.x to v2.0 naming changes
#
# This script automatically updates code to use the new naming conventions:
# - Rebind() → RebindPlaceholders()
# - GetRegistry() → Global()
# - GetFieldAccessor() → FieldAccessor()
# - GetPaginationParams() → ParsePaginationParams()
# - NewOrderField() → Asc() or Desc() (requires manual review)

set -e

echo "=== Forge Framework Naming Migration v2.0 ==="
echo ""
echo "This script will update your code to use the new naming conventions."
echo "Please review the changes before committing."
echo ""

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "Warning: Not in a git repository. Changes will be made directly."
    read -p "Continue? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
else
    echo "Git repository detected. Creating backup branch..."
    git checkout -b backup-before-naming-migration-$(date +%Y%m%d-%H%M%S) 2>/dev/null || true
    git checkout - 2>/dev/null || true
fi

# Find all Go files (excluding vendor and test files for now)
FILES=$(find . -type f -name "*.go" -not -path "./vendor/*" -not -path "./.git/*")

echo "Found $(echo "$FILES" | wc -l) Go files to process"
echo ""

# Fix 1: Rebind() → RebindPlaceholders()
echo "1. Updating Rebind() → RebindPlaceholders()..."
echo "$FILES" | xargs sed -i.bak \
    -e 's/\.Rebind(/\.RebindPlaceholders(/g' \
    -e 's/database\.Rebind(/database.RebindPlaceholders(/g' \
    -e 's/db\.Rebind(/db.RebindPlaceholders(/g' \
    -e 's/r\.db\.Rebind(/r.db.RebindPlaceholders(/g' 2>/dev/null || true
echo "   ✓ Done"

# Fix 2: GetRegistry() → Global() (only in registry package context)
echo "2. Updating registry.GetRegistry() → registry.Global()..."
echo "$FILES" | xargs sed -i.bak \
    -e 's/registry\.GetRegistry()/registry.Global()/g' 2>/dev/null || true
echo "   ✓ Done"

# Fix 3: GetFieldAccessor() → FieldAccessor()
echo "3. Updating GetFieldAccessor() → FieldAccessor()..."
echo "$FILES" | xargs sed -i.bak \
    -e 's/\.GetFieldAccessor()/\.FieldAccessor()/g' \
    -e 's/manager\.GetFieldAccessor()/manager.FieldAccessor()/g' 2>/dev/null || true
echo "   ✓ Done"

# Fix 4: GetPaginationParams() → ParsePaginationParams()
echo "4. Updating GetPaginationParams() → ParsePaginationParams()..."
echo "$FILES" | xargs sed -i.bak \
    -e 's/GetPaginationParams(/ParsePaginationParams(/g' \
    -e 's/api\.GetPaginationParams(/api.ParsePaginationParams(/g' 2>/dev/null || true
echo "   ✓ Done"

# Fix 5: NewOrderField() - Manual review needed
echo ""
echo "5. NewOrderField() requires manual review:"
echo "   Search for: NewOrderField("
echo "   Replace:"
echo "     - NewOrderField(field, true)  → Asc(field)"
echo "     - NewOrderField(field, false) → Desc(field)"
echo ""
echo "   Files that may need updates:"
grep -rn "NewOrderField(" --include="*.go" . 2>/dev/null | head -10 || echo "   (none found)"

# Clean up backup files
echo ""
echo "Cleaning up backup files..."
find . -name "*.bak" -type f -delete 2>/dev/null || true
echo "   ✓ Done"

echo ""
echo "=== Migration Complete ==="
echo ""
echo "Next steps:"
echo "1. Review the changes: git diff"
echo "2. Manually update NewOrderField() calls (see above)"
echo "3. Run tests: go test ./..."
echo "4. Build: go build ./..."
echo "5. Commit changes"
echo ""
echo "For detailed migration guide, see: docs/MIGRATION_V1_TO_V2.md"
