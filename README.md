[opmlfeed]
==========

opmlfeed is an opmlfeed subscription service. The name is somewhat
misleading now; while it originally served OPML files, now it serves
JSON containing feed information suitable for constructing an OPML
file on the client.


Dependencies
------------

* redis (although you can change the database backend)


Deployment
----------
* Install redis
* Edit env.sh to suit your needs and source it
* go build && ./opmlfeed


Notes
-----

* It would be trivial to adapt it other uses.


License
-------
opmlfeed is licensed under an ISC license. See the LICENSE file for the
full text of the license.
