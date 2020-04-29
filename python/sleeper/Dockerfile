FROM python:3

COPY parcs parcs/
RUN pip install --no-cache-dir docker
EXPOSE 4444

COPY sleeper sleeper/
CMD [ "python", "-m", "sleeper.main" ]
