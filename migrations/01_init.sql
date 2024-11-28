DO $$
    BEGIN
        CREATE TABLE IF NOT EXISTS public.account(
            id serial PRIMARY KEY,
            username VARCHAR (255) UNIQUE NOT NULL,
            password VARCHAR (255) NOT NULL,
            email VARCHAR (255),
            created_on TIMESTAMP NOT NULL default CURRENT_TIMESTAMP
        );

        CREATE TABLE IF NOT EXISTS public.pairs_data(
            id serial PRIMARY KEY,
            account_id integer not null,
            key VARCHAR (255) UNIQUE NOT NULL,
            pwd VARCHAR (255) NOT NULL,
            meta VARCHAR (255),
            uploaded_at TIMESTAMP NOT NULL default CURRENT_TIMESTAMP,
            FOREIGN KEY (account_id) REFERENCES public.account (id) ON DELETE CASCADE
        );

        CREATE TABLE IF NOT EXISTS public.simple_data(
            id serial PRIMARY KEY,
            account_id integer not null,
            text_data TEXT NOT NULL,
            meta VARCHAR (255),
            uploaded_at TIMESTAMP NOT NULL default CURRENT_TIMESTAMP,
            FOREIGN KEY (account_id) REFERENCES public.account (id) ON DELETE CASCADE
        );

        CREATE TABLE IF NOT EXISTS public.binary_data(
            id serial PRIMARY KEY,
            account_id integer not null,
            binary_data bytea NOT NULL,
            meta VARCHAR (255),
            uploaded_at TIMESTAMP NOT NULL default CURRENT_TIMESTAMP,
            FOREIGN KEY (account_id) REFERENCES public.account (id) ON DELETE CASCADE
        );

        CREATE TABLE IF NOT EXISTS public.card_data(
            id serial PRIMARY KEY,
            account_id integer not null,
            card_num VARCHAR (64) NOT NULL,
            card_owner VARCHAR (255) NOT NULL,
            card_exp TIMESTAMP NOT NULL,
            card_cvv smallint NOT NULL,
            meta VARCHAR (255),
            uploaded_at TIMESTAMP NOT NULL default CURRENT_TIMESTAMP,
            FOREIGN KEY (account_id) REFERENCES public.account (id) ON DELETE CASCADE
        );
    END
$$;