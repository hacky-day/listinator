# Listinator

A simple web-based list management application perfect for shopping lists, to-do lists, or any kind of item tracking. Backend built with Go, modern frontend powered by Vue.

## Features

- Web interface for managing shopping and to-do lists
- Mark items as bought/completed
- Search functionality
- Docker and Docker Compose support

## Requirements

- Go 1.24.3 or later
- Node.js 18+ and npm (for frontend development)

## Installation

### Option 1: Docker Compose (Recommended)

1. Clone the repository:

   ```bash
   git clone https://github.com/hacky-day/listinator.git
   cd listinator
   ```

2. Run with Docker Compose:

   ```bash
   docker-compose up
   ```

The application will be available at <http://localhost:8080>.

### Option 2: Build from Source

1. Clone the repository:

   ```bash
   git clone https://github.com/hacky-day/listinator.git
   cd listinator
   ```

2. **Build the frontend:**

   ```bash
   cd frontend
   npm install
   npm run build
   cd ..
   ```

   The built frontend will be placed in the `frontend/dist` directory and served as static files by the Go backend.

3. **Build the backend:**

   ```bash
   go build
   ```

4. Run the application:

   ```bash
   ./listinator
   ```

The application will be available at <http://localhost:8080>.

## Development

### Project Structure

- `main.go` – Application entry point and server setup (Go)
- `server/` – HTTP handlers and API routes (Go)
- `database/` – Database models and initialization (Go)
- `frontend/` – Vue frontend (source code, build, assets)

### Frontend Development (Vue)

1. Navigate to the frontend directory:

   ```bash
   cd frontend
   ```

2. Install dependencies:

   ```bash
   npm install
   ```

3. Start the development server:

   ```bash
   npm run dev
   ```

   The frontend will be available at <http://localhost:5173> (default). The backend needs also to be starten and will be proxied by the frontend.

### Backend Development

1. Make sure to have the frontend build and present in `frontend/dist`

2. Just build with

   ```bash
   go build
   ```

   The frontend and backend will be served at <http://localhost:8080>.

## Configuration

The application uses the following environment variable:

- `LISTINATOR_DATABASE_DIR` - Directory where the SQLite database file will be stored (required)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Commit your changes (`git commit -m 'Add some amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

## Support

If you encounter any issues or have questions, please open an issue on GitHub.
