---
sidebar_position: 2
description: Auto-generated Django-style admin interface for forge models. Full CRUD, filtering, search, and more - all type-safe.
keywords:
  - forge admin
  - django admin go
  - auto-generated admin
  - crud interface
  - admin panel
image: /img/forge-social-card.jpg
---

# Admin System

The admin system gives you a complete CRUD interface for your models, automatically. If you've used Django's admin, you know how powerful this is - but forge's version is type-safe and doesn't use reflection magic.

## Getting Started

Register a model and you're done:

```go
import "github.com/forgego/forge/admin"

func main() {
    // Register your models
    admin.RegisterModel(&models.User{})
    admin.RegisterModel(&models.Post{})
    admin.RegisterModel(&models.Category{})
    
    // Start your server
    server.Run()
}
```

Visit `/admin/` and you'll see a full admin interface for all your registered models.

## Basic Configuration

Customize how your model appears in the admin:

```go
func (User) AdminConfig() admin.Config {
    return admin.Config{
        ListDisplay: []string{"username", "email", "is_active", "created_at"},
        ListFilter: []string{"is_active", "created_at"},
        SearchFields: []string{"username", "email"},
        ListPerPage: 25,
        Ordering: []string{"-created_at"},
    }
}
```

## List Views

The list view is where you'll spend most of your time. forge gives you:

### Display Fields
Control what columns show up:

```go
ListDisplay: []string{
    "username",           // Simple field
    "email",             // Simple field
    "is_active",         // Boolean field (shows as checkbox)
    "post_count",        // Custom method
    "created_at",        // Date field
}
```

### Custom List Methods
Add computed columns:

```go
func (u *User) PostCount() int {
    count, _ := Post.Objects.Filter(Post.Fields.UserID.Equals(u.ID)).Count(context.Background())
    return count
}
```

### Search and Filtering
```go
ListFilter: []string{
    "is_active",         // Boolean filter
    "created_at",        // Date range filter
    "role",             // Choice filter
}

SearchFields: []string{
    "username",         // Search in username
    "email",           // Search in email
}
```

### Pagination and Ordering
```go
ListPerPage: 50,                    // Items per page
Ordering: []string{"-created_at"},  // Default ordering
```

## Form Views

### Fieldsets
Organize your form fields into logical groups:

```go
Fieldsets: []admin.Fieldset{
    {
        Title: "Basic Information",
        Fields: []string{"username", "email", "is_active"},
    },
    {
        Title: "Profile",
        Fields: []string{"bio", "avatar"},
        Classes: []string{"collapse"}, // Collapsed by default
    },
}
```

### Readonly Fields
Make fields read-only in the admin:

```go
ReadonlyFields: []string{"id", "created_at", "updated_at"},
```

### Prepopulated Fields
Auto-fill fields based on other fields:

```go
PrepopulatedFields: map[string][]string{
    "slug": {"title"},  // Fill slug from title
}
```

## Inlines

Edit related models inline:

```go
Inlines: []admin.Inline{
    {
        Model: &Post{},
        Extra: 3,           // Show 3 empty forms
        MinNum: 1,          // Require at least 1
        MaxNum: 10,         // Maximum 10
    },
}
```

Now when you edit a user, you'll see their posts right there on the same page.

## Actions

Add bulk actions to your list views:

```go
Actions: []admin.Action{
    {
        Name: "make_active",
        Description: "Mark selected users as active",
        Handler: func(queryset admin.QuerySet, form admin.FormData) error {
            return queryset.Update(map[string]interface{}{
                "is_active": true,
            })
        },
    },
    {
        Name: "send_welcome_email",
        Description: "Send welcome email to selected users",
        Handler: func(queryset admin.QuerySet, form admin.FormData) error {
            users, _ := queryset.All()
            for _, user := range users {
                sendWelcomeEmail(user.(*User))
            }
            return nil
        },
    },
}
```

## Custom Filters

Create your own filters:

```go
type ActiveUsersFilter struct {
    admin.BooleanFilter
}

func (f *ActiveUsersFilter) Filter(queryset admin.QuerySet, value interface{}) admin.QuerySet {
    if value.(bool) {
        return queryset.Filter(User.Fields.IsActive.Equals(true))
    }
    return queryset
}

func (User) AdminConfig() admin.Config {
    return admin.Config{
        ListFilter: []string{"is_active", "created_at"},
        CustomFilters: map[string]admin.Filter{
            "active_users": &ActiveUsersFilter{},
        },
    }
}
```

