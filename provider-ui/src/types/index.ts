export interface LTIContext {
  id_token: string;
  user: string;
  context: string;
  lineitem?: string;
}

export interface ExecuteRequest {
  code: string;
  language: string;
  user_id: string;
  lineitem?: string;
  id_token: string;
  max_score: number;
}

export interface Judge0Response {
  token?: string;
  status?: {
    id: number;
    description: string;
  };
  stdout?: string;
  stderr?: string;
  compile_output?: string;
  message?: string;
  exit_code?: number;
  exit_signal?: number;
  time?: string;
  memory?: number;
}

export interface ExecuteResponse {
  success: boolean;
  result?: Judge0Response;
  score?: number;
  error?: string;
}

export const SUPPORTED_LANGUAGES = [
  { label: "Go", value: "go" },
  { label: "Python", value: "python" },
  { label: "Java", value: "java" },
  { label: "JavaScript", value: "nodejs" },
  { label: "C++", value: "cpp" },
  { label: "C", value: "c" },
] as const;

export type SupportedLanguage = (typeof SUPPORTED_LANGUAGES)[number]["value"];
