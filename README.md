# Listinator - Vue Branch

A modern web-based list management application perfect for shopping lists, to-do lists, or any kind of item tracking. Built with **Go** backend and **Vue.js 3 + TypeScript** frontend, featuring a clean and intuitive interface with modern build tooling.

> **Note**: This is the **Vue branch** - a complete rewrite of the frontend from vanilla JavaScript to Vue.js 3 with TypeScript. For the original vanilla JS version, see the `main` branch.

## ✨ Features

- **Modern Vue.js Frontend**: Built with Vue 3, TypeScript, and Vite for fast development
- **Responsive Design**: Works seamlessly on desktop and mobile devices
- **List Management**: Create and manage multiple shopping lists and to-do lists
- **Smart Entry Management**: Add, edit, delete, and mark items as bought/completed
- **Category System**: Organize items with pre-defined types and emoji icons (🍎 fruit, 🥦 vegetable, etc.)
- **Real-time Search**: Filter items instantly as you type
- **Share Lists**: Generate shareable URLs for collaborative list management
- **Offline Capable**: Progressive Web App features for offline usage
- **Docker Support**: Easy deployment with Docker and docker-compose

## 🛠️ Technology Stack

### Backend
- **Go 1.24.3+** with Echo web framework
- **SQLite** database with GORM ORM
- **RESTful API** with comprehensive documentation
- **UUID-based** primary keys for all entities

### Frontend
- **Vue.js 3** with Composition API
- **TypeScript** for type safety
- **Vite** for fast builds and hot reload
- **CSS Grid/Flexbox** for responsive layouts
- **Progressive Web App** capabilities

### Build System
- **Comprehensive Makefile** with 15+ commands
- **Automated dependency management**
- **Production optimizations**
- **Embedded frontend** assets in Go binary

## 📋 Requirements

- **Go 1.24.3 or later**
- **Node.js 20.19.0+ or 22.12.0+**
- **NPM 10.8.2+**

## 🚀 Quick Start

### Option 1: Using Make (Recommended)

1. **Clone and setup:**
```bash
git clone https://github.com/hacky-day/listinator.git
cd listinator
git checkout vue  # Switch to Vue branch
```

2. **Install dependencies and build:**
```bash
make deps          # Install both Go and Node.js dependencies
make build         # Build frontend and backend
```

3. **Run the application:**
```bash
export LISTINATOR_DATABASE_DIR=./data
make dev          # Start development server
```

Visit http://localhost:8080 to use the application.

### Option 2: Manual Build

1. **Backend setup:**
```bash
go mod download
go build .
```

2. **Frontend setup:**
```bash
cd frontend
npm install
npm run build
cd ..
```

3. **Run:**
```bash
export LISTINATOR_DATABASE_DIR=./data
./listinator
```

### Option 3: Docker (Production)

```bash
docker-compose up --build
```

Or with Docker directly:
```bash
docker build -t listinator .
docker run -p 8080:8080 -v $(pwd)/data:/var/lib/listinator -e LISTINATOR_DATABASE_DIR=/var/lib/listinator listinator
```

## 🔧 Development

### Available Make Commands

```bash
make help          # Show all available commands
make deps          # Install all dependencies
make dev           # Start development server
make dev-frontend  # Start Vue dev server with hot reload
make build         # Build everything for production
make clean         # Clean build artifacts
make lint          # Run linters for Go and TypeScript
make test          # Run tests
make check-parity  # Compare features with main branch
```

### Project Structure

```
├── main.go              # Application entry point
├── api/v1/server/       # HTTP handlers and API routes
├── database/            # Database models and initialization  
├── frontend/            # Vue.js application
│   ├── src/            # Vue components and TypeScript
│   ├── public/         # Static assets
│   └── dist/           # Built frontend (auto-generated)
├── Makefile            # Build automation
└── docker-compose.yaml # Container orchestration
```

### API Endpoints

All endpoints are prefixed with `/api/v1/`:

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | `/entries?ListID={uuid}` | List entries for a specific list |
| POST   | `/entries` | Create a new entry |
| PUT    | `/entries/{id}` | Update an existing entry |
| DELETE | `/entries/{id}` | Delete an entry |
| POST   | `/lists` | Create a new list |
| GET    | `/types` | Get all available item types |

### Development Workflow

1. **Start backend with hot reload:**
```bash
make dev
```

2. **Start frontend dev server (separate terminal):**
```bash
make dev-frontend
```

3. **Make changes to Vue components** in `frontend/src/`
4. **Changes auto-reload** in your browser
5. **Build for production:**
```bash
make build-production
```

## 📊 Feature Parity with Main Branch

| Feature | Main Branch (Vanilla JS) | Vue Branch | Status |
|---------|-------------------------|------------|---------|
| List Management | ✓ | ✓ | ✅ Full parity |
| Entry CRUD Operations | ✓ | ✓ | ✅ Full parity |
| Search/Filter | ✓ | ✓ | ✅ Full parity |
| Item Categories | ✓ | ✓ | ✅ Full parity |
| Share Functionality | ✓ | ✓ | ✅ Full parity |
| Responsive Design | ✓ | ✓ | ✅ Enhanced |
| Build System | ❌ | ✓ | ✅ New feature |
| Type Safety | ❌ | ✓ | ✅ New feature |
| Hot Reload | ❌ | ✓ | ✅ New feature |
| Component Architecture | ❌ | ✓ | ✅ New feature |

## 🔧 Configuration

### Environment Variables

- `LISTINATOR_DATABASE_DIR` - **Required**. Directory where the SQLite database file will be stored

### Database

- **SQLite** database with automatic migrations
- **UUID primary keys** for all entities
- **Soft deletes** with GORM
- **Pre-seeded item types** with emoji icons

## 🐳 Docker

### Development
```bash
docker-compose up --build
```

### Production
```bash
docker build -t listinator .
docker run -p 8080:8080 -e LISTINATOR_DATABASE_DIR=/data listinator
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch from `vue` branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Run tests and linters: `make lint test`
5. Commit your changes: `git commit -m 'Add some amazing feature'`
6. Push to the branch: `git push origin feature/amazing-feature`
7. Open a Pull Request against the `vue` branch

### Development Guidelines

- Use TypeScript for all frontend code
- Follow Vue 3 Composition API patterns
- Add appropriate Go documentation comments
- Update tests for new functionality
- Use the provided Makefile for consistent builds

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 Support

If you encounter any issues or have questions:

1. Check the [Issues](https://github.com/hacky-day/listinator/issues) page
2. Run `make help` for available commands
3. Open a new issue with detailed information

---

**Vue Branch Enhancements**: Modern frontend, TypeScript safety, component architecture, comprehensive build system, and improved developer experience.