## Widgets

Control how fields are displayed in forms:

```go
FormWidgets: map[string]admin.Widget{
    "bio": admin.TextareaWidget{
        Rows: 10,
        Cols: 80,
    },
    "birth_date": admin.DateWidget{
        Format: "2006-01-02",
    },
    "avatar": admin.ImageField{
        UploadTo: "avatars/",
    },
    "role": admin.SelectWidget{
        Choices: []admin.Choice{
            {Value: "admin", Label: "Administrator"},
            {Value: "user", Label: "Regular User"},
            {Value: "moderator", Label: "Moderator"},
        },
    },
}
```

## Validation

Add custom validation to your admin forms:

```go
func (User) AdminClean(form admin.FormData) error {
    password := form.Get("password")
    confirm := form.Get("password_confirm")
    
    if password != confirm {
        return errors.New("Passwords don't match")
    }
    
    if len(password) < 8 {
        return errors.New("Password must be at least 8 characters")
    }
    
    return nil
}
```

## Permissions

Control who can do what:

```go
func (User) AdminConfig() admin.Config {
    return admin.Config{
        Permissions: admin.Permissions{
            Add:    "users.add_user",
            Change: "users.change_user",
            Delete: "users.delete_user",
            View:   "users.view_user",
        },
    }
}
```

## Custom Templates

Override admin templates for full control:

```go
Templates: map[string]string{
    "change_list.html":    "admin/user_change_list.html",
    "change_form.html":    "admin/user_change_form.html",
    "delete_confirmation.html": "admin/user_delete.html",
}
```

## Export

Add export functionality:

```go
ExportFormats: []string{"csv", "json", "xlsx"},
ExportFields: []string{"username", "email", "is_active", "created_at"},
```

Now users can export filtered data from the list view.

## Advanced Features

### Custom QuerySets
Override the default queryset:

```go
func (User) AdminGetQueryset(request *http.Request) admin.QuerySet {
    qs := User.Objects.All()
    
    // Only show active users to non-admins
    if !isAdmin(request) {
        qs = qs.Filter(User.Fields.IsActive.Equals(true))
    }
    
    return qs
}
```

### Custom Save Methods
Control how objects are saved:

```go
func (User) AdminSaveModel(request *http.Request, obj interface{}, form admin.FormData, change bool) error {
    user := obj.(*User)
    
    // Hash password if it's being set
    if form.HasChanged("password") {
        hashed, _ := bcrypt.GenerateFromPassword([]byte(form.Get("password")), bcrypt.DefaultCost)
        user.PasswordHash = string(hashed)
    }
    
    return User.Objects.Save(user)
}
```

### Custom URLs
Add custom admin URLs:

```go
func (User) AdminGetUrls() []admin.Url {
    return []admin.Url{
        {
            Path:     "send-email/",
            View:     sendEmailView,
            Name:     "send_email",
        },
        {
            Path:     "stats/",
            View:     userStatsView,
            Name:     "user_stats",
        },
    }
}
```

## Security

The admin system includes security features out of the box:

- **CSRF Protection** - All forms are CSRF protected
- **XSS Prevention** - Output is properly escaped
- **Permission Checking** - Users must have appropriate permissions
- **Audit Logging** - All admin actions are logged
- **Secure Defaults** - Safe defaults for all configurations

## Performance

The admin is built to perform well:

- **Efficient Queries** - Uses select_related and prefetch_related
- **Pagination** - Only loads the data you need
- **Caching** - Caches expensive operations
- **Lazy Loading** - Loads data only when needed

## Best Practices

1. **Keep it simple** - Don't over-customize your admin
2. **Use permissions** - Properly secure your admin
3. **Optimize queries** - Use select_related for foreign keys
4. **Test your admin** - Write tests for custom actions and validation
5. **Document custom features** - Help your users understand custom functionality

## Next Steps

- [REST API Framework](/docs/features/api-framework) - Build APIs for your models
- [Filter System](/docs/features/filter-system) - Advanced filtering capabilities
- [Examples](/docs/examples/blog) - See the admin in action
