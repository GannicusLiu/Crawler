

RM= rm -rf

all: clean release


release: Weibo

clean:
	$(RM) ../logs/*
	$(RM) Weibo

Weibo:
	go build -ldflags "-s -w" -o $@ ../$@/*.go