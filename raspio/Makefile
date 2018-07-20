MAKE=make
CC=gcc
CFLAGS=-Wall -ggdb
DIST_CFLAGS=-O2

WIRINGPI_PATH=./wiringPi/wiringPi
WIRINGPI_SONAME=wiringPi
WIRINGPI_LIB=$(WIRINGPI_PATH)/lib$(WIRINGPI_SONAME).so

.PHONY: dist

goldio: goldio.c $(WIRINGPI_LIB)
	$(CC) $(CFLAGS) -L$(WIRINGPI_PATH) -I$(WIRINGPI_PATH) -lrt -l$(WIRINGPI_SONAME) -o $@ $<

$(WIRINGPI_LIB): $(WIRINGPI_PATH)/*.c
	$(MAKE) -C $(WIRINGPI_PATH)
	cd $(WIRINGPI_PATH) && ln -s lib$(WIRINGPI_SONAME).so.* lib$(WIRINGPI_SONAME).so

dist:
	@$(CC) $(DIST_CFLAGS) -L$(WIRINGPI_PATH) -I$(WIRINGPI_PATH) -lrt -l$(WIRINGPI_SONAME) -o goldio goldio.c

clean:
	@rm goldio
