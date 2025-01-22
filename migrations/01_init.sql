DO $$
    BEGIN
        CREATE TABLE IF NOT EXISTS public.account(
            id serial PRIMARY KEY,
            username VARCHAR (255) UNIQUE NOT NULL,
            password VARCHAR (255) NOT NULL,
            created_on TIMESTAMP NOT NULL default CURRENT_TIMESTAMP
        );

        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'datatype') THEN
            CREATE TYPE datatype AS ENUM('PAIR','TEXT', 'BINARY', 'CARD');
        END IF;

        CREATE TABLE IF NOT EXISTS public.pass_keeper_data(
            id serial PRIMARY KEY,
            account_id integer not null,
            data_type datatype not null,
            data_info TEXT NOT NULL,
            meta VARCHAR (255),
            uploaded_at TIMESTAMP NOT NULL default CURRENT_TIMESTAMP,
            FOREIGN KEY (account_id) REFERENCES public.account (id) ON DELETE CASCADE
        );
    END
$$;