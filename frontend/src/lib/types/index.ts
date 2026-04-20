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

// ============================================================
// Etapa 2 — Cadastro Completo de Clientes e Profissionais
// ============================================================

// Cliente
export interface ClienteRequest {
	nome: string;
	cpf: string;
	telefone: string;
}

export interface ClienteResponse {
	id: string;
	usuario_id: string;
	nome: string;
	cpf: string;
	telefone: string;
	score: number;
}

// Endereco
export interface EnderecoRequest {
	logradouro: string;
	numero: string;
	complemento: string;
	bairro: string;
	cidade: string;
	estado: string;
	cep: string;
	num_quartos: number;
	num_banheiros: number;
	num_salas: number;
	num_cozinhas: number;
	area_m2: number;
	principal: boolean;
}

export interface EnderecoResponse {
	id: string;
	cliente_id: string;
	logradouro: string;
	numero: string;
	complemento: string;
	bairro: string;
	cidade: string;
	estado: string;
	cep: string;
	num_quartos: number;
	num_banheiros: number;
	num_salas: number;
	num_cozinhas: number;
	area_m2: number;
	principal: boolean;
}

// Profissional
export type StatusProfissional = 'PENDENTE' | 'APROVADA' | 'REPROVADA' | 'SUSPENSA' | 'DESCREDENCIADA';

export interface ProfissionalRequest {
	nome: string;
	cpf: string;
	rg: string;
	telefone: string;
	foto_url: string;
	mei: boolean;
}

export interface ProfissionalResponse {
	id: string;
	usuario_id: string;
	nome: string;
	cpf: string;
	rg: string;
	telefone: string;
	foto_url: string;
	mei: boolean;
	status: StatusProfissional;
	nota_media: number;
	total_servicos: number;
}

// Documento
export type TipoDocumento = 'RG_FRENTE' | 'RG_VERSO' | 'CPF' | 'COMPROVANTE' | 'FOTO' | 'OUTRO';
export type StatusDocumento = 'PENDENTE' | 'APROVADO' | 'REPROVADO';

export interface DocumentoRequest {
	tipo: TipoDocumento;
	url: string;
}

export interface DocumentoResponse {
	id: string;
	tipo: TipoDocumento;
	url: string;
	status: StatusDocumento;
}

// Referencia
export type StatusReferencia = 'PENDENTE' | 'CONFIRMADA' | 'NAO_CONFIRMADA';

export interface ReferenciaRequest {
	nome_contato: string;
	telefone_contato: string;
}

export interface ReferenciaResponse {
	id: string;
	nome_contato: string;
	telefone_contato: string;
	status: StatusReferencia;
}

// Regiao
export interface RegiaoResponse {
	id: string;
	nome: string;
	cidade: string;
	estado: string;
	cep_inicio: string;
	cep_fim: string;
}

export interface DefinirRegioesRequest {
	regiao_ids: string[];
}

// Disponibilidade
export interface DisponibilidadeSlot {
	dia_semana: number;
	hora_inicio: string;
	hora_fim: string;
}

export interface DefinirDisponibilidadesRequest {
	slots: DisponibilidadeSlot[];
}

export interface DisponibilidadeResponse {
	id: string;
	dia_semana: number;
	hora_inicio: string;
	hora_fim: string;
}

// Catálogo (diferencial C)
export type FrequenciaServico = 'UNICA' | 'SEMANAL' | 'QUINZENAL' | 'DUAS_POR_SEMANA';

export interface CategoriaResponse {
	id: string;
	nome: string;
	descricao: string;
	duracao_minima_min: number;
	ativa: boolean;
}

export interface OpcionalResponse {
	id: string;
	nome: string;
	descricao: string;
	valor_extra: number;
	tempo_extra_min: number;
	ativo: boolean;
}

export interface CalculoPrecoRequest {
	categoria_id: string;
	regiao_id: string;
	num_quartos: number;
	num_banheiros?: number;
	num_salas?: number;
	num_cozinhas?: number;
	opcionais_ids?: string[];
	frequencia?: FrequenciaServico;
	data_servico?: string;
}

export type TipoItemCalculo = 'BASE' | 'COMODO' | 'OPCIONAL' | 'ACRESCIMO' | 'DESCONTO' | 'TOTAL';

export interface ItemCalculo {
	tipo: TipoItemCalculo;
	label: string;
	valor: number;
}

export interface CalculoPrecoResponse {
	categoria_id: string;
	regiao_id: string;
	duracao_min: number;
	valor_total: number;
	itens: ItemCalculo[];
}

// Preferências (diferencial B)
export type TipoPreferencia = 'FAVORITA' | 'BLOQUEADA';

export interface PreferenciaResponse {
	id: string;
	profissional_id: string;
	tipo: TipoPreferencia;
	motivo: string;
	nome: string;
	foto_url: string;
	nota_media: number;
	criado_em: string;
}

// ============================================================
// Etapa 4 — Solicitações de Serviço
// ============================================================

export type StatusSolicitacao =
	| 'AGUARDANDO'
	| 'ATRIBUIDA'
	| 'CONFIRMADA'
	| 'EM_ANDAMENTO'
	| 'CONCLUIDA'
	| 'CANCELADA';

export interface CriarSolicitacaoRequest {
	endereco_id: string;
	categoria_id: string;
	regiao_id: string;
	num_quartos: number;
	num_banheiros: number;
	num_salas: number;
	num_cozinhas: number;
	opcionais_ids: string[];
	frequencia: FrequenciaServico;
	data_servico: string; // ISO 8601
	observacao?: string;
}

export interface SolicitacaoResponse {
	id: string;
	endereco_id: string;
	categoria_id: string;
	regiao_id: string;
	num_quartos: number;
	num_banheiros: number;
	num_salas: number;
	num_cozinhas: number;
	frequencia: FrequenciaServico;
	data_servico: string;
	observacao: string;
	valor_total: number;
	duracao_min: number;
	breakdown: ItemCalculo[];
	status: StatusSolicitacao;
	criada_em: string;
	cancelada_em?: string;
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
