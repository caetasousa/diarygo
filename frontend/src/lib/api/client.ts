import type {
	LoginRequest,
	RegistroRequest,
	RegistroResponse,
	TokenResponse,
	SolicitarRecuperacaoRequest,
	SolicitarRecuperacaoResponse,
	RedefinirSenhaRequest,
	TokenPayload,
	ApiError,
	ClienteRequest,
	ClienteResponse,
	EnderecoRequest,
	EnderecoResponse,
	ProfissionalRequest,
	ProfissionalResponse,
	DocumentoRequest,
	DocumentoResponse,
	ReferenciaRequest,
	ReferenciaResponse,
	RegiaoResponse,
	DefinirRegioesRequest,
	DefinirDisponibilidadesRequest,
	DisponibilidadeResponse
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

	// Cliente
	async buscarPerfilCliente(): Promise<ClienteResponse> {
		return this.request('GET', '/clientes/me');
	}

	async criarPerfilCliente(req: ClienteRequest): Promise<ClienteResponse> {
		return this.request('POST', '/clientes/me', req);
	}

	async atualizarPerfilCliente(req: ClienteRequest): Promise<ClienteResponse> {
		return this.request('PUT', '/clientes/me', req);
	}

	// Enderecos
	async listarEnderecos(): Promise<EnderecoResponse[]> {
		return this.request('GET', '/clientes/me/enderecos');
	}

	async criarEndereco(req: EnderecoRequest): Promise<EnderecoResponse> {
		return this.request('POST', '/clientes/me/enderecos', req);
	}

	async atualizarEndereco(id: string, req: EnderecoRequest): Promise<EnderecoResponse> {
		return this.request('PUT', `/clientes/me/enderecos/${id}`, req);
	}

	async removerEndereco(id: string): Promise<void> {
		return this.request('DELETE', `/clientes/me/enderecos/${id}`);
	}

	async definirEnderecoPrincipal(id: string): Promise<EnderecoResponse> {
		return this.request('PUT', `/clientes/me/enderecos/${id}/principal`);
	}

	// Profissional
	async buscarPerfilProfissional(): Promise<ProfissionalResponse> {
		return this.request('GET', '/profissionais/me');
	}

	async criarPerfilProfissional(req: ProfissionalRequest): Promise<ProfissionalResponse> {
		return this.request('POST', '/profissionais/me', req);
	}

	async atualizarPerfilProfissional(req: ProfissionalRequest): Promise<ProfissionalResponse> {
		return this.request('PUT', '/profissionais/me', req);
	}

	// Documentos
	async enviarDocumento(req: DocumentoRequest): Promise<DocumentoResponse> {
		return this.request('POST', '/profissionais/me/documentos', req);
	}

	async listarDocumentos(): Promise<DocumentoResponse[]> {
		return this.request('GET', '/profissionais/me/documentos');
	}

	// Referencias
	async adicionarReferencia(req: ReferenciaRequest): Promise<ReferenciaResponse> {
		return this.request('POST', '/profissionais/me/referencias', req);
	}

	async listarReferencias(): Promise<ReferenciaResponse[]> {
		return this.request('GET', '/profissionais/me/referencias');
	}

	// Regioes (publica)
	async listarRegioes(): Promise<RegiaoResponse[]> {
		return this.request('GET', '/regioes');
	}

	// Regioes de atuacao (profissional)
	async definirRegioesAtuacao(req: DefinirRegioesRequest): Promise<RegiaoResponse[]> {
		return this.request('PUT', '/profissionais/me/regioes', req);
	}

	async listarRegioesAtuacao(): Promise<RegiaoResponse[]> {
		return this.request('GET', '/profissionais/me/regioes');
	}

	// Disponibilidade
	async definirDisponibilidades(req: DefinirDisponibilidadesRequest): Promise<DisponibilidadeResponse[]> {
		return this.request('PUT', '/profissionais/me/disponibilidades', req);
	}

	async listarDisponibilidades(): Promise<DisponibilidadeResponse[]> {
		return this.request('GET', '/profissionais/me/disponibilidades');
	}
}

export const api = new ApiClient();
