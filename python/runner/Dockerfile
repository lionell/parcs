FROM python:3

COPY parcs parcs/
EXPOSE 4444
RUN pip install --no-cache-dir docker

COPY runner runner/
CMD [ "python", "-m", "runner.main" ]
