all: utest test

CC=gcc
LDFLAGS+=
LDLIBS+=  -lrt -lgcov
CFLAGS= -g3 -O0 -W -Wall -Wextra -Werror -fprofile-arcs -ftest-coverage  -DUTEST ##-fsanitize=leak -fsanitize=address -fno-omit-frame-pointer -D_FORTIFY_SOURCE=2
SOURCEDIR=.
SOURCES = $(wildcard $(SOURCEDIR)/*.c)
OBJECTS = $(SOURCES:.c=.o)

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

utest: $(OBJECTS)
	gcc $(OBJECTS) $(LDFLAGS) $(LDLIBS) -o $@ 

test:utest
	./utest
clean:
	-rm -f $(OBJECTS)
	-rm -f utest
	-rm -f *.gcno
	-rm -f *.gcda