import { useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";
import { LTIContext } from "@/types";

export function useLTIContext() {
  const searchParams = useSearchParams();
  const [context, setContext] = useState<LTIContext | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    try {
      // Get required parameters
      const id_token = searchParams.get("id_token");
      const user = searchParams.get("user");
      const contextId = searchParams.get("context");

      if (!id_token || !user || !contextId) {
        throw new Error("Missing required LTI parameters");
      }

      // Get optional parameters
      const lineitem = searchParams.get("lineitem") || undefined;

      // Set context
      setContext({
        id_token,
        user,
        context: contextId,
        lineitem,
      });
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to get LTI context");
    }
  }, [searchParams]);

  return { context, error };
}
