import AdminLayout from "../components/layout/AdminLayout";
import { StatsCard } from "../components/widgets/StatsCard";
import {
  Users,
  Package,
  ShoppingCart,
  BarChart3,
  ArrowRight,
  ChevronRight,
  Calendar
} from "lucide-react";
import { Button } from "../components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "../components/ui/card";
import { cn } from "../lib/utils";
import { widgetRegistry } from "../lib/widgets";
import { useAdminMetaStore } from "../store/adminStore";
import type { WidgetMetadata } from "../api/types";

export default function DashboardPage() {
  const { dashboard, plugins } = useAdminMetaStore();
  const widgetList = buildWidgetList(dashboard?.widgets, plugins);

  // Mock data for charts
  const salesData = [
    { value: 400 }, { value: 300 }, { value: 600 },
    { value: 800 }, { value: 500 }, { value: 900 }, { value: 1100 }
  ];

  return (
    <AdminLayout>
      <div className="space-y-8 animate-in fade-in duration-500">
        {/* Header */}
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <h1 className="text-4xl font-black tracking-tight text-foreground/90">
              Insights <span className="text-primary/60">Overview</span>
            </h1>
            <p className="text-muted-foreground flex items-center gap-2 mt-1">
              <Calendar className="h-3.5 w-3.5" />
              Checking performance for the last 30 days
            </p>
          </div>
          <div className="flex items-center gap-3">
             <Button variant="outline" size="sm" className="glass-lite border-border/50">
               <BarChart3 className="h-4 w-4 mr-2" />
               View Reports
             </Button>
             <Button size="sm" className="shadow-lg shadow-primary/20">
               <Plus className="h-4 w-4 mr-2" />
               Add Entry
             </Button>
          </div>
        </div>

        {/* Widgets */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {widgetList.length > 0 ? (
            widgetList.map((widget) => (
              <DashboardWidget key={widget.type + widget.title} widget={widget} />
            ))
          ) : (
            <>
              <StatsCard
                title="Total Revenue"
                value="$128,430"
                trend={12.5}
                icon={BarChart3}
                color="primary"
                chartData={salesData}
              />
              <StatsCard
                title="Active Customers"
                value="2,450"
                trend={8.2}
                icon={Users}
                color="success"
                chartData={salesData.map(v => ({ value: v.value * 0.8 }))}
              />
              <StatsCard
                title="Pending Orders"
                value="43"
                trend={-2.4}
                icon={ShoppingCart}
                color="warning"
                chartData={salesData.map(v => ({ value: 1000 - v.value }))}
              />
              <StatsCard
                title="Inventory Value"
                value="$840k"
                trend={4.1}
                icon={Package}
                color="primary"
                chartData={salesData}
              />
            </>
          )}
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Main Chart Area */}
          <Card className="lg:col-span-2 glass-lite border-border/50 shadow-sm overflow-hidden">
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-lg font-bold">Growth Trends</CardTitle>
              <Button variant="ghost" size="sm" className="text-xs">
                Details <ArrowRight className="h-3 w-3 ml-1" />
              </Button>
            </CardHeader>
            <CardContent>
              <div className="h-[300px] flex items-center justify-center text-muted-foreground border-2 border-dashed border-border/20 rounded-xl bg-muted/5">
                 <div className="text-center">
                    <BarChart3 className="h-12 w-12 mx-auto mb-2 opacity-20" />
                    <p className="text-sm font-medium">Main Performance Visualization</p>
                    <p className="text-xs opacity-60">Interactive Recharts Layer</p>
                 </div>
              </div>
            </CardContent>
          </Card>

          {/* Activity Sidebar */}
          <Card className="glass-lite border-border/50 shadow-sm overflow-hidden">
             <CardHeader>
                <CardTitle className="text-lg font-bold">Recent Events</CardTitle>
             </CardHeader>
             <CardContent className="p-0">
                <div className="divide-y divide-border/20">
                  {[
                    { label: "New Order", time: "2m ago", desc: "Order #8492 by Sarah J.", icon: ShoppingCart, color: "text-primary" },
                     { label: "Stock Alert", time: "15m ago", desc: "Low stock for 'MacBook Air'", icon: Package, color: "text-amber-500" },
                    { label: "User Joined", time: "1h ago", desc: "A new customer has registered", icon: Users, color: "text-emerald-500" },
                    { label: "Report Ready", time: "3h ago", desc: "Monthly export is now available", icon: BarChart3, color: "text-primary" },
                  ].map((item, idx) => (
                    <div key={idx} className="p-4 flex gap-4 hover:bg-muted/30 transition-colors group cursor-pointer">
                      <div className={cn("mt-1", item.color)}>
                        <item.icon className="h-4 w-4" />
                      </div>
                      <div className="flex-1">
                        <div className="flex items-center justify-between">
                          <p className="text-xs font-black uppercase tracking-widest">{item.label}</p>
                          <span className="text-[10px] text-muted-foreground">{item.time}</span>
                        </div>
                        <p className="text-sm text-foreground/80 font-medium leading-tight mt-0.5">{item.desc}</p>
                      </div>
                      <ChevronRight className="h-4 w-4 self-center text-muted-foreground opacity-0 group-hover:opacity-100 transition-all -translate-x-2 group-hover:translate-x-0" />
                    </div>
                  ))}
                </div>
                <div className="p-4 bg-muted/20 border-t border-border/20 text-center">
                   <Button variant="link" size="sm" className="text-xs font-bold text-primary">
                     View All Activity
                   </Button>
                </div>
             </CardContent>
          </Card>
        </div>
      </div>
    </AdminLayout>
  );
}

function Plus(props: any) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <path d="M5 12h14" />
      <path d="M12 5v14" />
    </svg>
  )
}

function buildWidgetList(dashboardWidgets?: WidgetMetadata[], plugins: any[] = []) {
  const pluginWidgets = plugins.flatMap((plugin) => plugin.widgets || []);
  const combined = [...(dashboardWidgets || []), ...pluginWidgets];
  const seen = new Set<string>();
  return combined.filter((widget) => {
    const key = `${widget.type}:${widget.title}`;
    if (seen.has(key)) return false;
    seen.add(key);
    return true;
  });
}

function DashboardWidget({ widget }: { widget: WidgetMetadata }) {
  const Component = widgetRegistry.get(widget.type);
  if (!Component) {
    return (
      <Card className="glass-lite border-border/50 shadow-sm">
        <CardHeader>
          <CardTitle className="text-sm font-bold">{widget.title}</CardTitle>
        </CardHeader>
        <CardContent className="text-xs text-muted-foreground">
          Missing widget component: {widget.type}
        </CardContent>
      </Card>
    );
  }

  return <Component data={null} config={{ id: widget.type, ...widget }} />;
}
