/*
Copyright (c) 2016 Nick Potts
Licensed to You under the GNU GPLv3
See the LICENSE file at github.com/npotts/homehub/LICENSE

This file is part of the HomeHub project
*/

/*Package attendants holds various wrappers around various protocols
in order to expose easier access to storage mechanisms.

Mostly attendants monitor some sort of input from remote clients and
via callbacks, send data to some sort of backend, which would be SQL,
NOSQL, binary blobs, or whatever else is desired.

Built Ins

Sub-package http provide a HTTP server that listens on a particular port
and optionally requires a username/password via Basic Auth before forwarding
data on.  Likewise, udp performs a similar service, but over UDP spockets.  It
is planned to also have a ZMQv4 setup
*/
package attendants
