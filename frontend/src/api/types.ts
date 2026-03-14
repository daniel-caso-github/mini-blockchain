export interface Block {
  index: number;
  timestamp: string;
  data: string;
  prev_hash: string;
  hash: string;
  nonce: number;
}

export interface ChainResponse {
  chain: Block[];
  length: number;
  difficulty: number;
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface ValidateResponse {
  valid: boolean;
}

export interface HealthResponse {
  status: string;
}

export interface ErrorResponse {
  error: string;
}

export interface WSDifficultyMessage {
  type: 'difficulty_adjusted';
  difficulty: number;
}
