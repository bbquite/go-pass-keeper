DO $$
    BEGIN
        CREATE TABLE IF NOT EXISTS public.account(
            id serial PRIMARY KEY,
            username VARCHAR (255) UNIQUE NOT NULL,
            password VARCHAR (255) NOT NULL,
            created_on TIMESTAMP NOT NULL default CURRENT_TIMESTAMP
        );

        CREATE TABLE IF NOT EXISTS public.pairs(
            id serial PRIMARY KEY,
            login VARCHAR (255) NOT NULL,
            password VARCHAR (255) NOT NULL,
            meta VARCHAR (255),
            created_on TIMESTAMP NOT NULL default CURRENT_TIMESTAMP,
            FOREIGN KEY (account_id) REFERENCES public.account (id) ON DELETE CASCADE
        );

        CREATE TABLE IF NOT EXISTS public.simple_data(
            id serial PRIMARY KEY,
            text TEXT NOT NULL,
            meta VARCHAR (255),
            created_on TIMESTAMP NOT NULL default CURRENT_TIMESTAMP,
            FOREIGN KEY (account_id) REFERENCES public.account (id) ON DELETE CASCADE
        );

        CREATE TABLE IF NOT EXISTS public.binary_data(
            id serial PRIMARY KEY,
            binary bytea NOT NULL,
            meta VARCHAR (255),
            created_on TIMESTAMP NOT NULL default CURRENT_TIMESTAMP,
            FOREIGN KEY (account_id) REFERENCES public.account (id) ON DELETE CASCADE
        );

        CREATE TABLE IF NOT EXISTS public.card_data(
            id serial PRIMARY KEY,
            card_num VARCHAR (64) NOT NULL,
            card_owner VARCHAR (255) NOT NULL,
            card_exp TIMESTAMP NOT NULL,
            card_cvv smallint NOT NULL,
            meta VARCHAR (255),
            created_on TIMESTAMP NOT NULL default CURRENT_TIMESTAMP,
            FOREIGN KEY (account_id) REFERENCES public.account (id) ON DELETE CASCADE
        );
    END
$$;