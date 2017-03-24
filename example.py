#!/usr/bin/python

import urllib
import urllib2

url = "http://127.0.01:4027/"
key = "xxx-secret-api-key"
username = "username"
headers = {
    'apikey': key, 'Starfish': ' client 1.0', 'Accept': 'application/json',
    'Accept-Charset': 'ISO-8859-1,utf-8;q=0.7,*;q=0.7'
}


def send_sms(num, msg):
    load = {
        'username': username, 'to': num,
        'message': msg, "from": "Test"
    }
    request = urllib2.Request(
        url + 'aft', urllib.urlencode(load), headers)
    return urllib2.urlopen(request).read()


def send_rms_sms(num, msg):
    load = {
        'username': username, 'destination': num,
        'message': msg, 'source': 'Test', 'password': 'pass',
        'dlr': 1, 'type': 0,
    }
    request = urllib2.Request(
        url + 'routesms', urllib.urlencode(load), headers)
    return urllib2.urlopen(request).read()


if __name__ == '__main__':
    import sys
    import json
    try:
        num, msg = sys.argv[1], sys.argv[2]
    except Exception:
        print "Usage: %s <phone> <message>" % (sys.argv[0])
    else:
        print json.loads(send_sms(num, msg))
        print send_rms_sms(num, msg)
