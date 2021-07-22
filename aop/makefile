all: test

CC=gcc
OBJECTS=LDasm.o pinpoint.o
LDFLAGS+=
LDLIBS+=  -lasan
CFLAGS=-g -W -Wall -Wextra -Werror -fsanitize=leak -fsanitize=address -fno-omit-frame-pointer -D_FORTIFY_SOURCE=2

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

test: $(OBJECTS)
	gcc $(OBJECTS) $(LDFLAGS) $(LDLIBS) -o $@ 

clean:
	-rm -f $(OBJECTS)
	-rm -f test