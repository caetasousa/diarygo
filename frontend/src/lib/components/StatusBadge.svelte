<script lang="ts">
	interface Props {
		status: string;
	}

	let { status }: Props = $props();

	const LABELS: Record<string, string> = {
		PENDENTE: 'Pendente',
		APROVADO: 'Aprovado',
		APROVADA: 'Aprovada',
		REPROVADO: 'Reprovado',
		REPROVADA: 'Reprovada',
		CONFIRMADA: 'Confirmada',
		NAO_CONFIRMADA: 'Não confirmada',
		SUSPENSA: 'Suspensa',
		DESCREDENCIADA: 'Descredenciada',
		// Solicitação
		AGUARDANDO: 'Aguardando',
		ATRIBUIDA: 'Atribuída',
		EM_ANDAMENTO: 'Em andamento',
		CONCLUIDA: 'Concluída',
		CANCELADA: 'Cancelada'
	};

	const CLASSES: Record<string, string> = {
		PENDENTE: 'badge-orange',
		APROVADO: 'badge-green',
		APROVADA: 'badge-green',
		REPROVADO: 'badge-red',
		REPROVADA: 'badge-red',
		CONFIRMADA: 'badge-green',
		NAO_CONFIRMADA: 'badge-red',
		SUSPENSA: 'badge-red',
		DESCREDENCIADA: 'badge-red',
		// Solicitação
		AGUARDANDO: 'badge-orange',
		ATRIBUIDA: 'badge-blue',
		EM_ANDAMENTO: 'badge-blue',
		CONCLUIDA: 'badge-green',
		CANCELADA: 'badge-red'
	};
</script>

<span class="badge {CLASSES[status] ?? 'badge-blue'}">
	{LABELS[status] ?? status}
</span>
