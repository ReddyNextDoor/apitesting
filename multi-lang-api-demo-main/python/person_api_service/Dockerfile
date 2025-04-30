FROM python:3.13-slim
WORKDIR /app
COPY . .
RUN pip install poetry && poetry config virtualenvs.create false && poetry install --no-interaction --no-ansi
EXPOSE 8000
CMD ["poetry", "run", "uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
