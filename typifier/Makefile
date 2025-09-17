.PHONY: all build run clean

all: build

venv: requirements.txt
	python -m venv venv
	./venv/bin/pip install -r requirements.txt

build: model.pkl

model.pkl: venv train_model.py training.csv
	./venv/bin/python train_model.py


run: venv model.pkl
	./venv/bin/uvicorn service:app --reload --port 8000

clean:
	rm -rf venv __pycache__ *.pyc model.pkl
