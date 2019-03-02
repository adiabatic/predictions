.PHONY: check-copyright-headers clean

check-copyright-headers:
	find . -name \*.go -print -exec head -n 1 {} \;

clean:
	rm -rf dist
