# alienvaulturls
Accept line-delimited domains on stdin, fetch known URLs from the AlienVault OTX for `*.domain` and output them on stdout.

Usage example:

```
▶ cat domains.txt | alienvaulturls > urls
```

Install:

```
▶ go get github.com/antuache/alienvaulturls
```

## Credit

This tool was inspired by @tomnomnom's [waybackurls](https://github.com/tomnomnom/waybackurls) script.
