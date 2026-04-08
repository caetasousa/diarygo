# Regras de Negocio

## Entidades Principais

- **Cliente:** pessoa fisica/juridica que contrata servicos de limpeza
- **Profissional (Diarista):** trabalhadora autonoma, sem vinculo empregaticio
- **Servico:** agendamento de limpeza com tipo, data, hora, endereco, opcionais
- **Avaliacao:** nota (1-5) + comentario, mutua (cliente avalia profissional e vice-versa)
- **Endereco:** com CEP para geolocalizacao

## Regras Criticas

1. **Limite legal (LC 150/2015):** maximo 2 visitas/semana da mesma profissional no mesmo endereco. 3+ dias configura vinculo empregaticio — sistema DEVE bloquear e rotacionar profissionais automaticamente
2. **Agendamento minimo:** 24h de antecedencia
3. **Cancelamento < 24h:** penalizacao no score do cliente
4. **Avaliacao:** profissionais com nota < 4.0 recebem alerta; nota < 3.5 por 3 servicos consecutivos = suspensao temporaria
5. **Atribuicao:** priorizar profissionais por melhor avaliacao (minimo configuravel, ex: 4.6/5.0) e proximidade geografica
6. **Profissional tem autonomia:** escolhe regioes, horarios e quais servicos aceitar
7. **Tempo de resposta:** profissional tem 30 min para aceitar/recusar; apos isso, redireciona para proxima do ranking

## Tipos de Servico

- Limpeza Padrao, Pesada, Pos-Obra, Pre-Mudanca, Comercial, Express, Passadoria
- Add-ons: geladeira, armarios, tapetes, janelas, lavagem/passadoria de roupas, area externa
