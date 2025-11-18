-- +goose Up
CREATE TABLE IF NOT EXISTS public."notify" (
   event_id uuid NOT NULL,
   send_time timestamptz NULL,
   status int2 NULL,
   create_time timestamptz NULL,
   CONSTRAINT notify_pk PRIMARY KEY (event_id)
);

-- +goose Down
DROP TABLE public.notify;
