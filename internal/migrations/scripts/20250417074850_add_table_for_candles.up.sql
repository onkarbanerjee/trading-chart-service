CREATE TABLE IF NOT EXISTS candlesticks (
     symbol      TEXT        NOT NULL,
     interval    TEXT        NOT NULL,
     open_time   TIMESTAMPTZ NOT NULL,
     open        NUMERIC,
     high        NUMERIC,
     low         NUMERIC,
     close       NUMERIC,
     volume      NUMERIC,
     PRIMARY KEY (symbol, interval, open_time)
);
