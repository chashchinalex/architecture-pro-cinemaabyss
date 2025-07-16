
# Python образ
FROM python:3.11-slim

# Установка зависимостей
RUN pip install flask requests

# Копирование приложения
COPY proxy.py /app/proxy.py
WORKDIR /app

# Прослушиваемый порт
ENV PORT=8080
EXPOSE 8080

# Запуск приложения
CMD ["python", "proxy.py"]
