# Listinator

A simple web-based list management application perfect for shopping lists,
to-do lists, or any kind of item tracking. Backends are built with Go and
Python, modern frontend powered by Vue.

---

## Features

- Go Backend with persistent storage via SQLite
- Optional machine learning python backend for product type prediction
- Web interface for managing lists

## Setup (Production)

For production, use the pre-built Container Images. The
[docker-compose.yaml](./infra/docker-compose.yaml) can be used as a reference for
production deployment.

## Setup (Development)

### Requirements

- Go 1.24.3 or later
- Node.js 18+ and npm
- Python 3
- `make`

### Quickstart

You can build everything clean with

```bash
make clean build
```

If you need the different parts running for development and testing, you can
just run `make run-*` for the different parts. Just run in different terminals

```bash
make run-core
make run-typifier
make run-frontend
```

For more information, take a look at the `Readme.md` files for the different parts.

## Project Structure

- [`core`](./core/) - Core Backend with embeded Frontend
- [`typifier`](./typifier/) - Microservice for product type prediction

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE)
file for details.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Commit your changes (`git commit -m 'Add some amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

## Support

If you encounter any issues or have questions, please open an issue on GitHub.
