// ============================================================
// Domain Types — espelham os DTOs do backend Go
// ============================================================

export type TipoUsuario = 'CLIENTE' | 'PROFISSIONAL' | 'ADMIN';

// Auth
export interface LoginRequest {
	email: string;
	senha: string;
}

export interface RegistroRequest {
	email: string;
	senha: string;
}

export interface RegistroResponse {
	id: string;
	email: string;
	tipo: TipoUsuario;
}

export interface TokenResponse {
	access_token: string;
	token_type: string;
	expires_in: number;
}

export interface TokenPayload {
	sub: string;
	email: string;
	tipo: TipoUsuario;
	exp: number;
	iat: number;
	iss: string;
}

export interface SolicitarRecuperacaoRequest {
	email: string;
}

export interface SolicitarRecuperacaoResponse {
	mensagem: string;
	token?: string; // apenas em development
}

export interface RedefinirSenhaRequest {
	token: string;
	nova_senha: string;
}

// API Error
export interface ApiError {
	erro: string;
}

// Auth State
export interface AuthState {
	token: string | null;
	payload: TokenPayload | null;
	loading: boolean;
}
