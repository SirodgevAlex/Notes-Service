#!bin/sh

psql -U postgres -d notes_service -f 001_init_users.sql
psql -U postgres -d notes_service -f 002_get_users.sql
psql -U postgres -d notes_service -f 003_init_notes.sql
psql -U postgres -d notes_service -f 004_get_notes.sql