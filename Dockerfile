FROM python:3

WORKDIR /usr/src/app
RUN mkdir -p /usr/src/app
COPY requirements.txt /usr/src/app/
RUN pip install --no-cache-dir -r requirements.txt
COPY . /usr/src/app
RUN pip install .

RUN mkdir -p /data
WORKDIR /data

CMD asnlookup-server
