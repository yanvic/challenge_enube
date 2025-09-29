# Challenge Enube

**Sistema completo com autenticação JWT, importação de dados e interface moderna**

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go&logoColor=white)

![React](https://img.shields.io/badge/React-20+-61DAFB?style=for-the-badge&logo=react&logoColor=black)

![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16+-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)

![Docker](https://img.shields.io/badge/Docker-24+-2496ED?style=for-the-badge&logo=docker&logoColor=white)


---

## 📋 Sobre o Projeto

Sistema full-stack com backend em **Golang**, frontend em **React** e banco de dados **PostgreSQL**. Implementa autenticação JWT, importação de dados e interface responsiva.

### ✨ Principais Funcionalidades

- 🔐 **Autenticação JWT** - Sistema seguro de login e autorização
- 📊 **Importação de Dados** - Suporte para upload e processamento de arquivos
- 🎨 **Interface Moderna** - UI responsiva e intuitiva em React
- 🐳 **Docker Ready** - Deploy simplificado com Docker Compose
- 🔄 **CORS Habilitado** - Comunicação fluida entre frontend e backend

---

## 🚀 Início Rápido

### Opção 1: Docker Compose (Recomendado)
```bash
# Clone o repositório
git clone <seu-repositorio>
cd challenge-enube

# Configure as variáveis de ambiente
cp .env.example .env

# Suba todos os serviços
docker-compose up --build
Pronto! Acesse:

🌐 Frontend: http://localhost:3000
⚙️ Backend: http://localhost:8080
🗄️ Adminer: http://localhost:8000
```

### Opção 2: Execução Local
Veja a seção 📦 Instalação Local abaixo.

📋 Pré-requisitos
Para Execução Local
FerramentaVersão Mínima Go1.25+ 
Node.js20+
NPM9+
PostgreSQL16+
Docker24+
Docker Compose2.20+

--- 

⚙️ Configuração
### 1. Variáveis de Ambiente
Crie um arquivo .env na raiz do projeto:
env# Banco de Dados
DATABASE_URL=host=localhost user=admin password=admin dbname=challenge sslmode=disable

# Autenticação
JWT_SECRET=supersecretjwtkey123456

# Servidor
PORT=8080

💡 Dica: Altere o JWT_SECRET para um valor único em produção!

### 2. Otimização Docker (Opcional)
Para builds mais rápidos, ative o BuildKit:
bashexport DOCKER_BUILDKIT=1

📦 Instalação Local
Backend (Golang)
bash# Navegue até a pasta do backend
cd back

# Baixe as dependências
go mod download

# Execute o servidor
go run cmd/api/main.go
✅ Backend rodando em http://localhost:8080
Frontend (React)
bash# Navegue até a pasta do frontend
cd front

# Instale as dependências
npm ci

# Inicie o servidor de desenvolvimento
npm start
✅ Frontend rodando em http://localhost:3000
Banco de Dados (PostgreSQL)
Configure um banco PostgreSQL com as seguintes credenciais (ou ajuste o .env):
Database: challenge
User: admin
Password: admin
Usando Docker para apenas o banco:
bashdocker run -d \
  --name postgres_challenge \
  -e POSTGRES_DB=challenge \
  -e POSTGRES_USER=admin \
  -e POSTGRES_PASSWORD=admin \
  -p 5432:5432 \
  postgres:16-alpine

🐳 Docker Compose
Arquitetura dos Serviços
O docker-compose.yml provisiona 4 serviços:
ServiçoDescriçãoPortapostgres_dbBanco de dados PostgreSQL 165432go-appBackend API em Golang8080react-appFrontend React3000adminer_uiInterface web para gerenciar o banco8000
Comandos Docker
bash# Iniciar todos os serviços
docker-compose up --build

# Iniciar em background
docker-compose up -d

# Parar todos os serviços
docker-compose down

# Rebuild completo (limpa cache)
docker-compose build --no-cache

# Ver logs de um serviço específico
docker-compose logs -f go-app

# Remover volumes (reset completo do banco)
docker-compose down -v
Acessando o Adminer
Acesse http://localhost:8000 e use:
Sistema: PostgreSQL
Servidor: postgres_db
Usuário: admin
Senha: admin
Base de dados: challenge

```bash
📁 Estrutura do Projeto
challenge-enube/
│
├── 📂 back/                    # Backend Golang
│   ├── 📂 cmd/
│   │   └── 📂 api/
│   │       └── main.go         # Entry point da aplicação
│   ├── 📂 internal/
│   │   ├── 📂 database/        # Conexão e migrations
│   │   ├── 📂 handlers/        # Controllers HTTP
│   │   ├── 📂 models/          # Estruturas de dados
│   │   ├── 📂 service/         # Regras de negócio
│   │   └── 📂 usecase/         # Casos de uso
│   ├── go.mod
│   └── go.sum
│
├── 📂 front/                   # Frontend React
│   ├── 📂 public/
│   ├── 📂 src/
│   │   ├── 📂 components/
│   │   ├── 📂 pages/
│   │   └── App.js
│   ├── package.json
│   └── package-lock.json
│
├── 🐳 Dockerfile.go            # Dockerfile otimizado do Go
├── 🐳 Dockerfile.react         # Dockerfile otimizado do React
├── 🐳 docker-compose.yml       # Orquestração de containers
├── 📄 .env                     # Variáveis de ambiente
└── 📖 README.md
```

🔧 Desenvolvimento
Hot Reload
Backend:
O Go não tem hot reload nativo. Recomenda-se usar Air:
bash# Instalar Air
go install github.com/cosmtrek/air@latest

# Executar com hot reload
cd back
air
Frontend:
O React já vem com hot reload habilitado via npm start.
Executando Testes
bash# Backend
cd back
go test ./...

# Frontend
cd front
npm test

🌐 Endpoints da API
Autenticação
httpPOST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "senha123"
}
Recursos Protegidos
httpGET /api/protected
Authorization: Bearer <seu-jwt-token>

📚 Documentação completa: A API possui documentação Swagger (se implementada) em /api/docs


🛠️ Troubleshooting
Porta já em uso
bash# Verificar processos usando a porta
lsof -i :8080  # ou :3000, :5432

# Matar o processo
kill -9 <PID>
Erro de conexão com banco

Verifique se o PostgreSQL está rodando
Confirme as credenciais no .env
Teste a conexão:

bashpsql -h localhost -U admin -d challenge
Container não inicia
bash# Ver logs detalhados
docker-compose logs <nome-do-serviço>

# Rebuild forçado
docker-compose down -v
docker-compose build --no-cache
docker-compose up

📝 Notas Importantes

⚠️ CORS: O backend já vem configurado para aceitar requisições de localhost:3000
🔒 Segurança: Altere o JWT_SECRET para produção
💾 Persistência: Os dados do PostgreSQL são salvos em volume Docker
🚀 Performance: O BuildKit do Docker acelera significativamente os builds
📦 Cache: Docker reutiliza camadas para builds incrementais

```
