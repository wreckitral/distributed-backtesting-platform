import yfinance as yf
import sys
import os

if len(sys.argv) < 2:
    print("Usage: python download_data.py SYMBOL")
    sys.exit(1)

symbol = sys.argv[1]

os.makedirs("data/sample", exist_ok=True)

print(f"Downloading {symbol}...")
ticker = yf.Ticker(symbol)
data = ticker.history(period="2y", interval="1d")

data.reset_index(inplace=True)

data = data[['Date', 'Open', 'High', 'Low', 'Close', 'Volume']]

output_path = f"data/sample/{symbol}_daily.csv"
data.to_csv(output_path, index=False)

print(f"✓ Saved to {output_path}")
print(f"  Rows: {len(data)}")
