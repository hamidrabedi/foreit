package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/forgego/forge/admin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"examples/ecommerce/app/catalog"
	"examples/ecommerce/app/customers"
	"examples/ecommerce/app/inventory"
	"examples/ecommerce/app/marketing"
	"examples/ecommerce/app/orders"
)

//go:generate forge generate

func main() {
	ctx := context.Background()
	fmt.Println("🛒 Forge Ecommerce System")
	fmt.Println("=" + string(make([]byte, 50)))

	// Initialize Admin Site
	adminSite := admin.DefaultSite
	adminSite.Title = "Forge Ecommerce Admin"

	// Register application admins
	catalog.RegisterAdmin(ctx)
	orders.RegisterAdmin(ctx)
	customers.RegisterAdmin(ctx)
	inventory.RegisterAdmin(ctx)
	marketing.RegisterAdmin(ctx)

	// Create a new Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Admin handler (handles both REST API and UI)
	r.Mount("/admin", adminSite.Handler())

	// Health check
	r.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"healthy","message":"Forge Ecommerce Example is running"}`)
	})

	// Homepage
	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `
<!DOCTYPE html>
<html>
<head>
	<title>Forge Ecommerce</title>
	<style>
		body { font-family: Inter, system-ui, sans-serif; max-width: 900px; margin: 50px auto; padding: 20px; background: #0f172a; color: #f8fafc; }
		h1 { color: #38bdf8; border-bottom: 2px solid #334155; padding-bottom: 10px; }
		.section { margin: 20px 0; padding: 20px; background: #1e293b; border-radius: 12px; border: 1px solid #334155; box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1); }
		.section h2 { color: #7dd3fc; margin-top: 0; }
		a { color: #38bdf8; text-decoration: none; font-weight: 500; }
		a:hover { color: #7dd3fc; }
		ul { list-style: none; padding: 0; }
		li { margin: 10px 0; padding: 12px; background: #334155; border-radius: 8px; border-left: 4px solid #38bdf8; }
		.badge { background: #0ea5e9; color: white; padding: 2px 8px; border-radius: 6px; font-size: 11px; margin-left: 10px; font-weight: bold; }
		.btn-admin { display: inline-block; background: #0ea5e9; color: white; padding: 12px 24px; border-radius: 8px; font-weight: bold; margin-top: 20px; transition: all 0.2s; }
		.btn-admin:hover { background: #0284c7; transform: translateY(-1px); }
	</style>
</head>
<body>
	<h1>🛒 Forge Ecommerce</h1>
	<p><strong>Experience the next generation of Admin Frameworks</strong></p>
	
	<div class="section">
		<h2>🚀 Premium Admin Experience</h2>
		<p>The entire admin interface is now powered by our new, world-class system featuring:</p>
		<ul>
			<li>✨ <strong>High-Performance UI</strong> - Built with React 18, Vite, and Shadcn UI</li>
			<li>📊 <strong>Dynamic Dashboards</strong> - Auto-generated analytics and charts</li>
			<li>🛠️ <strong>Deep Customization</strong> - Full override engine for every component</li>
			<li>🔗 <strong>Smart Relations</strong> - Advanced M2M and FK handling</li>
		</ul>
		<a href="/admin/" class="btn-admin">Launch Admin Panel</a>
	</div>
	
	<div class="section">
		<h2>📊 Project Statistics</h2>
		<ul>
			<li>📈 <strong>Total Models:</strong> 29 (Across 5 applications)</li>
			<li>🔧 <strong>Deep Customization:</strong> 100% Component & Template Overrides</li>
			<li>🛠️ <strong>Automated Logic:</strong> AI-driven smart filtering and search</li>
		</ul>
	</div>

	<footer style="margin-top: 40px; text-align: center; color: #64748b; font-size: 0.875rem;">
		<p>Built with ❤️ using the new Forge World-Class Admin Framework</p>
	</footer>
</body>
</html>
		`)
	})

	port := ":8003"
	fmt.Printf("\n✨ Forge Ecommerce is Alive ✨\n")
	fmt.Printf("------------------------------\n")
	fmt.Printf("🏠 Homepage: http://localhost%s\n", port)
	fmt.Printf("🛠️  Premium Admin: http://localhost%s/admin/\n", port)
	fmt.Printf("------------------------------\n\n")

	log.Fatal(http.ListenAndServe(port, r))
}
