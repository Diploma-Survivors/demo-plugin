import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Card } from "@/components/ui/card";
import {
  LTIContext,
  ExecuteResponse,
  SUPPORTED_LANGUAGES,
  SupportedLanguage,
} from "@/types";

interface Props {
  context: LTIContext;
}

export function CodeExecutor({ context }: Props) {
  const [code, setCode] = useState("");
  const [lang, setLang] = useState<SupportedLanguage>("python");
  const [result, setResult] = useState<ExecuteResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit() {
    setLoading(true);
    setResult(null);
    setError(null);

    try {
      const res = await fetch("/api/execute", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          code,
          language: lang,
          user_id: context.user,
          lineitem: context.lineitem,
          id_token: context.id_token,
          max_score: 100,
        }),
      });

      if (!res.ok) {
        throw new Error(`HTTP error! status: ${res.status}`);
      }

      const data: ExecuteResponse = await res.json();
      setResult(data);

      if (!data.success) {
        setError(data.error || "Unknown error occurred");
      }
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to execute code");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Select
          value={lang}
          onValueChange={(v) => setLang(v as SupportedLanguage)}
        >
          <SelectTrigger className="w-40">
            <SelectValue placeholder="Select language" />
          </SelectTrigger>
          <SelectContent>
            {SUPPORTED_LANGUAGES.map((l) => (
              <SelectItem key={l.value} value={l.value}>
                {l.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        <Button onClick={handleSubmit} disabled={loading || !code}>
          {loading ? "Running..." : "Run Code"}
        </Button>
      </div>

      <Textarea
        className="min-h-[200px] font-mono text-sm"
        value={code}
        onChange={(e) => setCode(e.target.value)}
        placeholder={`Enter your ${lang} code here...`}
      />

      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {result?.result && (
        <Card className="p-4 space-y-4">
          {result.result.stdout && (
            <div>
              <h3 className="font-semibold mb-2">Output:</h3>
              <pre className="bg-gray-100 dark:bg-gray-900 p-3 rounded text-sm overflow-x-auto">
                {result.result.stdout}
              </pre>
            </div>
          )}

          {result.result.stderr && (
            <div>
              <h3 className="font-semibold mb-2 text-red-500">Errors:</h3>
              <pre className="bg-red-50 dark:bg-red-900/20 p-3 rounded text-sm overflow-x-auto text-red-600 dark:text-red-400">
                {result.result.stderr}
              </pre>
            </div>
          )}

          {result.result.compile_output && (
            <div>
              <h3 className="font-semibold mb-2">Compiler Output:</h3>
              <pre className="bg-gray-100 dark:bg-gray-900 p-3 rounded text-sm overflow-x-auto">
                {result.result.compile_output}
              </pre>
            </div>
          )}

          {typeof result.score === "number" && (
            <div className="text-sm text-gray-600 dark:text-gray-400">
              Score: {result.score}/100
            </div>
          )}
        </Card>
      )}
    </div>
  );
}
