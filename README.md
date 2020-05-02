# fix-nmea
Program to fix nma file of TripMate 850 that went wrong with GPS Week Number Rollover.

This program fixes the date by adding 7168 days (1024 weeks) to the date of the GPRMC line in the NMEA file. This program is valid until around November 21, 2038.
