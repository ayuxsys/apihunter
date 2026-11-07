## Usage

```console
$ ./apihunter -h
Usage: apihunter [flags]
  -d, --delay int       delay timeout in seconds to wait for while multiple scraping orgs (default 5)
  -i, --hide-browser    hide browser GUI while scraping (leads to captcha err better use 'xvfb')
  -j, --json            enable jsonl output
  -l, --list string     file containing a list of organisations name to hunt api collection for
  -m, --md              enable markdown output
  -c, --org string      organisation name to hunt api collection for
  -o, --output string   path to the output file (note: can only use one output mode at a time)
  -r, --random string   use random delay between requests in seconds (min,max)
  -v, --verbose         enable verbose logging
```
