---
description: Purge legacy Admin directories safely
---

1. Remove legacy directories in `forge/admin`:
   - `advanced`
   - `fields`
   - `filter`
   - `handlers`
   - `orm`
   - `permissions`
   - `relations`
   - `schema`
   - `validation`
   - `TODO.md`

// turbo
2. Execute deletion command
    ```bash
    rm -rf forge/admin/advanced forge/admin/fields forge/admin/filter forge/admin/handlers forge/admin/orm forge/admin/permissions forge/admin/relations forge/admin/schema forge/admin/validation forge/admin/TODO.md
    ```
