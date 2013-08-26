# voicemail

`voicemail` intercepts voicemails send by a FRITZ!Box and provides a
web interface for playing them.

The FRITZ!Box provides an option to send voicemails to an email
address. `voicemail` provides a SMTP server, so setup your FRITZ!Box
to send voicemails to `voicemail@<host>`. The integrated SMTP service
is fake and is tested to only work with a FRITZ!Box 7270. The
voicemail's metadata is stored inside a SQLite database and the actual
audio message is converted to a MP3.

The web interface currently works on Android's Browser, Google Chrome
and Safari on iOS.

`voicemail` is currently hardcoded to only provide correct dates
on voicemails in timezone GMT+1 (winter time) and GMT+2 (summer time).
The dates that the FRITZ!Box 7270 sends in a voicemail email lack a
timezone.

## Command line options

There is no config file, all parameters are available as command line
options. Note that to actually receive voicemails from your FRITZ!Box
you need to use `-smtp-port=25`. Output of `voicemail -h`:

    -database="./voicemail.sqlite": Database file location
    -host="localhost": Hostname or IP to bind to
    -http-port="8080": Port for the HTTP service
    -limit=-1: Only display this many voicemails in the web interface
    -smtp-port="2500": Port for the SMTP service
    -user="nobody": User to drop to after binding
    -voicemail="./mp3/": Voicemail storage directory

## License

Copyright © 2012-2013 Tobias Kortkamp

Distributed under the Eclipse Public License.
