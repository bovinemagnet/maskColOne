# Mask Column One

will take an input file and output the file with the first column encrypted.

I wanted a simple 2 way hash via aes256 to encrypt a column of data. This is the result.


```bash
maskColOne --mode=e --in=test/input.tsv --out=test/output.tsv

maskColOne --mode=d --in=test/output.tsv --out=test/decrypted.tsv

```

You can specify a key with the `--key=` option. If you don't specify a key, a built in key will be used (not recommended).

