import type {
	LoginRequest,
	RegistroRequest,
	RegistroResponse,
	TokenResponse,
	SolicitarRecuperacaoRequest,
	SolicitarRecuperacaoResponse,
	RedefinirSenhaRequest,
	TokenPayload,
	ApiError
} from '$lib/types';

const BASE_URL = '/api/v1';

class ApiClient {
	private token: string | null = null;

	setToken(token: string | null) {
		this.token = token;
	}

	private headers(extra: Record<string, string> = {}): Record<string, string> {
		const h: Record<string, string> = {
			'Content-Type': 'application/json',
			...extra
		};
		if (this.token) {
			h['Authorization'] = `Bearer ${this.token}`;
		}
		return h;
	}

	private async request<T>(
		method: string,
		path: string,
		body?: unknown
	): Promise<T> {
		const res = await fetch(`${BASE_URL}${path}`, {
			method,
			headers: this.headers(),
			body: body ? JSON.stringify(body) : undefined
		});

		const data = await res.json();

		if (!res.ok) {
			const err = data as ApiError;
			throw new Error(err.erro ?? `Erro ${res.status}`);
		}

		return data as T;
	}

	// Auth
	async registrarCliente(req: RegistroRequest): Promise<RegistroResponse> {
		return this.request('POST', '/auth/registro/cliente', req);
	}

	async registrarProfissional(req: RegistroRequest): Promise<RegistroResponse> {
		return this.request('POST', '/auth/registro/profissional', req);
	}

	async login(req: LoginRequest): Promise<TokenResponse> {
		return this.request('POST', '/auth/login', req);
	}

	async solicitarRecuperacao(req: SolicitarRecuperacaoRequest): Promise<SolicitarRecuperacaoResponse> {
		return this.request('POST', '/auth/solicitar-recuperacao-senha', req);
	}

	async redefinirSenha(req: RedefinirSenhaRequest): Promise<{ mensagem: string }> {
		return this.request('POST', '/auth/redefinir-senha', req);
	}

	async me(): Promise<TokenPayload> {
		return this.request('GET', '/me');
	}
}

export const api = new ApiClient();
