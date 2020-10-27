# alienvaulturls
Accept line-delimited domains on stdin, fetch known URLs from the AlienVault OTX for `*.domain` and output them on stdout.

## Installation:

```
▶ go get github.com/antuache/alienvaulturls
```

## Usage:

Go to https://otx.alienvault.com/settings and grab the "OTX Key".

```
▶ export OTX="insert_api_key_here"

▶ cat domains.txt | alienvaulturls > urls
```

## Credit

This tool was inspired by @tomnomnom's [waybackurls](https://github.com/tomnomnom/waybackurls) script.
