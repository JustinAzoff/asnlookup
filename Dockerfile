FROM python:3

WORKDIR /usr/src/app
RUN mkdir -p /usr/src/app
RUN pip install 'https://github.com/JustinAzoff/asnlookup-client-python/archive/e76ede9e041571eeee80eada71347fc886b80b0e.zip#egg=asnlookup-client'
COPY requirements.txt /usr/src/app/
RUN pip install --no-cache-dir -r requirements.txt
COPY . /usr/src/app
RUN pip install .

RUN mkdir -p /data
WORKDIR /data

EXPOSE 5555
CMD asnlookup-server
