FROM python:3.9-slim

RUN pip install --upgrade pip

RUN pip install astor psutil

WORKDIR /app

COPY ./pkg/analyzer/languages/python/static/analyze_script.py /app/analyze_script.py

ENTRYPOINT [ "python3", "/app/analyze_script.py" ]
