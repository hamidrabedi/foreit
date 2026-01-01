# Naming Documentation Index

Complete naming conventions, standards, and guidelines for the Forge framework.

---

## 📚 Documents Overview

### 1. **[NAMING_ARCHITECTURE.md](./NAMING_ARCHITECTURE.md)** - The Complete Guide
**Status**: Draft  
**Purpose**: Comprehensive naming conventions and design philosophy

**Contains**:
- Core naming principles
- Package, type, function, variable, interface, constant, and method naming
- Query Expression API design (Q, F, Where)
- Field Expression API design
- Framework-specific conventions
- Anti-patterns and best practices
- Detailed examples and justifications

**Read this**: For understanding the "why" behind naming decisions

---

### 2. **[NAMING_AUDIT.md](./NAMING_AUDIT.md)** - Current State & Migration
**Status**: Draft  
**Purpose**: Audit of current naming issues and migration plan

**Contains**:
- Critical issues (NewQ, NewFieldQueryExpr)
- Medium priority improvements
- File-by-file audit results
- 6-phase migration plan
- Compatibility layer design
- Testing strategy
- Risk assessment

**Read this**: For understanding what needs to change and how

---

### 3. **[NAMING_QUICK_REFERENCE.md](./NAMING_QUICK_REFERENCE.md)** - Developer Cheat Sheet
**Status**: Complete  
**Purpose**: Fast lookup for daily development

