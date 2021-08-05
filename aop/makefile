all: test

CC=gcc
LDFLAGS+=
LDLIBS+=  -lrt -lgcov
CFLAGS=-g -O0 -W -Wall -Wextra -Werror -fprofile-arcs -ftest-coverage  ##-fsanitize=leak -fsanitize=address -fno-omit-frame-pointer -D_FORTIFY_SOURCE=2
SOURCEDIR=.
SOURCES = $(wildcard $(SOURCEDIR)/*.c)
OBJECTS = $(SOURCES:.c=.o)

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

test: $(OBJECTS)
	gcc $(OBJECTS) $(LDFLAGS) $(LDLIBS) -o $@ 

clean:
	-rm -f $(OBJECTS)
	-rm -f test