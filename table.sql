CREATE TABLE klines (
    symbol VARCHAR(255),
    interval VARCHAR(255),
    open_time BIGINT,
	open_price VARCHAR(255),
	high_price VARCHAR(255),
	low_price VARCHAR(255),
	close_price VARCHAR(255),
	volume VARCHAR(255),
	close_time BIGINT,
	quote_asset_volume VARCHAR(255),
	number_of_trades INTEGER,
	taker_buy_base_asset_volume VARCHAR(255),
	taker_buy_quote_asset_volume VARCHAR(255)
);

CREATE INDEX idx_symbol_interval ON klines (symbol, interval);