**Contains**:
- Quick decision tree
- Common patterns
- Anti-patterns (DON'T list)
- Package-specific conventions
- Real-world examples
- Migration checklist
- Quick lookup table

**Read this**: For day-to-day naming decisions

---

## 🎯 Quick Navigation

### By Role

**I'm a framework contributor**:
1. Start with [NAMING_QUICK_REFERENCE.md](./NAMING_QUICK_REFERENCE.md)
2. Review [NAMING_ARCHITECTURE.md](./NAMING_ARCHITECTURE.md) for details
3. Check [NAMING_AUDIT.md](./NAMING_AUDIT.md) before making changes

**I'm implementing a new feature**:
1. Use [NAMING_QUICK_REFERENCE.md](./NAMING_QUICK_REFERENCE.md) as guide
2. Reference [NAMING_ARCHITECTURE.md](./NAMING_ARCHITECTURE.md) for edge cases
3. Follow the decision tree in Quick Reference

**I'm reviewing a PR**:
1. Check [NAMING_QUICK_REFERENCE.md](./NAMING_QUICK_REFERENCE.md) anti-patterns
2. Verify consistency with [NAMING_ARCHITECTURE.md](./NAMING_ARCHITECTURE.md)
3. Ensure no new issues from [NAMING_AUDIT.md](./NAMING_AUDIT.md)

**I'm migrating old code**:
1. Read migration section in [NAMING_AUDIT.md](./NAMING_AUDIT.md)
2. Use compatibility layer
3. Follow migration checklist in [NAMING_QUICK_REFERENCE.md](./NAMING_QUICK_REFERENCE.md)

---

## 🔥 Critical Issues Summary

### Issue #1: `NewQ()` Function
**Problem**: Meaningless name, unnecessary wrapper  
**Solution**: Replace with `Q()` factory function  
**Files Affected**: ~15  
**Priority**: 🔴 Critical  
**Status**: Pending implementation

### Issue #2: `NewFieldQueryExpr()` Function
**Problem**: Too verbose, not user-friendly  
**Solution**: Replace with `Where()` or field methods  
**Files Affected**: ~25  
**Priority**: 🔴 Critical  
**Status**: Pending implementation

### Issue #3: `NewF()` Function
**Problem**: Planned but poorly named  
**Solution**: Use `F()` directly  
**Files Affected**: ~10  
**Priority**: 🔴 Critical  
**Status**: Prevention (not yet widely used)

---

## 📋 Key Principles (TL;DR)

### The Big 5

1. **Clarity over brevity** - Spell it out if there's any ambiguity
2. **Framework consistency** - Match Django/Laravel when it makes sense
3. **Go idioms first** - Follow Go conventions for exports, interfaces
4. **Self-documenting** - Code should read like prose
5. **Type safety** - Use generics, not name encoding

### Golden Rules

```
✅ DO:
- Use "New" for constructors that return concrete structs
- Use direct names for factories (Q, F, CharField)
- Follow Django naming for query operations
- Use full words in public APIs
- Keep receiver names consistent (qs, m, f)

❌ DON'T:
- Use meaningless abbreviations (NewQ, NewF)
- Add redundant prefixes (QuerySetQuerySet)
- Use Hungarian notation (strName, intCount)
- Create generic utils/helpers packages
- Abbreviate public API unnecessarily
```

---

## 📊 Migration Status

### Phase 1: Preparation (Current)
- [x] NAMING_ARCHITECTURE.md created
- [x] NAMING_AUDIT.md created  
- [x] NAMING_QUICK_REFERENCE.md created
- [ ] Team review
- [ ] Grep all usages
- [ ] Create compatibility shims

### Phase 2-6: Implementation (Planned)
See [NAMING_AUDIT.md - Migration Plan](./NAMING_AUDIT.md#migration-plan)

**Timeline**:
- v1.5.0: New API available, old API deprecated (Week 6)
- v1.6.0-1.9.0: Migration period (6 months)
- v2.0.0: Old API removed

---

## 🔍 Common Questions

### Q: Why replace `NewQ()` with `Q()`?

**A**: Three reasons:
1. **Framework consistency** - Django uses `Q()`, not `NewQ()`
2. **Clarity** - Q is meaningless, but Django developers know it
3. **Simplicity** - No need to wrap expressions that already work

See: [NAMING_ARCHITECTURE.md - Query Expression API](./NAMING_ARCHITECTURE.md#query-expression-api)

---

### Q: When should I use "New" prefix?

**A**: Use "New" when:
- Creating concrete struct instances
- Initialization requires validation/setup
- Returns the concrete struct (not just interface)

**Don't use "New" when**:
- Creating expressions/builders (factory pattern)
- Simple value construction
- Mimicking well-known API (Q, F from Django)

See: [NAMING_ARCHITECTURE.md - Function Naming](./NAMING_ARCHITECTURE.md#function-naming)

---

### Q: What about backwards compatibility?

**A**: We provide:
- 6-month deprecation period
- Compatibility layer (old functions call new ones)
- Clear migration guide
- Deprecation warnings in logs

See: [NAMING_AUDIT.md - Compatibility Layer](./NAMING_AUDIT.md#compatibility-layer)

---

### Q: How do I name a new feature?

**A**: Follow this decision tree:

```
1. Is it similar to Django/Laravel feature?
   → Use their naming (QuerySet, Manager, ViewSet)

2. Is it pure Go (no equivalent)?
   → Follow Go conventions (Reader, Writer, Builder)

3. Is it a query operation?
   → Use Django-style: Filter, Exclude, OrderBy

4. Is it a constructor?
   → Use New prefix: NewQuerySet, NewManager

5. Is it a factory?
   → Direct name: Q, F, CharField

6. Still unsure?
   → Ask in PR or check Quick Reference
```

See: [NAMING_QUICK_REFERENCE.md - Decision Tree](./NAMING_QUICK_REFERENCE.md#quick-decision-tree)

---

## 📖 Examples

### Before (Bad)

```go
// ❌ Meaningless abbreviations
q := orm.NewQ(nameField.Eq("John"))
f := orm.NewFieldQueryExpr("age", orm.OpGreater, 18)

// ❌ Too verbose
users, err := userManager.
    Filter(orm.NewQ(
        orm.NewFieldQueryExpr("name", orm.OpEquals, "John"),
    )).
    All(ctx)
```

### After (Good)

```go
// ✅ Clear and concise
users, err := userManager.
    Filter(User.Name.Eq("John")).
    All(ctx)

// ✅ Complex queries are readable
users, err := userManager.
    Filter(And(
        User.Name.Eq("John"),
        User.Age.Gt(18),
    )).
    All(ctx)

// ✅ SQL-like alternative
users, err := userManager.
    Filter(Where("name", OpEquals, "John")).
    All(ctx)
```

---

## 🎓 Learning Path

### Beginner (Just starting with Forge)
1. Read [NAMING_QUICK_REFERENCE.md](./NAMING_QUICK_REFERENCE.md)
2. Look at examples in your feature area
3. Ask questions in PR

### Intermediate (Contributing regularly)
1. Master [NAMING_QUICK_REFERENCE.md](./NAMING_QUICK_REFERENCE.md)
2. Read [NAMING_ARCHITECTURE.md](./NAMING_ARCHITECTURE.md) relevant sections
3. Help review PRs for naming consistency

### Advanced (Core maintainer)
1. Know all three documents thoroughly
2. Update [NAMING_AUDIT.md](./NAMING_AUDIT.md) as issues are found
3. Help design naming for new features

---

## 🛠️ Tools & Automation

### Linting Rules (Planned)

```go
// Add to .golangci.yml
linters-settings:
  gocritic:
    enabled-checks:
      - unnecessaryNew      # Flag unnecessary New prefix
      - abbreviations       # Flag unclear abbreviations
      - receiverNames       # Check receiver consistency
```

### Pre-commit Hooks (Planned)

```bash
# Check for deprecated functions
grep -r "NewQ(" . && echo "Use Q() instead of NewQ()"
grep -r "NewFieldQueryExpr(" . && echo "Use Where() instead"
```

### Migration Script (Planned)

```bash
# Automated migration helper
./scripts/migrate-naming.sh --check   # Dry run
./scripts/migrate-naming.sh --fix     # Auto-fix simple cases
```

---

## 📝 Contributing

### Adding New Names

Before adding any new exported name:

1. Check [NAMING_QUICK_REFERENCE.md](./NAMING_QUICK_REFERENCE.md) decision tree
2. Ensure it follows [NAMING_ARCHITECTURE.md](./NAMING_ARCHITECTURE.md) principles
3. Verify it doesn't conflict with [NAMING_AUDIT.md](./NAMING_AUDIT.md) known issues
4. Add examples to documentation
5. Get PR review focusing on naming

### Reporting Naming Issues

Found a naming problem?

1. Check if it's already in [NAMING_AUDIT.md](./NAMING_AUDIT.md)
2. If not, create an issue with:
   - Current name and usage
   - Why it's problematic
   - Proposed alternative
   - Impact assessment (how widely used)
3. Tag with `naming` label

### Proposing Changes

Want to change naming conventions?

1. Propose change in [NAMING_ARCHITECTURE.md](./NAMING_ARCHITECTURE.md)
2. Add to [NAMING_AUDIT.md](./NAMING_AUDIT.md) with impact assessment
3. Create RFC (Request for Comments) issue
4. Get consensus before implementation

---

## 🔗 Related Documentation

- [API_REFERENCE.md](./API_REFERENCE.md) - Complete API documentation
- [ARCHITECTURE.md](./ARCHITECTURE.md) - Framework architecture
- [GETTING_STARTED.md](./GETTING_STARTED.md) - Getting started guide
- [CONTRIBUTING.md](./CONTRIBUTING.md) - Contribution guidelines

---

## 📅 Changelog

### 2026-01-01 - Initial Release
- Created NAMING_ARCHITECTURE.md
- Created NAMING_AUDIT.md
- Created NAMING_QUICK_REFERENCE.md
- Created NAMING_INDEX.md
- Identified critical issues (NewQ, NewFieldQueryExpr)
- Proposed migration plan

---

## ✅ Status Dashboard

| Area | Status | Priority | Owner |
|------|--------|----------|-------|
| Documentation | ✅ Complete | - | @team |
| Q() Implementation | ⏳ Pending | 🔴 Critical | TBD |
| F() Implementation | ⏳ Pending | 🔴 Critical | TBD |
| Where() Implementation | ⏳ Pending | 🔴 Critical | TBD |
| Compatibility Layer | ⏳ Pending | 🟡 Medium | TBD |
| Internal Migration | ⏳ Pending | 🟡 Medium | TBD |
| Documentation Update | ⏳ Pending | 🟡 Medium | TBD |
| v1.5.0 Release | ⏳ Pending | 🟡 Medium | TBD |

---

## 🎯 Next Actions

### Immediate (This Week)
- [ ] Team review of naming documents
- [ ] Decide on migration timeline
- [ ] Assign owners for implementation

### Short Term (This Month)
- [ ] Implement Q(), F(), Where() functions
- [ ] Add deprecation warnings
- [ ] Create compatibility tests
- [ ] Begin internal migration

### Long Term (This Quarter)
- [ ] Complete internal migration
- [ ] Release v1.5.0 with new API
- [ ] Update all documentation
- [ ] Announce to community

---

**Last Updated**: 2026-01-01  
**Maintained By**: @team  
**Next Review**: Before v1.5.0 release
