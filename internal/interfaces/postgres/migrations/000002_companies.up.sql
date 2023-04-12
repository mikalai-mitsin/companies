CREATE TABLE public.companies
(
    id                  uuid                 DEFAULT uuid_generate_v4()
        CONSTRAINT companies_pk PRIMARY KEY,
    name                varchar(15) NOT NULL UNIQUE,
    description         text        NOT NULL,
    amount_of_employees int         NOT NULL,
    registered          boolean     NOT NULL,
    type                smallint    NOT NULL,
    updated_at          timestamp   NOT NULL DEFAULT (now() at time zone 'utc'),
    created_at          timestamp   NOT NULL DEFAULT (now() at time zone 'utc')
);

CREATE INDEX companies_registered_filter ON public.companies USING HASH (registered);
CREATE INDEX companies_type_filter ON public.companies USING HASH (type);

CREATE INDEX companies_search
    ON public.companies
        USING GIN (to_tsvector('english', name || description));

CREATE TRIGGER update_companies_updated_at
    BEFORE UPDATE
    ON
        public.companies
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_task();
