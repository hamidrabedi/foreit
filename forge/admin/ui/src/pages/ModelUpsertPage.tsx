import { useParams, useNavigate } from "@tanstack/react-router";
import {
  useModelMetadata,
  useModelDetail,
  useCreateObject,
  useUpdateObject,
  useModelHistory,
} from "../api/hooks/adminHooks";
import { Button } from "../components/ui/button";
import { Input } from "../components/ui/input";
import { Checkbox } from "../components/ui/checkbox";
import { SearchableSelect } from "../components/ui/searchable-select";
import {
  Card,
  CardContent,
} from "../components/ui/card";
import { Loader2, ArrowLeft, Save, X, UploadCloud, Eye } from "lucide-react";
import { useState, useEffect } from "react";
import AdminLayout from "../components/layout/AdminLayout";
import { useUIComponent } from "../hooks/useUIComponent";
import { cn } from "../lib/utils";
import { useToast } from "../hooks/use-toast";
import { adminAPI } from "../api/client";

interface ModelFormPageProps {
  mode: "create" | "edit";
}

export default function ModelFormPage({ mode }: ModelFormPageProps) {
  const params = useParams({ strict: false }) as any;
  const { model, id } = params; // 'model' and 'id' are used to derive modelName and objectId
  const navigate = useNavigate();
  const { toast } = useToast();
  const [formData, setFormData] = useState<Record<string, any>>({});
  const [uploadingFields, setUploadingFields] = useState<Record<string, boolean>>({});

  const modelName = model as string;
  const objectId = id as string || "";

  const { data: metadata, isLoading: metaLoading } =
    useModelMetadata(modelName);

  const { data: objectData, isLoading: objectLoading } = useModelDetail(
    modelName,
    objectId,
    {
      enabled: mode === "edit" && !!objectId,
    } as any
  );

  const createMutation = useCreateObject(modelName);
  const updateMutation = useUpdateObject(modelName);
  const isViewOnly = mode === "edit" && !!metadata && !metadata.permissions.change;

  useEffect(() => {
    if (objectData && mode === "edit") {
      setFormData(objectData);
    }
  }, [objectData, mode]);

  // Resolve Overrides
  const FormHeader = useUIComponent(metadata?.ui_overrides?.['form.header'] || "", 'div');
  const FormFooter = useUIComponent(metadata?.ui_overrides?.['form.footer'] || "", 'div');

  const normalizeFormData = () => {
    if (!metadata) return formData;
    const normalized = { ...formData };
    metadata.fields
      .filter((field: any) => field.type === "json")
      .forEach((field: any) => {
        const raw = normalized[field.name];
        if (typeof raw === "string") {
          const trimmed = raw.trim();
          if (trimmed === "") {
            normalized[field.name] = null;
            return;
          }
          try {
            normalized[field.name] = JSON.parse(trimmed);
          } catch (err) {
            throw new Error(`Invalid JSON in ${field.label || field.name}`);
          }
        }
      });
    return normalized;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (isViewOnly) {
      return;
    }

    try {
      const normalizedData = normalizeFormData();
      if (mode === "create") {
        await createMutation.mutateAsync(normalizedData);
        toast({
          title: "Success",
          description: `${metadata?.verbose_name} created successfully`,
        });
      } else {
        await updateMutation.mutateAsync({ id: objectId, data: normalizedData });
        toast({
          title: "Success",
          description: `${metadata?.verbose_name} updated successfully`,
        });
      }
      navigate({ to: `/${modelName}` });
    } catch (error) {
      console.error("Failed to save:", error);
      const errorMsg =
        (error as any)?.response?.data?.message ||
        (error instanceof Error ? error.message : "") ||
        "Failed to save changes. Please try again.";
      toast({
        title: "Error",
        description: errorMsg,
        variant: "destructive",
      });
    }
  };

  const handleChange = (name: string, value: any) => {
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleFileUpload = async (field: any, file: File | null) => {
    if (!file) return;
    setUploadingFields((prev) => ({ ...prev, [field.name]: true }));
    try {
      const response = await adminAPI.uploadFile(modelName, file);
      handleChange(field.name, response.url);
      toast({
        title: "Upload complete",
        description: `${file.name} uploaded successfully`,
      });
    } catch (err) {
      toast({
        title: "Upload failed",
        description: "Unable to upload file. Please try again.",
        variant: "destructive",
      });
    } finally {
      setUploadingFields((prev) => ({ ...prev, [field.name]: false }));
    }
  };

  if (metaLoading || (mode === "edit" && objectLoading)) {
    return (
      <AdminLayout>
        <div className="flex items-center justify-center h-full">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
        </div>
      </AdminLayout>
    );
  }

  if (!metadata) return null;

  if (mode === "create" && !metadata.permissions.add) {
    return (
      <AdminLayout>
        <div className="flex items-center justify-center h-full">
          <div className="text-center">
            <h2 className="text-2xl font-bold text-destructive">Permission Denied</h2>
            <p className="text-muted-foreground mt-2">You do not have permission to add {metadata.verbose_name}.</p>
          </div>
        </div>
      </AdminLayout>
    );
  }

  if (mode === "edit" && !metadata.permissions.view) {
    return (
      <AdminLayout>
        <div className="flex items-center justify-center h-full">
          <div className="text-center">
            <h2 className="text-2xl font-bold text-destructive">Permission Denied</h2>
            <p className="text-muted-foreground mt-2">You do not have permission to view {metadata.verbose_name}.</p>
          </div>
        </div>
      </AdminLayout>
    );
  }

  const readOnlyFields = new Set<string>(metadata?.read_only_fields || []);

  const renderField = (field: any) => {
    if (field.read_only && mode === 'create') return null;

    const value = formData[field.name] || "";
    const isReadOnly = field.read_only || readOnlyFields.has(field.name);
    const isDisabled = isReadOnly || isViewOnly;
    const fieldOverride = metadata.ui_overrides?.[`field.${field.name}`];
    const CustomField = useUIComponent(fieldOverride || "", null as any);

    if (CustomField) {
      return (
        <CustomField 
          field={field} 
          value={value} 
          onChange={(val: any) => handleChange(field.name, val)} 
          metadata={metadata}
        />
      );
    }

    if (field.widget === "password") {
      return (
        <Input
          id={field.name}
          type="password"
          value={value}
          onChange={(e) => handleChange(field.name, e.target.value)}
          className="rounded-lg border-border/50 bg-background/50 focus-visible:ring-primary/20 focus-visible:border-primary transition-all"
          required={field.required}
          disabled={isDisabled}
        />
      );
    }

    if (field.widget === "rich_text") {
      return (
        <textarea
          id={field.name}
          value={value}
          onChange={(e) => handleChange(field.name, e.target.value)}
          className="flex min-h-[160px] w-full rounded-lg border border-border/50 bg-background/50 px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/20 focus-visible:border-primary disabled:cursor-not-allowed disabled:opacity-50 transition-all shadow-inner"
          required={field.required}
          disabled={isDisabled}
        />
      );
    }

    if (field.widget === "file" || field.widget === "image") {
      const previewUrl = typeof value === "string" ? value : "";
      return (
        <div className="space-y-3">
          {previewUrl && field.widget === "image" && (
            <div className="rounded-lg border border-border/50 overflow-hidden max-w-xs">
              <img src={previewUrl} alt={field.label} className="w-full h-auto" />
            </div>
          )}
          {previewUrl && field.widget !== "image" && (
            <a
              href={previewUrl}
              target="_blank"
              rel="noreferrer"
              className="text-xs text-primary underline"
            >
              View uploaded file
            </a>
          )}
          <div className="flex items-center gap-3">
            <Input
              id={field.name}
              type="file"
              accept={field.accept || (field.widget === "image" ? "image/*" : undefined)}
              onChange={(e) => handleFileUpload(field, e.target.files?.[0] || null)}
              className="rounded-lg border-border/50 bg-background/50 focus-visible:ring-primary/20 focus-visible:border-primary transition-all"
              disabled={isDisabled}
            />
            {uploadingFields[field.name] && (
              <UploadCloud className="h-4 w-4 animate-pulse text-muted-foreground" />
            )}
          </div>
        </div>
      );
    }

    switch (field.type) {
      case "boolean":
        return (
          <div className="flex items-center space-x-3 p-3 rounded-lg border border-border/50 bg-muted/20 hover:bg-muted/30 transition-colors">
            <Checkbox
              id={field.name}
              checked={!!value}
              onChange={(e) => handleChange(field.name, e.target.checked)}
              disabled={isDisabled}
            />
            <label
              htmlFor={field.name}
              className="text-sm font-semibold leading-none cursor-pointer select-none"
            >
              {field.label}
            </label>
          </div>
        );

      case "text":
        if (field.widget === "textarea") {
          return (
            <textarea
              id={field.name}
              value={value}
              onChange={(e) => handleChange(field.name, e.target.value)}
              className="flex min-h-[120px] w-full rounded-lg border border-border/50 bg-background/50 px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/20 focus-visible:border-primary disabled:cursor-not-allowed disabled:opacity-50 transition-all"
              required={field.required}
            />
          );
        }
        return (
          <Input
            id={field.name}
            value={value}
            onChange={(e) => handleChange(field.name, e.target.value)}
            className="rounded-lg border-border/50 bg-background/50 focus-visible:ring-primary/20 focus-visible:border-primary transition-all"
            required={field.required}
            disabled={isDisabled}
          />
        );

      case "integer":
      case "float":
      case "decimal":
        return (
          <Input
            type="number"
            id={field.name}
            value={value}
            onChange={(e) =>
              handleChange(field.name, parseFloat(e.target.value))
            }
            className="rounded-lg border-border/50 bg-background/50 focus-visible:ring-primary/20 focus-visible:border-primary transition-all"
            required={field.required}
            disabled={isDisabled}
          />
        );

      case "date":
        return (
          <Input
            type="date"
            id={field.name}
            value={value}
            onChange={(e) => handleChange(field.name, e.target.value)}
            className="rounded-lg border-border/50 bg-background/50 focus-visible:ring-primary/20 focus-visible:border-primary transition-all"
            required={field.required}
            disabled={isDisabled}
          />
        );

      case "datetime":
        return (
          <Input
            type="datetime-local"
            id={field.name}
            value={value}
            onChange={(e) => handleChange(field.name, e.target.value)}
            className="rounded-lg border-border/50 bg-background/50 focus-visible:ring-primary/20 focus-visible:border-primary transition-all"
            required={field.required}
            disabled={isDisabled}
          />
        );

      case "choice":
        return (
          <select
            id={field.name}
            value={value}
            onChange={(e) => handleChange(field.name, e.target.value)}
            className="flex h-10 w-full rounded-lg border border-border/50 bg-background/50 px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/20 focus-visible:border-primary disabled:cursor-not-allowed disabled:opacity-50 transition-all appearance-none"
            required={field.required}
              disabled={isDisabled}
          >
            <option value="">Select {field.label}</option>
            {field.choices?.map((choice: any) => (
              <option key={choice.value} value={choice.value}>
                {choice.label}
              </option>
            ))}
          </select>
        );

      case "password":
        return (
          <Input
            type="password"
            id={field.name}
            value={value}
            onChange={(e) => handleChange(field.name, e.target.value)}
            className="rounded-lg border-border/50 bg-background/50 focus-visible:ring-primary/20 focus-visible:border-primary transition-all"
            required={field.required}
            autoComplete="new-password"
            disabled={isDisabled}
          />
        );

      case "json":
        return (
          <textarea
            id={field.name}
            value={typeof value === 'object' ? JSON.stringify(value, null, 2) : value}
            onChange={(e) => {
              const val = e.target.value;
              handleChange(field.name, val);
              // Basic validation hint could go here
            }}
            className="flex min-h-[150px] w-full rounded-lg border border-border/50 bg-background/50 px-3 py-2 text-xs font-mono ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/20 focus-visible:border-primary disabled:cursor-not-allowed disabled:opacity-50 transition-all shadow-inner"
            placeholder="{}"
            required={field.required}
            disabled={isDisabled}
          />
        );

      case "foreign_key":
      case "one_to_one":
        return (
          <SearchableSelect
            model={field.related_model}
            value={value}
            onChange={(val) => handleChange(field.name, val)}
            placeholder={`Select ${field.label}...`}
            required={field.required}
            disabled={isDisabled}
          />
        );

      case "many_to_many":
        const m2mValue = Array.isArray(value) ? value : [];
        return (
          <div className="space-y-4">
            <div className="flex flex-wrap gap-2 mb-2">
              {m2mValue.map((item: any, idx: number) => (
                <div 
                  key={item.id || idx} 
                  className="flex items-center gap-1.5 px-3 py-1 rounded-full bg-primary/10 text-primary border border-primary/20 text-xs font-medium animate-in fade-in zoom-in-95"
                >
                  <span>{item.name || item.title || item.label || `ID: ${item.id || item}`}</span>
                  <button
                    type="button"
                    onClick={() => {
                      const newValue = m2mValue.filter((_: any, i: number) => i !== idx);
                      handleChange(field.name, newValue);
                    }}
                    className="hover:text-destructive transition-colors"
                  >
                    <X className="h-3 w-3" />
                  </button>
                </div>
              ))}
            </div>
              <SearchableSelect
                model={field.related_model}
                value={null}
                onChange={(val) => {
                  if (!val) return;
                  // Avoid duplicates
                  if (m2mValue.some((v: any) => (typeof v === 'object' ? v.id : v) === val)) return;
                  handleChange(field.name, [...m2mValue, val]);
                }}
                placeholder={`Add ${field.label}...`}
                disabled={isDisabled}
              />
          </div>
        );

      default:
        return (
          <Input
            id={field.name}
            value={value}
            onChange={(e) => handleChange(field.name, e.target.value)}
            className="rounded-lg border-border/50 bg-background/50 focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all outline-none"
            required={field.required}
              disabled={isDisabled}
          />
        );
    }
  };

  const fieldsetGroups = metadata.fieldsets || [];
  const fieldsetFields = new Set<string>(
    fieldsetGroups.flatMap((fieldset: any) => fieldset.fields || [])
  );

  const { data: historyData } = useModelHistory(modelName, objectId, {
    enabled: mode === "edit" && !!objectId,
  } as any);
  const ungroupedFields = metadata.fields.filter(
    (field: any) => !fieldsetFields.has(field.name)
  );

  return (
    <AdminLayout>
      <div className="space-y-6">
        <FormHeader className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div className="flex items-center gap-4">
            <Button
              variant="ghost"
              size="icon"
              className="rounded-full hover:bg-muted/50 transition-colors"
              onClick={() => navigate({ to: `/${modelName}` })}
            >
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <div>
              <h1 className="text-3xl font-bold tracking-tight text-foreground/90">
                {mode === "create" ? "Add" : "Edit"} {metadata.verbose_name}
              </h1>
              <p className="text-muted-foreground text-sm">
                {mode === "create"
                  ? `Create a new instance of ${metadata.verbose_name}`
                  : `Updating ${metadata.verbose_name} #${objectId}`}
              </p>
              {isViewOnly && (
                <div className="inline-flex items-center gap-2 mt-2 px-2 py-1 rounded-full bg-muted text-muted-foreground text-[10px] font-bold uppercase tracking-widest">
                  <Eye className="h-3 w-3" />
                  View only
                </div>
              )}
            </div>
          </div>
          <div className="flex items-center gap-2">
            <Button
              type="button"
              variant="ghost"
              onClick={() => navigate({ to: `/${modelName}` })}
              className="text-muted-foreground hover:text-foreground"
            >
              Cancel
            </Button>
            {!isViewOnly && (
              <Button
                data-testid="submit-button"
                onClick={handleSubmit}
                disabled={createMutation.isPending || updateMutation.isPending}
                className="bg-primary hover:bg-primary/90 shadow-lg shadow-primary/20 min-w-[120px]"
              >
                {createMutation.isPending || updateMutation.isPending ? (
                  <Loader2 className="h-4 w-4 animate-spin mr-2" />
                ) : (
                  <Save className="h-4 w-4 mr-2" />
                )}
                {mode === "create" ? "Create" : "Save Changes"}
              </Button>
            )}
          </div>
        </FormHeader>

        <Card className="glass-lite border-border/50 shadow-xl shadow-black/5 overflow-hidden max-w-4xl mx-auto">
          <CardContent className="p-8">
            <form onSubmit={handleSubmit} className="space-y-8">
              {fieldsetGroups.map((fieldset: any, idx: number) => (
                <div
                  key={`${fieldset.name}-${idx}`}
                  className="border border-border/50 rounded-xl p-6 space-y-4"
                >
                  <div className="flex items-center justify-between">
                    <div>
                      <h3 className="text-sm font-bold uppercase tracking-widest text-muted-foreground/80">
                        {fieldset.name}
                      </h3>
                      {fieldset.description && (
                        <p className="text-xs text-muted-foreground mt-1">
                          {fieldset.description}
                        </p>
                      )}
                    </div>
                  </div>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-x-8 gap-y-6">
                    {(fieldset.fields || []).map((fieldName: string) => {
                      const field = metadata.fields.find((f: any) => f.name === fieldName);
                      if (!field) return null;
                      if (field.read_only && mode === "create") return null;
                      if (field.name === "id") return null;

                      const isFullWidth = field.type === 'text' && field.widget === 'textarea';

                      return (
                        <div key={field.name} className={cn("space-y-2", isFullWidth && "md:col-span-2")}>
                          <div className="flex items-center justify-between px-1">
                            <label
                              htmlFor={field.name}
                              className="text-xs font-bold uppercase tracking-widest text-muted-foreground/80"
                            >
                              {field.label}{" "}
                              {field.required && (
                                <span className="text-destructive font-normal">*</span>
                              )}
                            </label>
                          </div>
                          {renderField(field)}
                          {field.help_text && (
                            <p className="text-[11px] text-muted-foreground px-1 leading-relaxed">
                              {field.help_text}
                            </p>
                          )}
                        </div>
                      );
                    })}
                  </div>
                </div>
              ))}

              {ungroupedFields.length > 0 && (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-x-8 gap-y-6">
                  {ungroupedFields.map((field: any) => {
                    if (field.read_only && mode === "create") return null;
                    if (field.name === "id") return null;

                    const isFullWidth = field.type === 'text' && field.widget === 'textarea';

                    return (
                      <div key={field.name} className={cn("space-y-2", isFullWidth && "md:col-span-2")}>
                        <div className="flex items-center justify-between px-1">
                          <label
                            htmlFor={field.name}
                            className="text-xs font-bold uppercase tracking-widest text-muted-foreground/80"
                          >
                            {field.label}{" "}
                            {field.required && (
                              <span className="text-destructive font-normal">*</span>
                            )}
                          </label>
                        </div>
                        {renderField(field)}
                        {field.help_text && (
                          <p className="text-[11px] text-muted-foreground px-1 leading-relaxed">
                            {field.help_text}
                          </p>
                        )}
                      </div>
                    );
                  })}
                </div>
              )}

              <FormFooter className="pt-4 border-t border-border/50 flex justify-end">
                {!isViewOnly && (
                  <Button
                    type="submit"
                    disabled={createMutation.isPending || updateMutation.isPending}
                    className="bg-primary hover:bg-primary/90 shadow-lg shadow-primary/20 min-w-[140px]"
                  >
                    <Save className="h-4 w-4 mr-2" />
                    {mode === "create" ? "Create " + metadata.verbose_name : "Save Changes"}
                  </Button>
                )}
              </FormFooter>
            </form>
            {mode === "edit" && historyData && historyData.length > 0 && (
              <div className="mt-10 border-t border-border/50 pt-6">
                <h3 className="text-sm font-bold uppercase tracking-widest text-muted-foreground/80 mb-4">
                  Change History
                </h3>
                <div className="space-y-3">
                  {historyData.map((entry: any) => (
                    <div
                      key={entry.id || `${entry.timestamp}-${entry.action}`}
                      className="p-3 rounded-lg border border-border/50 bg-muted/20"
                    >
                      <div className="flex items-center justify-between text-xs">
                        <div className="font-semibold uppercase tracking-wider">
                          {entry.action}
                        </div>
                        <div className="text-muted-foreground">
                          {new Date(entry.timestamp).toLocaleString()}
                        </div>
                      </div>
                      <div className="text-xs text-muted-foreground mt-2">
                        {entry.user_name || entry.user_id || "System"}
                      </div>
                      {entry.change_stats && (
                        <pre className="mt-2 text-[11px] bg-background/60 border border-border/40 rounded-md p-2 overflow-auto">
                          {entry.change_stats}
                        </pre>
                      )}
                    </div>
                  ))}
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </AdminLayout>
  );
}
