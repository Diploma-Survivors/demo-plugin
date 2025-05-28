"use client";

import { useLTIContext } from "@/hooks/useLTIContext";
import { CodeExecutor } from "@/components/CodeExecutor";
import { Alert, AlertDescription } from "@/components/ui/alert";

export default function Home() {
  const { context, error } = useLTIContext();

  if (error) {
    return (
      <div className="max-w-xl mx-auto py-12">
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      </div>
    );
  }

  if (!context) {
    return (
      <div className="max-w-xl mx-auto py-12">
        <Alert>
          <AlertDescription>Loading LTI context...</AlertDescription>
        </Alert>
      </div>
    );
  }

  return (
    <div className="max-w-xl mx-auto py-12">
      <div className="mb-8">
        <h1 className="text-2xl font-bold mb-2">Code Execution</h1>
        <p className="text-sm text-gray-600 dark:text-gray-400">
          User ID: {context.user}
          {context.lineitem && " • Grades will be submitted to Moodle"}
        </p>
      </div>

      <CodeExecutor context={context} />
    </div>
  );
}
