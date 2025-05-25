FROM python:3.12-alpine

WORKDIR /crypto

COPY requirements.txt /crypto

RUN pip install --no-cache-dir -r /crypto/requirements.txt

COPY crypto.py /crypto

COPY server.py /crypto

RUN adduser -D crypto \
    && chown -R crypto:crypto /crypto \
    && chmod -R 500 /crypto

USER crypto
