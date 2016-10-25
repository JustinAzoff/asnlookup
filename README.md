# Standalone Usage

    $ echo -e "8.8.8.8\n64.233.177.113\n31.13.65.1\n72.30.202.51"|asnlookup | column -t -s $'\t'
    DEBUG:asnlookup.data_manager:file_age_in_hours: filename=asn.db hours=16
    DEBUG:asnlookup.data_manager:maybe_update: func=db output_filename=asn.db max_age_in_hours=24 result=cached age=16
    DEBUG:asnlookup.data_manager:file_age_in_hours: filename=asnames.json hours=16
    DEBUG:asnlookup.data_manager:maybe_update: func=names output_filename=asnames.json max_age_in_hours=24 result=cached age=16
    DEBUG:asnlookup.backend:reloading databases...
    DEBUG:asnlookup.backend:reloading databases complete seconds=0.6
    8.8.8.8         15169  8.8.8.0/24       GOOGLE - Google Inc.       US
    64.233.177.113  15169  64.233.177.0/24  GOOGLE - Google Inc.       US
    31.13.65.1      32934  31.13.65.0/24    FACEBOOK - Facebook, Inc.  US
    72.30.202.51    26101  72.30.192.0/20   YAHOO-3 - Yahoo!           US

# Server usage

    $ asnlookup-server

# Client usage

Install https://github.com/JustinAzoff/asnlookup-client-python

    echo 8.8.8.8 | asnlookup-client

# Docker

## Using my public image:

    docker run --rm -t -i -v `pwd`/data:/data -p 5555:5555 justinazoff/asnlookup

## Building your own image

    docker build -t asnlookup .
    docker run --rm -t -i -v `pwd`/data:/data -p 5555:5555 asnlookup
