FROM alpine:3.13.0

RUN apk upgrade && \
    apk add bash gcc libc-dev python3-dev libffi-dev openssl-dev

COPY . /app
WORKDIR /app
RUN python3 -m ensurepip
RUN python3 -m pip install --upgrade pip setuptools
RUN python3 -m pip install -r requirements.txt

CMD ["python3", "swapi.py"]
