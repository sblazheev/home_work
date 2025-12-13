-- +goose Up
CREATE TABLE IF NOT EXISTS public.events (
   id uuid NOT NULL,
   title varchar NULL,
   date_time timestamptz NULL,
   duration int4 NULL,
   description varchar NULL,
   "user" int4 NULL,
   notify_time int4 NULL,
   CONSTRAINT events2_pk PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE public.events;
