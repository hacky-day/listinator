FROM python:3 AS train-model
COPY requirements.txt /
RUN pip install --no-cache-dir -r /requirements.txt
COPY train_model.py training.csv /
RUN python /train_model.py

FROM python:3
COPY requirements.txt /app/
RUN pip install --no-cache-dir -r /app/requirements.txt
COPY service.py /app/
COPY --from=train-model /model.pkl /app/

WORKDIR /app
ENTRYPOINT [ "uvicorn", "service:app" ]


