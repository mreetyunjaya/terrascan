FROM python:alpine

LABEL maintainer="Accurics <support@accurics.com>"

COPY setup.py .
COPY HISTORY.rst .
COPY docker_test.sh .
COPY terrascan terrascan
RUN pip install --no-cache-dir -e .

ENTRYPOINT ["terrascan"]
