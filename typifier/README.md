# Typifier

A simple microservice to predict the types of products using classic machine
learning algorithms.

---

## Setup (Production)

For production, use the pre-built Container Images.

## Setup (Development)

### Requirements

- Python 3
- `make`

### Quick Start

Simply run:

```bash
make clean run
```

This will:

1. Clean old artifacts
2. Create a virtual environment and install dependencies
3. Train the model
4. Start the microservice locally

## API Usage

Once running, you can call the API with:

```bash
$ curl "http://127.0.0.1:8000" \
  -H "Content-Type: application/json" \
  -d '{"product": "shampoo"}'
{"product":"shampoo","type":"household & baby & pets","uuid":"a14bca10-13b7-4a9c-a663-75a5203c3f09"}
```
