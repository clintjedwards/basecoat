CREATE TABLE IF NOT EXISTS accounts (
    id          TEXT    NOT NULL,
    name        TEXT    NOT NULL,
    hash        TEXT    NOT NULL,
    state       TEXT    NOT NULL,
    created     INTEGER NOT NULL,
    modified    INTEGER NOT NULL,
    PRIMARY KEY (id)
) STRICT;

CREATE TABLE IF NOT EXISTS colorants (
    account      TEXT NOT NULL,
    id           TEXT NOT NULL,
    label        TEXT NOT NULL,
    manufacturer TEXT NOT NULL,
    created      INTEGER NOT NULL,
    PRIMARY KEY (account, id),
    FOREIGN KEY (account) REFERENCES accounts(id) ON DELETE CASCADE
) STRICT;


CREATE TABLE IF NOT EXISTS bases (
    account      TEXT NOT NULL,
    id           TEXT NOT NULL,
    label        TEXT NOT NULL,
    manufacturer TEXT NOT NULL,
    created      INTEGER NOT NULL,
    PRIMARY KEY (account, id),
    FOREIGN KEY (account) REFERENCES accounts(id) ON DELETE CASCADE
) STRICT;

CREATE TABLE IF NOT EXISTS formulas (
    account     TEXT    NOT NULL,
    id          TEXT    NOT NULL,
    name        TEXT    NOT NULL,
    number      TEXT    NOT NULL,
    notes       TEXT    NOT NULL,
    created     INTEGER NOT NULL,
    modified    INTEGER NOT NULL,
    PRIMARY KEY (account, id),
    FOREIGN KEY (account) REFERENCES accounts(id) ON DELETE CASCADE
) STRICT;

CREATE TABLE IF NOT EXISTS formula_colorants (
    account     TEXT NOT NULL,
    formula     TEXT NOT NULL,
    colorant    TEXT NOT NULL,
    amount      TEXT NOT NULL,
    PRIMARY KEY (account, formula, colorant),
    FOREIGN KEY (account) REFERENCES accounts(id) ON DELETE CASCADE,
    FOREIGN KEY (account, formula) REFERENCES formulas(account, id) ON DELETE CASCADE,
    FOREIGN KEY (account, colorant) REFERENCES colorants(account, id) ON DELETE CASCADE
) STRICT;

CREATE TABLE IF NOT EXISTS formula_bases (
    account     TEXT NOT NULL,
    formula     TEXT NOT NULL,
    base        TEXT NOT NULL,
    amount      TEXT NOT NULL,
    PRIMARY KEY (account, formula, base),
    FOREIGN KEY (account) REFERENCES accounts(id) ON DELETE CASCADE,
    FOREIGN KEY (account, formula) REFERENCES formulas(account, id) ON DELETE CASCADE,
    FOREIGN KEY (account, base) REFERENCES bases(account, id) ON DELETE CASCADE
) STRICT;

CREATE TABLE IF NOT EXISTS contractors  (
    account    TEXT NOT NULL,
    id         TEXT NOT NULL,
    company    TEXT NOT NULL,
    contact    TEXT,
    created    INTEGER NOT NULL,
    modified   INTEGER NOT NULL,
    PRIMARY KEY (account, id),
    FOREIGN KEY (account) REFERENCES accounts(id) ON DELETE CASCADE,
    FOREIGN KEY (account, contact) REFERENCES contacts(account, id) ON DELETE SET NULL
) STRICT;

CREATE TABLE IF NOT EXISTS jobs  (
    account    TEXT    NOT NULL,
    id         TEXT    NOT NULL,
    contractor TEXT    NOT NULL,
    name       TEXT    NOT NULL,
    address    TEXT    NOT NULL,
    notes      TEXT    NOT NULL,
    created    INTEGER NOT NULL,
    modified   INTEGER NOT NULL,
    contact    TEXT,
    PRIMARY KEY (account, id),
    FOREIGN KEY (account) REFERENCES accounts(id) ON DELETE CASCADE,
    FOREIGN KEY (account, contractor) REFERENCES contractors(account, id) ON DELETE SET NULL,
    FOREIGN KEY (account, contact) REFERENCES contacts(account, id) ON DELETE SET NULL
) STRICT;

CREATE TABLE IF NOT EXISTS contacts (
    account    TEXT NOT NULL,
    id         TEXT NOT NULL,
    name       TEXT NOT NULL,
    email      TEXT NOT NULL,
    phone      TEXT NOT NULL,
    created    INTEGER NOT NULL,
    modified   INTEGER NOT NULL,
    PRIMARY KEY (account, id),
    FOREIGN KEY (account) REFERENCES accounts(id) ON DELETE CASCADE
) STRICT;

CREATE TABLE IF NOT EXISTS formula_jobs (
    account     TEXT NOT NULL,
    formula     TEXT NOT NULL,
    job         TEXT NOT NULL,
    PRIMARY KEY (account, job, formula),
    FOREIGN KEY (account) REFERENCES accounts(id) ON DELETE CASCADE,
    FOREIGN KEY (account, job) REFERENCES jobs(account, id) ON DELETE CASCADE,
    FOREIGN KEY (account, formula) REFERENCES formulas(account, id) ON DELETE CASCADE
) STRICT;

