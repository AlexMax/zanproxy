zanproxy
========
A proxy detector for Zandronum.

Compiling
---------
To compile this program, you need a working Go environment and a `GOPATH` set up.  Once you do, this project is `go get`table.

    go get github.com/AlexMax/zanproxy

Note that you do not necessarily need a Go environment on the server you wish to use the program on - you can simply compile the program on one machine and copy the binary.  For more information about how to cross-compile with Go, read [this](http://dave.cheney.net/2015/08/22/cross-compilation-with-go-1-5).

Configuration
-------------
To use this program, you first need a [TOML format](https://github.com/toml-lang/toml) configuration file.  Here is a sample configuration:

    Banlist = "/opt/zandronum/list/banlist.txt"
    Email = "alexmax@example.com"
    Logfiles = ["/opt/zandronum/log/*.log"]
    MinScore = 1.0
    BanMesssage = "You have been banned on suspicion of proxy use.  If you believe this is in error, please contact the administrators."

### Banlist
**Banlist** is the path to the banlist that your Zandronum servers use.  The program will automatically append any bans to the bottom of this file, first checking for duplicates.  The program should have write access to this file.

### Email
**Email** is a valid, working, contact e-mail address.  Without one, the API driving this software will refuse to work.

### Logfiles
**Logfiles** is an array of logfiles the program should monitor.  [Globbing](https://golang.org/pkg/path/filepath/#Match) is allowed, but the program will not automatically monitor any logfile that is created after the program starts.

### MinScore
**MinScore** is a score between 0.0 and 1.0 that the program will ban proxies greater than or equal to this number.  Each IP is rated on a scale from 0.0 to 1.0: scores of 1.0 are confirmed proxies, anything less is a likelyhood guesstimate based on a machine learning algorithm.  If you want to be safe, set this to 1.0.  **Use extreme caution when setting this to anything less than 0.99.**  I would recommend starting at 1.0 and then keeping an eye on the output of this program to see if any of your regular players would get caught by your desired score.  For more information on how this score is determined, see [this website](http://getipintel.net/).

### BanMessage
**BanMessage** is the reason the ban will be appended with.  The message will be visible to the suspected player when they try to connect.

Use
---
To run this program, simply pass the path to the configuration file.  It is also recommended that you pipe standard output to a file, as detailed logs are kept of every IP examined, their score, and any action taken.

    $ ./zanproxy zanproxy.cfg > zanproxy.log

License
-------
This program is licensed under the [GNU Affero General Public License v3](https://www.gnu.org/licenses/agpl-3.0.en.html).  If this license is not acceptable to you for some reason, please let me know and perhaps we can work something out.
